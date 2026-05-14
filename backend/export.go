package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func exportEncrypted(c *fiber.Ctx) error {
	passphrase := strings.TrimSpace(c.Get("X-Export-Passphrase"))
	if len(passphrase) < 12 {
		return c.Status(400).JSON(fiber.Map{"error": "passphrase must be at least 12 characters"})
	}

	rows, err := db.Query("SELECT id, title, content, folder, created_at, updated_at FROM notes ORDER BY updated_at DESC")
	if err != nil {
		logAppEvent("error", "export_encrypted_failed", c, map[string]interface{}{"reason": "db_query_failed", "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Folder, &n.CreatedAt, &n.UpdatedAt); err != nil {
			logAppEvent("error", "export_encrypted_failed", c, map[string]interface{}{"reason": "db_scan_failed", "error": err.Error()})
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		}
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		logAppEvent("error", "export_encrypted_failed", c, map[string]interface{}{"reason": "db_rows_failed", "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	payload, err := json.Marshal(fiber.Map{
		"version":     1,
		"exported_at": time.Now().UTC().Format(time.RFC3339),
		"notes":       notes,
	})
	if err != nil {
		logAppEvent("error", "export_encrypted_failed", c, map[string]interface{}{"reason": "json_marshal_failed", "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		logAppEvent("error", "export_encrypted_failed", c, map[string]interface{}{"reason": "salt_generation_failed", "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		logAppEvent("error", "export_encrypted_failed", c, map[string]interface{}{"reason": "nonce_generation_failed", "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	iterations := 150000
	key := deriveKey(passphrase, salt, iterations)
	block, err := aes.NewCipher(key)
	if err != nil {
		logAppEvent("error", "export_encrypted_failed", c, map[string]interface{}{"reason": "cipher_init_failed", "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logAppEvent("error", "export_encrypted_failed", c, map[string]interface{}{"reason": "gcm_init_failed", "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	ciphertext := gcm.Seal(nil, nonce, payload, nil)
	return c.JSON(fiber.Map{
		"version":    1,
		"algo":       "AES-256-GCM",
		"kdf":        "SHA256-ITER",
		"iterations": iterations,
		"salt_b64":   base64.StdEncoding.EncodeToString(salt),
		"nonce_b64":  base64.StdEncoding.EncodeToString(nonce),
		"data_b64":   base64.StdEncoding.EncodeToString(ciphertext),
	})
}

func deriveKey(passphrase string, salt []byte, iterations int) []byte {
	buf := append([]byte(passphrase), salt...)
	result := sha256.Sum256(buf)
	for i := 1; i < iterations; i++ {
		result = sha256.Sum256(result[:])
	}
	key := make([]byte, 32)
	copy(key, result[:])
	return key
}

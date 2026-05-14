package main

import (
	"crypto/rand"
	"fmt"
	"encoding/hex"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func uploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "no file"})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		return c.Status(400).JSON(fiber.Map{"error": "file type not allowed"})
	}
	if file.Size <= 0 || file.Size > 10*1024*1024 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid file size"})
	}

	opened, err := file.Open()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid file"})
	}
	defer opened.Close()

	buf := make([]byte, 512)
	n, _ := opened.Read(buf)
	mime := http.DetectContentType(buf[:n])
	if !isAllowedImageMime(ext, mime) {
		return c.Status(400).JSON(fiber.Map{"error": "invalid mime type"})
	}

	filename, err := randomFileName(ext)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	dest := filepath.Join("../uploads", filename)
	if err := c.SaveFile(file, dest); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	url := "/uploads/" + filename
	incUploads()
	logAppEvent("info", "image_uploaded", c, map[string]interface{}{
		"file_ext":  ext,
		"file_size": file.Size,
		"url":       url,
	})
	return c.JSON(fiber.Map{"url": url})
}

func randomFileName(ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", hex.EncodeToString(b), ext), nil
}

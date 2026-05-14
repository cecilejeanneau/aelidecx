package main

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func createNote(c *fiber.Ctx) error {
	var input noteInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	if input.Folder == "" {
		input.Folder = "inbox"
	}
	if !validateFolder(input.Folder) {
		return c.Status(400).JSON(fiber.Map{"error": "invalid folder"})
	}

	input.Title = sanitizeTitle(input.Title)
	input.Content = sanitizeContent(input.Content)
	if len(input.Content) > 1_000_000 {
		return c.Status(400).JSON(fiber.Map{"error": "content too large"})
	}

	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	result, err := db.Exec(
		"INSERT INTO notes (title, content, folder, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		input.Title, input.Content, input.Folder, now, now,
	)
	if err != nil {
		logAppEvent("error", "note_create_failed", c, map[string]interface{}{"reason": "db_insert_failed", "folder": input.Folder, "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	id, _ := result.LastInsertId()
	logAppEvent("info", "note_created", c, map[string]interface{}{"note_id": id, "folder": input.Folder})
	return c.Status(201).JSON(fiber.Map{"id": id})
}

func updateNote(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	var input noteInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	if !validateFolder(input.Folder) {
		return c.Status(400).JSON(fiber.Map{"error": "invalid folder"})
	}

	input.Title = sanitizeTitle(input.Title)
	input.Content = sanitizeContent(input.Content)
	if len(input.Content) > 1_000_000 {
		return c.Status(400).JSON(fiber.Map{"error": "content too large"})
	}

	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	result, err := db.Exec(
		"UPDATE notes SET title = ?, content = ?, folder = ?, updated_at = ? WHERE id = ?",
		input.Title, input.Content, input.Folder, now, id,
	)
	if err != nil {
		logAppEvent("error", "note_update_failed", c, map[string]interface{}{"reason": "db_update_failed", "note_id": id, "folder": input.Folder, "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	logAppEvent("info", "note_updated", c, map[string]interface{}{"note_id": id, "folder": input.Folder})
	return c.JSON(fiber.Map{"ok": true})
}

func deleteNote(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	result, err := db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		logAppEvent("error", "note_delete_failed", c, map[string]interface{}{"reason": "db_delete_failed", "note_id": id, "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	logAppEvent("info", "note_deleted", c, map[string]interface{}{"note_id": id})
	return c.JSON(fiber.Map{"ok": true})
}

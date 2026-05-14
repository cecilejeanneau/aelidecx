package main

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func listNotes(c *fiber.Ctx) error {
	folder := c.Query("folder", "")
	var rows *sql.Rows
	var err error

	if folder != "" {
		if !validateFolder(folder) {
			return c.Status(400).JSON(fiber.Map{"error": "invalid folder"})
		}
		rows, err = db.Query("SELECT id, title, content, folder, created_at, updated_at FROM notes WHERE folder = ? ORDER BY updated_at DESC", folder)
	} else {
		rows, err = db.Query("SELECT id, title, content, folder, created_at, updated_at FROM notes ORDER BY updated_at DESC")
	}
	if err != nil {
		logAppEvent("error", "notes_list_failed", c, map[string]interface{}{"reason": "db_query_failed", "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Folder, &n.CreatedAt, &n.UpdatedAt); err != nil {
			logAppEvent("error", "notes_list_failed", c, map[string]interface{}{"reason": "db_scan_failed", "error": err.Error()})
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		}
		n.Content = sanitizeContent(n.Content)
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		logAppEvent("error", "notes_list_failed", c, map[string]interface{}{"reason": "db_rows_failed", "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	return c.JSON(notes)
}

func getNote(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	var n Note
	err = db.QueryRow("SELECT id, title, content, folder, created_at, updated_at FROM notes WHERE id = ?", id).
		Scan(&n.ID, &n.Title, &n.Content, &n.Folder, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		}
		logAppEvent("error", "note_get_failed", c, map[string]interface{}{"reason": "db_query_failed", "note_id": id, "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	n.Content = sanitizeContent(n.Content)
	return c.JSON(n)
}

func searchNotes(c *fiber.Ctx) error {
	q := strings.TrimSpace(c.Query("q", ""))
	if q == "" {
		return c.JSON([]Note{})
	}
	if len(q) > 128 || !searchQueryPattern.MatchString(q) {
		return c.Status(400).JSON(fiber.Map{"error": "invalid search query"})
	}

	rows, err := db.Query(`
		SELECT n.id, n.title, n.content, n.folder, n.created_at, n.updated_at
		FROM notes n
		JOIN notes_fts fts ON n.id = fts.rowid
		WHERE notes_fts MATCH ?
		ORDER BY rank
		LIMIT 50
	`, q)
	if err != nil {
		logAppEvent("warn", "search_failed", c, map[string]interface{}{"reason": "query_or_db_error", "error": err.Error()})
		return c.Status(400).JSON(fiber.Map{"error": "invalid search query"})
	}
	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Folder, &n.CreatedAt, &n.UpdatedAt); err != nil {
			logAppEvent("error", "search_failed", c, map[string]interface{}{"reason": "db_scan_failed", "error": err.Error()})
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		}
		n.Content = sanitizeContent(n.Content)
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		logAppEvent("error", "search_failed", c, map[string]interface{}{"reason": "db_rows_failed", "error": err.Error()})
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	return c.JSON(notes)
}

package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// initDB opens the SQLite database at ../data/notes.db (WAL mode).
// Creates the notes table, FTS5 virtual table, and insert/update/delete sync triggers.
// Calls log.Fatal on any error.
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "../data/notes.db?_journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll("../data", 0755)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			folder TEXT NOT NULL DEFAULT 'inbox',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(title, content, content=notes, content_rowid=id);

		CREATE TRIGGER IF NOT EXISTS notes_ai AFTER INSERT ON notes BEGIN
			INSERT INTO notes_fts(rowid, title, content) VALUES (new.id, new.title, new.content);
		END;
		CREATE TRIGGER IF NOT EXISTS notes_ad AFTER DELETE ON notes BEGIN
			INSERT INTO notes_fts(notes_fts, rowid, title, content) VALUES('delete', old.id, old.title, old.content);
		END;
		CREATE TRIGGER IF NOT EXISTS notes_au AFTER UPDATE ON notes BEGIN
			INSERT INTO notes_fts(notes_fts, rowid, title, content) VALUES('delete', old.id, old.title, old.content);
			INSERT INTO notes_fts(rowid, title, content) VALUES (new.id, new.title, new.content);
		END;
	`)
	if err != nil {
		log.Fatal(err)
	}
}

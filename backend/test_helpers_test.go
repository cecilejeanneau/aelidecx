package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/microcosm-cc/bluemonday"
	_ "github.com/mattn/go-sqlite3"
)

const testKey = "test_secure_key_abcdefghijklmnopqrstuvwxyz012345"

func resetTestSecurityState() {
	authState.mu.Lock()
	authState.m = make(map[string]*authAttemptState)
	authState.mu.Unlock()

	secStats.mu.Lock()
	secStats.APIRequests = 0
	secStats.RateLimited = 0
	secStats.AuthFailures = 0
	secStats.AuthLocked = 0
	secStats.Uploads = 0
	secStats.mu.Unlock()

	obsStats.mu.Lock()
	obsStats.APIRequests = 0
	obsStats.Responses2xx = 0
	obsStats.Responses4xx = 0
	obsStats.Responses5xx = 0
	obsStats.TotalResponseTimeMS = 0
	obsStats.RecentLatenciesMS = nil
	obsStats.Routes = make(map[string]*routeStats)
	obsStats.mu.Unlock()
}

func resetRateLimiter() {
	rateLimiter.mu.Lock()
	rateLimiter.hits = make(map[string][]time.Time)
	rateLimiter.mu.Unlock()
}

func setupTestDB(t *testing.T) {
	t.Helper()
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			folder TEXT NOT NULL DEFAULT 'inbox',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(title, content, content=notes, content_rowid=id)`,
		`CREATE TRIGGER IF NOT EXISTS notes_ai AFTER INSERT ON notes BEGIN
			INSERT INTO notes_fts(rowid, title, content) VALUES (new.id, new.title, new.content);
		END`,
		`CREATE TRIGGER IF NOT EXISTS notes_ad AFTER DELETE ON notes BEGIN
			INSERT INTO notes_fts(notes_fts, rowid, title, content) VALUES('delete', old.id, old.title, old.content);
		END`,
		`CREATE TRIGGER IF NOT EXISTS notes_au AFTER UPDATE ON notes BEGIN
			INSERT INTO notes_fts(notes_fts, rowid, title, content) VALUES('delete', old.id, old.title, old.content);
			INSERT INTO notes_fts(rowid, title, content) VALUES (new.id, new.title, new.content);
		END`,
	}
	for _, s := range stmts {
		if _, err = db.Exec(s); err != nil {
			t.Fatalf("schema exec failed: %v\nSQL: %s", err, s)
		}
	}
	
	t.Cleanup(func() { db.Close() })
}

func newTestApp(t *testing.T) *fiber.App {
	t.Helper()
	setupTestDB(t)
	resetTestSecurityState()
	resetRateLimiter()
	htmlPolicy = bluemonday.UGCPolicy()
	apiKeys = []string{testKey}

	app := fiber.New()
	app.Get("/api/health", healthHandler)
	app.Get("/api/ready", readyHandler)

	api := app.Group("/api")
	api.Use(observabilityMiddleware)
	api.Use(noStoreAPIResponses)
	api.Use(authMiddleware)
	api.Get("/notes", listNotes)
	api.Get("/notes/:id", getNote)
	api.Post("/notes", createNote)
	api.Put("/notes/:id", updateNote)
	api.Delete("/notes/:id", deleteNote)
	api.Get("/search", searchNotes)
	api.Get("/security/metrics", getSecurityMetrics)
	api.Get("/observability/metrics", getObservabilityMetrics)
	api.Get("/export/encrypted", exportEncrypted)
	return app
}

func authedReq(method, url, body string) *http.Request {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, url, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	req.Header.Set("Authorization", "Bearer "+testKey)
	return req
}

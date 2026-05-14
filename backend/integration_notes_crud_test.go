package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestCRUDNotes(t *testing.T) {
	app := newTestApp(t)

	resp, err := app.Test(authedReq(http.MethodPost, "/api/notes", `{"title":"Test Note","content":"<p>Hello</p>","folder":"inbox"}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", resp.StatusCode)
	}
	var created map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&created)
	id := int64(created["id"].(float64))
	if id <= 0 {
		t.Fatalf("create: expected valid id, got %v", id)
	}

	resp, err = app.Test(authedReq(http.MethodGet, "/api/notes", ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", resp.StatusCode)
	}
	var notes []Note
	json.NewDecoder(resp.Body).Decode(&notes)
	if len(notes) != 1 {
		t.Fatalf("list: expected 1 note, got %d", len(notes))
	}

	resp, err = app.Test(authedReq(http.MethodGet, fmt.Sprintf("/api/notes/%d", id), ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", resp.StatusCode)
	}
	var got Note
	json.NewDecoder(resp.Body).Decode(&got)
	if got.Title != "Test Note" {
		t.Fatalf("get: expected title 'Test Note', got %q", got.Title)
	}

	resp, err = app.Test(authedReq(http.MethodPut, fmt.Sprintf("/api/notes/%d", id), `{"title":"Updated","content":"<p>World</p>","folder":"concepts"}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update: expected 200, got %d", resp.StatusCode)
	}

	resp, err = app.Test(authedReq(http.MethodGet, fmt.Sprintf("/api/notes/%d", id), ""))
	if err != nil {
		t.Fatal(err)
	}
	var updated Note
	json.NewDecoder(resp.Body).Decode(&updated)
	if updated.Title != "Updated" || updated.Folder != "concepts" {
		t.Fatalf("get after update: unexpected note %+v", updated)
	}

	resp, err = app.Test(authedReq(http.MethodDelete, fmt.Sprintf("/api/notes/%d", id), ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", resp.StatusCode)
	}

	resp, err = app.Test(authedReq(http.MethodGet, fmt.Sprintf("/api/notes/%d", id), ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("get after delete: expected 404, got %d", resp.StatusCode)
	}
}

func TestCreateNoteDefaultFolder(t *testing.T) {
	app := newTestApp(t)

	resp, err := app.Test(authedReq(http.MethodPost, "/api/notes", `{"title":"No folder","content":""}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var created map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&created)
	id := int64(created["id"].(float64))

	resp, err = app.Test(authedReq(http.MethodGet, fmt.Sprintf("/api/notes/%d", id), ""))
	if err != nil {
		t.Fatal(err)
	}
	var note Note
	json.NewDecoder(resp.Body).Decode(&note)
	if note.Folder != "inbox" {
		t.Fatalf("expected default folder 'inbox', got %q", note.Folder)
	}
}

func TestNotesInvalidInputs(t *testing.T) {
	app := newTestApp(t)

	cases := []struct {
		name   string
		method string
		url    string
		body   string
		want   int
	}{
		{"get invalid id", http.MethodGet, "/api/notes/abc", "", http.StatusBadRequest},
		{"get not found", http.MethodGet, "/api/notes/9999", "", http.StatusNotFound},
		{"delete invalid id", http.MethodDelete, "/api/notes/abc", "", http.StatusBadRequest},
		{"delete not found", http.MethodDelete, "/api/notes/9999", "", http.StatusNotFound},
		{"update invalid id", http.MethodPut, "/api/notes/abc", `{"title":"x","content":"","folder":"inbox"}`, http.StatusBadRequest},
		{"update not found", http.MethodPut, "/api/notes/9999", `{"title":"x","content":"","folder":"inbox"}`, http.StatusNotFound},
		{"create invalid folder", http.MethodPost, "/api/notes", `{"title":"x","content":"","folder":"evil"}`, http.StatusBadRequest},
		{"list invalid folder", http.MethodGet, "/api/notes?folder=admin", "", http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := app.Test(authedReq(tc.method, tc.url, tc.body))
			if err != nil {
				t.Fatalf("app.Test failed: %v", err)
			}
			if resp.StatusCode != tc.want {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("expected %d, got %d (body: %s)", tc.want, resp.StatusCode, body)
			}
		})
	}
}

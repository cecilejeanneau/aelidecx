package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestListNotesByFolder(t *testing.T) {
	app := newTestApp(t)

	app.Test(authedReq(http.MethodPost, "/api/notes", `{"title":"A","content":"","folder":"inbox"}`))
	app.Test(authedReq(http.MethodPost, "/api/notes", `{"title":"B","content":"","folder":"concepts"}`))

	resp, err := app.Test(authedReq(http.MethodGet, "/api/notes?folder=inbox", ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var notes []Note
	json.NewDecoder(resp.Body).Decode(&notes)
	if len(notes) != 1 || notes[0].Title != "A" {
		t.Fatalf("expected 1 inbox note named A, got %+v", notes)
	}
}

func TestSearchNotes(t *testing.T) {
	app := newTestApp(t)

	app.Test(authedReq(http.MethodPost, "/api/notes", `{"title":"Golang tutorial","content":"<p>Learn Go programming</p>","folder":"inbox"}`))

	resp, err := app.Test(authedReq(http.MethodGet, "/api/search?q=", ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("empty search: expected 200, got %d", resp.StatusCode)
	}

	resp, err = app.Test(authedReq(http.MethodGet, "/api/search?q=%3Cscript%3Ealert(1)%3C%2Fscript%3E", ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid search: expected 400, got %d", resp.StatusCode)
	}

	longQ := strings.Repeat("a", 129)
	resp, err = app.Test(authedReq(http.MethodGet, "/api/search?q="+longQ, ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("long search: expected 400, got %d", resp.StatusCode)
	}

	resp, err = app.Test(authedReq(http.MethodGet, "/api/search?q=Golang", ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("valid search: expected 200, got %d", resp.StatusCode)
	}
}

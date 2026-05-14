package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHealthAndReadinessProbes(t *testing.T) {
	app := newTestApp(t)

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("health: expected 200, got %d", resp.StatusCode)
	}
	var health map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&health)
	if health["status"] != "ok" {
		t.Fatalf("health: expected status=ok, got %v", health["status"])
	}

	resp, err = app.Test(httptest.NewRequest(http.MethodGet, "/api/ready", nil))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("ready: expected 200, got %d", resp.StatusCode)
	}
	var ready map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&ready)
	if ready["status"] != "ready" {
		t.Fatalf("ready: expected status=ready, got %v", ready["status"])
	}
}

func TestReadinessProbeWhenDBMissing(t *testing.T) {
	app := fiber.New()
	app.Get("/api/ready", readyHandler)

	originalDB := db
	db = nil
	t.Cleanup(func() { db = originalDB })

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/ready", nil))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", resp.StatusCode)
	}
	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	if body["status"] != "not_ready" {
		t.Fatalf("expected status=not_ready, got %v", body["status"])
	}
}

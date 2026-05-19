package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetSecurityMetrics(t *testing.T) {
	app := newTestApp(t)

	resp, err := app.Test(authedReq(http.MethodGet, "/api/security/metrics", ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	for _, field := range []string{"api_requests", "rate_limited", "auth_failures", "auth_locked", "uploads", "timestamp"} {
		if _, ok := result[field]; !ok {
			t.Errorf("expected field %q in security metrics response", field)
		}
	}
}

func TestObservabilityMetrics(t *testing.T) {
	app := newTestApp(t)

	resp, err := app.Test(authedReq(http.MethodGet, "/api/notes", ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from notes list, got %d", resp.StatusCode)
	}
	if resp.Header.Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID header on API response")
	}

	resp, err = app.Test(authedReq(http.MethodGet, "/api/observability/metrics", ""))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	for _, field := range []string{"api_requests", "responses_2xx", "responses_4xx", "responses_5xx", "avg_response_time_ms", "p50_response_time_ms", "p95_response_time_ms", "p99_response_time_ms", "routes", "system", "timestamp"} {
		if _, ok := result[field]; !ok {
			t.Errorf("expected field %q in observability metrics response", field)
		}
	}
	if result["api_requests"].(float64) < 1 {
		t.Fatalf("expected at least 1 observed API request, got %v", result["api_requests"])
	}
	if result["responses_2xx"].(float64) < 1 {
		t.Fatalf("expected at least 1 successful response, got %v", result["responses_2xx"])
	}

	routes, ok := result["routes"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected routes map in observability metrics, got %T", result["routes"])
	}
	if _, ok := routes["GET /api/notes"]; !ok {
		t.Fatalf("expected route metrics for GET /api/notes, got %v", routes)
	}
	system, ok := result["system"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected system metrics map in observability metrics, got %T", result["system"])
	}
	for _, field := range []string{"uptime_seconds", "goroutines", "go_arch", "go_os", "db_path", "db_size_bytes", "configured_api_keys"} {
		if _, ok := system[field]; !ok {
			t.Fatalf("expected system field %q in observability metrics, got %v", field, system)
		}
	}
}

func TestExportEncrypted(t *testing.T) {
	app := newTestApp(t)

	req := authedReq(http.MethodGet, "/api/export/encrypted", "")
	req.Header.Set("X-Export-Passphrase", "short")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("short passphrase: expected 400, got %d", resp.StatusCode)
	}

	req = authedReq(http.MethodGet, "/api/export/encrypted", "")
	req.Header.Set("X-Export-Passphrase", "my-long-secure-passphrase-123!")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("valid export: expected 200, got %d", resp.StatusCode)
	}

	var payload map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&payload)
	for _, field := range []string{"salt_b64", "nonce_b64", "data_b64", "algo", "kdf", "iterations"} {
		if payload[field] == nil {
			t.Errorf("expected field %q in encrypted export, got nil", field)
		}
	}
	if payload["algo"] != "AES-256-GCM" {
		t.Fatalf("expected algo AES-256-GCM, got %v", payload["algo"])
	}
}

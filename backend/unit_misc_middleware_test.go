package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(securityHeaders())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.Header.Get("X-Frame-Options") != "SAMEORIGIN" {
		t.Fatalf("expected X-Frame-Options to be set")
	}
	if !strings.Contains(resp.Header.Get("Content-Security-Policy"), "script-src 'self'") {
		t.Fatalf("expected CSP script-src self policy")
	}
}

func TestRandomFileNameFormat(t *testing.T) {
	for _, ext := range []string{".png", ".jpg", ".jpeg", ".gif", ".webp"} {
		name, err := randomFileName(ext)
		if err != nil {
			t.Fatalf("ext %s: unexpected error: %v", ext, err)
		}
		if !strings.HasSuffix(name, ext) {
			t.Fatalf("ext %s: expected suffix, got %s", ext, name)
		}
		hexPart := strings.TrimSuffix(name, ext)
		if len(hexPart) != 32 {
			t.Fatalf("ext %s: expected 32 hex chars, got %d (%s)", ext, len(hexPart), hexPart)
		}
	}
	n1, _ := randomFileName(".png")
	n2, _ := randomFileName(".png")
	if n1 == n2 {
		t.Fatal("expected unique filenames, got identical")
	}
}

func TestNoStoreAPIResponsesHeader(t *testing.T) {
	app := fiber.New()
	app.Use(noStoreAPIResponses)
	app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("expected Cache-Control: no-store, got %q", resp.Header.Get("Cache-Control"))
	}
}

func TestRateLimitMiddlewareBlocks(t *testing.T) {
	resetRateLimiter()
	app := fiber.New()
	app.Use(rateLimitMiddleware)
	app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req0 := httptest.NewRequest(http.MethodGet, "/", nil)
	if _, err := app.Test(req0); err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	rateLimiter.mu.Lock()
	now := time.Now()
	for ip := range rateLimiter.hits {
		full := make([]time.Time, rateLimiter.limit)
		for i := range full {
			full[i] = now
		}
		rateLimiter.hits[ip] = full
	}
	rateLimiter.mu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("expected 429 after rate limit exceeded, got %d", resp.StatusCode)
	}
}

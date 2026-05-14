package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestAuthMiddlewareUnauthorizedForbiddenAndSuccess(t *testing.T) {
	resetTestSecurityState()
	apiKeys = []string{"super_secure_key_abcdefghijklmnopqrstuvwxyz"}

	app := fiber.New()
	app.Use(authMiddleware)
	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req1 := httptest.NewRequest(http.MethodGet, "/ok", nil)
	resp1, err := app.Test(req1)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp1.StatusCode != fiber.Status401Unauthorized {
		t.Fatalf("expected 401, got %d", resp1.StatusCode)
	}

	resetTestSecurityState()
	req2 := httptest.NewRequest(http.MethodGet, "/ok", nil)
	req2.Header.Set("Authorization", "Bearer wrong_key")
	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp2.StatusCode != fiber.Status403Forbidden {
		t.Fatalf("expected 403, got %d", resp2.StatusCode)
	}

	resetTestSecurityState()
	req3 := httptest.NewRequest(http.MethodGet, "/ok", nil)
	req3.Header.Set("Authorization", "Bearer super_secure_key_abcdefghijklmnopqrstuvwxyz")
	resp3, err := app.Test(req3)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp3.StatusCode != fiber.Status200OK {
		body, _ := io.ReadAll(resp3.Body)
		t.Fatalf("expected 200, got %d with body %s", resp3.StatusCode, string(body))
	}
}

func TestAuthBackoffLocksAfterFailure(t *testing.T) {
	resetTestSecurityState()
	apiKeys = []string{"super_secure_key_abcdefghijklmnopqrstuvwxyz"}

	app := fiber.New()
	app.Use(authMiddleware)
	app.Get("/ok", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req1 := httptest.NewRequest(http.MethodGet, "/ok", nil)
	req1.Header.Set("Authorization", "Bearer wrong_key")
	resp1, err := app.Test(req1)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp1.StatusCode != fiber.Status403Forbidden {
		t.Fatalf("expected first invalid token response 403, got %d", resp1.StatusCode)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/ok", nil)
	req2.Header.Set("Authorization", "Bearer wrong_key")
	resp2, err := app.Test(req2)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}
	if resp2.StatusCode != fiber.Status429TooManyRequests {
		t.Fatalf("expected lockout response 429, got %d", resp2.StatusCode)
	}
}

func TestAuthBackoffDuration(t *testing.T) {
	cases := []struct {
		failures int
		expected time.Duration
	}{
		{0, 0},
		{1, 1 * time.Second},
		{2, 2 * time.Second},
		{3, 4 * time.Second},
		{4, 8 * time.Second},
		{5, 16 * time.Second},
		{6, 32 * time.Second},
		{7, 64 * time.Second},
		{8, 128 * time.Second},
		{9, 128 * time.Second},
		{100, 128 * time.Second},
	}
	for _, tc := range cases {
		got := authBackoffDuration(tc.failures)
		if got != tc.expected {
			t.Errorf("failures=%d: expected %v, got %v", tc.failures, tc.expected, got)
		}
	}
}

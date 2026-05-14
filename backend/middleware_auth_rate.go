package main

import (
	"crypto/subtle"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func authMiddleware(c *fiber.Ctx) error {
	ip := c.IP()
	incAPIRequests()

	if retryAfter, locked := isAuthLocked(ip); locked {
		incAuthLocked()
		c.Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
		logAppEvent("warn", "auth_locked", c, map[string]interface{}{
			"ip":              ip,
			"retry_after_sec": int(retryAfter.Seconds()),
		})
		return c.Status(429).JSON(fiber.Map{"error": "too many authentication failures"})
	}

	authHeader := strings.TrimSpace(c.Get("Authorization"))
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		logAppEvent("warn", "auth_failure", c, map[string]interface{}{"ip": ip, "reason": "missing_or_malformed_bearer_header"})
		recordAuthFailure(ip)
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if token == "" {
		logAppEvent("warn", "auth_failure", c, map[string]interface{}{"ip": ip, "reason": "empty_bearer_token"})
		recordAuthFailure(ip)
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	if !isValidAPIKey(token) {
		logAppEvent("warn", "auth_failure", c, map[string]interface{}{"ip": ip, "reason": "invalid_token"})
		recordAuthFailure(ip)
		return c.Status(403).JSON(fiber.Map{"error": "invalid token"})
	}
	resetAuthFailure(ip)
	return c.Next()
}

func rateLimitMiddleware(c *fiber.Ctx) error {
	now := time.Now()
	ip := c.IP()

	rateLimiter.mu.Lock()
	defer rateLimiter.mu.Unlock()

	history := rateLimiter.hits[ip]
	cutoff := now.Add(-rateLimiter.win)
	filtered := history[:0]
	for _, t := range history {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= rateLimiter.limit {
		incRateLimited()
		return c.Status(429).JSON(fiber.Map{"error": "too many requests"})
	}

	rateLimiter.hits[ip] = append(filtered, now)
	return c.Next()
}

func isValidAPIKey(token string) bool {
	for _, k := range apiKeys {
		if subtle.ConstantTimeCompare([]byte(token), []byte(k)) == 1 {
			return true
		}
	}
	return false
}

func recordAuthFailure(ip string) {
	incAuthFailures()
	authState.mu.Lock()
	defer authState.mu.Unlock()

	s, ok := authState.m[ip]
	if !ok {
		s = &authAttemptState{}
		authState.m[ip] = s
	}

	s.Failures++
	s.LockedUntil = time.Now().Add(authBackoffDuration(s.Failures))
}

func resetAuthFailure(ip string) {
	authState.mu.Lock()
	defer authState.mu.Unlock()
	delete(authState.m, ip)
}

func isAuthLocked(ip string) (time.Duration, bool) {
	authState.mu.Lock()
	defer authState.mu.Unlock()

	s, ok := authState.m[ip]
	if !ok {
		return 0, false
	}
	if now := time.Now(); now.Before(s.LockedUntil) {
		return s.LockedUntil.Sub(now), true
	}
	return 0, false
}

func authBackoffDuration(failures int) time.Duration {
	if failures < 1 {
		return 0
	}
	if failures > 8 {
		failures = 8
	}
	return time.Duration(1<<(failures-1)) * time.Second
}

func getConfiguredAPIKeys() []string {
	if raw := strings.TrimSpace(os.Getenv("API_KEYS")); raw != "" {
		parts := strings.Split(raw, ",")
		keys := make([]string, 0, len(parts))
		for _, p := range parts {
			k := strings.TrimSpace(p)
			if k != "" {
				keys = append(keys, k)
			}
		}
		return keys
	}

	if single := strings.TrimSpace(os.Getenv("API_KEY")); single != "" {
		return []string{single}
	}
	return nil
}

package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func getSecurityMetrics(c *fiber.Ctx) error {
	secStats.mu.Lock()
	defer secStats.mu.Unlock()

	return c.JSON(fiber.Map{
		"api_requests":  secStats.APIRequests,
		"rate_limited":  secStats.RateLimited,
		"auth_failures": secStats.AuthFailures,
		"auth_locked":   secStats.AuthLocked,
		"uploads":       secStats.Uploads,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
	})
}

func incAPIRequests() {
	secStats.mu.Lock()
	secStats.APIRequests++
	secStats.mu.Unlock()
}

func incRateLimited() {
	secStats.mu.Lock()
	secStats.RateLimited++
	secStats.mu.Unlock()
}

func incAuthFailures() {
	secStats.mu.Lock()
	secStats.AuthFailures++
	secStats.mu.Unlock()
}

func incAuthLocked() {
	secStats.mu.Lock()
	secStats.AuthLocked++
	secStats.mu.Unlock()
}

func incUploads() {
	secStats.mu.Lock()
	secStats.Uploads++
	secStats.mu.Unlock()
}

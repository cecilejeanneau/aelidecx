package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// healthHandler returns a minimal liveness response used by probes.
func healthHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// readyHandler returns readiness status based on critical local dependencies.
func readyHandler(c *fiber.Ctx) error {
	if db == nil {
		logAppEvent("error", "readiness_failed", c, map[string]interface{}{"reason": "db_nil"})
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "not_ready"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logAppEvent("error", "readiness_failed", c, map[string]interface{}{"reason": "db_ping_failed", "error": err.Error()})
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "not_ready"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ready"})
}

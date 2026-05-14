package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func securityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; object-src 'none'; base-uri 'self'; frame-ancestors 'self'; form-action 'self'")
		c.Set("Referrer-Policy", "no-referrer")
		c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		if c.Protocol() == "https" {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		return c.Next()
	}
}

func buildCORSConfig() cors.Config {
	origin := os.Getenv("ALLOWED_ORIGINS")
	if origin == "" {
		origin = "http://localhost:3000,http://127.0.0.1:3000"
	}

	return cors.Config{
		AllowOrigins: origin,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Authorization, Content-Type",
	}
}

func noStoreAPIResponses(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")
	return c.Next()
}

// @title          Aelidecx API
// @version        1.0
// @description    Self-hosted personal knowledge base. All routes under /api require Bearer authentication.
// @host           localhost:3000
// @BasePath       /api
// @securityDefinitions.apikey BearerAuth
// @in             header
// @name           Authorization
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/microcosm-cc/bluemonday"
	fiberswagger "github.com/swaggo/fiber-swagger"
	_ "aelidecx/docs"
)

func main() {
	initDB()
	defer db.Close()
	htmlPolicy = bluemonday.UGCPolicy()

	apiKeys = getConfiguredAPIKeys()
	if len(apiKeys) == 0 {
		if strings.EqualFold(os.Getenv("APP_ENV"), "development") || strings.EqualFold(os.Getenv("APP_ENV"), "dev") {
			apiKeys = []string{"dev-key-change-me-now"}
			log.Println("WARNING: API_KEY not set. Development fallback key is active. Do not use this outside local dev.")
		} else {
			log.Fatal("API_KEY or API_KEYS must be set in environment")
		}
	}
	for _, key := range apiKeys {
		if len(key) < 24 {
			log.Fatal("All API keys must be at least 24 characters")
		}
	}

	os.MkdirAll("../uploads", 0755)
	app := newApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	bindAddr := os.Getenv("BIND_ADDR")
	if bindAddr == "" {
		bindAddr = "127.0.0.1"
	}

	listenAddr := bindAddr + ":" + port
	tlsCert := os.Getenv("TLS_CERT_FILE")
	tlsKey := os.Getenv("TLS_KEY_FILE")
	if tlsCert != "" && tlsKey != "" {
		fmt.Printf("Server running on https://%s:%s\n", bindAddr, port)
		log.Fatal(app.ListenTLS(listenAddr, tlsCert, tlsKey))
		return
	}

	fmt.Printf("Server running on http://%s:%s\n", bindAddr, port)
	log.Fatal(app.Listen(listenAddr))
}

func newApp() *fiber.App {
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			if code >= 500 {
				return c.Status(code).JSON(fiber.Map{"error": "internal server error"})
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(logger.New())
	app.Use(securityHeaders())
	app.Use(cors.New(buildCORSConfig()))
	app.Static("/", "../frontend")
	app.Static("/uploads", "../uploads")
	app.Get("/swagger/*", fiberswagger.WrapHandler)
	app.Get("/api/health", healthHandler)
	app.Get("/api/ready", readyHandler)

	api := app.Group("/api")
	api.Use(observabilityMiddleware)
	api.Use(rateLimitMiddleware)
	api.Use(noStoreAPIResponses)
	api.Use(authMiddleware)
	api.Get("/notes", listNotes)
	api.Get("/notes/:id", getNote)
	api.Post("/notes", createNote)
	api.Put("/notes/:id", updateNote)
	api.Delete("/notes/:id", deleteNote)
	api.Get("/search", searchNotes)
	api.Post("/upload", uploadImage)
	api.Get("/security/metrics", getSecurityMetrics)
	api.Get("/observability/metrics", getObservabilityMetrics)
	api.Get("/export/encrypted", exportEncrypted)

	return app
}

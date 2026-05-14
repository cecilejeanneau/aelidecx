package main

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
)

func observabilityMiddleware(c *fiber.Ctx) error {
	requestID := randomRequestID()
	c.Set("X-Request-ID", requestID)
	c.Locals("request_id", requestID)

	start := time.Now()
	err := c.Next()
	durationMS := float64(time.Since(start).Microseconds()) / 1000.0
	statusCode := c.Response().StatusCode()
	if statusCode == 0 {
		statusCode = fiber.StatusOK
	}

	route := routeLabel(c)
	recordObservability(c.Method(), route, statusCode, durationMS)
	logRequestEvent(requestID, c.IP(), c.Method(), c.Path(), route, statusCode, durationMS)
	return err
}

func routeLabel(c *fiber.Ctx) string {
	if route := c.Route(); route != nil && route.Path != "" {
		return route.Path
	}
	return c.Path()
}

func getObservabilityMetrics(c *fiber.Ctx) error {
	obsStats.mu.Lock()
	defer obsStats.mu.Unlock()

	avg := 0.0
	if obsStats.APIRequests > 0 {
		avg = obsStats.TotalResponseTimeMS / float64(obsStats.APIRequests)
	}

	latencies := append([]float64(nil), obsStats.RecentLatenciesMS...)
	routes := make(map[string]fiber.Map, len(obsStats.Routes))
	for route, stats := range obsStats.Routes {
		routeAvg := 0.0
		if stats.Count > 0 {
			routeAvg = stats.TotalResponseTimeMS / float64(stats.Count)
		}
		routes[route] = fiber.Map{
			"count":                stats.Count,
			"responses_2xx":        stats.Responses2xx,
			"responses_4xx":        stats.Responses4xx,
			"responses_5xx":        stats.Responses5xx,
			"avg_response_time_ms": routeAvg,
		}
	}

	return c.JSON(fiber.Map{
		"api_requests":         obsStats.APIRequests,
		"responses_2xx":        obsStats.Responses2xx,
		"responses_4xx":        obsStats.Responses4xx,
		"responses_5xx":        obsStats.Responses5xx,
		"avg_response_time_ms": avg,
		"p50_response_time_ms": percentile(latencies, 0.50),
		"p95_response_time_ms": percentile(latencies, 0.95),
		"p99_response_time_ms": percentile(latencies, 0.99),
		"routes":               routes,
		"system":               collectSystemMetrics(),
		"timestamp":            time.Now().UTC().Format(time.RFC3339),
	})
}

func collectSystemMetrics() fiber.Map {
	uptimeSeconds := int64(time.Since(startedAt).Seconds())
	if uptimeSeconds < 0 {
		uptimeSeconds = 0
	}

	dbPath := filepath.Clean("../data/notes.db")
	dbSizeBytes := int64(0)
	if info, err := os.Stat(dbPath); err == nil {
		dbSizeBytes = info.Size()
	}

	return fiber.Map{
		"uptime_seconds":      uptimeSeconds,
		"goroutines":          runtime.NumGoroutine(),
		"go_arch":             runtime.GOARCH,
		"go_os":               runtime.GOOS,
		"db_path":             dbPath,
		"db_size_bytes":       dbSizeBytes,
		"configured_api_keys": len(apiKeys),
	}
}

func recordObservability(method, route string, statusCode int, durationMS float64) {
	obsStats.mu.Lock()
	obsStats.APIRequests++
	obsStats.TotalResponseTimeMS += durationMS
	obsStats.RecentLatenciesMS = append(obsStats.RecentLatenciesMS, durationMS)
	if len(obsStats.RecentLatenciesMS) > 512 {
		obsStats.RecentLatenciesMS = append([]float64(nil), obsStats.RecentLatenciesMS[len(obsStats.RecentLatenciesMS)-512:]...)
	}

	key := method + " " + route
	stats, ok := obsStats.Routes[key]
	if !ok {
		stats = &routeStats{}
		obsStats.Routes[key] = stats
	}
	stats.Count++
	stats.TotalResponseTimeMS += durationMS

	switch {
	case statusCode >= 200 && statusCode < 300:
		obsStats.Responses2xx++
		stats.Responses2xx++
	case statusCode >= 400 && statusCode < 500:
		obsStats.Responses4xx++
		stats.Responses4xx++
	case statusCode >= 500:
		obsStats.Responses5xx++
		stats.Responses5xx++
	}
	obsStats.mu.Unlock()
}

func percentile(values []float64, ratio float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)
	index := int(float64(len(sorted)-1) * ratio)
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

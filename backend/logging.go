package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func logRequestEvent(requestID, ip, method, path, route string, statusCode int, durationMS float64) {
	level := "info"
	if statusCode >= 500 {
		level = "error"
	} else if statusCode >= 400 {
		level = "warn"
	}
	logJSON(level, "api_request", map[string]interface{}{
		"request_id":  requestID,
		"ip":          ip,
		"method":      method,
		"path":        path,
		"route":       route,
		"status":      statusCode,
		"duration_ms": durationMS,
	})
}

func logAppEvent(level, event string, c *fiber.Ctx, fields map[string]interface{}) {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	fields["request_id"] = requestIDFromCtx(c)
	fields["method"] = c.Method()
	fields["path"] = c.Path()
	fields["route"] = routeLabel(c)
	logJSON(level, event, fields)
}

func requestIDFromCtx(c *fiber.Ctx) string {
	if c == nil {
		return ""
	}
	if requestID, ok := c.Locals("request_id").(string); ok {
		return requestID
	}
	return c.Get("X-Request-ID")
}

func logJSON(level, event string, fields map[string]interface{}) {
	payload := map[string]interface{}{
		"level":     level,
		"event":     event,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	for key, value := range fields {
		payload[key] = value
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		log.Printf("level=error event=log_encode_failed original_event=%s err=%v", event, err)
		return
	}
	log.Println(string(encoded))
}

func randomRequestID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	return hex.EncodeToString(buf)
}

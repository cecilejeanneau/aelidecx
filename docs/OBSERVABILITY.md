# Observability - Aelidecx

## Objective

This document describes what the project currently exposes for logs and metrics, how to use it for debugging and lightweight operations, and what can be improved next.

## 1. Existing logs

### 1.1 HTTP logs

The backend enables Fiber's `logger` middleware for all HTTP requests.

At minimum, it logs:

- HTTP method
- requested URL
- response status
- latency
- client address (Fiber default behavior)

These logs are emitted to standard output.

### 1.2 Structured API request logs

API requests also emit a structured JSON log entry to standard output.

Current fields:

- `level`: `info`, `warn`, or `error` based on response class
- `event`: event type (currently `api_request`)
- `timestamp`: UTC RFC3339 timestamp
- `request_id`: correlation id
- `ip`: client IP
- `method`: HTTP method
- `path`: requested path
- `route`: logical Fiber route
- `status`: final HTTP status code
- `duration_ms`: request duration

### 1.3 Structured application events

In addition to generic request logs, sensitive actions emit dedicated events.

Currently covered examples:

- `note_created`
- `note_updated`
- `note_deleted`
- `image_uploaded`
- `auth_failure`
- `auth_locked`
- `note_get_failed`, `notes_list_failed`, `search_failed`
- `note_create_failed`, `note_update_failed`, `note_delete_failed`
- `export_encrypted_failed`

These events stay intentionally compact and avoid sensitive note content or secrets.

Example:

```json
{
  "level": "info",
  "event": "api_request",
  "timestamp": "2026-05-07T19:10:00Z",
  "request_id": "5c1d2a9f4b7e31aa",
  "ip": "127.0.0.1",
  "method": "GET",
  "path": "/api/notes",
  "route": "/notes",
  "status": 200,
  "duration_ms": 2.31
}
```

### 1.4 Startup logs

At startup, the backend logs:

- warning when a development fallback API key is used
- fatal error when no API key is configured outside dev mode
- fatal error when an API key is too short
- HTTP or HTTPS bind address announcement
- fatal SQLite initialization errors

### 1.5 Current limits

Logs are not yet:

- persisted with file rotation
- centralized to an external aggregator

Current observability is intentionally simple and local-first.

## 2. Existing metrics

### 2.1 Security endpoint

The backend exposes:

- `GET /api/security/metrics`

This endpoint is protected by Bearer authentication.

### 2.2 Exposed security counters

- `api_requests`: total API requests that reached auth middleware
- `rate_limited`: requests rejected by rate limiting
- `auth_failures`: authentication failures
- `auth_locked`: temporary lock rejections
- `uploads`: successful image uploads
- `timestamp`: UTC snapshot timestamp

These counters are:

- in-memory
- reset on process restart
- useful for lightweight operations
- not persisted by default

### 2.3 API observability metrics endpoint

The backend also exposes:

- `GET /api/observability/metrics`

Current response includes:

- `api_requests`
- `responses_2xx`
- `responses_4xx`
- `responses_5xx`
- `avg_response_time_ms`
- `p50_response_time_ms`
- `p95_response_time_ms`
- `p99_response_time_ms`
- `routes`: counters and average latency per logical route
- `system`: basic runtime/system metrics
- `timestamp`

Percentiles are computed on a recent in-memory latency window.

Current `system` block includes:

- `uptime_seconds`
- `goroutines`
- `go_arch`
- `go_os`
- `db_path`
- `db_size_bytes`
- `configured_api_keys`

### 2.4 Health probes

For quick operational checks, the backend exposes two unauthenticated probes:

- `GET /api/health`: liveness probe (`status: ok`)
- `GET /api/ready`: readiness probe based on SQLite (`status: ready` or `status: not_ready`)

### 2.5 Quick curl examples

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" http://localhost:3000/api/security/metrics
curl -H "Authorization: Bearer YOUR_API_KEY" http://localhost:3000/api/observability/metrics
curl -i http://localhost:3000/api/health
curl -i http://localhost:3000/api/ready
```

## 3. Recommended usage

### 3.1 Local debugging

- watch backend logs in your terminal
- track sensitive app events without logging secret data
- query metrics endpoints after auth, upload, and load tests
- monitor route counters and p50/p95/p99 latency snapshots

### 3.2 Lightweight self-hosted operations

- run under `tmux`, `screen`, or a simple service manager
- periodically capture metrics if manual trend tracking is needed
- keep observability simple unless a real scaling need appears

## 4. Missing pieces

This is not yet a full observability platform. Missing pieces include:

- dedicated Prometheus exposition
- broader runtime/system coverage
- persistent event history

## 5. Suggested roadmap

1. Extend runtime/system metrics that are useful for operations.
2. Improve sampling strategy if percentiles become business critical.
3. Add Prometheus output if external monitoring is required.
4. Add log persistence/rotation for long-running deployments.
5. Expand structured business-failure events where needed.

## 6. Project positioning

Aelidecx does not aim to replicate a full observability stack at this stage. The target is:

- enough logs to troubleshoot quickly
- enough metrics to validate security and usage behavior
- no heavy dependencies that make local setup harder
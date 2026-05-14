# Backend Development Guide

## Goal

Keep backend code easy to change by enforcing:

- Small files and focused responsibilities
- Clear separation between routing, handlers, middleware, and infrastructure
- Stable tests split by scope (unit vs integration)

## Production code layout

- `main.go`: process bootstrap and route wiring only
- `state.go`: shared runtime state and metrics structures
- `models.go`: API/domain structs
- `db.go`: SQLite initialization and schema setup
- `handlers_platform.go`: health/readiness handlers
- `handlers_notes_read.go`: notes read + search handlers
- `handlers_notes_write.go`: notes write/delete handlers
- `handlers_upload.go`: image upload handler
- `middleware_http.go`: CORS/security/cache headers middleware
- `middleware_auth_rate.go`: auth, API keys, rate limit, backoff
- `observability.go`: request metrics and system metrics
- `metrics_security.go`: security counters endpoint + counter helpers
- `export.go`: encrypted export + key derivation
- `logging.go`: structured logging helpers
- `sanitize.go`: sanitization and validation helpers

## Testing layout

- `test_helpers_test.go`: test app builder + shared helpers
- `unit_sanitize_config_test.go`: sanitize/folder/config unit tests
- `unit_auth_middleware_test.go`: auth and backoff unit tests
- `unit_misc_middleware_test.go`: headers/rate-limit/file-name unit tests
- `integration_platform_test.go`: health/readiness integration tests
- `integration_notes_crud_test.go`: CRUD/invalidation integration tests
- `integration_notes_search_test.go`: folder filtering + search integration tests
- `integration_metrics_export_test.go`: metrics and encrypted export integration tests

## Conventions

- Prefer one responsibility per file.
- Keep files under ~200 lines when possible.
- Place shared test setup in one helper test file, not duplicated in each test file.
- Add tests next to the concern they verify (unit or integration split).

## Typical commands

From `backend/`:

```bash
go test ./...
go test -race -covermode=atomic -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

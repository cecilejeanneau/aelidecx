# Security - Aelidecx

## Overview

Aelidecx applies pragmatic security controls suitable for a self-hosted personal application:

- mandatory authentication on API routes
- abuse protection through rate limiting and auth backoff
- HTML sanitization
- strict upload validation
- protective HTTP headers
- encrypted export support

The project is not designed to be exposed directly to the public internet without HTTPS, reverse proxying, and explicit network hardening.

## 1. API authentication

All routes under `/api` are protected by Bearer token authentication.

Expected format:

```http
Authorization: Bearer <token>
```

### Current behavior

- API keys are loaded from `API_KEYS` or `API_KEY`.
- `API_KEYS` supports a comma-separated list.
- Each key must be at least 24 characters.
- Key comparisons use constant-time checks.
- In development only, a local fallback key can be enabled when no key is configured.

### Expected responses

- `401` when header is missing or malformed
- `403` when token is invalid
- `429` when client is temporarily locked after repeated failures

## 2. Abuse protection

### Rate limiting

The backend applies sliding-window rate limiting per IP:

- window: 60 seconds
- default limit: 120 requests per minute

### Authentication backoff

Authentication failures are tracked per IP with exponential lock duration:

- 1st failure: 1s
- 2nd failure: 2s
- 3rd failure: 4s
- and so on
- cap: 128s

## 3. Security headers

The `securityHeaders()` middleware sets:

- `X-Frame-Options: SAMEORIGIN`
- `X-Content-Type-Options: nosniff`
- restrictive `Content-Security-Policy` for local assets
- `Referrer-Policy: no-referrer`
- `Permissions-Policy: camera=(), microphone=(), geolocation=()`
- `Strict-Transport-Security` only when HTTPS is enabled

API responses also include:

- `Cache-Control: no-store`

## 4. Data validation

### Notes

- title is trimmed and capped at 255 characters
- HTML content is sanitized with `bluemonday.UGCPolicy()`
- folder values are restricted to a strict allow-list
- content size is capped at 1 MB on create and update

### Search

- query length is limited to 128 characters
- allowed-character pattern reduces invalid/heavy requests

## 5. Image uploads

Uploads are validated with the following constraints:

- allowed extensions: `.png`, `.jpg`, `.jpeg`, `.gif`, `.webp`
- `.svg` is blocked
- max file size: 10 MB
- MIME type is verified from file bytes
- stored file name is random

Files are stored locally in `uploads/`.

## 6. Encrypted export

`GET /api/export/encrypted` exports notes in encrypted form.

Current characteristics:

- passphrase required via `X-Export-Passphrase`
- minimum passphrase length: 12 characters
- encryption: `AES-256-GCM`
- iterative key derivation using `SHA-256` for 150000 rounds
- JSON output includes salt, nonce, parameters, and base64 ciphertext

## 7. CORS and network exposure

- server binds to `127.0.0.1` by default
- allowed CORS origins come from `ALLOWED_ORIGINS`
- if unset, only local origins on port 3000 are allowed
- HTTPS is supported via `TLS_CERT_FILE` and `TLS_KEY_FILE`

## 8. Backup scope and local data

- database: `data/notes.db`
- uploads: `uploads/`

A minimal backup must include both locations.

## 9. Quick checks

### Protected API

```bash
curl http://localhost:3000/api/notes
curl -H "Authorization: Bearer wrong" http://localhost:3000/api/notes
curl -H "Authorization: Bearer your-api-key" http://localhost:3000/api/notes
```

### Security headers

```bash
curl -I http://localhost:3000
```

### Upload validation

```bash
curl -H "Authorization: Bearer your-api-key" \
     -F "image=@test.png" \
     http://localhost:3000/api/upload
```

## 10. Current limitations

The project does not yet provide:

- multi-user management
- full audit logging
- at-rest encryption of the full SQLite database
- identity federation
- centralized security monitoring

## 11. Suggested improvements

1. Continue expanding structured audit events for sensitive operations.
2. Evaluate SQLCipher if at-rest encryption becomes mandatory.
3. Add more detailed HTTPS operations guidance for remote exposure scenarios.

## 12. References

- `backend/main.go`
- `docs/OBSERVABILITY.md`
- `docs/REQUIREMENTS.md`
- `.env.example`
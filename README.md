# Aelidecx — Self-hosted personal knowledge base

[![CI](https://github.com/cecilejeanneau/aelidecx/actions/workflows/ci.yml/badge.svg)](https://github.com/cecilejeanneau/aelidecx/actions/workflows/ci.yml)
[![Coverage](https://codecov.io/gh/cecilejeanneau/aelidecx/branch/main/graph/badge.svg)](https://codecov.io/gh/cecilejeanneau/aelidecx)

> Self-hosted personal knowledge base — encrypted, offline-first, runs on your phone.

Personal encrypted wiki. Runs on Android via Termux, accessible from any browser.

## Stack

- **Backend**: Go + Fiber + SQLite
- **Frontend**: Vanilla JS (local rich-text editor, no CDN)
- **Runtime**: Termux (Android)
- **Cost**: free

## Features

- Create / edit / delete notes
- Rich-text editor (bold, italic, headings, lists, code blocks, links)
- Paste images directly into the editor
- Image upload via drag & drop or toolbar button
- Folders (inbox, concepts, practices, ui-ux, postmortems, references)
- Full-text search (SQLite FTS5)
- Swagger UI for API documentation
- Dedicated observability metrics endpoint for API traffic
- Responsive (mobile + desktop)

## Documentation

- Requirements: [docs/REQUIREMENTS.md](docs/REQUIREMENTS.md)
- Security: [docs/SECURITY.md](docs/SECURITY.md)
- Observability: [docs/OBSERVABILITY.md](docs/OBSERVABILITY.md)
- Backend development: [backend/DEVELOPMENT.md](backend/DEVELOPMENT.md)
- Change history: [CHANGELOG.md](CHANGELOG.md)

## Tests & Coverage

From the [backend/](backend/) directory:

```bash
go test ./...
go test -race -covermode=atomic -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Line-by-line HTML report:

```bash
go tool cover -html=coverage.out
```

Tests are split by scope in dedicated files under [backend/](backend/) and documented in [backend/DEVELOPMENT.md](backend/DEVELOPMENT.md).

## CI / CD

CI workflow: [.github/workflows/ci.yml](.github/workflows/ci.yml)

The pipeline runs:

1. `go mod tidy`
2. `go test -race -covermode=atomic -coverprofile=coverage.out ./...`
3. `go tool cover -func=coverage.out`
4. Global coverage threshold (35%)
5. Per-package coverage threshold (30%)
6. Coverage upload to Codecov
7. Backend build (`go build ./...`)

CD workflow: [.github/workflows/release.yml](.github/workflows/release.yml)

- Triggers on manual dispatch or a `v*` tag push
- Builds Linux binaries for `amd64` and `arm64`
- Publishes a GitHub Release with binaries attached

## Structure

```
aelidecx/
├── backend/
│   ├── main.go
│   └── go.mod
├── frontend/
│   ├── index.html
│   ├── css/style.css
│   └── js/
│       ├── app.js       ← main UI logic
│       ├── api.js       ← HTTP client
│       └── editor.js    ← rich-text editor + image upload
├── uploads/             ← uploaded images
├── data/                ← SQLite database (auto-created)
└── README.md
```

For the full backend file map (production + tests), see [backend/DEVELOPMENT.md](backend/DEVELOPMENT.md).

## Setup on Termux (Android)

### 1. Install Termux

Download from **F-Droid** (not the Play Store — the Play Store version is outdated):
https://f-droid.org/packages/com.termux/

### 2. Install Go

```bash
pkg update
pkg install golang git
```

### 3. Clone the project

```bash
cd ~
git clone <your-repo> aelidecx
# or copy files manually
```

### 4. Install Go dependencies

```bash
cd ~/aelidecx/backend
go mod tidy
```

### 5. Build and run

```bash
go build -o server .
./server
```

On Windows PowerShell, you can use the shorter launcher from the repo root:

```powershell
.\start-dev.ps1
```

If port `3000` is already taken, use a different one:

```powershell
.\start-dev.ps1 -Port 3001
```

You should see:
```
Server running on http://0.0.0.0:3000
```

### 6. Access from a browser

**From your phone:**
```
http://localhost:3000
```

**From your PC (same Wi-Fi):**
1. Find your phone's IP: `ifconfig wlan0` (look for `inet`)
2. Open on your PC: `http://192.168.x.x:3000`

**From your PC via hotspot:**
1. Enable Wi-Fi hotspot on your phone
2. Connect your PC to the hotspot
3. Phone IP is usually `192.168.43.1`
4. Open: `http://192.168.43.1:3000`

## Auto-start (optional)

To start the server automatically when Termux launches:

```bash
mkdir -p ~/.termux
echo "cd ~/aelidecx/backend && ./server" >> ~/.termux/boot.sh
chmod +x ~/.termux/boot.sh
```

Then install the **Termux:Boot** plugin from F-Droid.

## Backup

The database is at `data/notes.db`. Images are in `uploads/`.

Quick backup:
```bash
cp data/notes.db data/notes-backup-$(date +%Y%m%d).db
```

Or push everything to a private Git repository.

## Security

- Mandatory Bearer token auth on all API routes
- CORS restricted via `ALLOWED_ORIGINS`
- Server binds to localhost by default (`BIND_ADDR=127.0.0.1`)
- Exponential backoff on auth failures + per-IP rate limiting
- Strict image upload validation (size, extension, MIME)
- Do not expose to the internet without HTTPS (`TLS_CERT_FILE` + `TLS_KEY_FILE`)

## Observability

- HTTP request logs are emitted by the Fiber logger middleware
- API requests also emit structured JSON logs with `level`, `event`, `request_id`, route, status and latency
- Sensitive actions also emit structured events such as note creation, update, deletion, upload, auth failures and export activity
- Security counters are available at `GET /api/security/metrics`
- API traffic metrics are available at `GET /api/observability/metrics`
- The observability endpoint exposes average latency, simple `p50` / `p95` / `p99` latency snapshots and per-route counters
- The observability endpoint also exposes simple system metrics such as uptime, goroutine count and database file size
- Each API response now includes an `X-Request-ID` header for request correlation

## Probes

For quick operational checks, the backend exposes two unauthenticated probe endpoints:

- `GET /api/health`: liveness probe (returns `200` with `{"status":"ok"}` when the process is alive)
- `GET /api/ready`: readiness probe (returns `200` with `{"status":"ready"}` when SQLite is available, otherwise `503` with `{"status":"not_ready"}`)

Quick check:

```bash
curl -i http://localhost:3000/api/health
curl -i http://localhost:3000/api/ready
```

## Export

Notes are stored in SQLite. Raw text export:
```bash
sqlite3 data/notes.db "SELECT title, content FROM notes" > export.txt
```

Encrypted JSON export is available via `GET /api/export/encrypted` (AES-256-GCM).

# Requirements Specification - Aelidecx

## 1. Project purpose

Aelidecx is a self-hosted personal knowledge base, mobile-first, designed to run locally on Android via Termux and be accessible from a modern web browser.

The goal is to provide a simple, lightweight, sovereign tool to store, organize, search, and export rich notes with a security level suitable for personal or semi-private usage.

This specification is inferred from the current codebase, documentation, and CI/CD workflows.

## 2. Context and need

Aelidecx addresses a personal wiki / second-brain use case that is:

- independent from SaaS providers
- hosted by the user on their own device
- accessible from phone and desktop on the same network
- simple enough for autonomous installation
- robust enough to protect potentially sensitive notes

The project prioritizes autonomy, low cost, portability, and user control over data.

## 3. Objectives

### 3.1 Functional objectives

- Create, view, update, and delete notes.
- Edit rich HTML content in an integrated editor.
- Insert and upload images inside notes.
- Organize notes using predefined folders.
- Search notes with full-text search.
- Export data, including encrypted export.

### 3.2 Non-functional objectives

- Run locally on Android through Termux.
- Avoid CDN and third-party frontend runtime dependencies.
- Stay simple to deploy, maintain, and back up.
- Provide a reasonably hardened security baseline.
- Be automatically testable in CI and releasable as binaries.

## 4. Scope

### 4.1 In scope

- Go HTTP backend.
- Static frontend in HTML, CSS, and vanilla JavaScript.
- Local SQLite database for note storage.
- Local image upload storage.
- Bearer API key authentication for HTTP API.
- Swagger/OpenAPI API documentation.
- GitHub Actions CI/CD workflows.

### 4.2 Out of scope

- Real-time multi-device sync.
- Multi-user roles and permissions.
- Public note sharing.
- Collaborative concurrent editing.
- Native Android/iOS mobile apps.
- Direct internet exposure without explicit HTTPS/reverse-proxy setup.

## 5. Target users

- Individual users who want to self-host their notes.
- Technical or semi-technical users able to run software in Termux or Linux.
- Primary usage on Android phone, secondary usage from desktop browser.

## 6. Functional requirements

### 6.1 Note management

The system must provide:

- note creation
- single note retrieval
- list of notes
- note update
- note deletion
- sorting by update date
- folder filtering

Each note must include at least:

- id
- title
- rich content
- folder
- created timestamp
- updated timestamp

### 6.2 Folder organization

Notes must be constrained to an allow-list of folders to keep a simple and consistent structure.

Current folders:

- `inbox`
- `concepts`
- `practices`
- `ui-ux`
- `postmortems`
- `references`

Default folder on creation must be `inbox`.

### 6.3 Rich editor

The frontend must provide at least:

- bold and italic
- headings
- lists
- code blocks
- links
- paste support
- image insertion via paste, drag-and-drop, or upload button

### 6.4 Search

The system must support full-text search over note title and content.

Search must:

- use SQLite FTS5
- return relevant results quickly
- reject or limit invalid/excessive queries
- work on mobile and desktop

### 6.5 Image uploads

The system must support image upload for notes.

Minimum constraints:

- allowed extensions: PNG, JPG, JPEG, GIF, WEBP
- SVG forbidden
- max size: 10 MB
- MIME check from actual file bytes
- randomized stored filename
- local storage in a dedicated folder

### 6.6 Export

The system must provide:

- raw SQLite-level export for local backup
- encrypted export through API
- passphrase-based encryption flow

Encrypted export must use an AES-256-GCM style mechanism with iterative key derivation.

### 6.7 API documentation

The backend must expose Swagger/OpenAPI documentation through a web UI.

Documentation must describe at least:

- available routes
- input parameters
- nominal responses
- possible errors
- Bearer authentication mechanism

## 7. Technical requirements

### 7.1 Required stack

- Backend: Go 1.21+ with Fiber
- Database: SQLite with WAL mode
- Search: SQLite FTS5
- Frontend: HTML/CSS/vanilla JavaScript, no mandatory external framework
- HTML sanitization: bluemonday

### 7.2 Architecture

The app should remain a simple monolith:

- one backend binary
- static frontend served by backend
- local SQLite database
- local uploads directory
- local assets and data

### 7.3 Portability

Project must remain runnable on:

- Android via Termux
- Linux amd64
- Linux arm64

Release pipeline must produce binaries at least for `linux/amd64` and `linux/arm64`.

## 8. Security requirements

### 8.1 Authentication

- All `/api` routes must require Bearer authentication.
- API keys must come from `API_KEY` or `API_KEYS`.
- Key comparison must be constant-time.

### 8.2 Abuse protection

- Per-IP rate limiting must be enabled.
- Exponential backoff must be applied on auth failures.
- Clients may be temporarily blocked after repeated failures.

### 8.3 Content protection

- HTML content must be sanitized to reduce stored XSS risk.
- API responses must be non-cacheable.
- Security HTTP headers must be set by default.

### 8.4 Network exposure

- Default bind address must be `127.0.0.1`.
- HTTPS support must be available through cert/key files.
- CORS origin opening must be strictly configuration-driven.

## 9. Operations requirements

### 9.1 Installation

Installation should stay concise:

- install Go
- clone repository
- `go mod tidy`
- `go build`
- run binary

### 9.2 Backup

Minimum backup coverage:

- `data/notes.db`
- `uploads/`

### 9.3 Configuration

Main runtime options must be environment variables, including:

- API key(s)
- bind address and port
- allowed CORS origins
- TLS certificate and key
- development-mode options

## 10. Quality requirements

### 10.1 Testing

Backend automated tests should cover at least:

- critical business rules
- input validation
- baseline security behavior
- main CRUD paths
- search
- encrypted export

### 10.2 CI

Continuous integration pipeline must run at least:

- dependency resolution
- backend tests
- code coverage reporting
- global coverage threshold check
- per-package coverage threshold check
- backend build
- Swagger documentation generation

### 10.3 CD

Release pipeline must:

- trigger on version tag or manual dispatch
- produce target binaries
- attach artifacts to GitHub Release

### 10.4 CI/CD supply-chain hardening

GitHub Actions workflows should avoid floating references for critical actions and tools.

Minimum expectations:

- pinned versions or SHAs for GitHub Actions
- fixed versions for CLIs downloaded during workflows
- reduced implicit/uncontrolled updates

## 11. UX constraints

- Interface must be usable on mobile screens.
- Interface must remain responsive on desktop.
- Main interactions must stay fast and clear.
- Note-taking flow should prioritize simplicity over excessive feature depth.

## 12. Acceptance criteria

The project is considered compliant at minimum if:

- a user can run the server locally and access UI from a browser
- API authentication blocks unauthorized access
- a note can be created, updated, deleted, and found via search
- a valid image can be uploaded and reused in a note
- allowed folders are enforced
- encrypted export works with a valid passphrase
- Swagger documentation is reachable
- CI passes successfully
- release workflow produces expected binaries

## 13. Potential future enhancements

The following are possible but not required at this stage:

- cross-folder tags
- richer import/export
- stronger client-side offline mode
- automated backups
- multi-device sync
- full at-rest database encryption
- multi-user support

## 14. Document status

This requirements specification is reconstructed from the existing product state. It is a framing document used to formalize expectations and align documentation, code, and future evolution.
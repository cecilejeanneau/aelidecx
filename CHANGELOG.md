# Changelog

All notable changes to Aelidecx are documented in this file.

The format is inspired by Keep a Changelog and uses a simple project-oriented history.

## [Unreleased]

### Added

- Swagger/OpenAPI documentation for backend routes with Swagger UI exposure.
- Project requirements document in `docs/REQUIREMENTS.md`.
- Observability document in `docs/OBSERVABILITY.md`.
- Explicit CI/CD supply-chain hardening with pinned action/tool versions.

### Changed

- Project renamed from `my-notes` to `Aelidecx`.
- Full source translation to English for user-facing code and main technical documentation.
- Go doc comments and JSDoc coverage substantially improved across backend and frontend.
- CI coverage pipeline now includes per-package enforcement and Codecov upload.
- Release workflow now produces Linux `amd64` and `arm64` binaries for GitHub Releases.
- Backend API documentation comments were simplified after adding Swagger annotations.

### Security

- Bearer auth supports multiple configured API keys via `API_KEYS` with `API_KEY` fallback.
- Constant-time API key comparison enforced.
- Authentication backoff and lockout metrics exposed through the security metrics endpoint.
- CI and release workflows now avoid floating action references and unpinned CLI downloads.

### Fixed

- `sanitizeTitle()` now trims whitespace before truncating to 255 characters.
- CI coverage upload path corrected.
- Security section comments and several documentation inconsistencies aligned with the current codebase.

## [2026-05] Documentation and delivery baseline

### Added

- Full backend test suite covering CRUD, validation, search, export, auth, rate limiting and sanitization.
- GitHub Actions CI workflow for tests, coverage and build verification.
- GitHub Actions release workflow for binary publication.
- Swagger bootstrap for generated API docs.

### Changed

- README updated with badges, setup guidance and deployment notes.
- Project distribution aligned around self-hosted Android/Termux usage and browser access.

## [2026-05] Initial application scope

### Added

- Go backend using Fiber and SQLite.
- Vanilla JS frontend with local rich-text editor.
- Folder-based note organisation.
- Image upload support.
- SQLite FTS5 full-text search.
- Encrypted export endpoint.

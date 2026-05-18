# Project Structure

## Workspace Layout

```
SWU_OSR/                              # Root workspace
├── .kiro/
│   ├── specs/swu-osr-platform/       # Feature specs (requirements, design, tasks)
│   └── steering/                     # AI steering rules (this directory)
├── plan/                             # Planning & skill documents
├── SWAAP/                            # Main application repository
│   ├── cmd/wrapper-api/main.go       # Go API entrypoint
│   ├── legacy/                       # SIAKAD proxy client library
│   │   ├── client.go                 # HTTP client with PHP session simulation
│   │   └── client_test.go
│   ├── flutter_app/                  # Flutter web/mobile frontend
│   │   ├── lib/
│   │   │   ├── main.dart
│   │   │   ├── models/              # Data models (jadwal, presensi, zoom)
│   │   │   ├── screens/             # UI screens
│   │   │   └── services/            # API services, credential storage
│   │   ├── Dockerfile               # Flutter web build + nginx serve
│   │   └── pubspec.yaml
│   ├── nginx/nginx.conf              # Internal nginx (rate limiting, proxy rules)
│   ├── docker-compose.yml            # Container orchestration
│   ├── Dockerfile.backend            # Go multi-stage build (distroless)
│   ├── go.mod                        # Go module (module name: swaap)
│   ├── .env.example                  # Environment variable template
│   └── docs/DEPLOY.md               # VPS deployment guide (Indonesian)
```

## Planned Structure (SWU OSR Platform)

```
backend/
├── cmd/api/main.go                   # Server entrypoint
├── internal/
│   ├── config/                       # Environment-based configuration
│   ├── domain/                       # Business entities & interfaces
│   ├── handler/                      # HTTP handlers (thin, validation only)
│   ├── service/                      # Business logic layer
│   ├── repository/                   # DB queries (sqlc-generated)
│   ├── middleware/                   # Auth, logging, CORS, rate-limit
│   ├── siakad/                       # SIAKAD proxy (adapted from SWAAP legacy/)
│   └── github/                       # GitHub OAuth + API client
├── migrations/                       # golang-migrate SQL files
├── sqlc/                             # sqlc queries & config
├── Dockerfile
└── go.mod

frontend/
├── src/
│   ├── app/                          # Next.js App Router pages & layouts
│   ├── components/                   # Reusable UI (shadcn/ui)
│   ├── features/                     # Feature modules (auth, feed, forum, etc.)
│   ├── hooks/                        # Custom React hooks
│   ├── lib/                          # API client, utilities
│   └── types/                        # TypeScript type definitions
├── Dockerfile
└── package.json
```

## Architecture Conventions

- **Go backend**: `cmd/` for entrypoints, `internal/` for private packages, package-per-concern
- **Clean Architecture layers**: Handler → Service → Repository (dependencies point inward)
- **Frontend**: Feature-based module organization under `src/features/`
- **Docker**: Each service has its own Dockerfile; orchestrated via docker-compose.yml
- **Backend not exposed to host**: Only accessible via internal Docker network through nginx
- **Frontend bound to localhost**: Host nginx handles TLS termination and proxies to container

## Key Patterns

- JSON response envelope: `{"ok": bool, "data": ..., "error": ...}`
- Interface-driven services (all business logic behind interfaces)
- UUID primary keys with soft deletes (`deleted_at` column)
- Cursor-based pagination (base64-encoded timestamps)
- RBAC via JWT claims: student, faculty, admin
- Domain types use `json:"-"` to hide private fields from public responses
- Atomic operations for multi-step mutations (showcase updates)
- HMAC-SHA256 for webhook signature verification

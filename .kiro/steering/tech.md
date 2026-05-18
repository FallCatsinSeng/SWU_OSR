# Tech Stack & Build System

## Backend (Current — SWAAP)

- **Language**: Go 1.22
- **Router**: Standard library `net/http` with `http.NewServeMux()`
- **Dependencies**: None (stdlib only)
- **Container**: Multi-stage Docker (golang:1.22-alpine → gcr.io/distroless/static-debian12:nonroot)

## Backend (Planned — SWU OSR Platform)

- **Language**: Go (monolith)
- **Router**: go-chi/chi/v5
- **Architecture**: Clean Architecture (Handler → Service → Repository → PostgreSQL)
- **Database**: PostgreSQL 16 + sqlc for type-safe query generation
- **Cache/Sessions**: Redis 7
- **Migrations**: golang-migrate
- **Auth**: JWT (golang-jwt/jwt/v5) + single-use refresh tokens + GitHub OAuth
- **Validation**: go-playground/validator/v10
- **Logging**: uber-go/zap
- **Config**: spf13/viper
- **Testing**: stretchr/testify + leanovate/gopter (property-based)
- **UUID**: google/uuid
- **Encryption**: AES-256-GCM for GitHub token storage at rest

## Frontend (Current — SWAAP)

- **Framework**: Flutter (Dart SDK ^3.10.4)
- **Dependencies**: http, google_fonts, shared_preferences
- **Served via**: Nginx Alpine container on port 8080

## Frontend (Planned — SWU OSR Platform)

- **Framework**: Next.js 14+ (App Router)
- **Language**: TypeScript 5.x
- **Styling**: Tailwind CSS 3.x
- **UI Components**: shadcn/ui
- **Data Fetching**: @tanstack/react-query v5
- **Forms**: react-hook-form + zod
- **HTTP Client**: axios
- **Icons**: lucide-react

## Infrastructure

- Docker Compose for orchestration
- Nginx reverse proxy (container-level + host-level for TLS termination)
- Cloudflare Tunnel for public exposure
- Let's Encrypt / Certbot for HTTPS
- UFW firewall (ports 22, 80, 443 only)
- Security hardening: `cap_drop: ALL`, `read_only`, `no-new-privileges`

## Common Commands

### Backend

```bash
# Run locally
go run ./cmd/wrapper-api

# Run tests
go test ./...

# Build binary
CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o wrapper-api ./cmd/wrapper-api
```

### Frontend (Flutter — current)

```bash
cd flutter_app
flutter pub get
flutter run
```

### Docker

```bash
# Build all images
docker compose build

# Deploy (background)
docker compose up -d

# View logs
docker compose logs -f backend
docker compose logs -f frontend

# Restart
docker compose restart

# Clean old images
docker image prune -f
```

### Deployment (VPS)

```bash
cd /opt/swaap
git pull origin main
docker compose build
docker compose up -d
```

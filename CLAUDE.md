# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build the API binary
go build -o main ./cmd/api/main.go

# Run all tests
go test ./...

# Run tests in a specific package (with verbose output)
go test -v ./internal/service/...

# Run a single test function
go test -v -run TestRegister_Success ./internal/service/

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run the full stack (PostgreSQL + app + nginx)
docker compose up --build

# Run only the database for local development
docker compose up db -d
```

## Architecture

A **Go + Gin + GORM** backend API for hospital staff to search patient records across hospitals. Uses a clean layered architecture with dependency injection in `main.go`.

### Layered structure (bottom → top)

| Layer | Directory | Role |
|-------|-----------|------|
| **Domain** | `internal/domain/` | Data models (Hospital, Staff, Patient) + repository interfaces |
| **Repository** | `internal/repository/` | GORM-based DB access; implements domain interfaces |
| **Infrastructure** | `internal/infrastructure/` | External API client (Hospital A patient lookup) |
| **Service** | `internal/service/` | Business logic (auth, patient search with external fallback) |
| **Middleware** | `internal/middleware/` | JWT auth middleware (extracts `hospital_id` claim into context) |
| **Handler** | `internal/handler/` | Gin HTTP handlers (JSON binding, response) |
| **Entrypoint** | `cmd/api/main.go` | Wires everything together, runs migrations, seeds hospitals |

### Key patterns

- **Repository interfaces live in `internal/domain/`**, implementations in `internal/repository/` — this keeps domain clean of GORM imports.
- **Services depend on interfaces**, not concrete repos — tests mock at the interface level.
- **Handler layer is thin** — validates input, calls service, returns JSON.
- **JWT contains `hospital_id`** — middleware extracts it and sets it in the gin context. All patient queries are scoped by `hospital_id` for data isolation.
- **External API fallback**: When Hospital A staff search with no local results, the service calls an external API and caches the result locally via upsert.
- **Config** via Viper — reads `.env` file with env var overrides (critical: env vars are explicitly `BindEnv`'d).

### Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/staff/create` | No | Register staff (username, password, hospital) |
| POST | `/staff/login` | No | Login, returns JWT |
| GET | `/patient/search` | JWT | Search patients by various filters, scoped to user's hospital |
| GET | `/` | No | Serves frontend HTML page |

### Testing pattern

Mocks are hand-written structs in `_test.go` files (no mock library). Each test file defines its own mocks implementing the domain interfaces. Tests use `gin.TestMode` and `httptest.NewRecorder()` for HTTP-level testing.

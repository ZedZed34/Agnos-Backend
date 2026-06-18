---
name: go-dev
description: Go development toolkit for Agnos-Backend. Run tests, format code, lint, build, and manage the Go toolchain. Activates on requests like "run tests", "format code", "lint", "build the project", "check coverage", "run Postman collections".
---

# Agnos-Backend Go Dev Toolkit

## Project Info

| Field | Value |
|-------|-------|
| Module | `agnos-gin` |
| Go version | 1.24.5 |
| Entrypoint | `cmd/api/main.go` |
| Framework | Gin + GORM + Viper |
| Test framework | Standard `testing` + `httptest` |

## Test Suites

### Run All Tests

```bash
cd /mnt/d/Agnos-Backend && go test ./...
```

### Run Tests for Specific Package

```bash
go test ./internal/handler/...
go test ./internal/service/...
go test ./internal/middleware/...
```

### Verbose Output

```bash
go test -v ./...
go test -v ./internal/service/...
```

### Coverage

```bash
# Coverage report
go test -cover ./...

# Detailed coverage profile
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# HTML coverage view
go tool cover -html=coverage.out -o coverage.html
```

### Run Single Test Function

```bash
go test -run TestRegister_Success ./internal/service/...
go test -run TestLoginHandler ./internal/handler/...
go test -v -run "TestLogin" ./...
```

### Race Detection

```bash
go test -race ./...
```

### Current Test Files (5 total)

| File | Package |
|------|---------|
| `internal/handler/auth_handler_test.go` | handler |
| `internal/handler/patient_handler_test.go` | handler |
| `internal/middleware/auth_middleware_test.go` | middleware |
| `internal/service/auth_service_test.go` | service |
| `internal/service/patient_service_test.go` | service |

Tests use mock repositories (in-memory maps) — no database required.

## Code Formatting

### Format All Go Files

```bash
cd /mnt/d/Agnos-Backend && go fmt ./...
```

### Vet (static analysis)

```bash
go vet ./...
```

### goimports (format + organize imports)

```bash
# Install if missing
go install golang.org/x/tools/cmd/goimports@latest

# Run
goimports -w .
```

### Combined Format + Vet Check

```bash
go fmt ./... && go vet ./...
```

## Build

### Local Build

```bash
cd /mnt/d/Agnos-Backend && go build -o bin/api ./cmd/api/main.go
```

### Docker Build

```bash
docker compose build app
```

### Tidy Dependencies

```bash
go mod tidy
go mod verify
```

## Postman Collections

No Postman collections found in this repository. API endpoints:

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/staff/create` | None | Register staff |
| POST | `/staff/login` | None | Login, returns JWT |
| GET | `/patient/search` | JWT | Search patients |
| GET | `/` | None | Static frontend |

To create a Postman collection for these endpoints, use the `newman` runner:

```bash
# Install newman (if not present)
npm install -g newman

# Run a collection file
newman run collection.json
```

## Common Workflows

### Quick Verify (format + vet + test)

```bash
cd /mnt/d/Agnos-Backend && go fmt ./... && go vet ./... && go test ./...
```

### Full CI Check (format + vet + test + coverage)

```bash
cd /mnt/d/Agnos-Backend && \
  go fmt ./... && \
  go vet ./... && \
  go test -race -coverprofile=coverage.out ./... && \
  go tool cover -func=coverage.out
```

### After Docker Rebuild — Smoke Test

```bash
curl -s http://localhost:8080/ && echo "OK" || echo "FAIL"
curl -s -X POST http://localhost:8080/staff/create \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123","hospital":"Hospital A"}'
```

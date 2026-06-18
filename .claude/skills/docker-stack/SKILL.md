---
name: docker-stack
description: Manage the Agnos-Backend Docker Compose stack. Use for starting, stopping, restarting, rebuilding, or troubleshooting the full Docker stack (app, db, nginx). Activates on requests like "start docker", "restart the stack", "rebuild containers", "check logs", "reset the database volume".
---

# Agnos-Backend Docker Stack

## Stack Overview

| Service | Image/Build | Port | Depends On |
|---------|-------------|------|------------|
| `app` | Build from `Dockerfile` (Go 1.24) | `8080:8080` | `db` |
| `db` | `postgres:15-alpine` | `5432:5432` | — |
| `nginx` | `nginx:alpine` | `80:80` | `app` |

- **Volume**: `pgdata` (PostgreSQL data, persists across restarts)
- **Network**: All services on default Compose network, internal DNS works via service names

## Commands

### Start / Spin Up

```bash
docker compose up -d
```

First run or after code changes:

```bash
docker compose up -d --build
```

Rebuild only the app (not db or nginx):

```bash
docker compose up -d --build app
```

### Check Status

```bash
docker compose ps
```

All three services (`app`, `db`, `nginx`) should show `Up`.

### Restart Services

Restart everything:

```bash
docker compose restart
```

Restart only one service:

```bash
docker compose restart app
docker compose restart db
docker compose restart nginx
```

Restart with fresh code (rebuild + restart app):

```bash
docker compose up -d --build app
```

### Stop

```bash
docker compose stop
```

Stop and remove containers (network also removed):

```bash
docker compose down
```

### Logs

```bash
docker compose logs -f          # follow all
docker compose logs -f app      # follow app only
docker compose logs --tail=50   # last 50 lines
```

### Reset Database (destroy + recreate)

```bash
docker compose down -v          # removes container + pgdata volume
docker compose up -d            # creates fresh, empty DB
```

### Full Clean Teardown

```bash
docker compose down -v --rmi all
```

## Troubleshooting

### App won't start — database connection refused

App tries to connect before PostgreSQL is ready. Fix:

```bash
docker compose restart app
```

If persistent, ensure `depends_on` is present (it is) and check DB healthcheck. Add a healthcheck to `db` in `docker-compose.yml` if needed:

```yaml
db:
  healthcheck:
    test: ["CMD-SHELL", "pg_isready -U agnos"]
    interval: 5s
    timeout: 5s
    retries: 5
```

Then change `depends_on` on `app`:

```yaml
depends_on:
  db:
    condition: service_healthy
```

### Port already in use

Check what's on the port:

```bash
sudo lsof -i :8080
sudo lsof -i :5432
sudo lsof -i :80
```

### Container stuck restarting

```bash
docker compose down
docker compose up -d
```

### nginx returns 502

Usually means `app` isn't ready yet or crashed:

```bash
docker compose ps app
docker compose logs app --tail=20
docker compose restart app
```

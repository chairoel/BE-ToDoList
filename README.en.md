# To-Do List API

Backend API (Go + Gin + PostgreSQL) with a **handler → usecase → repository** architecture, Goose migrations, and Docker.

**Languages:** [English](README.en.md) | [Bahasa Indonesia](README.md)

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [Clone & environment setup](#2-clone--environment-setup)
3. [Start the database](#3-start-the-database)
4. [Run migrations](#4-run-migrations)
5. [Run the API](#5-run-the-api)
6. [How to use the API](#6-how-to-use-the-api)
7. [Make commands](#7-make-commands)
8. [Project structure](#8-project-structure)
9. [Troubleshooting](#9-troubleshooting)

---

## 1. Prerequisites

Install these tools before you start:

| Tool | Purpose | Check version |
|------|---------|---------------|
| [Docker](https://docs.docker.com/get-docker/) + Docker Compose | PostgreSQL & API containers | `docker --version` |
| [Go](https://go.dev/dl/) 1.26+ | Run the API locally | `go version` |
| [Goose](https://github.com/pressly/goose) (optional) | Local migrations; Make falls back to Docker if missing | `goose --version` |
| Make | Helper commands | `make --version` |

Install Goose (macOS):

```bash
brew install goose
```

---

## 2. Clone & environment setup

```bash
# go to the project folder
cd to-do-list

# copy environment file
cp .env.example .env
```

`.env` contents (defaults are ready for local use):

```env
APP_PORT=8080

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=todolist
POSTGRES_SSLMODE=disable

GOOSE_DRIVER=postgres
GOOSE_DBSTRING=host=localhost port=5432 user=postgres password=postgres dbname=todolist sslmode=disable
GOOSE_MIGRATION_DIR=./migrations
```

> **Note:** When the API runs via Docker Compose, `POSTGRES_HOST` is set automatically to `postgres` (the service name), not `localhost`.

Install Go dependencies:

```bash
go mod download
```

---

## 3. Start the database

Start PostgreSQL only:

```bash
make db-up
```

Check status:

```bash
docker compose ps
```

Postgres is available at `localhost:5432`.

Stop all Docker services:

```bash
make db-down
# or
make docker-down
```

---

## 4. Run migrations

Make sure the database is running, then apply the schema:

```bash
make migrate-up
```

Check migration status:

```bash
make migrate-status
```

Other migration commands:

```bash
make migrate-down                    # rollback 1 step
make migrate-reset                   # rollback all
make migrate-version                 # current DB version
make migrate-create name=add_column  # create a new migration file
```

Existing migrations:

- `00001_create_users_table.sql` — `users` table
- `00002_create_notes_table.sql` — `notes` table

---

## 5. Run the API

Choose one option.

### Option A — Local (Go)

```bash
make db-up
make migrate-up
make run
```

API: `http://localhost:8080`

### Option B — Docker (API + Postgres)

```bash
make migrate-up   # run migrations first (local Goose / Docker tools)
make docker-up    # build image + start api & postgres
```

Follow API logs:

```bash
make docker-logs
```

Restart the API after code changes:

```bash
make docker-restart
```

Stop everything:

```bash
make docker-down
```

### Health check

```bash
curl http://localhost:8080/health
```

Response:

```json
{"status":"ok"}
```

---

## 6. How to use the API

Base URL: `http://localhost:8080/api/v1`

### Create user

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Irul",
    "email": "irul@example.com",
    "password": "secret123"
  }'
```

### List users

```bash
curl http://localhost:8080/api/v1/users
```

### Get user by ID

```bash
curl http://localhost:8080/api/v1/users/<user-id>
```

### Update user

```bash
curl -X PUT http://localhost:8080/api/v1/users/<user-id> \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Irul Updated"
  }'
```

### Delete user

```bash
curl -X DELETE http://localhost:8080/api/v1/users/<user-id>
```

### Endpoint summary

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/users` | Create user |
| `GET` | `/api/v1/users` | List users |
| `GET` | `/api/v1/users/:id` | Get user |
| `PUT` | `/api/v1/users/:id` | Update user |
| `DELETE` | `/api/v1/users/:id` | Delete user |

Passwords are hashed with bcrypt and are **not** returned in JSON responses.

---

## 7. Make commands

```bash
make help
```

| Command | Description |
|---------|-------------|
| `make db-up` | Start PostgreSQL |
| `make db-down` | Stop all Docker services |
| `make migrate-up` | Apply migrations |
| `make migrate-down` | Rollback 1 migration |
| `make migrate-status` | Migration status |
| `make migrate-reset` | Rollback all migrations |
| `make migrate-create name=...` | Create a new migration file |
| `make run` | Run API locally |
| `make build` | Build local binary (`bin/api`) |
| `make docker-up` | Build & start API + Postgres |
| `make docker-down` | Stop API + Postgres |
| `make docker-build` | Build API image only |
| `make docker-restart` | Rebuild & recreate API |
| `make docker-logs` | Follow API logs |

---

## 8. Project structure

```text
.
├── cmd/api/main.go              # entrypoint + graceful shutdown
├── internal/
│   ├── config/                  # env loading
│   ├── database/                # PostgreSQL connection (pgx)
│   ├── domain/                  # entities & interfaces
│   ├── repository/              # database access
│   ├── usecase/                 # business logic
│   ├── handler/                 # HTTP handlers
│   └── router/                  # Gin routes
├── migrations/                  # SQL migrations (Goose)
├── Dockerfile                   # multi-stage build
├── docker-compose.yml           # postgres + api + goose
├── Makefile
├── .env.example
├── README.md                    # Bahasa Indonesia
└── README.en.md                 # English
```

Request flow:

```text
HTTP → Handler → Usecase → Repository → PostgreSQL
```

---

## 9. Troubleshooting

### Port 5432 / 8080 already in use

Change the port in `.env` (`POSTGRES_PORT` / `APP_PORT`) or stop the process using that port.

### Migration cannot connect to the DB

Make sure Postgres is healthy:

```bash
make db-up
docker compose ps
```

Then retry `make migrate-up`.

### Docker API cannot connect to the DB

Use `make docker-up` (Compose sets `POSTGRES_HOST=postgres`). Do not use `localhost` from inside the API container.

### Goose is not installed

Not required. `make migrate-*` automatically uses the Goose Docker image if the local binary is missing.

---

## Quick start (end to end)

```bash
# 1. Setup
cp .env.example .env
go mod download

# 2. Database
make db-up

# 3. Migration
make migrate-up

# 4. Run API (pick one)
make run          # local
# or
make docker-up    # docker

# 5. Test
curl http://localhost:8080/health
```

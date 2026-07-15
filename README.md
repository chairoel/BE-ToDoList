# To-Do List API

Backend API (Go + Gin + PostgreSQL) dengan arsitektur **handler → usecase → repository**, migration Goose, dan Docker.

**Bahasa:** [Bahasa Indonesia](README.md) | [English](README.en.md)

## Daftar Isi

1. [Prasyarat](#1-prasyarat)
2. [Clone & setup environment](#2-clone--setup-environment)
3. [Jalankan database](#3-jalankan-database)
4. [Jalankan migration](#4-jalankan-migration)
5. [Jalankan API](#5-jalankan-api)
6. [Cara menggunakan API](#6-cara-menggunakan-api)
7. [Perintah Make](#7-perintah-make)
8. [Struktur project](#8-struktur-project)
9. [Troubleshooting](#9-troubleshooting)

---

## 1. Prasyarat

Install tools berikut sebelum mulai:

| Tool | Keterangan | Cek versi |
|------|------------|-----------|
| [Docker](https://docs.docker.com/get-docker/) + Docker Compose | PostgreSQL & API container | `docker --version` |
| [Go](https://go.dev/dl/) 1.26+ | Jalankan API secara lokal | `go version` |
| [Goose](https://github.com/pressly/goose) (opsional) | Migration lokal; tanpa ini Make pakai Docker | `goose --version` |
| Make | Helper command | `make --version` |

Install Goose (macOS):

```bash
brew install goose
```

---

## 2. Clone & setup environment

```bash
# masuk ke folder project
cd to-do-list

# salin file environment
cp .env.example .env
```

Isi `.env` (default sudah siap dipakai lokal):

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

> **Catatan:** Saat API dijalankan lewat Docker Compose, `POSTGRES_HOST` di-set otomatis ke `postgres` (nama service), bukan `localhost`.

Install dependency Go:

```bash
go mod download
```

---

## 3. Jalankan database

Nyalakan PostgreSQL saja:

```bash
make db-up
```

Cek status:

```bash
docker compose ps
```

Postgres tersedia di `localhost:5432`.

Hentikan semua service Docker:

```bash
make db-down
# atau
make docker-down
```

---

## 4. Jalankan migration

Pastikan database sudah jalan, lalu terapkan schema:

```bash
make migrate-up
```

Cek status migration:

```bash
make migrate-status
```

Perintah migration lain:

```bash
make migrate-down                    # rollback 1 step
make migrate-reset                   # rollback semua
make migrate-version                 # versi DB saat ini
make migrate-create name=add_column  # buat file migration baru
```

Migration yang ada:

- `00001_create_users_table.sql` — tabel `users`
- `00002_create_notes_table.sql` — tabel `notes`

---

## 5. Jalankan API

Pilih salah satu cara.

### Opsi A — Lokal (Go)

```bash
make db-up
make migrate-up
make run
```

API: `http://localhost:8080`

### Opsi B — Docker (API + Postgres)

```bash
make migrate-up   # migration dulu (Goose lokal / Docker tools)
make docker-up    # build image + start api & postgres
```

Cek log API:

```bash
make docker-logs
```

Restart API setelah ubah kode:

```bash
make docker-restart
```

Stop semua:

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

## 6. Cara menggunakan API

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

### Ringkasan endpoint

| Method | Path | Keterangan |
|--------|------|------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/v1/users` | Create user |
| `GET` | `/api/v1/users` | List users |
| `GET` | `/api/v1/users/:id` | Detail user |
| `PUT` | `/api/v1/users/:id` | Update user |
| `DELETE` | `/api/v1/users/:id` | Delete user |

Password di-hash (bcrypt) dan **tidak** dikembalikan di response JSON.

---

## 7. Perintah Make

```bash
make help
```

| Command | Fungsi |
|---------|--------|
| `make db-up` | Start PostgreSQL |
| `make db-down` | Stop semua Docker service |
| `make migrate-up` | Apply migration |
| `make migrate-down` | Rollback 1 migration |
| `make migrate-status` | Status migration |
| `make migrate-reset` | Rollback semua migration |
| `make migrate-create name=...` | Buat file migration baru |
| `make run` | Jalankan API lokal |
| `make build` | Build binary lokal (`bin/api`) |
| `make docker-up` | Build & start API + Postgres |
| `make docker-down` | Stop API + Postgres |
| `make docker-build` | Build image API saja |
| `make docker-restart` | Rebuild & recreate API |
| `make docker-logs` | Follow log API |

---

## 8. Struktur project

```text
.
├── cmd/api/main.go              # entrypoint + graceful shutdown
├── internal/
│   ├── config/                  # load env
│   ├── database/                # koneksi PostgreSQL (pgx)
│   ├── domain/                  # entity & interface
│   ├── repository/              # akses database
│   ├── usecase/                 # business logic
│   ├── handler/                 # HTTP handler
│   └── router/                  # route Gin
├── migrations/                  # SQL migration (Goose)
├── Dockerfile                   # multi-stage build
├── docker-compose.yml           # postgres + api + goose
├── Makefile
├── .env.example
├── README.md                    # Bahasa Indonesia
└── README.en.md                 # English
```

Alur request:

```text
HTTP → Handler → Usecase → Repository → PostgreSQL
```

---

## 9. Troubleshooting

### Port 5432 / 8080 sudah dipakai

Ubah port di `.env` (`POSTGRES_PORT` / `APP_PORT`) atau hentikan proses yang memakai port tersebut.

### Migration gagal konek DB

Pastikan Postgres sudah sehat:

```bash
make db-up
docker compose ps
```

Lalu coba lagi `make migrate-up`.

### API Docker gagal connect ke DB

Pastikan memakai `make docker-up` (compose set `POSTGRES_HOST=postgres`). Jangan pakai `localhost` dari dalam container API.

### Goose tidak terpasang

Tidak wajib. `make migrate-*` otomatis memakai image Docker Goose jika binary lokal tidak ada.

---

## Alur cepat (dari awal sampai jalan)

```bash
# 1. Setup
cp .env.example .env
go mod download

# 2. Database
make db-up

# 3. Migration
make migrate-up

# 4. Jalankan API (pilih satu)
make run          # lokal
# atau
make docker-up    # docker

# 5. Test
curl http://localhost:8080/health
```

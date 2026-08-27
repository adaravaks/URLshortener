# URL Shortener

A small REST API for shortening URLs, built in Go as a learning project. Boasts no ORM or web framework, just the standard library plus a couple of focused dependencies. Built to practice idiomatic Go.

## Features

- `POST /links` — shorten a URL
- `GET /{code}` — redirect to the original URL, with async click tracking
- `GET /links/{code}/stats` — view click count and creation date for a link

## Tech stack

- **Go 1.27**, using the standard library's `net/http` router
- **PostgreSQL**, accessed via `database/sql` + `lib/pq`
- **golang-migrate** for schema migrations
- **testify** for (obviously) testing
- **Docker Compose** for local orchestration (app + database)
- **GitHub Actions** for CI (build + test on every push)

## Architecture

```
cmd/api/          entry point, wires everything together
internal/
  handler/        HTTP handlers, middleware, request/response logic
  model/          plain data structs
  repository/     data access — interface + Postgres implementation
migrations/       versioned SQL migrations
```

The handler layer depends only on the `LinkRepository` interface, never on the concrete Postgres implementation — this keeps HTTP concerns and persistence concerns fully decoupled and makes the handlers testable without a real database (see `internal/handler/handler_test.go`, which uses an in-memory fake repository).

## Running locally

**Requirements:** Docker Desktop, Go 1.27+ (only needed if running outside Docker)

Start everything with Docker Compose:

```bash
docker compose up -d --build
```

This starts Postgres, runs the app in a container, and exposes the API on `http://localhost:8080`.

Check it's healthy:

```bash
curl http://localhost:8080/health
```

### Running outside Docker

If you'd rather run the Go app directly against a local Postgres:

```bash
docker compose up -d postgres
go run ./cmd/api
```

The app falls back to a local connection string (`localhost:5433`) when the `DATABASE_URL` environment variable isn't set — see `cmd/api/main.go`.

### Migrations

Migrations live in `migrations/` and are applied with [golang-migrate](https://github.com/golang-migrate/migrate):

```bash
migrate -database "postgres://urlshortener:urlshortener@localhost:5433/urlshortener?sslmode=disable" -path migrations up
```

## API examples

**Create a short link:**

```bash
curl -X POST http://localhost:8080/links -d '{"url":"https://example.com"}'
```

```json
{
  "id": 1,
  "short_code": "aZ3xQ1p",
  "original_url": "https://example.com",
  "created_at": "2026-08-26T10:00:00Z",
  "click_count": 0
}
```

**Follow a short link:**

```bash
curl -v http://localhost:8080/aZ3xQ1p
```

**Check stats:**

```bash
curl http://localhost:8080/links/aZ3xQ1p/stats
```

## Tests

```bash
go test ./... -v
```

Unit tests cover the handler layer using an in-memory fake repository — no database required to run them. CI runs the same command on every push via GitHub Actions (`.github/workflows/ci.yml`).

### The entire README except for this sentence and the select few lines I edited is AI-written, sorry.
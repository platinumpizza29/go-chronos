# go-chronos

Go task prioritisation server using Gin + PostgreSQL + Gemini AI.

## Quick start

```bash
# build & run via Docker (only supported path)
make build   # docker build -t go-chronos .
make run     # docker run --env-file .env -p 8080:8080 go-chronos
make logs    # tail container logs
make stop    # stop & remove container
make restart # stop + build + run
```

The binary can also be built natively: `go build -o go-chronos .`

## Environment

`.env` is **not** auto-loaded — `godotenv.Load()` is commented out in `main.go`. Set `DATABASE_URL` and `GEMINI_API_KEY` in the environment (the Makefile passes `--env-file .env` to Docker).

| Variable         | Required for            |
|------------------|-------------------------|
| `DATABASE_URL`   | all endpoints           |
| `GEMINI_API_KEY` | `POST /api/tasks/optimise` |

## API

| Method | Path                 | Description                           |
|--------|----------------------|---------------------------------------|
| POST   | `/api/tasks/`        | Create a task                         |
| GET    | `/api/tasks/`        | List all tasks                        |
| POST   | `/api/tasks/optimise`| AI-prioritise tasks (calls Gemini)    |

The optimise endpoint uses model `gemini-3-flash-preview` and requires a valid `GEMINI_API_KEY`.

## Horizon (day-of-week scheduling)

The `horizon` field stores which day a task is scheduled on:
`0=Monday, 1=Tuesday, 2=Wednesday, 3=Thursday, 4=Friday, 5=Saturday, 6=Sunday`.

Gemini sorts tasks by `priority_score` descending, slots the highest-priority task onto today, then distributes remaining tasks across weekdays with overflow. Weekend days are only used when weekdays are saturated.

## Layout

```
main.go                          — entrypoint (Gin router, DB connect, DI)
internal/
  db/db.go                       — pgxpool connection (global Pool var)
  db/taskDb.go                   — TaskDB CRUD
  handlers/tasksHandlers.go      — Gin handlers
  models/tasks.go                — Task struct + request/response types
  services/taskService.go        — business logic + Gemini API call
  utils/generatePrompt.go        — prompt builder for Gemini
schema.sql                       — PostgreSQL schema
```

DB uses `github.com/jackc/pgx/v5/pgxpool`. Global `db.Pool` is set in `Connect()`.

## No tests / no CI / no linter config

No test files, no lint or typecheck commands, no CI workflow. Add any before relying on quality gates.

## .env secrets

The `.env` file contains `DATABASE_URL` and `GEMINI_API_KEY` — do not commit. `.gitignore` already excludes `.env`.

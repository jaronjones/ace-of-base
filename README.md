# ace-of-base

A Go web application serving **NEON_OS**, a synthwave/cyberpunk dashboard
inspired by 1980s sci-fi terminals — neon palettes, scan-lines, glowing
borders, and high-density "terminal block" widgets.

## Stack

- **Go** `1.26` — `net/http` standard mux, `log/slog` structured logging
- **[templ](https://templ.guide)** `v0.3` — typed HTML templates
- **[Tailwind CSS v4](https://tailwindcss.com)** + **[daisyUI v5](https://daisyui.com)** — loaded via CDN (`@tailwindcss/browser@4` + `daisyui@5`)
- **[HTMX](https://htmx.org)** `2.0.8` — vendored at `static/js/`
- **OKLCH** color space throughout the seven themes

## Running the server

### One-time setup

```bash
cp .env.example .env
# Edit .env to set TOMORROW_IO_API_KEY (free key at https://app.tomorrow.io/)
```

`.env` is auto-loaded in dev mode (TTY-detected, or `DEV_MODE=true`). Real
environment variables always take precedence. See the
[Configuration](#configuration) table for every supported var.

### Path A — Docker compose (full app + postgres)

The simplest first run. Brings up `app` on `:8081` and `db` (postgres 17)
on `127.0.0.1:5432`. Compose synthesises `DATABASE_URL` for the app from
the `POSTGRES_*` vars (see `docker-compose.yml`), so you don't set it
yourself.

```bash
docker compose up --build       # build + start app and db
docker compose logs -f app      # tail app logs
docker compose down             # stop, keep the pgdata volume
docker compose down -v          # stop and wipe the pgdata volume
```

Migrations are embedded in the binary (`//go:embed db/migrations/*.sql` in
`main.go`) and run automatically on app boot.

### Path B — Local Go against the compose db

Fastest dev loop with a real database. Bring up just the postgres service,
point the app at it, and run `go run .` natively for instant rebuilds.

```bash
docker compose up -d db
export DATABASE_URL=postgres://aceofbase:aceofbase@127.0.0.1:5432/aceofbase?sslmode=disable
go run .
```

The DSN above matches the commented `DATABASE_URL` line in `.env.example`
— uncomment it there if you'd rather load it from `.env`.

Omit `DATABASE_URL` entirely (or leave it blank) to skip postgres
altogether: the app boots, the weather widget works, and `/health` simply
doesn't include the `database` check.

### Path C — Production-style binary

What CI/prod does — stamp the version, build a static binary, run it with
real env vars.

```bash
go run ./cmd/version                          # stamps internal/version/version.go from git
go build -o bin/ace-of-base .
DATABASE_URL=... TOMORROW_IO_API_KEY=... ./bin/ace-of-base
```

With `DEV_MODE` unset and stdout not a TTY, logs are JSON at `LOG_LEVEL`
(default `info`). Verify the server is up:

```bash
curl localhost:8081/health           # {"running":"ok",...}
curl localhost:8081/api/v1/version   # short git sha matching `git rev-parse --short HEAD`
```

### Regenerating templates

If you edit any `.templ` files you'll need the `templ` CLI:

```bash
go install github.com/a-h/templ/cmd/templ@latest
templ generate   # regenerates *_templ.go from *.templ
```

## Configuration

All configuration is via environment variables. A `.env` file in the project
root is read automatically when running in dev mode (real env vars always win).

| Variable                | Default     | Description                              |
| ----------------------- | ----------- | ---------------------------------------- |
| `HTTP_HOST`             | `0.0.0.0`   | Listen address                           |
| `HTTP_PORT`             | `8081`      | Listen port                              |
| `HTTP_READ_TIMEOUT`     | `5s`        | Request read timeout                     |
| `HTTP_WRITE_TIMEOUT`    | `10s`       | Response write timeout                   |
| `HTTP_SHUTDOWN_TIMEOUT` | `30s`       | Graceful shutdown deadline               |
| `LOG_LEVEL`             | `info`      | `debug`, `info`, `warn`, `error`         |
| `DEV_MODE`              | TTY-detect  | Pretty/text logs and `.env` loading      |
| `TOMORROW_IO_API_KEY`   | _(empty)_   | Required for the LOCAL_CLIMATE widget; without it the widget renders `API_KEY_MISSING` |
| `WEATHER_DEFAULT_LAT`   | `35.6762`   | Fallback latitude on first visit (Tokyo) |
| `WEATHER_DEFAULT_LON`   | `139.6503`  | Fallback longitude                       |
| `WEATHER_DEFAULT_LABEL` | `Tokyo, Japan` | Display label for the fallback location |
| `DATABASE_URL`          | _(empty)_   | Postgres DSN. When unset the app starts without a DB; when set, the app pings on boot, runs pending migrations, and registers a `database` health check |

## Project layout

```
.
├── main.go                  # Server bootstrap, signal handling, graceful shutdown
├── static-handler.go        # Embedded /static/ file server
├── api/                     # JSON API
│   ├── health.go            # GET /health (with pluggable health checks)
│   ├── routes.go
│   └── v1/
│       ├── routes.go
│       └── version.go       # GET /api/v1/version
├── views/                   # HTML views (templ)
│   ├── layout.templ         # <html>, CDN imports, theme bootstrap script
│   ├── dashboard.templ      # NEON_OS dashboard + theme picker
│   ├── dashboard.go         # GET /, GET /views/status (htmx fragment)
│   └── routes.go
├── static/
│   ├── css/app.css          # daisyUI v5 themes + NEON_OS custom utilities
│   └── js/htmx-2.0.8.min.js
├── internal/
│   ├── config/              # Env-var driven config + .env loader
│   ├── logging/             # slog setup (JSON in prod, text in dev)
│   ├── version/             # Build-stamped version string
│   └── weather/             # Tomorrow.io realtime client + Open-Meteo geocoder
├── cmd/version/             # Stamps internal/version/version.go from git
└── design/stitch/           # Design system spec (NEON_OS_DESIGN.md + screens)
```

## Theme picker

The dashboard navbar contains a palette icon. Clicking it opens a dropdown
with all seven NEON_OS palettes from `design/stitch/NEON_OS_DESIGN.md`:

- `synthgrid` — Electric Magenta + Cyber Cyan (default)
- `neon-horizon` — Magenta Glow + Cyan Data
- `solaris-terminal` — Solar Red + Amber Warning
- `retro-future` — Purple + Orange Sunset
- `tech-noir` — Crimson + Emerald
- `cyan-sunset` — Vibrant Cyan + Sunset Amber
- `sky-blue` — Sky Blue + Violet

Each theme is a `[data-theme="..."]` block in `static/css/app.css` defining the
daisyUI v5 `--color-*` variables (in `oklch()`) plus a `--neon-grid-rgb`
custom property used by the body-grid background. Selection is applied to
`<html data-theme>` and persisted in `localStorage` under `neon-os-theme`;
a small inline script in `<head>` restores it before paint to avoid FOUC.

## Database

Postgres runs as the `db` service in `docker compose` with a persistent
`pgdata` volume — see [Path A](#path-a--docker-compose-full-app--postgres)
or [Path B](#path-b--local-go-against-the-compose-db) above for how to
bring it up. Schema lives in `db/migrations/` (sql-migrate format) and
hand-written queries in `db/queries/` are compiled to type-safe Go via
sqlc.

Install the CLIs once:

```bash
go install github.com/rubenv/sql-migrate/sql-migrate@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Migrations (sql-migrate)

The app runs pending migrations automatically on startup (embedded via
`//go:embed db/migrations/*.sql` in `main.go`). The `sql-migrate` CLI is
still useful for scaffolding new migrations and rolling back:

```bash
sql-migrate new add_users_table   # scaffold db/migrations/<ts>-add_users_table.sql
sql-migrate up                    # apply pending migrations (also done on app startup)
sql-migrate down                  # roll back the most recent migration
sql-migrate status                # show applied / pending
```

`dbconfig.yml` defines two envs:

- `development` (default) — points at the compose db on `127.0.0.1:5432`
- `production` — reads `DATABASE_URL` from the environment; run with `-env=production`

Migration files use sql-migrate's `-- +migrate Up` / `-- +migrate Down`
markers. Multiple statements are split on `;`; use `-- +migrate StatementBegin`
/ `StatementEnd` for functions or anything containing semicolons.

### Type-safe queries (sqlc)

Add SQL to `db/queries/*.sql` with sqlc annotation comments, then:

```bash
sqlc generate    # writes Go to internal/db/sqlc/
```

`sqlc.yaml` is pre-configured for `engine: postgresql`, `sql_package:
database/sql`, and reads the schema directly from the sql-migrate files in
`db/migrations/`. Generated code is committed so CI doesn't need the sqlc
binary.

## Build & version stamping

```bash
# Stamp the build version from the current git branch + sha (and optional
# BUILD_NUMBER env var). Writes internal/version/version.go.
go run ./cmd/version

# Build the binary
go build -o bin/ace-of-base .
```

The current version is exposed at `GET /api/v1/version` and rendered into
the navbar via an htmx `hx-get` on page load.

## Endpoints

| Method | Path                  | Description                                |
| ------ | --------------------- | ------------------------------------------ |
| GET    | `/`                   | NEON_OS dashboard (HTML)                   |
| GET    | `/views/status`       | Status badge htmx fragment                 |
| GET    | `/views/weather`      | Weather widget htmx fragment (`?lat=&lon=&label=`) |
| GET    | `/views/weather/search` | Geocoding autocomplete htmx fragment (`?q=`) |
| GET    | `/health`             | JSON health check (200 OK / 503 degraded)  |
| GET    | `/api/v1/version`     | Plain-text build version                   |
| GET    | `/static/*`           | CSS, JS, fonts                             |

## Weather widget

The dashboard's `LOCAL_CLIMATE` block fetches real conditions from
[Tomorrow.io's Realtime API](https://docs.tomorrow.io/reference/realtime-weather).
Set `TOMORROW_IO_API_KEY` in `.env` (or your real environment) to enable it.

Users pick a city via the gear-icon dropdown on the widget. Typing into the
input hits `/views/weather/search`, which proxies to the no-auth
[Open-Meteo geocoding API](https://open-meteo.com/en/docs/geocoding-api).
The chosen location and the °C/°F preference are stored in browser
`localStorage` (`neon-os-weather-{lat,lon,label,units}`) and survive reloads.

The server caches Tomorrow.io responses for 5 minutes per coordinate to stay
well within the free-tier limits (25/hour, 500/day). The browser refreshes
every 30 minutes.

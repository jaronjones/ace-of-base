# Widgets

Tracking doc for NEON_OS dashboard widgets — what's live, what's faux, and
what we want to build next. Effort is a rough T-shirt size: **S** ≈ a few
hours, **M** ≈ a day, **L** ≈ multi-day with a backend integration.

## Status legend

- 🟢 **Live** — wired to real data via htmx/SSE
- 🟡 **Static** — block exists in `views/dashboard.templ` with placeholder
  copy; cosmetic only
- ⚪ **Proposed** — not built yet

---

## Currently in the dashboard

| Widget          | Status   | Notes                                                  |
| --------------- | -------- | ------------------------------------------------------ |
| `LOCAL_CLIMATE` | 🟢 Live  | Tomorrow.io realtime + Open-Meteo geocoding            |
| Status badge    | 🟢 Live  | `GET /views/status` every 5s                           |
| Version (nav)   | 🟢 Live  | `GET /api/v1/version` on load                          |
| `CPU_USAGE`     | 🟡 Static | Faux gauge — promote to real host stats               |
| `VRAM_ALLOC`    | 🟡 Static | Faux gauge — pair with CPU_USAGE under SYSTEM_VITALS  |
| `KERNEL_LOGS`   | 🟡 Static | Faux scrollback — promote to slog SSE tail            |
| `MARKET_PULSE`  | 🟡 Static | Faux ticker — wire to a real quote API                |
| `AETHER_LINK`   | 🟡 Static | Faux network meter                                     |
| `NETWORK_SCAN`  | 🟡 Static | Faux ping/scan readout                                 |
| `SYS_STABILITY` | 🟡 Static | Faux nominal/warn indicator                            |
| `UPTIME`        | 🟡 Static | Trivial to wire — server boot time + delta            |
| `QUICK_OPS`     | 🟡 Static | Decorative button row                                  |

---

## Proposed — server-introspection

We own the data, so these are the cheapest to build and double as ops
tooling.

| Widget            | Effort | Description                                                           |
| ----------------- | ------ | --------------------------------------------------------------------- |
| `TERMINAL_LOG`    | M      | SSE-streamed `slog` tail with level filter — promotes `KERNEL_LOGS`   |
| `SYSTEM_VITALS`   | S      | Host CPU / mem / load via `gopsutil` — promotes `CPU_USAGE` + `VRAM_ALLOC` |
| `REQ_FEED`        | M      | Live request log: method, path, status, ms — sparkline of req/sec    |
| `DB_PULSE`        | M      | Postgres connection count, slow-query top-N, table sizes              |
| `UPTIME_HUD`      | S      | Boot time, uptime delta, build sha, goroutine count                   |
| `BUILD_STATUS`    | S      | Last GitHub Actions run for this repo (status + commit + duration)    |

## Proposed — outward-facing data

What makes a personal dashboard feel alive day-to-day.

| Widget          | Effort | Description                                                              |
| --------------- | ------ | ------------------------------------------------------------------------ |
| `CHRONOMETER`   | S      | Multi-timezone clock + sunrise/sunset (uses weather widget's lat/lon)   |
| `MARKET_TICKER` | M      | Real stock/crypto quotes — promotes `MARKET_PULSE`                       |
| `NEWS_FEED`     | S      | RSS / HN / lobste.rs scrolling headlines                                 |
| `GH_PULSE`      | M      | GitHub repo activity: recent commits, open PRs, issues                   |
| `PACKET_TRACE`  | S      | Ping/HTTP probe a configurable list of endpoints — promotes `NETWORK_SCAN` |
| `GEO_TRACE`     | L      | World map with pins for saved weather locations / visitor IPs            |
| `CALENDAR_GRID` | M      | Month view with events (Google Calendar MCP integration)                 |
| `TASK_QUEUE`    | M      | Persistent TODO list backed by postgres (first real use of the DB)       |

## Proposed — ambient / decorative

Pure aesthetic — synthwave instruments, no real data.

| Widget          | Effort | Description                                                       |
| --------------- | ------ | ----------------------------------------------------------------- |
| `PHRASEBOOK`    | S      | Rotating quote / koan / song lyric                                |
| `AUDIO_VIZ`     | M      | Web Audio visualizer styled as a console scope                    |
| `STARFIELD`     | S      | Background canvas — slow-moving 80s vector starfield              |
| `BOOT_SEQUENCE` | S      | One-time on-load fake POST/boot scroll before the dashboard fades in |

---

## Notes

- Promote 🟡 → 🟢 by replacing the faux content with an htmx fragment route
  under `views/` (mirror the `LOCAL_CLIMATE` pattern: `views/weather.go` +
  cache + browser refresh interval).
- For SSE-based widgets (`TERMINAL_LOG`, `REQ_FEED`), add a small
  `internal/stream/` helper rather than open-coding `EventSource` per route.
- First real use of postgres is a good moment to pick `TASK_QUEUE` or
  `REQ_FEED` (request log persistence) so the schema isn't speculative.

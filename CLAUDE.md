# IslandTime Dashboard — Claude Agent Context

## Project Identity
**IslandTime** is a personal, self-hosted Go dashboard with configurable, grid-based widgets.
It serves as a single-pane-of-glass for daily utilities: weather, stocks, bookmarks, home automation, dev tools, reminders, and more.

## Tech Stack
| Layer | Choice | Notes |
|---|---|---|
| Language | Go 1.22+ | Standard library preferred |
| Frontend | HTMX | Hypermedia-first, no SPA complexity |
| CSS | Tailwind CSS (CDN) + DaisyUI (CDN) | No build step; custom tropical themes |
| Templates | `html/template` (stdlib) | Embedded via `embed.FS` |
| Config | JSON | `encoding/json` stdlib |
| Storage | In-memory → file → SQLite | Phased; start simple |

## Architecture Overview

```
cmd/server/main.go          → Entry point, wires everything
internal/
  config/       → Dashboard JSON config (layout, widget defs, theme)
  server/       → HTTP server setup, middleware, routing
  handlers/     → Dashboard page + widget HTMX fragment handlers
  widgets/      → Widget interface, registry, per-widget packages
templates/      → html/template files (base layout + per-widget)
static/         → CSS, JS, images
config/         → dashboard.json (user's layout config)
```

## Widget System

### Interface (internal/widgets/widget.go)
Every widget implements:
```go
type Widget interface {
    Type() string
    Render(r *http.Request, def Definition) (template.HTML, error)
}
```

### Registry Pattern
Widgets self-register via `init()`:
```go
func init() { widgets.Register(&ClockWidget{}) }
```

### Grid System
- 12-column CSS Grid
- Each widget has `col`, `row`, `col_span`, `row_span` in config
- Rendered as inline `grid-column` / `grid-row` styles

### HTMX Widget Loading Pattern
1. Dashboard page renders shell divs with `hx-get="/widget/{id}"`
2. `hx-trigger="load"` for initial load; add `every Ns` for live data
3. `/widget/{id}` handler looks up config → registry → calls `Render()`
4. Returns HTML fragment (no full page)

## Themes
Two custom DaisyUI themes:
- `tropical-dark` — Deep navy, hot pink primary, cyan secondary, gold accent
- `tropical-light` — Warm sand background, coral primary, teal secondary

Theme toggled via `data-theme` attribute on `<html>`, persisted in `localStorage`.

## Coding Standards
- **No external dependencies** unless stdlib genuinely cannot do it
- `errors.New` / `fmt.Errorf` with `%w` for error wrapping
- All handlers return proper HTTP status codes
- Widget errors render an error card, not a 500 page (resilience)
- Config validation on startup; crash fast with clear error messages
- Use `embed.FS` for templates in production; `os.DirFS` dev flag for live reload

## Adding a New Widget Type
1. Create `internal/widgets/{name}/{name}.go`
2. Implement `Widget` interface
3. Add `init()` that calls `widgets.Register()`
4. Create `templates/widgets/{name}.html`
5. Import the widget package (blank import) in `cmd/server/main.go`
6. Add widget definition to `config/dashboard.json`

## File Naming Conventions
- Go files: `snake_case.go`
- Templates: `snake_case.html`
- CSS classes: DaisyUI components + Tailwind utilities
- Widget type strings: `kebab-case` (e.g., `"clock"`, `"weather"`, `"url-encoder"`)

## Dev Workflow
```bash
# Run with live template reload
go run ./cmd/server -dev

# Run production (embedded templates)
go run ./cmd/server

# Build
go build -o islandtime ./cmd/server
```

## Future Considerations
- SQLite for TODO/bookmark persistence
- API key management (env vars or encrypted config field)
- Widget drag-and-drop grid editor UI
- Per-widget refresh intervals in config
- WebSocket for push updates (home automation)
- Auth (single-user htpasswd or similar)

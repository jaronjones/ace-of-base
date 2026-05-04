package views

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"jaronjones/ace-of-base/internal/config"
	"jaronjones/ace-of-base/internal/weather"
)

var (
	weatherClient   *weather.Client
	weatherDefaults config.WeatherConfig
)

// SetWeather wires the package-level weather client and defaults from main().
// Mirrors the codebase's existing init/registerRoute pattern — no DI framework.
func SetWeather(c *weather.Client, defaults config.WeatherConfig) {
	weatherClient = c
	weatherDefaults = defaults
}

// WeatherDefaults exposes the configured fallback location to templates.
func WeatherDefaults() config.WeatherConfig {
	return weatherDefaults
}

func init() {
	registerRoute(func(router *http.ServeMux) {
		router.HandleFunc("GET /views/weather", handleWeather)
		router.HandleFunc("GET /views/weather/search", handleWeatherSearch)
	})
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	latStr := q.Get("lat")
	lonStr := q.Get("lon")
	label := q.Get("label")

	if latStr == "" || lonStr == "" {
		weatherEmpty().Render(r.Context(), w)
		return
	}

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lon, err2 := strconv.ParseFloat(lonStr, 64)
	if err1 != nil || err2 != nil {
		weatherEmpty().Render(r.Context(), w)
		return
	}

	if weatherClient == nil {
		weatherError("WEATHER_DISABLED", "client not initialized").Render(r.Context(), w)
		return
	}

	data, err := weatherClient.GetRealtime(r.Context(), lat, lon)
	if err != nil {
		if errors.Is(err, weather.ErrNoAPIKey) {
			weatherError("API_KEY_MISSING", "set TOMORROW_IO_API_KEY").Render(r.Context(), w)
			return
		}
		slog.Warn("tomorrow.io fetch failed", "err", err, "lat", lat, "lon", lon)
		weatherError("UPSTREAM_OFFLINE", "").Render(r.Context(), w)
		return
	}

	if label == "" {
		label = weatherDefaults.DefaultLabel
	}
	weatherContent(data, label).Render(r.Context(), w)
}

func handleWeatherSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	places, err := weather.Search(r.Context(), q)
	if err != nil {
		slog.Warn("geocode search failed", "err", err, "q", q)
		// Render an empty list rather than 5xx so HTMX swaps cleanly.
		weatherSearchResults(nil).Render(r.Context(), w)
		return
	}
	weatherSearchResults(places).Render(r.Context(), w)
}

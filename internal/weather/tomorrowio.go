package weather

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	realtimeURL    = "https://api.tomorrow.io/v4/weather/realtime"
	requestTimeout = 5 * time.Second
	cacheTTL       = 5 * time.Minute
)

// Realtime is a normalized snapshot ready for the dashboard widget.
// Both unit systems are pre-computed so the °C/°F toggle is a pure client-side flip.
type Realtime struct {
	TempC             float64
	TempF             float64
	Humidity          float64
	PrecipProbability float64
	WindSpeedKph      float64
	WindSpeedMph      float64
	WeatherCode       int
	Condition         string
	ObservedAt        time.Time
}

// ErrNoAPIKey signals the caller didn't configure Tomorrow.io. Callers should
// render an "API_KEY_MISSING" UI state rather than a generic error.
var ErrNoAPIKey = errors.New("tomorrow.io API key not configured")

type Client struct {
	APIKey string
	HTTP   *http.Client

	mu    sync.Mutex
	cache map[string]cacheEntry
}

type cacheEntry struct {
	data      Realtime
	expiresAt time.Time
}

func (c *Client) GetRealtime(ctx context.Context, lat, lon float64) (Realtime, error) {
	if c.APIKey == "" {
		return Realtime{}, ErrNoAPIKey
	}

	key := cacheKey(lat, lon)
	if hit, ok := c.cacheGet(key); ok {
		return hit, nil
	}

	q := url.Values{}
	q.Set("location", fmt.Sprintf("%.4f,%.4f", lat, lon))
	q.Set("units", "metric")
	q.Set("apikey", c.APIKey)

	reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, realtimeURL+"?"+q.Encode(), nil)
	if err != nil {
		return Realtime{}, err
	}
	req.Header.Set("Accept", "application/json")

	httpClient := c.HTTP
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return Realtime{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Realtime{}, fmt.Errorf("tomorrow.io returned %s", resp.Status)
	}

	var raw struct {
		Data struct {
			Time   time.Time `json:"time"`
			Values struct {
				Temperature              float64 `json:"temperature"`
				Humidity                 float64 `json:"humidity"`
				PrecipitationProbability float64 `json:"precipitationProbability"`
				WindSpeed                float64 `json:"windSpeed"`
				WeatherCode              int     `json:"weatherCode"`
			} `json:"values"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return Realtime{}, fmt.Errorf("decode tomorrow.io response: %w", err)
	}

	v := raw.Data.Values
	rt := Realtime{
		TempC:             round1(v.Temperature),
		TempF:             round1(v.Temperature*9/5 + 32),
		Humidity:          round1(v.Humidity),
		PrecipProbability: round1(v.PrecipitationProbability),
		WindSpeedKph:      round1(v.WindSpeed * 3.6),
		WindSpeedMph:      round1(v.WindSpeed * 2.23694),
		WeatherCode:       v.WeatherCode,
		Condition:         describeCode(v.WeatherCode),
		ObservedAt:        raw.Data.Time,
	}

	c.cachePut(key, rt)
	return rt, nil
}

func (c *Client) cacheGet(key string) (Realtime, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil {
		return Realtime{}, false
	}
	e, ok := c.cache[key]
	if !ok || time.Now().After(e.expiresAt) {
		return Realtime{}, false
	}
	return e.data, true
}

func (c *Client) cachePut(key string, rt Realtime) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cache == nil {
		c.cache = make(map[string]cacheEntry)
	}
	c.cache[key] = cacheEntry{data: rt, expiresAt: time.Now().Add(cacheTTL)}
}

func cacheKey(lat, lon float64) string {
	return strconv.FormatFloat(round2(lat), 'f', 2, 64) + "," + strconv.FormatFloat(round2(lon), 'f', 2, 64)
}

func round1(v float64) float64 { return math.Round(v*10) / 10 }
func round2(v float64) float64 { return math.Round(v*100) / 100 }

// describeCode maps Tomorrow.io's `weatherCode` enumeration to a short
// human-readable label. Values per the public Data Layers reference.
func describeCode(code int) string {
	if s, ok := weatherCodeMap[code]; ok {
		return s
	}
	return "Unknown"
}

var weatherCodeMap = map[int]string{
	0:    "Unknown",
	1000: "Clear",
	1001: "Cloudy",
	1100: "Mostly Clear",
	1101: "Partly Cloudy",
	1102: "Mostly Cloudy",
	2000: "Fog",
	2100: "Light Fog",
	3000: "Light Wind",
	3001: "Wind",
	3002: "Strong Wind",
	4000: "Drizzle",
	4001: "Rain",
	4200: "Light Rain",
	4201: "Heavy Rain",
	5000: "Snow",
	5001: "Flurries",
	5100: "Light Snow",
	5101: "Heavy Snow",
	6000: "Freezing Drizzle",
	6001: "Freezing Rain",
	6200: "Light Freezing Rain",
	6201: "Heavy Freezing Rain",
	7000: "Ice Pellets",
	7101: "Heavy Ice Pellets",
	7102: "Light Ice Pellets",
	8000: "Thunderstorm",
}

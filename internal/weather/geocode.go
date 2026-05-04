package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const geocodeURL = "https://geocoding-api.open-meteo.com/v1/search"

type Place struct {
	Name    string
	Admin1  string
	Country string
	Lat     float64
	Lon     float64
}

// Label returns a single-line "City, Region, Country" label suitable for the UI.
// Region is omitted when it duplicates the city name (e.g. "Tokyo, Tokyo").
func (p Place) Label() string {
	parts := []string{p.Name}
	if p.Admin1 != "" && p.Admin1 != p.Name {
		parts = append(parts, p.Admin1)
	}
	if p.Country != "" {
		parts = append(parts, p.Country)
	}
	return strings.Join(parts, ", ")
}

// Search queries Open-Meteo's free geocoding API. No key required.
func Search(ctx context.Context, q string) ([]Place, error) {
	q = strings.TrimSpace(q)
	if len(q) < 2 {
		return nil, nil
	}

	// Open-Meteo matches plain city names, not "City, Region" — strip after the
	// first comma so a query like "Lehi, Utah" still returns matches.
	name := q
	if i := strings.Index(name, ","); i > 0 {
		name = strings.TrimSpace(name[:i])
	}
	if len(name) < 2 {
		return nil, nil
	}

	params := url.Values{}
	params.Set("name", name)
	params.Set("count", "5")
	params.Set("language", "en")
	params.Set("format", "json")

	reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, geocodeURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("open-meteo returned %s", resp.Status)
	}

	var raw struct {
		Results []struct {
			Name      string  `json:"name"`
			Admin1    string  `json:"admin1"`
			Country   string  `json:"country"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode open-meteo response: %w", err)
	}

	out := make([]Place, 0, len(raw.Results))
	for _, r := range raw.Results {
		out = append(out, Place{
			Name:    r.Name,
			Admin1:  r.Admin1,
			Country: r.Country,
			Lat:     r.Latitude,
			Lon:     r.Longitude,
		})
	}
	return out, nil
}


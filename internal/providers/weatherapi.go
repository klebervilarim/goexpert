package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type WeatherAPIClient struct {
	http *http.Client
	base string
	key  string
}

func NewWeatherAPIClient(h *http.Client, base, key string) *WeatherAPIClient {
	return &WeatherAPIClient{http: h, base: base, key: key}
}

type weatherCurrentResp struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (c *WeatherAPIClient) CurrentTempC(ctx context.Context, query string) (float64, error) {
	if c.key == "" {
		return 0, fmt.Errorf("missing WEATHERAPI_KEY")
	}
	u := fmt.Sprintf("%s/v1/current.json?key=%s&q=%s&aqi=no",
		c.base, url.QueryEscape(c.key), url.QueryEscape(query))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	res, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi status %d", res.StatusCode)
	}

	var wr weatherCurrentResp
	if err := json.NewDecoder(res.Body).Decode(&wr); err != nil {
		return 0, err
	}
	return wr.Current.TempC, nil
}

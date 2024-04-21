package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/softstone1/woc/domain"
)

const (
	// client request timeout
	timeout               = 30 * time.Second
	maxIdleConns          = 5
	maxIdleConnsPerHost   = 3
	dialTimeout           = time.Second
	keepAlive             = time.Minute
	responseHeaderTimeout = time.Second
	tlsHandshakeTimeout   = 2 * time.Second
)

type OpenMeteo struct {
	baseUrl string
	client  *http.Client
}

func NewOpenMeteo(url string) *OpenMeteo {
	return &OpenMeteo{
		baseUrl: url,
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        maxIdleConns,
				MaxIdleConnsPerHost: maxIdleConnsPerHost,
				DialContext: (&net.Dialer{
					Timeout:   dialTimeout,
					KeepAlive: keepAlive,
				}).DialContext,
				ResponseHeaderTimeout: responseHeaderTimeout,
				TLSHandshakeTimeout:   tlsHandshakeTimeout,
			},
		},
	}
}

type WeatherReponse struct {
	Hourly struct {
        Temperature2m []float64 `json:"temperature_2m"`
		WindSpeed10m  []float64 `json:"wind_speed_10m"`
    } `json:"hourly"`
}

func (c *OpenMeteo) FetchWeatherByCity(ctx context.Context, city domain.City) (*domain.Weather, error) {
	url := fmt.Sprintf("%s/v1/forecast?latitude=%s&longitude=%s&hourly=temperature_2m,wind_speed_10m", c.baseUrl, city.Latitude, city.Longitude)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var data WeatherReponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &domain.Weather{
		City:        city.Name,
		Temperature: data.Hourly.Temperature2m[0],
		WindSpeed:   data.Hourly.WindSpeed10m[0],
	}, nil

}

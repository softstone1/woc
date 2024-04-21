// Business logic and data model weather
package domain

import "context"

type Weather struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"windSpeed"`
}

type WeatherClient interface {
	FetchWeatherByCity(ctx context.Context, city City) (*Weather, error)
}

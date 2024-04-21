// application services
package app

import (
	"context"
	"github.com/softstone1/woc/domain"
)
type WeatherService interface {
	GetWeatherByCity(ctx context.Context, cityName string) (*domain.Weather, error)
	GetAllCities() ([]domain.City, error)
}

type weatherService struct {
	client domain.WeatherClient
	cityRepository domain.CityRepository
}

func NewWeatherService(weatherClient domain.WeatherClient, cityRepository domain.CityRepository) *weatherService {
	return &weatherService{
		client: weatherClient,
		cityRepository: cityRepository,
	}
}

func (s *weatherService) GetWeatherByCity(ctx context.Context, cityName string) (*domain.Weather, error) {
	city, err := s.cityRepository.GetCity(cityName)
	if err != nil {
		return nil, err
	}
	return s.client.FetchWeatherByCity(ctx, *city)
}

func (s *weatherService) GetAllCities() ([]domain.City, error) {
	return s.cityRepository.GetAllCities()
}
package app

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/softstone1/woc/domain"
	gomock "go.uber.org/mock/gomock"
)

func TestWeatherService_GetWeatherByCity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockWeatherClient := domain.NewMockWeatherClient(mockCtrl)
	mockCityRepository := domain.NewMockCityRepository(mockCtrl)
	service := NewWeatherService(mockWeatherClient, mockCityRepository)

	tests := []struct {
		name          string
		cityName      string
		setupMocks    func()
		expectedWeather *domain.Weather
		expectedErr   error
	}{
		{
			name:     "successful weather fetch",
			cityName: "Berlin",
			setupMocks: func() {
				mockCity := &domain.City{Name: "Berlin", Latitude: "52.5200", Longitude: "13.4050"}
				mockWeather := &domain.Weather{City: "Berlin", Temperature: 20.5, WindSpeed: 5.0}
				mockCityRepository.EXPECT().GetCity("Berlin").Return(mockCity, nil)
				mockWeatherClient.EXPECT().FetchWeatherByCity(gomock.Any(), *mockCity).Return(mockWeather, nil)
			},
			expectedWeather: &domain.Weather{City: "Berlin", Temperature: 20.5, WindSpeed: 5.0},
			expectedErr:     nil,
		},
		{
			name:     "city not found error",
			cityName: "Unknown",
			setupMocks: func() {
				mockCityRepository.EXPECT().GetCity("Unknown").Return(nil, errors.New("city not found"))
			},
			expectedWeather: nil,
			expectedErr:     errors.New("city not found"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			weather, err := service.GetWeatherByCity(context.Background(), tc.cityName)
			if !reflect.DeepEqual(err, tc.expectedErr) {
				t.Errorf("%s: expected error %v, got %v", tc.name, tc.expectedErr, err)
			}
			if !reflect.DeepEqual(weather, tc.expectedWeather) {
				t.Errorf("%s: expected weather %v, got %v", tc.name, tc.expectedWeather, weather)
			}
		})
	}
}

func TestWeatherService_GetAllCities(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCityRepository := domain.NewMockCityRepository(mockCtrl)
	service := NewWeatherService(nil, mockCityRepository) // WeatherClient is not used here

	tests := []struct {
		name           string
		setupMocks     func()
		expectedCities []domain.City
		expectedErr    error
	}{
		{
			name: "get all cities successfully",
			setupMocks: func() {
				expectedCities := []domain.City{{Name: "Berlin", Latitude: "52.5200", Longitude: "13.4050"}}
				mockCityRepository.EXPECT().GetAllCities().Return(expectedCities, nil)
			},
			expectedCities: []domain.City{{Name: "Berlin", Latitude: "52.5200", Longitude: "13.4050"}},
			expectedErr:    nil,
		},
		{
			name: "failure in fetching cities",
			setupMocks: func() {
				mockCityRepository.EXPECT().GetAllCities().Return(nil, errors.New("database error"))
			},
			expectedCities: nil,
			expectedErr:    errors.New("database error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks()
			}

			cities, err := service.GetAllCities()
			if !reflect.DeepEqual(err, tc.expectedErr) {
				t.Errorf("%s: expected error %v, got %v", tc.name, tc.expectedErr, err)
			}
			if !reflect.DeepEqual(cities, tc.expectedCities) {
				t.Errorf("%s: expected cities %v, got %v", tc.name, tc.expectedCities, cities)
			}
		})
	}
}


package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/softstone1/woc/domain"
)

func TestFetchWeatherByCity(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the URL and respond with mock data
		if r.URL.Path == "/v1/forecast" {
			fmt.Fprint(w, `{"hourly": {"temperature_2m": [25.5], "wind_speed_10m": [10.2]}}`)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Create a new OpenMeteo client with the test server URL
	client := NewOpenMeteo(server.URL)

	// Define test cases
	testCases := []struct {
		name           string
		city           domain.City
		expectedWeather *domain.Weather
	}{
		{
			name: "Test case 1",
			city: domain.City{
				Name:      "New York",
				Latitude:  "40.7128",
				Longitude: "-74.0060",
			},
			expectedWeather: &domain.Weather{
				City:        "New York",
				Temperature: 25.5,
				WindSpeed:   10.2,
			},
		},
		{
			name: "Test case 2",
			city: domain.City{
				Name:      "TestCity",
				Latitude:  "123.456",
				Longitude: "789.012",
			},
			expectedWeather: &domain.Weather{
				City:        "TestCity",
				Temperature: 25.5,
				WindSpeed:   10.2,
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the FetchWeatherByCity function
			weather, err := client.FetchWeatherByCity(context.Background(), tc.city)

			// Check the result
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if weather == nil {
				t.Errorf("Expected weather to be fetched, but it was not")
			} else {
				if !reflect.DeepEqual(weather, tc.expectedWeather) {
					t.Errorf("Expected weather %v, but got %v", tc.expectedWeather, weather)
				}
			}
		})
	}
}

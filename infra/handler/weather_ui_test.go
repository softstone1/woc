package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/softstone1/woc/app"
	"github.com/softstone1/woc/domain"
	"go.uber.org/mock/gomock"
)

func TestHome(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockWeatherService := app.NewMockWeatherService(mockCtrl)
	weatherHandler := NewWeather(mockWeatherService)

	tests := []struct {
		name           string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful Home Load",
			mockSetup: func() {
				mockWeatherService.EXPECT().GetAllCities().Return([]domain.City{
					{
						Name:      "New York",
						Latitude:  "40.7128",
						Longitude: "-74.0060",
					},
					{
						Name:      "London",
						Latitude:  "51.5074",
						Longitude: "-0.1278",
					},
					{
						Name:      "Paris",
						Latitude:  "48.8566",
						Longitude: "2.3522",
					},
					{
						Name:      "Tokyo",
						Latitude:  "35.6895",
						Longitude: "139.6917",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "New York,London,Paris,Tokyo",
		},
		{
			name: "Failed Home Load",
			mockSetup: func() {
				mockWeatherService.EXPECT().GetAllCities().Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(weatherHandler.Home)
			if tc.mockSetup != nil {
				tc.mockSetup()
			}
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			//Only check for city names if the status code is 200 (OK)
			if rr.Code == http.StatusOK {
				// Check that the body contains the city names, rather than being equal to them
				body := rr.Body.String()
				for _, city := range []string{"New York", "London", "Paris", "Tokyo"} {
					if !strings.Contains(body, city) {
						t.Errorf("Expected body to contain %q, but it did not", city)
					}
				}
			} else if rr.Code == http.StatusInternalServerError {
				// Check for an error message in the body
				body := rr.Body.String()
				if !strings.Contains(body, "error") {
					t.Errorf("Expected body to contain %q, but it did not", "error")
				}
			}
		})
	}
}

func TestGetWeatherByCity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockWeatherService := app.NewMockWeatherService(mockCtrl)
	weatherHandler := NewWeather(mockWeatherService)

	tests := []struct {
		name           string
		city           string
		mockSetup      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid City Request",
			city: "London",
			mockSetup: func() {
				mockWeatherService.EXPECT().
					GetWeatherByCity(gomock.Any(), "London").
					Return(&domain.Weather{City: "London", Temperature: 15.5}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "<div>\n    <h2>Weather for London</h2>\n    <p>Temperature: 15.5°C</p>\n    <p>Windspeed: 0 km/h</p>\n</div>", 
		},
		{
			name:           "City Missing in Request",
			city:           "",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing city query parameter",
		},
		{
			name: "Weather Service Error",
			city: "Unknown",
			mockSetup: func() {
				mockWeatherService.EXPECT().
					GetWeatherByCity(gomock.Any(), "Unknown").
					Return(nil, errors.New("city not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "city not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/weather?city="+tc.city, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(weatherHandler.GetWeatherByCity)
			if tc.mockSetup != nil {
				tc.mockSetup()
			}
			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			body := strings.TrimSuffix(rr.Body.String(), "\n")
			if body != tc.expectedBody {
				t.Errorf("Expected body %q, got %q", tc.expectedBody, body)
			}
		})
	}
}

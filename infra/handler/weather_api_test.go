package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/softstone1/woc/app"
	"github.com/softstone1/woc/domain"
)

func TestGetWeatherByCityAPI(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockWeatherService := app.NewMockWeatherService(mockCtrl)
	weatherHandler := NewWeather(mockWeatherService)

	tests := []struct {
		name           string
		city           string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid City",
			city: "London",
			setupMock: func() {
				mockWeatherService.EXPECT().
					GetWeatherByCity(gomock.Any(), "London").
					Return(&domain.Weather{City: "London", Temperature: 6.7, WindSpeed: 5.5}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"city": "London", "temperature": 6.7, "windSpeed": 5.5}`,
		},
		{
			name:           "City Missing",
			city:           "",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing city query parameter\n",
		},
		{
			name: "Service Error",
			city: "Unknown",
			setupMock: func() {
				mockWeatherService.EXPECT().
					GetWeatherByCity(gomock.Any(), "Unknown").
					Return(nil, errors.New("city not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "city not found\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "/weather?city="+tc.city, nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			if tc.setupMock != nil {
				tc.setupMock()
			}

			weatherHandler.GetWeatherByCityAPI(recorder, request)

			if recorder.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, recorder.Code)
			}

			// Normalize JSON strings by removing spaces for comparison
			expectedBody := tc.expectedBody
			actualBody := recorder.Body.String()

			var buf1, buf2 bytes.Buffer
			json.Compact(&buf1, []byte(expectedBody))
			json.Compact(&buf2, []byte(actualBody))

			if buf1.String() != buf2.String() {
				t.Errorf("Expected body %q, got %q", buf1.String(), buf2.String())
			}
		})
	}
}

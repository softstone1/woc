package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/softstone1/woc/app"
)

type Weather struct {
	weatherService app.WeatherService
}

func NewWeather(weatherService app.WeatherService) *Weather {
	return &Weather{
		weatherService: weatherService,
	}
}

// GetWeatherByCityAPI returns weather information for a given city.
func (h *Weather) GetWeatherByCityAPI(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()
	cityName := r.URL.Query().Get("city")
	if cityName == "" {
		http.Error(w, "missing city query parameter", http.StatusBadRequest)
		return
	}
	weather, err := h.weatherService.GetWeatherByCity(ctx, cityName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, weather)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

package handler

import (
	"context"
	"embed"
	"net/http"
	"text/template"
	"time"
)
var (
	//go:embed templates/*.gohtml
	FS embed.FS
	tmpl = template.Must(template.ParseFS(FS, "templates/*.gohtml"))
)

// Home is the handler for the home page
func (h *Weather) Home(w http.ResponseWriter, r *http.Request) {
	cities, err := h.weatherService.GetAllCities()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "home.gohtml", cities); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetWeatherByCity is the handler for the weather page
func (h *Weather) GetWeatherByCity(w http.ResponseWriter, r *http.Request) {
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
	if err := tmpl.ExecuteTemplate(w, "weather.gohtml", weather); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

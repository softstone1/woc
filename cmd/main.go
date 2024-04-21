package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/softstone1/woc/app"
	"github.com/softstone1/woc/config"
	"github.com/softstone1/woc/infra/client"
	"github.com/softstone1/woc/infra/db"
	"github.com/softstone1/woc/infra/handler"
	"github.com/softstone1/woc/infra/server"
)

func main() {
	// Set up default logger to slog with JSON format
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
	// Load the environment variables
	slog.Info("loading environment variables")
	config.LoadEnv()

	/* Dependency injection */

	// Create a new weather client using the OpenMeteo API
	weatherClient := client.NewOpenMeteo(config.GetEnv().WeatherBaseURL())
	// Create in-memory city repository
	cityRepo := db.NewInMemoryCityRepository()

	// Create a new weather service
	weatherService := app.NewWeatherService(weatherClient, cityRepo)

	// Create weather handler
	weatherHandler := handler.NewWeather(weatherService)

	// Create a new server
	server, err := server.NewMux(config.GetEnv(), weatherHandler)
	if err != nil {
		slog.Error("error creating server", "error", err)
		os.Exit(1)
	}

	// Run the server with graceful shutdown
	if err := server.Run(context.Background()); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}

}

package config

import "github.com/spf13/viper"

const (
	serverPort = "SERVER_PORT"
	weatherBaseURL = "WEATHER_BASE_URL"
	enableProfiling = "ENABLE_PROFILING"
)

type Env struct {
	ServerPort func() string	
	EnableProfiling func() bool
	WeatherBaseURL func() string
}

func GetEnv() Env {
	return Env{
		ServerPort: func() string {
			return viper.GetString(serverPort)
			},
		EnableProfiling: func() bool {
			return viper.GetBool(enableProfiling)
			},
		WeatherBaseURL: func() string {
			return viper.GetString(weatherBaseURL)
		},
	}
}

func LoadEnv() {
	//viber 
	viper.AutomaticEnv()
	viper.SetDefault(serverPort, "8080")
	viper.SetDefault(enableProfiling, false)
	viper.SetDefault(weatherBaseURL, "https://api.open-meteo.com")
}

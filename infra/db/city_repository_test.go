package db

import (
	"testing"

	"github.com/softstone1/woc/domain"
)

func TestGetCity(t *testing.T) {
	repo := NewInMemoryCityRepository()

	// Test case 1: City found
	cityName := "New York"
	expectedCity := &domain.City{
		Name:      "New York",
		Latitude:  "40.7128",
		Longitude: "-74.0060",
	}
	city, err := repo.GetCity(cityName)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if city == nil {
		t.Errorf("Expected city to be found, but it was not")
	} else if *city != *expectedCity {
		t.Errorf("Expected city %v, but got %v", expectedCity, city)
	}

	// Test case 2: City not found
	unknownCityName := "Unknown"
	_, err = repo.GetCity(unknownCityName)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	} else if err.Error() != "city not found" {
		t.Errorf("Expected error message 'city not found', but got '%v'", err.Error())
	}
}

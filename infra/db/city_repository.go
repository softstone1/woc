package db

import (
    "errors"
    "github.com/softstone1/woc/domain"
)

// InMemoryCityRepository is an in-memory implementation of CityRepository.
type InMemoryCityRepository struct {
    cities map[string]domain.City
}

// NewInMemoryCityRepository creates a new instance of InMemoryCityRepository with preloaded data.
func NewInMemoryCityRepository() *InMemoryCityRepository {
    return &InMemoryCityRepository{
        cities: map[string]domain.City{
            "Tokyo":    {Name: "Tokyo", Latitude: "35.6895", Longitude: "139.6917"},
            "New York": {Name: "New York", Latitude: "40.7128", Longitude: "-74.0060"},
            "London":   {Name: "London", Latitude: "51.5074", Longitude: "-0.1278"},
            "Paris":    {Name: "Paris", Latitude: "48.8566", Longitude: "2.3522"},
        },
    }
}

// GetCity retrieves city information by name.
func (repo *InMemoryCityRepository) GetCity(name string) (*domain.City, error) {
    if city, ok := repo.cities[name]; ok {
        return &city, nil
    }
    return nil, errors.New("city not found")
}

// Returns all cities in the repository.
func (repo *InMemoryCityRepository) GetAllCities() ([]domain.City, error) {
    allCities := make([]domain.City, 0, len(repo.cities))
    for _, city := range repo.cities {
        allCities = append(allCities, city)
    }
    return allCities, nil
}
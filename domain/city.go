package domain

// City represents city data with coordinates.
type City struct {
    Name      string
    Latitude  string
    Longitude string
}

// CityRepository defines the interface for accessing city data.
type CityRepository interface {
    GetCity(name string) (*City, error)
    GetAllCities() ([]City, error)
}
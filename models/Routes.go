package models

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/JojiiOfficial/gaw"
)

// Route a reverseproxy route
type Route struct {
	ServerName string
	Ports      []uint16
	Locations  []RouteLocation
}

// RouteLocation location for route
type RouteLocation struct {
	Location    string
	Destination string
}

// CreateExampleRoute creates an example route
func CreateExampleRoute(file string) error {
	// Check if already exists
	fStat, err := os.Stat(file)
	if err == nil && fStat.Size() > 0 {
		return nil
	}

	// Example route
	r := Route{
		ServerName: "localhost",
		Ports:      []uint16{80, 443},
		Locations: []RouteLocation{
			RouteLocation{
				Location:    "/",
				Destination: "http://127.0.0.1/",
			},
		},
	}

	// Create path if neccessary
	err = gaw.CreatePath(file, 0740)
	if err != nil {
		return err
	}

	// Create file
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	return toml.NewEncoder(f).Encode(r)
}

// LoadRoute loads route
func LoadRoute(file string) (*Route, error) {
	var route Route

	// Read file
	if _, err := toml.DecodeFile(file, &route); err != nil {
		return nil, err
	}

	return &route, nil
}

// Check checks a route for errors. Returns true on success
func (route Route) Check() bool {
	return true
}

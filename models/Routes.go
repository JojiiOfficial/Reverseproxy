package models

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/JojiiOfficial/gaw"
)

// Route a reverseproxy route
type Route struct {
	ServerName string
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
		ServerName: "test.de",
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

package models

import (
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/JojiiOfficial/gaw"
	log "github.com/sirupsen/logrus"
)

// Route a reverseproxy route
type Route struct {
	FileName        string `toml:"-"`
	ServerNames     []string
	Addresses       []string         `toml:"Addresses"`
	ListenAddresses []*ListenAddress `toml:"-"`
	SSL             TLSKeyCertPair
	Locations       []RouteLocation
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
		FileName: gaw.FileFromPath(file),
		ServerNames: []string{
			"localhost",
		},
		Addresses: []string{
			"127.0.0.1:80",
			"127.0.0.1:443",
		},
		Locations: []RouteLocation{
			RouteLocation{
				Location:    "/",
				Destination: "http://127.0.0.1/",
			},
			RouteLocation{
				Location:    "/subroute",
				Destination: "http://127.0.0.1/lol.html",
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

	// Set filename
	route.FileName = gaw.FileFromPath(file)

	return &route, nil
}

// Check checks a route for errors. Returns true on success
func (route Route) Check(config *Config) bool {
	// Check ssl config if used
	if route.NeedSSL() {
		if len(route.SSL.Cert) == 0 || len(route.SSL.Key) == 0 {
			log.Error("Missing SSL config: SSL cert or key")
			return false
		}

		// Check cert exists
		if !gaw.FileExists(route.SSL.Cert) {
			log.Error("SSL cert file not found")
			return false
		}

		// Check key exists
		if !gaw.FileExists(route.SSL.Key) {
			log.Error("SSL key file not found")
			return false
		}
	}

	return true
}

// NeedSSL return true if route needs ssl
func (route Route) NeedSSL() bool {
	for _, address := range route.ListenAddresses {
		if address.SSL {
			return true
		}
	}

	return false
}

// LoadAddress loads an address from config. Returns false if at least one address was not found
func (route *Route) LoadAddress(config *Config) bool {
	var addresses []*ListenAddress
	for _, add := range route.Addresses {
		address := config.GetAddress(add)
		if address.Port == 0 {
			return false
		}

		addresses = append(addresses, config.GetAddress(add))
	}

	// Set addresses
	route.ListenAddresses = addresses

	return true
}

//GetTLSCerts get all required certifitates/keys from routes
func GetTLSCerts(routes []Route, address *ListenAddress) []TLSKeyCertPair {
	var pairs []TLSKeyCertPair

	// Loop routes and addresses to find all matching keys/certs
	for _, route := range routes {
		for i := range route.ListenAddresses {
			if route.ListenAddresses[i] == address {
				pairs = append(pairs, TLSKeyCertPair{
					Cert: route.SSL.Cert,
					Key:  route.SSL.Key,
				})
			}
		}
	}

	return pairs
}

// Handle http handler function
func (route *Route) Handle(w http.ResponseWriter, r *http.Request) {

}

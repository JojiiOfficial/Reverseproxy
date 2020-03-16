package models

import (
	"net/http"
	"os"
	"strconv"
	"strings"

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
	Locations       []RouteLocation `toml:"Location"`
	DefaultLocation *RouteLocation  `toml:"-"`
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
			"localhost.de",
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

	// Init locations
	for i := range route.Locations {
		route.Locations[i].Init()
		if route.Locations[i].Location == "/" {
			route.DefaultLocation = &route.Locations[i]
		}
	}

	// Make servernames toLower
	sliceToLower(route.ServerNames)

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

	// Validate locations
	for _, location := range route.Locations {
		if !isURLValid(location.Destination) {
			log.Fatalf("Location in %s is malformed!\n", route.FileName)
			return false
		}

		// Check for http request loops
		for _, address := range config.ListenAddresses {
			var p uint64
			if len(location.DestinationURL.Port()) > 0 {
				p, _ = strconv.ParseUint(location.DestinationURL.Port(), 10, 16)
			} else {
				p = 80
			}

			if location.DestinationURL.Hostname() == address.Address && p == uint64(address.Port) {
				log.Fatal("Error Request loop detected")
				return false
			}
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

// GetRoutesFromAddress gets all routes assigned to an address
func GetRoutesFromAddress(routes []Route, address ListenAddress) []*Route {
	var retRoutes []*Route

	for i := range routes {
		if routes[i].HasAddress(address) {
			retRoutes = append(retRoutes, &routes[i])
		}
	}

	return retRoutes
}

// HasAddress return true if route has an address
func (route Route) HasAddress(address ListenAddress) bool {
	for i := range route.Addresses {
		if *route.ListenAddresses[i] == address {
			return true
		}
	}

	return false
}

// FindMatchingLocation finds location
func FindMatchingLocation(routes []*Route, r *http.Request) *RouteLocation {
	pathItems := trunSlice(strings.Split(r.URL.Path, "/"))

	for _, route := range routes {
		// Match hostname
		if !inStrSl(route.ServerNames, r.URL.Hostname()) {
			continue
		}

		// Find matching route
		found := findMatchingLocation(pathItems, route.Locations)
		if found != nil {
			return found
		}

		// Otherwise use default location
		if route.DefaultLocation != nil {
			return route.DefaultLocation
		}
	}

	// Return nil if nothing was found
	return nil
}

func inStrSl(ss []string, str string) bool {
	for _, s := range ss {
		if str == strings.ToLower(s) {
			return true
		}
	}
	return false
}

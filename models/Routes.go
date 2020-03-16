package models

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/JojiiOfficial/gaw"
	log "github.com/sirupsen/logrus"
)

//RegexpStore store regexes
var RegexpStore = NewRegexStore()

// Route a reverseproxy route
type Route struct {
	FileName        string `toml:"-"`
	ServerNames     []string
	Interfaces      []string
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
			"localhost",
			"127.0.0.1",
		},
		Interfaces: []string{
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
				Allow: []string{
					"127.0.0.1",
					"1.1.1.1",
				},
				Deny: "all",
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
		route.Locations[i].Init(&route)
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
			log.Errorf("SSL cert file '%s' not found", route.SSL.Cert)
			return false
		}

		// Check key exists
		if !gaw.FileExists(route.SSL.Key) {
			log.Errorf("SSL key file '%s' not found", route.SSL.Key)
			return false
		}
	}

	// Validate locations
	for _, location := range route.Locations {
		if !isURLValid(location.Destination) {
			log.Fatalf("Location in %s is malformed!\n", route.FileName)
			return false
		}

		// Check if location points to reverseproxies address
		port := "80"
		if len(location.DestinationURL.Port()) > 0 {
			port = location.DestinationURL.Port()
		}
		if isHostsAddress(location.DestinationURL.Hostname()) && gaw.IsInStringArray(port, location.Ports()) {
			log.Fatal("Error Request loop detected")
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
	for _, iface := range route.Interfaces {
		address := config.GetAddress(iface)
		if address.Address == "" {
			return false
		}

		addresses = append(addresses, config.GetAddress(iface))
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
	for i := range route.Interfaces {
		if route.Interfaces[i] == address.GetAddress() {
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
		log.Debug(r.URL.String(), " -> ", found.DestinationURL.String())

		if found != nil {
			found.Route = route
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

// Return true if given address belongs to host adresses
func isHostsAddress(address string) bool {
	if gaw.IsInStringArray(address, []string{"[::1]", "127.0.0.1"}) {
		return true
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return false
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return false
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil {
				continue
			}

			if ip.To4() != nil {
				if ip.To4().String() == address {
					return true
				}
			}

			if ip.To16() != nil {
				if ip.To16().String() == address {
					return true
				}
			}
		}
	}
	return false
}

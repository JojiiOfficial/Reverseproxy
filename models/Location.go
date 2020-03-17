package models

import (
	"net/http"
	"net/url"
	"path"
	"strings"
)

// RouteLocation location for route
type RouteLocation struct {
	// Toml config attributes
	Location    string
	Destination string
	SrcIPHeader string
	Regex       bool

	// Allow/deny hosts
	Allow []string
	Deny  string

	// Non toml attrs
	DestinationURL *url.URL `toml:"-"`
	Route          *Route   `toml:"-"`
	HasDenyRoule   bool     `toml:"-"`
}

// Init inits a location. Gets called on loading its assigned route
func (location *RouteLocation) Init(route *Route) {
	location.Route = route
	location.HasDenyRoule = strings.ToLower(location.Deny) == "all"
	location.DestinationURL, _ = url.Parse(location.Destination)
}

//Ports returns a list with ports used by the given RouteLocation
func (location *RouteLocation) Ports() []string {
	var ports []string
	for _, address := range location.Route.ListenAddresses {
		ports = append(ports, address.GetPort())
	}
	return ports
}

// ModifyProxyRequest modifies a request to a proxy forward request
func (location *RouteLocation) ModifyProxyRequest(req *http.Request) {
	// Set new host & scheme
	req.URL.Scheme = location.DestinationURL.Scheme
	req.URL.Host = location.DestinationURL.Host

	// Build Path
	if strings.HasSuffix(location.DestinationURL.Path, "/") {
		req.URL.Path = path.Join(location.DestinationURL.Path, (req.URL.Path[len(location.Location):]))
	} else {
		req.URL.Path = path.Join(location.DestinationURL.Path, req.URL.Path)
	}

	targetQuery := location.DestinationURL.RawQuery
	if targetQuery == "" || req.URL.RawQuery == "" {
		// Add Query
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		// Add locations custom-query if set
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}

	location.finalMods(req)
}

func (location *RouteLocation) finalMods(req *http.Request) {
	// explicitly disable User-Agent so it's not set to default value
	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", "")
	}
}

func findMatchingLocation(pathItems []string, locations []RouteLocation) *RouteLocation {
	for index := range locations {
		locationItems := trunSlice(strings.Split(locations[index].Location, "/"))

		// Try to find anything else than /
		if len(locationItems) == 0 {
			continue
		}

		// Calc match depth and compare to required depth. If >=, add to matching items
		if calcMatchDepth(pathItems, locationItems, locations[index].Regex) >= len(locationItems) {
			return &locations[index]
		}
	}

	return nil
}

func calcMatchDepth(path, location []string, regex bool) int {
	matchCount := 0
	for i := range location {
		if len(location) <= i || len(path) <= i {
			break
		}

		if !regex || !isRegexString(location[i]) {
			if location[i] != path[i] {
				return matchCount
			}
		} else if regex && isRegexString(location[i]) {
			r := RegexpStore.GetPattern(location[i][1 : len(location[i])-1])
			if r == nil {
				return 0
			}

			if !r.MatchString(path[i]) {
				return matchCount
			}
		}

		matchCount++
	}

	return matchCount
}

func isRegexString(str string) bool {
	return strings.HasSuffix(str, "}") && strings.HasPrefix(str, "{")
}

func trunSlice(sl []string) []string {
	j := 0
	for i := range sl {
		if sl[i] != "" {
			sl[j] = sl[i]
			j++
		}
	}
	return sl[:j]
}

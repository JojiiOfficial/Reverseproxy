package models

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/JojiiOfficial/gaw"
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
func (location RouteLocation) ModifyProxyRequest(req *http.Request) {
	destination := location.DestinationURL

	targetQuery := destination.RawQuery
	req.URL.Scheme = destination.Scheme
	req.URL.Host = destination.Host

	// Only join file name if location.Location ends with a /
	if strings.HasSuffix(destination.Path, "/") {
		// Remove the the target url prefix from the request url
		if len(destination.Path) > 1 && len(req.URL.Path) > 1 && strings.HasPrefix(req.URL.Path[1:], destination.Path[1:]) {
			req.URL.Path = req.URL.Path[len(destination.Path)-1:]
		}

		// Append path if both aren't the same
		if trimPath(destination.Path) != trimPath(req.URL.Path) {
			req.URL.Path = singleJoiningSlash(destination.Path, req.URL.Path)
		}
	} else {
		req.URL.Path = destination.Path
	}

	// Append HTTP queries of target and request
	// This allows to add custom queries into locations which will be joined
	// with the requested query
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}

	location.finalMods(req)
}

func trimPath(p string) string {
	return strings.Trim(p, "/")
}

func (location *RouteLocation) finalMods(req *http.Request) {
	// explicitly disable User-Agent so it's not set to default value
	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", "")
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func findMatchingLocation(pathItems []string, locations []RouteLocation) *RouteLocation {
	for index := range locations {
		locationItems := gaw.TrimEmptySlice(strings.Split(locations[index].Location, "/"))

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

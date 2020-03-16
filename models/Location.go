package models

import (
	"net/http"
	"net/url"
	"strings"
)

// RouteLocation location for route
type RouteLocation struct {
	// Toml config attributes
	Location    string
	Destination string
	Regex       bool

	// Allow/deny hosts
	Allow []string
	Deny  string

	// Non toml attrs
	DestinationURL *url.URL `toml:"-"`
	Route          *Route   `toml:"-"`
}

// Init inits a location
func (location *RouteLocation) Init() {
	location.DestinationURL, _ = url.Parse(location.Destination)
}

// ModifyProxyRequest modifies a request to a proxy forward request
func (location RouteLocation) ModifyProxyRequest(req *http.Request) {
	target := location.DestinationURL

	targetQuery := target.RawQuery
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host

	// Only join file name if location.Location ends with a /
	if strings.HasSuffix(target.Path, "/") {
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
	} else {
		req.URL.Path = target.Path
	}

	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
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

func singleJoiningSlash(a, b string) string {
	if strings.HasPrefix(b, a) {
		b = b[len(a)-1:]
	}
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
		locationItems := trunSlice(strings.Split(locations[index].Location, "/"))
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

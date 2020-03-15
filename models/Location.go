package models

import (
	"net/http"
	"net/url"
	"strings"
)

// RouteLocation location for route
type RouteLocation struct {
	Location    string
	Destination string

	DestinationURL *url.URL `toml:"-"`
}

// Init inits a location
func (location *RouteLocation) Init() {
	location.DestinationURL, _ = url.Parse(location.Destination)
}

// GetDetinationURL gets url from location
func (location RouteLocation) GetDetinationURL() *url.URL {
	u, _ := url.Parse(location.Destination)
	return u
}

// ModifyRequest modifies a request to a proxy forward request
func (location RouteLocation) ModifyRequest(req *http.Request) {
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

func findMatchingLocation(reqPath string, locations []RouteLocation) *RouteLocation {
	var matching []*RouteLocation

	for i := range locations {
		location := locations[i]

		if !strings.HasPrefix(location.Location, reqPath) {
			continue
		}

		matching = append(matching, &location)
	}

	// If only one found, use it
	if len(matching) == 1 {
		return matching[0]
	}

	if len(matching) > 1 {
		for i := range matching {
			if matching[i].Location == reqPath {
				return matching[i]
			}
		}
	}

	return nil
}

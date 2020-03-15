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
	req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
	if _, ok := req.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
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

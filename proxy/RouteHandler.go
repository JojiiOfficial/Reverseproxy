package proxy

import (
	"errors"
	"net/http"
	"time"

	"github.com/JojiiOfficial/ReverseProxy/models"
	log "github.com/sirupsen/logrus"
)

// Director directs
func (httpServer *HTTPServer) Director(req *http.Request) {}

// RoundTrip trips stuff around
func (httpServer *HTTPServer) RoundTrip(req *http.Request) (*http.Response, error) {
	var start time.Time

	if httpServer.Loglevel == log.DebugLevel {
		start = time.Now()
	}

	// Set host of new request
	var err error
	taskResponse := &http.Response{}
	req.URL.Host = req.Host

	if httpServer.ListenAddress.IsRedirectInterface {
		// Handle Redirection
		taskResponse = httpServer.redirectTask(req, httpServer.ListenAddress)
	} else {
		// Handle proxy route
		location := models.FindMatchingLocation(httpServer.Routes, req)
		if location == nil {
			log.Warnf("No matching route found for %s", req.URL.String())
			return nil, errors.New("Route not found")
		}
		taskResponse, err = httpServer.proxyTask(req, location)
	}

	// Prevent useless operations
	if httpServer.Loglevel == log.DebugLevel {
		log.Debug("Action took ", time.Since(start).String())
	}

	return taskResponse, err
}

// --- Tasks

// Proxy a request
func (httpServer *HTTPServer) proxyTask(req *http.Request, location *models.RouteLocation) (*http.Response, error) {
	// Handle access control
	if !isRequestAllowed(req, location) {
		log.Debugf("IP %s is not allowed", req.RemoteAddr)
		return getForbiddenResponse(req), nil
	}

	// Modifies the request
	location.ModifyProxyRequest(req)

	// Call default roundTrip to forward the request
	return http.DefaultTransport.RoundTrip(req)
}

// Send redirect request
func (httpServer *HTTPServer) redirectTask(req *http.Request, listenAddress *models.ListenAddress) *http.Response {
	body := listenAddress.TaskData.Redirect.GetBody()
	req.URL.Scheme = "https"
	to := req.URL.String()

	if len(to) == 0 {
		log.Fatalln("To redirect target specified for %s")
		return nil
	}

	// Build header
	header := make(http.Header)
	header.Set("Location", to)

	// Build and return response
	return buildResponse(req, http.StatusMovedPermanently, body, "301 Moved permanently", header)
}

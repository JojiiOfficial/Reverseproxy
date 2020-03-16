package proxy

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/JojiiOfficial/ReverseProxy/models"
	"github.com/JojiiOfficial/gaw"
	log "github.com/sirupsen/logrus"
)

// Director directs
func (httpServer *HTTPServer) Director(req *http.Request) {}

// RoundTrip trips stuff around
func (httpServer *HTTPServer) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Set host of new request
	req.URL.Host = req.Host

	// Find matching location
	location := models.FindMatchingLocation(httpServer.Routes, req)
	if location == nil {
		log.Warnf("No matching route found for %s", req.URL.String())
		return nil, errors.New("Route not found")
	}

	// Get interface

	rif := httpServer.Config.GetAddress(httpServer.Server.Addr)
	if rif == nil {
		return nil, errors.New("Interface not found")
	}

	var err error
	taskResponse := &http.Response{}

	// Do specified task
	switch rif.GetTask() {
	case models.ProxyTask:
		{
			// Forward request
			taskResponse, err = httpServer.proxyTask(req, location)
		}
	case models.HTTPRedirectTask:
		{
			// Send redirect
			taskResponse, err = httpServer.redirectTask(req, rif), nil
		}
	}

	log.Debug("Action took ", time.Since(start).String())
	return taskResponse, err
}

// --- Tasks

// Proxy a request
func (httpServer *HTTPServer) proxyTask(req *http.Request, location *models.RouteLocation) (*http.Response, error) {
	// Handle access control
	if location.Deny == "all" {
		ip := strings.Split(req.RemoteAddr, ":")[0]
		if !gaw.IsInStringArray(ip, location.Allow) {
			return getForbiddenResponse(req), nil
		}
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

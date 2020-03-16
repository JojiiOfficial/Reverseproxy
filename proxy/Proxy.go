package proxy

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/JojiiOfficial/ReverseProxy/models"
	log "github.com/sirupsen/logrus"
)

// Director directs
func (httpServer *HTTPServer) Director(req *http.Request) {}

// RoundTrip trips stuff around
func (httpServer *HTTPServer) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	req.URL.Host = req.Host

	location := models.FindMatchingLocation(httpServer.Routes, req)
	if location == nil {
		log.Warnf("No matching route found for %s", req.URL.String())
		req.URL = nil
		return nil, errors.New("Route not found")
	}

	rif := httpServer.GetInterfaceFromRoute(location.Route)
	if rif == nil {
		log.Error("Couldn't find interface!")
		req.URL = nil
		return nil, errors.New("Interface not found")
	}

	switch rif.GetTask() {
	case models.ProxyTask:
		{
			httpServer.proxyTask(req, location)
		}
	case models.HTTPRedirectTask:
		{
			return httpServer.redirectTask(req, rif), nil
		}
	}

	log.Info("Action took ", time.Since(start).String())

	return http.DefaultTransport.RoundTrip(req)
}

// --- Tasks

func (httpServer *HTTPServer) proxyTask(req *http.Request, location *models.RouteLocation) {
	// Modifies the request
	location.ModifyProxyRequest(req)

	if httpServer.Debug {
		log.Info("Forwarding to ", req.URL.String())
	}
}

func (httpServer *HTTPServer) redirectTask(req *http.Request, addressInterface *models.AddressInterface) *http.Response {
	data := addressInterface.TaskData.Redirect
	to := data.Location

	if len(to) == 0 {
		log.Fatalln("To redirect target specified for %s")
		return nil
	}

	header := make(http.Header)
	header.Set("Location", to)

	t := &http.Response{
		Status:        "301 Redirect",
		StatusCode:    data.GetHTTPCode(),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(data.GetBody())),
		ContentLength: int64(len(data.GetBody())),
		Request:       req,
		Header:        header,
	}
	return t
}

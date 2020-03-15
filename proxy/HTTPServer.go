package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/JojiiOfficial/ReverseProxy/models"
	log "github.com/sirupsen/logrus"
)

// HTTPServer http server
type HTTPServer struct {
	SSL    bool
	Routes []*models.Route
	Server *http.Server
}

// Start starts the server
func (httpServer *HTTPServer) Start() {
	httpServer.initRouter()
	go httpServer.run()
}

func (httpServer *HTTPServer) initRouter() {
	httpServer.Server.Handler = &httputil.ReverseProxy{
		Director: httpServer.Director,
	}
}

func (httpServer *HTTPServer) run() {
	log.Fatalln(httpServer.Server.ListenAndServe())
}

// GetScheme returns scheme
func (httpServer HTTPServer) GetScheme() string {
	if httpServer.SSL {
		return "https"
	}
	return "http"
}

// Director director
func (httpServer *HTTPServer) Director(req *http.Request) {
	req.URL.Scheme = httpServer.GetScheme()
	req.URL.Host = req.Host

	location := models.FindMatchingLocation(httpServer.Routes, req)
	if location == nil {
		log.Warnf("No matching route found for %s", req.URL.String())
		req.URL = nil
		return
	}

	// Modifies the request
	location.ModifyRequest(req)
}

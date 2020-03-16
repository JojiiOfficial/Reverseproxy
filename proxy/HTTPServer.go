package proxy

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/JojiiOfficial/ReverseProxy/models"
	log "github.com/sirupsen/logrus"
)

// HTTPServer http server
type HTTPServer struct {
	SSL    bool
	Debug  bool
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
	if httpServer.SSL {
		listener, err := tls.Listen("tcp", httpServer.Server.Addr, httpServer.Server.TLSConfig)
		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(httpServer.Server.Serve(listener))
	} else {
		log.Fatalln(httpServer.Server.ListenAndServe())
	}
}

// GetInterfaceFromRoute returns AddressInterface
func (httpServer HTTPServer) GetInterfaceFromRoute(r *models.Route) *models.AddressInterface {
	for _, rif := range r.Interfaces {
		if rif.Address == httpServer.Server.Addr {
			return &rif
		}
	}
	return nil
}

// Director directs
func (httpServer *HTTPServer) Director(req *http.Request) {
	start := time.Now()
	req.URL.Host = req.Host

	location := models.FindMatchingLocation(httpServer.Routes, req)
	if location == nil {
		log.Warnf("No matching route found for %s", req.URL.String())
		req.URL = nil
		return
	}

	rif := httpServer.GetInterfaceFromRoute(location.Route)
	if rif == nil {
		log.Error("Couldn't find interface!")
		req.URL = nil
		return
	}

	switch rif.GetTask() {
	case models.ProxyTask:
		httpServer.proxyTask(req, location)
	case models.HTTPRedirectTask:
		httpServer.redirectTask(req, location)
	}

	log.Info("Action took ", time.Since(start).String())
}

func (httpServer *HTTPServer) proxyTask(req *http.Request, location *models.RouteLocation) {
	// Modifies the request
	location.ModifyProxyRequest(req)

	if httpServer.Debug {
		log.Info("Forwarding to ", req.URL.String())
	}
}

func (httpServer *HTTPServer) redirectTask(req *http.Request, location *models.RouteLocation) {
	// Modifies the request
	location.ModifyRedirectRequest(req)

	if httpServer.Debug {
		log.Info("Redirecting")
	}
}

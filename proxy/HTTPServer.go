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

// Director director
func (httpServer *HTTPServer) Director(req *http.Request) {
	start := time.Now()
	req.URL.Host = req.Host

	location := models.FindMatchingLocation(httpServer.Routes, req)
	if location == nil {
		log.Warnf("No matching route found for %s", req.URL.String())
		req.URL = nil
		return
	}

	// Modifies the request
	location.ModifyRequest(req)

	if httpServer.Debug {
		log.Info("Forwarding to ", req.URL.String())
	}

	log.Info("Forwarding took ", time.Since(start).String())
}

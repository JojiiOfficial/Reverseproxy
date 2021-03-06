package proxy

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"

	"github.com/JojiiOfficial/ReverseProxy/models"
	log "github.com/sirupsen/logrus"
)

// HTTPServer http server
type HTTPServer struct {
	SSL           bool
	Debug         bool
	ListenAddress *models.ListenAddress
	Routes        []*models.Route
	Server        *http.Server
	Config        *models.Config
	Loglevel      log.Level
}

// Start starts the server
func (httpServer *HTTPServer) Start() {
	httpServer.Loglevel = log.GetLevel()
	httpServer.initRouter()
	go httpServer.run()
}

func (httpServer *HTTPServer) initRouter() {
	httpServer.Server.Handler = &httputil.ReverseProxy{
		Director:       httpServer.Director,
		Transport:      httpServer,
		ModifyResponse: httpServer.ModifyResponse,
	}
}

// Start the server
func (httpServer *HTTPServer) run() {
	if httpServer.SSL {
		log.Debugf("Starting HTTPS server on '%s' with %d certificates and %d routes",
			httpServer.Server.Addr,
			len(httpServer.Server.TLSConfig.Certificates),
			len(httpServer.Routes),
		)

		// Create TLS Listener for ... tls
		listener, err := tls.Listen("tcp", httpServer.Server.Addr, httpServer.Server.TLSConfig)
		if err != nil {
			log.Fatal(err)
		}

		// Start the server
		log.Fatal(httpServer.Server.Serve(listener))
	} else {
		log.Debugf("Starting HTTP server on '%s' with %d routes",
			httpServer.Server.Addr,
			len(httpServer.Routes),
		)

		// Start the http server
		log.Fatalln(httpServer.Server.ListenAndServe())
	}
}

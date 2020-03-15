package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/JojiiOfficial/ReverseProxy/models"
	"github.com/sirupsen/logrus"
)

// HTTPServer http server
type HTTPServer struct {
	SSL    bool
	Server *http.Server
}

// ReverseProxyServer a reverseproxy server
type ReverseProxyServer struct {
	Config *models.Config
	Routes []models.Route
	Server []HTTPServer
}

// NewReverseProxyServere create a new reverseproxy server
func NewReverseProxyServere(config *models.Config, routes []models.Route) *ReverseProxyServer {
	return &ReverseProxyServer{
		Config: config,
		Routes: routes,
	}
}

// InitHTTPServers inits http servers
func (server *ReverseProxyServer) InitHTTPServers() {
	for i, listenAddress := range server.Config.ListenAddresses {
		httpServer := &http.Server{
			Addr:      listenAddress.GetAddress(),
			TLSConfig: &tls.Config{},
			// TODO add more config
		}

		// If address is ssl address, add tls config
		if listenAddress.SSL {
			certKeyPairs := models.GetTLSCerts(server.Routes, &server.Config.ListenAddresses[i])
			if len(certKeyPairs) == 0 {
				logrus.Fatalln("Couldn't find correct certificate pairs!")
			}

			var tlsConfig tls.Config
			for _, pair := range certKeyPairs {
				// Load cert
				cert, err := pair.GetCertificate()
				if err != nil {
					log.Fatalln(err)
				}

				// Add cert to tls.config
				tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
			}

			// Set tls config
			httpServer.TLSConfig = &tlsConfig
		}

		// Append server
		server.Server = append(server.Server, HTTPServer{
			SSL:    listenAddress.SSL,
			Server: httpServer,
		})
	}
}

// Start starts the server
func (server *ReverseProxyServer) Start() {

	// Wait for shutting down
	server.WaitForShutdown()
}

// WaitForShutdown waiting for shutdown
func (server *ReverseProxyServer) WaitForShutdown() {
	for {
		// magic
	}
}

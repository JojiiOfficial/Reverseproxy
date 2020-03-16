package proxy

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JojiiOfficial/ReverseProxy/models"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// ReverseProxyServer a reverseproxy server
type ReverseProxyServer struct {
	Config *models.Config
	Routes []models.Route
	Server []HTTPServer
	Debug  bool
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
	var foundRoutes int

	for i, listenAddress := range server.Config.ListenAddresses {
		serverConf := server.Config.Server

		httpServer := http.Server{
			Addr:           listenAddress.GetAddress(),
			MaxHeaderBytes: int(serverConf.MaxHeaderSize.Bytes()),
			ReadTimeout:    time.Duration(serverConf.ReadTimeout),
			WriteTimeout:   time.Duration(serverConf.WriteTimeout),
		}

		// If address is ssl address, add tls config
		if listenAddress.SSL {
			certKeyPairs := models.GetTLSCerts(server.Routes, &server.Config.ListenAddresses[i])
			if len(certKeyPairs) == 0 {
				logrus.Warnf("Couldn't find any certificate pairs for Address '%s'. This Route/Server might be unavailable", listenAddress.Address)
				continue
			}

			tlsConfig := &tls.Config{}
			for _, pair := range certKeyPairs {
				log.Debug("Found cert: ", pair.Cert, " Key: ", pair.Key)

				// Load cert
				cert, err := pair.GetCertificate()
				if err != nil {
					log.Fatalln(err)
				}

				// Add cert to tls.config
				tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
			}

			tlsConfig.BuildNameToCertificate()

			// Set tls config
			httpServer.TLSConfig = tlsConfig
		}

		// Append server
		routes := models.GetRoutesFromAddress(server.Routes, server.Config.ListenAddresses[i])
		server.Server = append(server.Server, HTTPServer{
			SSL:    listenAddress.SSL,
			Server: &httpServer,
			Routes: routes,
			Debug:  server.Debug,
			Config: server.Config,
		})

		foundRoutes += len(routes)
	}

	// Exit if no route was found
	if foundRoutes == 0 {
		log.Fatal("No route found")
	}
}

// Start starts the server
func (server *ReverseProxyServer) Start() {
	for i := range server.Server {
		server.Server[i].Start()
	}

	// Wait for shutting down
	server.WaitForShutdown()
}

// WaitForShutdown waiting for shutdown
func (server *ReverseProxyServer) WaitForShutdown() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)

	// await os signal
	<-signalChan

	// Create a deadline for the await
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	log.Info("Shutting down server")

	for i := range server.Server {
		server.Server[i].Server.Shutdown(ctx)
	}

	log.Info("Shutting down complete")
	os.Exit(0)
}

package main

import (
	"github.com/JojiiOfficial/ReverseProxy/models"
	"github.com/JojiiOfficial/ReverseProxy/proxy"
	log "github.com/sirupsen/logrus"
)

const (
	//Version version of reverseproxy
	Version = "v0.1"
	// DefaultConfigFile default config file
	DefaultConfigFile = "./config/config.toml"
)

func main() {
	// Init config
	config := models.InitConfig(DefaultConfigFile)
	// Check route count
	if len(config.RouteFiles) == 0 {
		log.Error("No route found!")
		return
	}

	// Loading routes
	routes, err := config.LoadRoutes()
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Infof("Successfully loaded %d routes", len(routes))

	// Create and start the reverseproxy server
	server := proxy.NewReverseProxyServere(config, routes)
	server.InitHTTPServers()
	server.Start()

}

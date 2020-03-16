package main

import (
	"flag"

	"github.com/JojiiOfficial/ReverseProxy/models"
	"github.com/JojiiOfficial/ReverseProxy/proxy"
	log "github.com/sirupsen/logrus"
)

const (
	//Version version of reverseproxy
	Version = "v0.1"
	// DefaultConfigFile default config file
	DefaultConfigFile = "/etc/reverseproxy/config.toml"
)

var (
	configPath *string
)

func initFlags() {
	configPath = flag.String("config", "", "Specify the configfil")
	flag.Parse()
}

func main() {
	initFlags()

	// Determine configfile
	configFile := DefaultConfigFile
	if len(*configPath) > 0 {
		configFile = *configPath
	}

	// Init config
	config := models.InitConfig(configFile)
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

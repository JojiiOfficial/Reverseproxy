package main

import (
	"flag"
	"os"

	"github.com/JojiiOfficial/ReverseProxy/models"
	"github.com/JojiiOfficial/ReverseProxy/proxy"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const (
	//Version version of reverseproxy
	Version = "v0.1"
	// DefaultConfigPath default config path
	DefaultConfigPath = "/etc/reverseproxy/"
	// DefaultConfigFile default config file
	DefaultConfigFile = DefaultConfigPath + "config.toml"
)

var (
	configPath *string
	debug      *bool
)

func initFlags() {
	configPath = flag.String("config", "", "Specify the configfile")
	debug = flag.Bool("debug", false, "Debug")
	flag.Parse()
}

func main() {
	initFlags()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debugmode: on")
	}

	// Determine configfile
	configFile := DefaultConfigFile
	if len(*configPath) > 0 {
		configFile = *configPath
	}
	configFile = getEnvar("PROXY_CONFIG", configFile)

	// Init config
	config := models.InitConfig(configFile, DefaultConfigPath)

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
	server.Debug = *debug
	server.InitHTTPServers()
	server.Start()

}

func getEnvar(key, fallbackValue string) string {
	va, ok := os.LookupEnv(key)
	if ok {
		return va
	}
	return fallbackValue
}

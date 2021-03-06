package main

import (
	"flag"
	"os"
	"time"

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
	configPath            *string
	debug, forceColorLogs *bool
)

func initFlags() {
	configPath = flag.String("config", "", "Specify the configfile")
	debug = flag.Bool("debug", false, "Debug")
	forceColorLogs = flag.Bool("force-colors", false, "Force colors")
	flag.Parse()

	*configPath = getEnvar("PROXY_CONFIG", *configPath)

	if os.Getenv("PROXY_FORCE_COLORS") == "true" {
		*forceColorLogs = true
	}

	if os.Getenv("PROXY_DEBUG") == "true" {
		*debug = true
	}
}

func main() {
	initFlags()

	log.SetOutput(os.Stdout)
	if *forceColorLogs {
		setupForceColor()
		logrus.Info("Force Colors: on")
	}

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debugmode: on")
	}

	// Determine configfile
	configFile := DefaultConfigFile
	if len(*configPath) > 0 {
		configFile = *configPath
	}

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

	if len(routes) == 0 {
		log.Fatal("No route was found")
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

func setupForceColor() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.Stamp,
		FullTimestamp:    true,
		ForceColors:      true,
	})
}

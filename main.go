package main

import (
	"ReverseProxy/models"

	log "github.com/sirupsen/logrus"
)

const (
	//Version version of reverseproxy
	Version = "v0.1"
	// DefaultConfigFile default config file
	DefaultConfigFile = "./config/config.toml"
)

var (
	config *models.Config
)

func main() {
	// Create config if not exists
	created, err := models.CreateDefaultConfig(DefaultConfigFile)
	if err != nil {
		log.Fatalln(err)
	}

	// Exit if config was created
	if created {
		log.Infof("Config %s created successfully", DefaultConfigFile)
		return
	}

	// Read config
	con, err := models.ReadConfig(DefaultConfigFile)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Use config
	config = con

	// Check route count
	if len(config.RouteFiles) == 0 {
		log.Error("No route found!")
		return
	}

	log.Infof("Starting reverseproxy %s with %d routes", Version, len(config.RouteFiles))
}

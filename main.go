package main

import (
	"os"

	"github.com/JojiiOfficial/ReverseProxy/models"

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
	// Init config
	config = initConfig(DefaultConfigFile)

	// Check route count
	if len(config.RouteFiles) == 0 {
		log.Error("No route found!")
		return
	}

	log.Infof("Starting reverseproxy %s with %d routes", Version, len(config.RouteFiles))

	// Loading routes
	err := config.LoadRoutes()
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Infof("Successfully loaded %d routes", len(config.Routes))
}

func initConfig(file string) *models.Config {
	// Create config if not exists
	created, err := models.CreateDefaultConfig(file)
	if err != nil {
		log.Fatalln(err)
	}

	// Exit if config was created
	if created {
		log.Infof("Config %s created successfully", file)
		os.Exit(0)
		return nil
	}

	// Read config
	con, err := models.ReadConfig(file)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return con
}

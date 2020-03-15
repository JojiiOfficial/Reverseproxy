package models

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/JojiiOfficial/ReverseProxy/models/units"
	"github.com/JojiiOfficial/gaw"
	log "github.com/sirupsen/logrus"
)

// Config configuration file
type Config struct {
	ListenAddresses []ListenAddress
	RouteFiles      []string
	Server          ServerConfig `toml:"Server"`
}

// ServerConfig configuration for webserver
type ServerConfig struct {
	MaxHeaderSize units.Datasize
}

// ListenAddress config for ports
type ListenAddress struct {
	Address string
	Port    uint16
	SSL     bool
}

// GetAddress returns address of a listenAddress
func (address ListenAddress) GetAddress() string {
	return fmt.Sprintf("%s:%d", address.Address, address.Port)
}

// ReadConfig read the config file
func ReadConfig(file string) (*Config, error) {
	// Unmarshal config
	var conf Config
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

// CreateDefaultConfig creates the default config file
func CreateDefaultConfig(file string) (bool, error) {
	fStat, err := os.Stat(file)
	if err == nil && fStat.Size() > 0 {
		return false, nil
	}

	// Example route
	exampleRoute := "routes/route1.toml"

	// Create default config struct
	config := Config{
		ListenAddresses: []ListenAddress{
			ListenAddress{
				Address: "127.0.0.1",
				Port:    80,
				SSL:     false,
			},
			ListenAddress{
				Address: "127.0.0.1",
				Port:    443,
				SSL:     true,
			},
		},
		RouteFiles: []string{
			exampleRoute,
		},
	}

	// Create a route
	err = CreateExampleRoute(exampleRoute)
	if err != nil {
		return false, err
	}

	// Create config-path if neccessary
	err = gaw.CreatePath(file, 0740)
	if err != nil {
		return false, err
	}

	// Open file
	f, err := os.Create(file)
	if err != nil {
		return false, err
	}

	// Encode
	return true, toml.NewEncoder(f).Encode(config)
}

// InitConfig the config
func InitConfig(file string) *Config {
	// Create config if not exists
	created, err := CreateDefaultConfig(file)
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
	con, err := ReadConfig(file)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return con
}

// LoadRoutes loads the routes specified in config
func (config *Config) LoadRoutes() ([]Route, error) {
	var routes []Route

	for _, sRoute := range config.RouteFiles {
		// Load route
		route, err := LoadRoute(sRoute)
		if err != nil {
			return []Route{}, err
		}

		// Set route filename
		route.FileName = gaw.FileFromPath(sRoute)

		// Load addresses
		if !route.LoadAddress(config) {
			return []Route{}, ErrAddrNotFound
		}

		// Check route
		if !route.Check(config) {
			return []Route{}, errors.New("Config check failed")
		}

		// Append to existing routes
		routes = append(routes, *route)
	}

	return routes, nil
}

// GetAddress gets address from config
func (config Config) GetAddress(sAddress string) *ListenAddress {
	for i, address := range config.ListenAddresses {
		if address.GetAddress() == sAddress {
			return &config.ListenAddresses[i]
		}
	}

	return &ListenAddress{Port: 0}
}

// IsListeningOn return true if server is listening on port
func (config Config) IsListeningOn(port uint16) bool {
	for _, address := range config.ListenAddresses {
		if address.Port == port {
			return true
		}
	}

	return false
}

// IsSSLPort return true if server listens on port using ssl
func (config Config) IsSSLPort(port uint16) bool {
	for _, address := range config.ListenAddresses {
		if address.Port == port {
			return address.SSL
		}
	}

	return false
}

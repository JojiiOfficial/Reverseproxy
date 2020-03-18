package models

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/JojiiOfficial/ReverseProxy/models/units"
	"github.com/JojiiOfficial/gaw"
	log "github.com/sirupsen/logrus"
)

// Config configuration file
type Config struct {
	Server          ServerConfig `toml:"Server"`
	ListenAddresses []ListenAddress
	RouteFiles      []string
}

// ServerConfig configuration for webserver
type ServerConfig struct {
	MaxHeaderSize units.Datasize
	ReadTimeout   ConfigDuration
	WriteTimeout  ConfigDuration
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
func CreateDefaultConfig(file, defaultPath string) (bool, error) {
	fStat, err := os.Stat(file)
	if err == nil && fStat.Size() > 0 {
		return false, nil
	}

	// Example route
	exampleRoute := filepath.Join(gaw.PathFromFilepath(file), "routes/route1.toml")

	// Create default config struct
	config := Config{
		Server: ServerConfig{
			MaxHeaderSize: units.Kilobyte * 16,
			ReadTimeout:   ConfigDuration(10 * time.Second),
			WriteTimeout:  ConfigDuration(10 * time.Second),
		},
		ListenAddresses: []ListenAddress{
			ListenAddress{
				Address: "127.0.0.1:80",
				SSL:     false,
				Task:    HTTPRedirectTask,
				TaskData: TaskData{
					Redirect: RedirectData{
						HTTPCode: 301,
					},
				},
			},
			ListenAddress{
				Address: "127.0.0.1:443",
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

// ConfigDuration duration for config
type ConfigDuration time.Duration

// UnmarshalText implements encoding.TextUnmarshaler
func (d *ConfigDuration) UnmarshalText(data []byte) error {
	duration, err := time.ParseDuration(string(data))
	if err == nil {
		*d = ConfigDuration(duration)
	}
	return err
}

// MarshalText implements encoding.TextMarshaler
func (d ConfigDuration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

// InitConfig the config
func InitConfig(file, defaultPath string) *Config {
	// Create config if not exists
	created, err := CreateDefaultConfig(file, defaultPath)
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
			log.Warnf("At least one address used in '%s' was not found. This Route/Address might be unavailable", sRoute)
			continue
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

// GetPreferredSSLAddress returns preferred SSL address
func (config Config) GetPreferredSSLAddress() *ListenAddress {
	for i := range config.ListenAddresses {
		if config.ListenAddresses[i].SSL == true {
			return &config.ListenAddresses[i]
		}
	}
	return nil
}

// GetAddress gets address from config
func (config Config) GetAddress(sAddress string) *ListenAddress {
	for i, address := range config.ListenAddresses {
		if address.GetAddress() == sAddress {
			return &config.ListenAddresses[i]
		}
	}

	return &ListenAddress{Address: ""}
}

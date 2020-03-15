package models

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/JojiiOfficial/gaw"
)

//Config configuration file
type Config struct {
	RouteFiles []string
}

//ReadConfig read the config file
func ReadConfig(file string) (*Config, error) {
	// Unmarshal config
	var conf Config
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

//CreateDefaultConfig creates the default config file
func CreateDefaultConfig(file string) (bool, error) {
	fStat, err := os.Stat(file)
	if err == nil && fStat.Size() > 0 {
		return false, nil
	}

	// Example route
	exampleRoute := "routes/route1.toml"

	// Create default config struct
	config := Config{
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

package models

import (
	"os"

	"github.com/BurntSushi/toml"
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

	// Create default config struct
	config := Config{
		RouteFiles: []string{
			"routes/route1.toml",
		},
	}

	// Open file
	f, err := os.Create(file)
	if err != nil {
		return false, err
	}

	// Encode
	return true, toml.NewEncoder(f).Encode(config)
}

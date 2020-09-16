package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
	"trace/pkg/database"
)

type Config struct {
	DatabaseConfig database.Config

	// The time it will take for a student to be automatically be signed out
	// TODO: Make this per location?
	Timeout time.Duration `json:"timeout"`

	// The Username and Password to login to the website
	Username string `json:"username"`
	Password string `json:"password"`
}

// GlobalConfig is the config used throughout the API package
var GlobalConfig *Config

var DefaultConfig = Config{
	Timeout: 3 * time.Hour,
}

// LoadConfig loads a config from filename. If the file does not exist,
// it will create an empty config at the location and exit the program.
// If the config loaded is invalid, the error will be returned
func LoadConfig(filename string) (config Config, err error) {
	file, err := os.Open(filename)
	if err != nil {
		// Check if the error was from the file not existing
		if os.IsNotExist(err) {
			log.Infof("Config file %s does not exist. Creating an empty one", filename)
			if err := CreateDefaultConfig(filename); err != nil {
				log.WithField("filename", filename).Fatalf("Could not create the default config")
			}
			os.Exit(0)
		}
		return Config{}, err
	}

	configBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, nil
	}

	if err := json.Unmarshal(configBytes, &config); err != nil {
		return Config{}, err
	}

	log.Infof("Loaded config file %s", filename)

	// implicitly return config and err
	return
}

// CreateDefaultConfig creates a default config at the specified location
func CreateDefaultConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	err = enc.Encode(DefaultConfig)
	if err != nil {
		panic(err)
	}

	return nil
}
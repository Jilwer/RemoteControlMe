package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func LoadConfig(path string) (*Config, error) {

	// check if the file exists, if it doesn't create a new one and exit with a message
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Config file not found, creating a new one")
		file, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create config: %w", err)
		}
		defer file.Close()

		// create a new config with default values
		cfg := Config{
			Port:          "8080",
			StaticMessage: "Hello, World!",
			ChatEnabled:   true,
		}

		// encode the config to the file
		if err = toml.NewEncoder(file).Encode(cfg); err != nil {
			return nil, fmt.Errorf("failed to encode config: %w", err)
		}

		log.Fatal("Config file created, please fill in the values and restart the program")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %w", err)
	}

	var cfg Config
	if err = toml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Config struct {
	Port          string `toml:"port"`
	StaticMessage string `toml:"static_message"`
	ChatEnabled   bool   `toml:"chat_enabled"`
}

package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pelletier/go-toml/v2"
)

func LoadConfig(path string) (*StateConfig, error) {

	// check if the file exists, if it doesn't create a new one and exit with a message
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Config file not found, creating a new one")
		file, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create config: %w", err)
		}
		defer file.Close()

		cfg := StateConfig{
			StaticMessage: &StaticMessage{
				Send:    true,
				Message: "Hello, World!",
			},
			Chat: &Chat{
				Enabled: true,
			},
			Server: &Server{
				Port: "8080",
			},
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

	var cfg StateConfig
	if err = toml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func MustLoadConfig(path string) *StateConfig {
	cfg, err := LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

type StaticMessage struct {
	Send    bool         `toml:"send"`
	Timer   *time.Ticker `toml:"-"`
	Message string       `toml:"message"`
}

type Server struct {
	Port string `toml:"port"`
}

type Chat struct {
	LastMessageTime time.Time     `toml:"-"`
	RateLimit       time.Duration `toml:"-"`
	Mutex           *sync.Mutex   `toml:"-"`
	Enabled         bool          `toml:"enabled"`
}

type StateConfig struct {
	StaticMessage *StaticMessage `toml:"static_message"`
	Chat          *Chat          `toml:"chat"`
	Server        *Server        `toml:"server"`
}

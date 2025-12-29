package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string     `json:"env"`
	HTTPServer HTTPServer `json:"http_server"`
	Client     Client     `json:"client"`
}

type HTTPServer struct {
	Address        string `json:"address"`
	Timeout        time.Duration
	IdleTimeout    time.Duration
	TimeoutStr     string `json:"timeout"`      // temporary field to parse seconds and convert them later to time.Duration
	IdleTimeoutStr string `json:"idle_timeout"` // ~//~
}

type Client struct {
	Host       string `json:"host"`
	Timeout    time.Duration
	TimeoutStr string `json:"cleint_timeout"`
}

func MustLoad() *Config {
	// configPath := os.Getenv("CONFIG_PATH") // for production
	configPath := "./config/local/local.json" // simplification for review purposes
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist: %s", configPath)
	}

	var cfg Config

	data, _ := os.ReadFile(configPath)
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("can't unmarshal config file: %v", err)
	}

	var err error
	cfg.HTTPServer.Timeout, err = time.ParseDuration(cfg.HTTPServer.TimeoutStr)
	if err != nil {
		log.Fatalf("can't parse timeout: %v", err)
	}

	cfg.HTTPServer.IdleTimeout, err = time.ParseDuration(cfg.HTTPServer.IdleTimeoutStr)
	if err != nil {
		log.Fatalf("can't parse idle_timeout: %v", err)
	}

	cfg.Client.Timeout, err = time.ParseDuration(cfg.Client.TimeoutStr)
	if err != nil {
		log.Fatalf("can't parse client timeout: %v", err)
	}

	return &cfg
}

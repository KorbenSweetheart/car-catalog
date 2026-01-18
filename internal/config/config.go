package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"viewer/internal/lib/e"
)

type Config struct {
	Env        string     `json:"env"`
	HTTPServer HTTPServer `json:"http_server"`
	Client     Client     `json:"client"`
	Repo       Repository `json:"repository"`
}

type HTTPServer struct {
	Address        string `json:"address"`
	StaticPath     string `json:"static_path"`
	TemplatesPath  string `json:"templates_path"`
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

type Repository struct {
	MediaHost          string `json:"media_host"`
	RefreshInterval    time.Duration
	RefreshIntervalStr string `json:"refresh_interval"`
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

	if err := loadStatic(cfg.HTTPServer.StaticPath); err != nil {
		log.Fatalf("can't load static directory: %v", err)
	}

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

	cfg.Repo.RefreshInterval, err = time.ParseDuration(cfg.Repo.RefreshIntervalStr)
	if err != nil {
		log.Fatalf("can't parse client timeout: %v", err)
	}

	return &cfg
}

func loadStatic(path string) error {

	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return e.Wrap("static directory does not exist", err)
	}

	if !info.IsDir() {
		return e.Wrap("path is not a directory", err)
	}

	if err != nil {
		return e.Wrap("error checking static directory", err)
	}

	return nil
}

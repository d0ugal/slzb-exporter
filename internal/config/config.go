package config

import (
	"os"
	"strconv"
	"time"

	"github.com/d0ugal/slzb-exporter/internal/logging"
)

// Config represents the application configuration
type Config struct {
	Server  ServerConfig          `yaml:"server"`
	Logging logging.LoggingConfig `yaml:"logging"`
	SLZB    SLZBConfig            `yaml:"slzb"`
}

// ServerConfig represents the HTTP server configuration
type ServerConfig struct {
	Host string `yaml:"host" env:"SLZB_EXPORTER_SERVER_HOST"`
	Port int    `yaml:"port" env:"SLZB_EXPORTER_SERVER_PORT"`
}

// SLZBConfig represents the SLZB device configuration
type SLZBConfig struct {
	APIURL   string        `yaml:"api_url" env:"SLZB_EXPORTER_SLZB_API_URL"`
	Interval time.Duration `yaml:"interval" env:"SLZB_EXPORTER_METRICS_DEFAULT_INTERVAL"`
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 9110,
		},
		Logging: logging.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		SLZB: SLZBConfig{
			APIURL:   "http://slzb-device.local",
			Interval: 10 * time.Second, // 10 seconds
		},
	}

	// Always load from environment variables when available
	loadFromEnv(cfg)

	return cfg, nil
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(cfg *Config) {
	// Server config
	if host := os.Getenv("SLZB_EXPORTER_SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if portStr := os.Getenv("SLZB_EXPORTER_SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Server.Port = port
		}
	}

	// Logging config
	if level := os.Getenv("SLZB_EXPORTER_LOG_LEVEL"); level != "" {
		cfg.Logging.Level = level
	}
	if format := os.Getenv("SLZB_EXPORTER_LOG_FORMAT"); format != "" {
		cfg.Logging.Format = format
	}

	// SLZB config
	if apiURL := os.Getenv("SLZB_EXPORTER_SLZB_API_URL"); apiURL != "" {
		cfg.SLZB.APIURL = apiURL
	}
	if intervalStr := os.Getenv("SLZB_EXPORTER_METRICS_DEFAULT_INTERVAL"); intervalStr != "" {
		if interval, err := time.ParseDuration(intervalStr); err == nil {
			cfg.SLZB.Interval = interval
		}
	}
}

package config

import (
	"os"
	"strconv"
	"time"

	"github.com/d0ugal/promexporter/config"
)

// Config represents the application configuration
type Config struct {
	config.BaseConfig

	SLZB SLZBConfig `yaml:"slzb"`
}

// SLZBConfig represents the SLZB device configuration
type SLZBConfig struct {
	APIURL   string        `yaml:"api_url"`
	Interval time.Duration `yaml:"interval"`
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// Load base configuration from environment
	baseConfig := &config.BaseConfig{}

	// Server configuration
	if host := os.Getenv("SLZB_EXPORTER_SERVER_HOST"); host != "" {
		baseConfig.Server.Host = host
	} else {
		baseConfig.Server.Host = "0.0.0.0"
	}

	if portStr := os.Getenv("SLZB_EXPORTER_SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			baseConfig.Server.Port = port
		} else {
			baseConfig.Server.Port = 9110
		}
	} else {
		baseConfig.Server.Port = 9110
	}

	// Logging configuration
	if level := os.Getenv("SLZB_EXPORTER_LOG_LEVEL"); level != "" {
		baseConfig.Logging.Level = level
	} else {
		baseConfig.Logging.Level = "info"
	}

	if format := os.Getenv("SLZB_EXPORTER_LOG_FORMAT"); format != "" {
		baseConfig.Logging.Format = format
	} else {
		baseConfig.Logging.Format = "json"
	}

	// Metrics configuration
	baseConfig.Metrics.Collection.DefaultInterval = config.Duration{}
	baseConfig.Metrics.Collection.DefaultIntervalSet = false

	cfg.BaseConfig = *baseConfig

	// SLZB configuration
	if apiURL := os.Getenv("SLZB_EXPORTER_SLZB_API_URL"); apiURL != "" {
		cfg.SLZB.APIURL = apiURL
	} else {
		cfg.SLZB.APIURL = "http://slzb-device.local"
	}

	if intervalStr := os.Getenv("SLZB_EXPORTER_METRICS_DEFAULT_INTERVAL"); intervalStr != "" {
		if interval, err := time.ParseDuration(intervalStr); err == nil {
			cfg.SLZB.Interval = interval
		} else {
			cfg.SLZB.Interval = 10 * time.Second
		}
	} else {
		cfg.SLZB.Interval = 10 * time.Second
	}

	return cfg, nil
}

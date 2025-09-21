package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Check default values
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default host '0.0.0.0', got '%s'", cfg.Server.Host)
	}

	if cfg.Server.Port != 9110 {
		t.Errorf("Expected default port 9110, got %d", cfg.Server.Port)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", cfg.Logging.Level)
	}

	if cfg.Logging.Format != "json" {
		t.Errorf("Expected default log format 'json', got '%s'", cfg.Logging.Format)
	}

	if cfg.SLZB.APIURL != "http://slzb-device.local" {
		t.Errorf("Expected default API URL 'http://slzb-device.local', got '%s'", cfg.SLZB.APIURL)
	}

	if cfg.SLZB.Interval != 10*time.Second {
		t.Errorf("Expected default interval 10s, got %v", cfg.SLZB.Interval)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Set environment variables
	_ = os.Setenv("SLZB_EXPORTER_SERVER_HOST", "0.0.0.0")
	_ = os.Setenv("SLZB_EXPORTER_SERVER_PORT", "9090")
	_ = os.Setenv("SLZB_EXPORTER_LOG_LEVEL", "warn")
	_ = os.Setenv("SLZB_EXPORTER_LOG_FORMAT", "text")
	_ = os.Setenv("SLZB_EXPORTER_SLZB_API_URL", "http://test-device.local")
	_ = os.Setenv("SLZB_EXPORTER_METRICS_DEFAULT_INTERVAL", "45s")

	defer func() {
		_ = os.Unsetenv("SLZB_EXPORTER_SERVER_HOST")
		_ = os.Unsetenv("SLZB_EXPORTER_SERVER_PORT")
		_ = os.Unsetenv("SLZB_EXPORTER_LOG_LEVEL")
		_ = os.Unsetenv("SLZB_EXPORTER_LOG_FORMAT")
		_ = os.Unsetenv("SLZB_EXPORTER_SLZB_API_URL")
		_ = os.Unsetenv("SLZB_EXPORTER_METRICS_DEFAULT_INTERVAL")
	}()

	// Load config
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	// Check environment variable values
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected host '0.0.0.0', got '%s'", cfg.Server.Host)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected port 9090, got %d", cfg.Server.Port)
	}

	if cfg.Logging.Level != "warn" {
		t.Errorf("Expected log level 'warn', got '%s'", cfg.Logging.Level)
	}

	if cfg.Logging.Format != "text" {
		t.Errorf("Expected log format 'text', got '%s'", cfg.Logging.Format)
	}

	if cfg.SLZB.APIURL != "http://test-device.local" {
		t.Errorf("Expected API URL 'http://test-device.local', got '%s'", cfg.SLZB.APIURL)
	}

	if cfg.SLZB.Interval != 45*time.Second {
		t.Errorf("Expected interval 45s, got %v", cfg.SLZB.Interval)
	}
}

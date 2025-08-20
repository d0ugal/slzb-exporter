package server

import (
	"testing"

	"github.com/d0ugal/slzb-exporter/internal/config"
	"github.com/d0ugal/slzb-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func TestServer_UsesMetricsRegistry(t *testing.T) {
	cfg := &config.Config{}
	metricsRegistry := metrics.NewRegistry()
	server := New(cfg, metricsRegistry)

	// Verify that the server has the metrics registry
	if server.metrics == nil {
		t.Fatal("Expected server to have metrics registry, but got nil")
	}

	// Verify that the metrics registry has metrics info
	metricsInfo := metricsRegistry.GetMetricsInfo()
	if len(metricsInfo) == 0 {
		t.Fatal("Expected metrics registry to contain metrics info, but got empty slice")
	}

	// Check that each metric has required fields
	for i, metric := range metricsInfo {
		if metric.Name == "" {
			t.Errorf("Metric %d has empty name", i)
		}

		if metric.Help == "" {
			t.Errorf("Metric %d has empty help text", i)
		}

		if len(metric.Labels) == 0 && metric.Name != "slzb_exporter_version_info" {
			// Most metrics should have labels, except version info
			t.Errorf("Metric %d (%s) has no labels", i, metric.Name)
		}
	}

	// Check for specific metrics we expect
	foundConnected := false
	foundTemp := false

	for _, metric := range metricsInfo {
		if metric.Name == "slzb_device_connected" {
			foundConnected = true
		}

		if metric.Name == "slzb_device_temperature_celsius" {
			foundTemp = true
		}
	}

	if !foundConnected {
		t.Error("Expected to find slzb_device_connected metric")
	}

	if !foundTemp {
		t.Error("Expected to find slzb_device_temperature_celsius metric")
	}
}

func TestServer_HandleRoot(t *testing.T) {
	// Use a test registry to avoid registration conflicts
	testRegistry := prometheus.NewRegistry()
	originalRegistry := prometheus.DefaultRegisterer
	prometheus.DefaultRegisterer = testRegistry

	defer func() {
		prometheus.DefaultRegisterer = originalRegistry
	}()

	cfg := &config.Config{
		SLZB: config.SLZBConfig{
			APIURL:  "http://localhost:8080",
			Interval: 30,
		},
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
	}
	metricsRegistry := metrics.NewRegistry()
	server := New(cfg, metricsRegistry)

	// Test that the server can be created without errors
	if server == nil {
		t.Fatal("Expected server to be created, but got nil")
	}

	// Test that the server has the expected configuration
	if server.config == nil {
		t.Fatal("Expected server to have config, but got nil")
	}

	if server.config.SLZB.APIURL != "http://localhost:8080" {
		t.Errorf("Expected API URL to be http://localhost:8080, got %s", server.config.SLZB.APIURL)
	}
}

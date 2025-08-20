package server

import (
	"testing"

	"github.com/d0ugal/slzb-exporter/internal/config"
	"github.com/d0ugal/slzb-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func TestServer_GetMetricsInfo(t *testing.T) {
	cfg := &config.Config{}
	metricsRegistry := metrics.NewRegistry()
	server := New(cfg, metricsRegistry)

	metricsInfo := server.getMetricsInfo()

	// Check that we have metrics
	if len(metricsInfo) == 0 {
		t.Fatal("Expected metrics info to contain metrics, but got empty slice")
	}

	// Check that each metric has required fields
	for i, metric := range metricsInfo {
		if metric.Name == "" {
			t.Errorf("Metric %d has empty name", i)
		}

		if metric.Help == "" {
			t.Errorf("Metric %d has empty help text", i)
		}

		if metric.ExampleValue == "" {
			t.Errorf("Metric %d has empty example value", i)
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

func TestServer_GetMetricLabels(t *testing.T) {
	cfg := &config.Config{}
	// Create a test registry to avoid registration conflicts
	testRegistry := prometheus.NewRegistry()
	originalRegistry := prometheus.DefaultRegisterer
	prometheus.DefaultRegisterer = testRegistry

	defer func() {
		prometheus.DefaultRegisterer = originalRegistry
	}()

	metricsRegistry := metrics.NewRegistry()
	server := New(cfg, metricsRegistry)

	tests := []struct {
		metricName string
		expected   map[string]string
	}{
		{
			metricName: "SLZBConnected",
			expected:   map[string]string{"device": ""}, // Device ID will be empty in test
		},
		{
			metricName: "SLZBEthernetConnected",
			expected: map[string]string{
				"device":      "",
				"ip_address":  "unknown",
				"mac_address": "unknown",
				"gateway":     "unknown",
				"subnet_mask": "unknown",
				"dns_server":  "unknown",
				"speed_mbps":  "unknown",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.metricName, func(t *testing.T) {
			labels := server.getMetricLabels(tt.metricName)
			if len(labels) != len(tt.expected) {
				t.Errorf("Expected %d labels, got %d", len(tt.expected), len(labels))
			}

			for k, v := range tt.expected {
				if labels[k] != v {
					t.Errorf("Expected label %s=%s, got %s", k, v, labels[k])
				}
			}
		})
	}
}

func TestServer_GetMetricValue(t *testing.T) {
	cfg := &config.Config{}
	// Create a test registry to avoid registration conflicts
	testRegistry := prometheus.NewRegistry()
	originalRegistry := prometheus.DefaultRegisterer
	prometheus.DefaultRegisterer = testRegistry

	defer func() {
		prometheus.DefaultRegisterer = originalRegistry
	}()

	metricsRegistry := metrics.NewRegistry()
	server := New(cfg, metricsRegistry)

	tests := []struct {
		metricName string
		expected   string
	}{
		{"SLZBConnected", "0"},         // Will be 0 since no real data in test
		{"SLZBDeviceTemp", "0.0"},      // Will be 0.0 since no real data in test
		{"SLZBUptime", "0"},            // Will be 0 since no real data in test
		{"SLZBHeapFree", "0"},          // Will be 0 since no real data in test
		{"SLZBHTTPRequestsTotal", "0"}, // Will be 0 since no real data in test
	}

	for _, tt := range tests {
		t.Run(tt.metricName, func(t *testing.T) {
			value := server.getMetricValue(tt.metricName)
			if value != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, value)
			}
		})
	}
}

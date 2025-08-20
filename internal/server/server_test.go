package server

import (
	"testing"

	"github.com/d0ugal/slzb-exporter/internal/config"
	"github.com/d0ugal/slzb-exporter/internal/metrics"
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

func TestServer_GetExampleLabels(t *testing.T) {
	cfg := &config.Config{}
	metricsRegistry := &metrics.Registry{}
	server := New(cfg, metricsRegistry)

	tests := []struct {
		metricName string
		expected   map[string]string
	}{
		{
			metricName: "SLZBConnected",
			expected:   map[string]string{"device": "slzb-01"},
		},
		{
			metricName: "SLZBEthernetConnected",
			expected: map[string]string{
				"device":      "slzb-01",
				"ip_address":  "192.168.1.100",
				"mac_address": "00:11:22:33:44:55",
				"gateway":     "192.168.1.1",
				"subnet_mask": "255.255.255.0",
				"dns_server":  "8.8.8.8",
				"speed_mbps":  "1000",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.metricName, func(t *testing.T) {
			labels := server.getExampleLabels(tt.metricName)
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

func TestServer_GetExampleValue(t *testing.T) {
	cfg := &config.Config{}
	metricsRegistry := &metrics.Registry{}
	server := New(cfg, metricsRegistry)

	tests := []struct {
		metricName string
		expected   string
	}{
		{"SLZBConnected", "1"},
		{"SLZBDeviceTemp", "45.2"},
		{"SLZBUptime", "86400"},
		{"SLZBHeapFree", "512"},
		{"SLZBHTTPRequestsTotal", "42"},
	}

	for _, tt := range tests {
		t.Run(tt.metricName, func(t *testing.T) {
			value := server.getExampleValue(tt.metricName)
			if value != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, value)
			}
		})
	}
}

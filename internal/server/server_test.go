package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestServer_HandleMetricsInfo(t *testing.T) {
	cfg := &config.Config{}
	// Use a mock metrics registry to avoid duplicate registration issues
	metricsRegistry := &metrics.Registry{}
	server := New(cfg, metricsRegistry)

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/metrics-info", nil)
	w := httptest.NewRecorder()

	// Call the handler
	server.handleMetricsInfo(w, req)

	// Check response status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Parse JSON response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check response structure
	if _, ok := response["metrics"]; !ok {
		t.Error("Expected 'metrics' field in response")
	}

	if _, ok := response["total_count"]; !ok {
		t.Error("Expected 'total_count' field in response")
	}

	if _, ok := response["generated_at"]; !ok {
		t.Error("Expected 'generated_at' field in response")
	}

	// Check that total_count is a number and matches metrics length
	totalCount, ok := response["total_count"].(float64)
	if !ok {
		t.Error("Expected total_count to be a number")
	}

	metricsArray, ok := response["metrics"].([]interface{})
	if !ok {
		t.Error("Expected metrics to be an array")
	}

	if int(totalCount) != len(metricsArray) {
		t.Errorf("Expected total_count %d to match metrics array length %d", int(totalCount), len(metricsArray))
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

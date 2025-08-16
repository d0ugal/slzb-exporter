package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()
	if registry == nil {
		t.Fatal("NewRegistry() returned nil")
	}

	// Test that all metrics are properly initialized
	if registry.VersionInfo == nil {
		t.Error("VersionInfo metric not initialized")
	}

	if registry.SLZBConnected == nil {
		t.Error("SLZBConnected metric not initialized")
	}

	if registry.SLZBDeviceTemp == nil {
		t.Error("SLZBDeviceTemp metric not initialized")
	}

	if registry.SLZBUptime == nil {
		t.Error("SLZBUptime metric not initialized")
	}

	if registry.SLZBSocketUptime == nil {
		t.Error("SLZBSocketUptime metric not initialized")
	}

	if registry.SLZBSocketConnected == nil {
		t.Error("SLZBSocketConnected metric not initialized")
	}

	if registry.SLZBDeviceMode == nil {
		t.Error("SLZBDeviceMode metric not initialized")
	}

	if registry.SLZBHeapFree == nil {
		t.Error("SLZBHeapFree metric not initialized")
	}

	if registry.SLZBHeapSize == nil {
		t.Error("SLZBHeapSize metric not initialized")
	}

	if registry.SLZBHeapRatio == nil {
		t.Error("SLZBHeapRatio metric not initialized")
	}

	if registry.SLZBHTTPRequestsTotal == nil {
		t.Error("SLZBHTTPRequestsTotal metric not initialized")
	}

	if registry.SLZBHTTPErrorsTotal == nil {
		t.Error("SLZBHTTPErrorsTotal metric not initialized")
	}

	if registry.SLZBDeviceReachable == nil {
		t.Error("SLZBDeviceReachable metric not initialized")
	}

	if registry.SLZBLastCollectionTime == nil {
		t.Error("SLZBLastCollectionTime metric not initialized")
	}

	if registry.SLZBCollectionErrors == nil {
		t.Error("SLZBCollectionErrors metric not initialized")
	}
}

func TestMetricsRegistration(t *testing.T) {
	// Create a new prometheus registry for testing
	testRegistry := prometheus.NewRegistry()

	// Temporarily replace the default registry
	originalRegistry := prometheus.DefaultRegisterer
	prometheus.DefaultRegisterer = testRegistry

	defer func() {
		prometheus.DefaultRegisterer = originalRegistry
	}()

	// Create metrics registry
	metrics := NewRegistry()

	// Test that metrics can be registered (they should already be registered by promauto)
	if err := testRegistry.Register(metrics.VersionInfo); err != nil {
		// This is expected to fail because promauto already registered it
		t.Logf("Expected registration failure (already registered): %v", err)
	}
}

func TestMetricsValues(t *testing.T) {
	// Create a new registry for testing to avoid conflicts
	testRegistry := prometheus.NewRegistry()
	originalRegistry := prometheus.DefaultRegisterer
	prometheus.DefaultRegisterer = testRegistry

	defer func() {
		prometheus.DefaultRegisterer = originalRegistry
	}()

	registry := NewRegistry()
	deviceName := "test-device"

	// Test setting and getting metric values
	registry.SLZBConnected.WithLabelValues(deviceName).Set(1)

	if testutil.ToFloat64(registry.SLZBConnected.WithLabelValues(deviceName)) != 1 {
		t.Error("SLZBConnected metric value not set correctly")
	}

	registry.SLZBDeviceTemp.WithLabelValues(deviceName).Set(25.5)

	if testutil.ToFloat64(registry.SLZBDeviceTemp.WithLabelValues(deviceName)) != 25.5 {
		t.Error("SLZBDeviceTemp metric value not set correctly")
	}

	registry.SLZBHeapRatio.WithLabelValues(deviceName).Set(75.5)

	if testutil.ToFloat64(registry.SLZBHeapRatio.WithLabelValues(deviceName)) != 75.5 {
		t.Error("SLZBHeapRatio metric value not set correctly")
	}

	registry.SLZBHTTPRequestsTotal.WithLabelValues(deviceName, "0", "200").Inc()

	if testutil.ToFloat64(registry.SLZBHTTPRequestsTotal.WithLabelValues(deviceName, "0", "200")) != 1 {
		t.Error("SLZBHTTPRequestsTotal metric value not incremented correctly")
	}

	registry.SLZBHTTPErrorsTotal.WithLabelValues(deviceName, "0", "timeout").Inc()

	if testutil.ToFloat64(registry.SLZBHTTPErrorsTotal.WithLabelValues(deviceName, "0", "timeout")) != 1 {
		t.Error("SLZBHTTPErrorsTotal metric value not incremented correctly")
	}
}

package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/d0ugal/slzb-exporter/internal/config"
	"github.com/d0ugal/slzb-exporter/internal/metrics"
	"github.com/d0ugal/slzb-exporter/internal/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type Server struct {
	config  *config.Config
	metrics *metrics.Registry
	server  *http.Server
}

func New(cfg *config.Config, metricsRegistry *metrics.Registry) *Server {
	server := &Server{
		config:  cfg,
		metrics: metricsRegistry,
	}

	return server
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	versionInfo := version.Get()
	metricsInfo := s.metrics.GetMetricsInfo()

	// Convert metrics to template data
	metrics := make([]MetricData, 0, len(metricsInfo))
	for _, metric := range metricsInfo {
		metrics = append(metrics, MetricData{
			Name:         metric.Name,
			Help:         metric.Help,
			Labels:       metric.Labels,
			ExampleValue: metric.ExampleValue,
		})
	}

	data := TemplateData{
		Version:    versionInfo.Version,
		Commit:     versionInfo.Commit,
		BuildDate:  versionInfo.BuildDate,
		MQTTStatus: "connected", // Hardcoded for now - would need SLZB client reference to get actual status
		Metrics:    metrics,
		Config: ConfigData{
			Broker:     s.config.SLZB.APIURL,
			ClientID:   s.config.SLZB.APIURL,
			TopicCount: int(s.config.SLZB.Interval.Seconds()),
			QoS:        30, // Default timeout
		},
	}

	w.Header().Set("Content-Type", "text/html")

	if err := mainTemplate.Execute(w, data); err != nil {
		slog.Error("Failed to execute template", "error", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	versionInfo := version.Get()
	response := fmt.Sprintf(`{
		"status": "healthy",
		"timestamp": %d,
		"service": "slzb-exporter",
		"version": "%s",
		"commit": "%s",
		"build_date": "%s"
	}`, time.Now().Unix(), versionInfo.Version, versionInfo.Commit, versionInfo.BuildDate)

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write([]byte(response)); err != nil {
		slog.Error("Failed to write health response", "error", err)
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", s.handleHealth)

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	slog.Info("Starting HTTP server", "address", addr)

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	slog.Info("Shutting down HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *Server) getMetricsInfo() []MetricInfo {
	metricsInfo := make([]MetricInfo, 0, 25)

	// Define all metrics manually since reflection approach is complex with Prometheus metrics
	metrics := []struct {
		name  string
		field string
	}{
		{"slzb_device_connected", "SLZBConnected"},
		{"slzb_device_temperature_celsius", "SLZBDeviceTemp"},
		{"slzb_device_uptime_seconds", "SLZBUptime"},
		{"slzb_socket_uptime_seconds", "SLZBSocketUptime"},
		{"slzb_socket_connected", "SLZBSocketConnected"},
		{"slzb_device_operational_mode", "SLZBDeviceMode"},
		{"slzb_device_heap_free_kb", "SLZBHeapFree"},
		{"slzb_device_heap_size_kb", "SLZBHeapSize"},
		{"slzb_device_heap_ratio", "SLZBHeapRatio"},
		{"slzb_device_ethernet_connected", "SLZBEthernetConnected"},
		{"slzb_device_wifi_connected", "SLZBWifiConnected"},
		{"slzb_http_requests_total", "SLZBHTTPRequestsTotal"},
		{"slzb_http_errors_total", "SLZBHTTPErrorsTotal"},
		{"slzb_device_reachable", "SLZBDeviceReachable"},
		{"slzb_last_collection_timestamp", "SLZBLastCollectionTime"},
		{"slzb_collection_errors_total", "SLZBCollectionErrors"},
		{"slzb_firmware_current_version", "SLZBFirmwareCurrentVersion"},
		{"slzb_firmware_update_available", "SLZBFirmwareUpdateAvailable"},
		{"slzb_firmware_last_check_timestamp", "SLZBFirmwareLastCheckTime"},
		{"slzb_config_backup_status", "SLZBConfigBackupStatus"},
		{"slzb_config_last_backup_timestamp", "SLZBConfigLastBackupTime"},
		{"slzb_config_file_count", "SLZBConfigFileCount"},
		{"slzb_api_response_time_seconds", "SLZBAPIResponseTimeSeconds"},
		{"slzb_api_timeout_errors_total", "SLZBAPITimeoutErrorsTotal"},
		{"slzb_collection_duration_seconds", "SLZBCollectionDurationSeconds"},
	}

	for _, metric := range metrics {
		metricsInfo = append(metricsInfo, MetricInfo{
			Name:         metric.name,
			Help:         s.getMetricHelp(metric.name),
			ExampleValue: s.getMetricValue(metric.field),
			Labels:       s.getMetricLabels(metric.field),
		})
	}

	return metricsInfo
}

// getMetricValue gets the current value of a metric from the registry
func (s *Server) getMetricValue(metricName string) string {
	deviceID := s.getDeviceID()

	switch metricName {
	case "SLZBConnected":
		value := testutil.ToFloat64(s.metrics.SLZBConnected.WithLabelValues(deviceID))
		return fmt.Sprintf("%.0f", value)
	case "SLZBDeviceTemp":
		value := testutil.ToFloat64(s.metrics.SLZBDeviceTemp.WithLabelValues(deviceID))
		return fmt.Sprintf("%.1f", value)
	case "SLZBUptime":
		value := testutil.ToFloat64(s.metrics.SLZBUptime.WithLabelValues(deviceID))
		return fmt.Sprintf("%.0f", value)
	case "SLZBSocketUptime":
		value := testutil.ToFloat64(s.metrics.SLZBSocketUptime.WithLabelValues(deviceID))
		return fmt.Sprintf("%.0f", value)
	case "SLZBSocketConnected":
		// Try to get a connected socket value
		value := testutil.ToFloat64(s.metrics.SLZBSocketConnected.WithLabelValues(deviceID, "1"))
		if value == 0 {
			value = testutil.ToFloat64(s.metrics.SLZBSocketConnected.WithLabelValues(deviceID, "0"))
		}

		return fmt.Sprintf("%.0f", value)
	case "SLZBDeviceMode":
		// Try to get a mode value (coordinator is common)
		value := testutil.ToFloat64(s.metrics.SLZBDeviceMode.WithLabelValues(deviceID, "coordinator"))
		if value == 0 {
			value = testutil.ToFloat64(s.metrics.SLZBDeviceMode.WithLabelValues(deviceID, "router"))
		}

		return fmt.Sprintf("%.0f", value)
	case "SLZBHeapFree":
		value := testutil.ToFloat64(s.metrics.SLZBHeapFree.WithLabelValues(deviceID))
		return fmt.Sprintf("%.0f", value)
	case "SLZBHeapSize":
		value := testutil.ToFloat64(s.metrics.SLZBHeapSize.WithLabelValues(deviceID))
		return fmt.Sprintf("%.0f", value)
	case "SLZBHeapRatio":
		value := testutil.ToFloat64(s.metrics.SLZBHeapRatio.WithLabelValues(deviceID))
		return fmt.Sprintf("%.1f", value)
	case "SLZBEthernetConnected":
		// Try to get ethernet connection value
		value := testutil.ToFloat64(s.metrics.SLZBEthernetConnected.WithLabelValues(deviceID, "192.168.1.100", "00:11:22:33:44:55", "192.168.1.1", "255.255.255.0", "8.8.8.8", "1000"))
		if value == 0 {
			value = testutil.ToFloat64(s.metrics.SLZBEthernetConnected.WithLabelValues(deviceID, "unknown", "unknown", "unknown", "unknown", "unknown", "unknown"))
		}

		return fmt.Sprintf("%.0f", value)
	case "SLZBWifiConnected":
		// Try to get wifi connection value
		value := testutil.ToFloat64(s.metrics.SLZBWifiConnected.WithLabelValues(deviceID, "MyWiFi", "192.168.1.101", "00:11:22:33:44:56", "192.168.1.1", "255.255.255.0", "8.8.8.8"))
		if value == 0 {
			value = testutil.ToFloat64(s.metrics.SLZBWifiConnected.WithLabelValues(deviceID, "none", "none", "none", "none", "none", "none"))
		}

		return fmt.Sprintf("%.0f", value)
	case "SLZBHTTPRequestsTotal":
		value := testutil.ToFloat64(s.metrics.SLZBHTTPRequestsTotal.WithLabelValues(deviceID, "0", "200"))
		return fmt.Sprintf("%.0f", value)
	case "SLZBHTTPErrorsTotal":
		value := testutil.ToFloat64(s.metrics.SLZBHTTPErrorsTotal.WithLabelValues(deviceID, "0", "timeout"))
		return fmt.Sprintf("%.0f", value)
	case "SLZBDeviceReachable":
		value := testutil.ToFloat64(s.metrics.SLZBDeviceReachable.WithLabelValues(deviceID))
		return fmt.Sprintf("%.0f", value)
	case "SLZBLastCollectionTime":
		value := testutil.ToFloat64(s.metrics.SLZBLastCollectionTime.WithLabelValues(deviceID))
		return fmt.Sprintf("%.0f", value)
	case "SLZBCollectionErrors":
		value := testutil.ToFloat64(s.metrics.SLZBCollectionErrors.WithLabelValues(deviceID, "timeout"))
		return fmt.Sprintf("%.0f", value)
	case "SLZBFirmwareCurrentVersion":
		value := testutil.ToFloat64(s.metrics.SLZBFirmwareCurrentVersion.WithLabelValues(deviceID, "1.0.0", "2024-01-01"))
		return fmt.Sprintf("%.0f", value)
	case "SLZBFirmwareUpdateAvailable":
		value := testutil.ToFloat64(s.metrics.SLZBFirmwareUpdateAvailable.WithLabelValues(deviceID, "unknown"))
		return fmt.Sprintf("%.0f", value)
	case "SLZBFirmwareLastCheckTime":
		value := testutil.ToFloat64(s.metrics.SLZBFirmwareLastCheckTime.WithLabelValues(deviceID))
		return fmt.Sprintf("%.0f", value)
	case "SLZBConfigBackupStatus":
		value := testutil.ToFloat64(s.metrics.SLZBConfigBackupStatus.WithLabelValues(deviceID, "auto"))
		return fmt.Sprintf("%.0f", value)
	case "SLZBConfigLastBackupTime":
		value := testutil.ToFloat64(s.metrics.SLZBConfigLastBackupTime.WithLabelValues(deviceID, "auto"))
		return fmt.Sprintf("%.0f", value)
	case "SLZBConfigFileCount":
		value := testutil.ToFloat64(s.metrics.SLZBConfigFileCount.WithLabelValues(deviceID, "config"))
		return fmt.Sprintf("%.0f", value)
	case "SLZBAPIResponseTimeSeconds":
		// For histograms, we'll show a default value since we can't easily read the current value
		return "0.125"
	case "SLZBAPITimeoutErrorsTotal":
		value := testutil.ToFloat64(s.metrics.SLZBAPITimeoutErrorsTotal.WithLabelValues(deviceID, "0"))
		return fmt.Sprintf("%.0f", value)
	case "SLZBCollectionDurationSeconds":
		// For histograms, we'll show a default value since we can't easily read the current value
		return "0.125"
	default:
		return "0"
	}
}

// getMetricLabels gets the labels from the registry for a metric
func (s *Server) getMetricLabels(metricName string) map[string]string {
	deviceID := s.getDeviceID()

	switch metricName {
	case "SLZBConnected", "SLZBDeviceTemp", "SLZBUptime", "SLZBHeapFree", "SLZBHeapSize", "SLZBHeapRatio":
		return map[string]string{"device": deviceID}
	case "SLZBSocketConnected":
		// Check if we have any socket connection data
		if testutil.ToFloat64(s.metrics.SLZBSocketConnected.WithLabelValues(deviceID, "1")) > 0 {
			return map[string]string{"device": deviceID, "connections": "1"}
		}

		return map[string]string{"device": deviceID, "connections": "0"}
	case "SLZBDeviceMode":
		// Check for common modes
		if testutil.ToFloat64(s.metrics.SLZBDeviceMode.WithLabelValues(deviceID, "coordinator")) > 0 {
			return map[string]string{"device": deviceID, "mode": "coordinator"}
		}

		if testutil.ToFloat64(s.metrics.SLZBDeviceMode.WithLabelValues(deviceID, "router")) > 0 {
			return map[string]string{"device": deviceID, "mode": "router"}
		}

		return map[string]string{"device": deviceID, "mode": "unknown"}
	case "SLZBEthernetConnected":
		// Try to get real ethernet connection data
		value := testutil.ToFloat64(s.metrics.SLZBEthernetConnected.WithLabelValues(deviceID, "192.168.1.100", "00:11:22:33:44:55", "192.168.1.1", "255.255.255.0", "8.8.8.8", "1000"))
		if value > 0 {
			return map[string]string{
				"device":      deviceID,
				"ip_address":  "192.168.1.100",
				"mac_address": "00:11:22:33:44:55",
				"gateway":     "192.168.1.1",
				"subnet_mask": "255.255.255.0",
				"dns_server":  "8.8.8.8",
				"speed_mbps":  "1000",
			}
		}

		return map[string]string{
			"device":      deviceID,
			"ip_address":  "unknown",
			"mac_address": "unknown",
			"gateway":     "unknown",
			"subnet_mask": "unknown",
			"dns_server":  "unknown",
			"speed_mbps":  "unknown",
		}
	case "SLZBWifiConnected":
		// Try to get real wifi connection data
		value := testutil.ToFloat64(s.metrics.SLZBWifiConnected.WithLabelValues(deviceID, "MyWiFi", "192.168.1.101", "00:11:22:33:44:56", "192.168.1.1", "255.255.255.0", "8.8.8.8"))
		if value > 0 {
			return map[string]string{
				"device":      deviceID,
				"ssid":        "MyWiFi",
				"ip_address":  "192.168.1.101",
				"mac_address": "00:11:22:33:44:56",
				"gateway":     "192.168.1.1",
				"subnet_mask": "255.255.255.0",
				"dns_server":  "8.8.8.8",
			}
		}

		return map[string]string{
			"device":      deviceID,
			"ssid":        "none",
			"ip_address":  "none",
			"mac_address": "none",
			"gateway":     "none",
			"subnet_mask": "none",
			"dns_server":  "none",
		}
	case "SLZBHTTPRequestsTotal":
		// Check for real HTTP request data
		value := testutil.ToFloat64(s.metrics.SLZBHTTPRequestsTotal.WithLabelValues(deviceID, "0", "200"))
		if value > 0 {
			return map[string]string{"device": deviceID, "action": "0", "status": "200"}
		}

		return map[string]string{"device": deviceID, "action": "get_status", "status": "200"}
	case "SLZBHTTPErrorsTotal":
		// Check for real HTTP error data
		value := testutil.ToFloat64(s.metrics.SLZBHTTPErrorsTotal.WithLabelValues(deviceID, "0", "timeout"))
		if value > 0 {
			return map[string]string{"device": deviceID, "action": "0", "error_type": "timeout"}
		}

		return map[string]string{"device": deviceID, "action": "get_status", "error_type": "timeout"}
	case "SLZBDeviceReachable", "SLZBLastCollectionTime":
		return map[string]string{"device": deviceID}
	case "SLZBCollectionErrors":
		// Check for real collection error data
		value := testutil.ToFloat64(s.metrics.SLZBCollectionErrors.WithLabelValues(deviceID, "timeout"))
		if value > 0 {
			return map[string]string{"device": deviceID, "error_type": "timeout"}
		}

		return map[string]string{"device": deviceID, "error_type": "timeout"}
	case "SLZBFirmwareCurrentVersion":
		// Check for real firmware data
		value := testutil.ToFloat64(s.metrics.SLZBFirmwareCurrentVersion.WithLabelValues(deviceID, "1.0.0", "2024-01-01"))
		if value > 0 {
			return map[string]string{"device": deviceID, "version": "1.0.0", "build_date": "2024-01-01"}
		}

		return map[string]string{"device": deviceID, "version": "unknown", "build_date": "unknown"}
	case "SLZBFirmwareUpdateAvailable":
		// Check for real firmware update data
		value := testutil.ToFloat64(s.metrics.SLZBFirmwareUpdateAvailable.WithLabelValues(deviceID, "unknown"))
		if value > 0 {
			return map[string]string{"device": deviceID, "available_version": "unknown"}
		}

		return map[string]string{"device": deviceID, "available_version": "none"}
	case "SLZBFirmwareLastCheckTime":
		return map[string]string{"device": deviceID}
	case "SLZBConfigBackupStatus":
		// Check for real backup status data
		value := testutil.ToFloat64(s.metrics.SLZBConfigBackupStatus.WithLabelValues(deviceID, "auto"))
		if value > 0 {
			return map[string]string{"device": deviceID, "backup_type": "auto"}
		}

		return map[string]string{"device": deviceID, "backup_type": "manual"}
	case "SLZBConfigLastBackupTime":
		// Check for real backup time data
		value := testutil.ToFloat64(s.metrics.SLZBConfigLastBackupTime.WithLabelValues(deviceID, "auto"))
		if value > 0 {
			return map[string]string{"device": deviceID, "backup_type": "auto"}
		}

		return map[string]string{"device": deviceID, "backup_type": "manual"}
	case "SLZBConfigFileCount":
		// Check for real file count data
		value := testutil.ToFloat64(s.metrics.SLZBConfigFileCount.WithLabelValues(deviceID, "config"))
		if value > 0 {
			return map[string]string{"device": deviceID, "file_type": "config"}
		}

		return map[string]string{"device": deviceID, "file_type": "configuration"}
	case "SLZBAPIResponseTimeSeconds":
		// For histograms, we'll show a default label since we can't easily read the current value
		return map[string]string{"device": deviceID, "action": "get_status"}
	case "SLZBAPITimeoutErrorsTotal":
		// Check for real timeout error data
		value := testutil.ToFloat64(s.metrics.SLZBAPITimeoutErrorsTotal.WithLabelValues(deviceID, "0"))
		if value > 0 {
			return map[string]string{"device": deviceID, "action": "0"}
		}

		return map[string]string{"device": deviceID, "action": "get_status"}
	case "SLZBCollectionDurationSeconds":
		return map[string]string{"device": deviceID}
	default:
		return map[string]string{"device": deviceID}
	}
}

// getDeviceID derives the device ID from the configuration
func (s *Server) getDeviceID() string {
	// Extract hostname/IP from API URL
	apiURL := s.config.SLZB.APIURL
	if strings.HasPrefix(apiURL, "http://") {
		apiURL = strings.TrimPrefix(apiURL, "http://")
	} else if strings.HasPrefix(apiURL, "https://") {
		apiURL = strings.TrimPrefix(apiURL, "https://")
	}

	// Remove port if present
	if idx := strings.Index(apiURL, ":"); idx != -1 {
		apiURL = apiURL[:idx]
	}

	// Use the hostname/IP as device ID
	return apiURL
}

// getMetricHelp gets the help text from the Prometheus registry
func (s *Server) getMetricHelp(metricName string) string {
	// Get all metric families from the default registry
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		slog.Warn("Failed to gather metrics for help text", "error", err)
		return s.getFallbackMetricHelp(metricName)
	}

	// Search for the metric by name
	for _, family := range metricFamilies {
		if family.GetName() == metricName {
			return family.GetHelp()
		}
	}

	// If not found in default registry, fall back to hardcoded help
	return s.getFallbackMetricHelp(metricName)
}

// getFallbackMetricHelp provides fallback help text when dynamic extraction fails
func (s *Server) getFallbackMetricHelp(metricName string) string {
	switch metricName {
	case "slzb_device_connected":
		return "SLZB device connection status (1=connected, 0=disconnected)"
	case "slzb_device_temperature_celsius":
		return "SLZB device temperature in degrees Celsius"
	case "slzb_device_uptime_seconds":
		return "SLZB device uptime in seconds since last boot"
	case "slzb_socket_uptime_seconds":
		return "SLZB socket connection uptime in seconds since connection established"
	case "slzb_socket_connected":
		return "SLZB socket connection status (1=connected, 0=disconnected)"
	case "slzb_device_operational_mode":
		return "SLZB device operational mode (1=active, 0=inactive) with mode label"
	case "slzb_device_heap_free_kb":
		return "SLZB device free heap memory in kilobytes"
	case "slzb_device_heap_size_kb":
		return "SLZB device total heap memory in kilobytes"
	case "slzb_device_heap_ratio":
		return "SLZB device heap usage ratio as percentage (free heap / total heap * 100)"
	case "slzb_device_ethernet_connected":
		return "SLZB device Ethernet connection status (1=connected, 0=disconnected)"
	case "slzb_device_wifi_connected":
		return "SLZB device WiFi connection status (1=connected, 0=disconnected)"
	case "slzb_http_requests_total":
		return "Total number of HTTP requests made by exporter to SLZB device API"
	case "slzb_http_errors_total":
		return "Total number of HTTP errors when making requests to SLZB device API"
	case "slzb_device_reachable":
		return "SLZB device reachability status (1=reachable, 0=unreachable)"
	case "slzb_last_collection_timestamp":
		return "Unix timestamp of the last successful collection from SLZB device"
	case "slzb_collection_errors_total":
		return "Total number of collection errors for SLZB device by error type"
	case "slzb_firmware_current_version":
		return "Current firmware version (always 1, used for joining with labels)"
	case "slzb_firmware_update_available":
		return "Firmware update availability (1=available, 0=not_available)"
	case "slzb_firmware_last_check_timestamp":
		return "Unix timestamp of last firmware check"
	case "slzb_config_backup_status":
		return "Status of the last configuration backup (1=success, 0=failure)"
	case "slzb_config_last_backup_timestamp":
		return "Unix timestamp of the last successful configuration backup"
	case "slzb_config_file_count":
		return "Number of configuration files on the device"
	case "slzb_api_response_time_seconds":
		return "Histogram of API response times in seconds"
	case "slzb_api_timeout_errors_total":
		return "Total number of API timeout errors"
	case "slzb_collection_duration_seconds":
		return "Histogram of collection durations in seconds"
	default:
		return "SLZB device metric"
	}
}

type MetricInfo struct {
	Name         string
	Help         string
	ExampleValue string
	Labels       map[string]string
}

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Registry holds all the metrics
type Registry struct {
	// Version information
	VersionInfo *prometheus.GaugeVec

	// SLZB device metrics
	SLZBConnected         *prometheus.GaugeVec
	SLZBDeviceTemp        *prometheus.GaugeVec
	SLZBUptime            *prometheus.GaugeVec
	SLZBSocketUptime      *prometheus.GaugeVec
	SLZBSocketConnected   *prometheus.GaugeVec
	SLZBDeviceMode        *prometheus.GaugeVec
	SLZBHeapFree          *prometheus.GaugeVec
	SLZBHeapSize          *prometheus.GaugeVec
	SLZBHeapRatio         *prometheus.GaugeVec
	SLZBEthernetConnected *prometheus.GaugeVec
	SLZBWifiConnected     *prometheus.GaugeVec

	// HTTP request metrics (exporter's requests to SLZB device)
	SLZBHTTPRequestsTotal *prometheus.CounterVec
	SLZBHTTPErrorsTotal   *prometheus.CounterVec

	// Device health and availability metrics
	SLZBDeviceReachable    *prometheus.GaugeVec
	SLZBLastCollectionTime *prometheus.GaugeVec
	SLZBCollectionErrors   *prometheus.CounterVec

	// NEW: Firmware Update Status
	SLZBFirmwareCurrentVersion  *prometheus.GaugeVec
	SLZBFirmwareUpdateAvailable *prometheus.GaugeVec
	SLZBFirmwareLastCheckTime   *prometheus.GaugeVec

	// NEW: Configuration Management
	SLZBConfigBackupStatus   *prometheus.GaugeVec
	SLZBConfigLastBackupTime *prometheus.GaugeVec
	SLZBConfigFileCount      *prometheus.GaugeVec

	// NEW: Performance Benchmarks
	SLZBAPIResponseTimeSeconds    *prometheus.HistogramVec
	SLZBAPITimeoutErrorsTotal     *prometheus.CounterVec
	SLZBCollectionDurationSeconds *prometheus.HistogramVec

	// metricInfo holds metric metadata for the UI
	metricInfo []MetricInfo
}

// MetricInfo contains information about a metric for the UI
type MetricInfo struct {
	Name         string
	Help         string
	ExampleValue string
	Labels       []string
}

// addMetricInfo adds metric information to the slice
func (r *Registry) addMetricInfo(name, help string, labels []string) {
	r.metricInfo = append(r.metricInfo, MetricInfo{
		Name:         name,
		Help:         help,
		Labels:       labels,
		ExampleValue: "",
	})
}

// GetMetricsInfo returns information about all metrics for the UI
func (r *Registry) GetMetricsInfo() []MetricInfo {
	return r.metricInfo
}

// NewRegistry creates a new metrics registry
func NewRegistry() *Registry {
	r := &Registry{
		VersionInfo: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_exporter_version_info",
				Help: "Version information about the SLZB exporter",
			},
			[]string{"version", "commit", "build_date"},
		),

		SLZBConnected: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_device_connected",
				Help: "SLZB device connection status (1=connected, 0=disconnected)",
			},
			[]string{"device"},
		),

		SLZBDeviceTemp: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_device_temperature_celsius",
				Help: "SLZB device temperature in degrees Celsius",
			},
			[]string{"device"},
		),

		SLZBUptime: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_device_uptime_seconds",
				Help: "SLZB device uptime in seconds since last boot",
			},
			[]string{"device"},
		),

		SLZBSocketUptime: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_socket_uptime_seconds",
				Help: "SLZB socket connection uptime in seconds since connection established",
			},
			[]string{"device"},
		),

		SLZBSocketConnected: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_socket_connected",
				Help: "SLZB socket connection status (1=connected, 0=disconnected)",
			},
			[]string{"device", "connections"},
		),

		SLZBDeviceMode: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_device_operational_mode",
				Help: "SLZB device operational mode (1=active, 0=inactive) with mode label",
			},
			[]string{"device", "mode"},
		),

		SLZBHeapFree: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_device_heap_free_kb",
				Help: "SLZB device free heap memory in kilobytes",
			},
			[]string{"device"},
		),

		SLZBHeapSize: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_device_heap_size_kb",
				Help: "SLZB device total heap memory in kilobytes",
			},
			[]string{"device"},
		),

		SLZBHeapRatio: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_device_heap_ratio",
				Help: "SLZB device heap usage ratio as percentage (free heap / total heap * 100)",
			},
			[]string{"device"},
		),

		SLZBEthernetConnected: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_device_ethernet_connected",
				Help: "SLZB device Ethernet connection status (1=connected, 0=disconnected)",
			},
			[]string{"device", "ip_address", "mac_address", "gateway", "subnet_mask", "dns_server", "speed_mbps"},
		),

		SLZBWifiConnected: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_device_wifi_connected",
				Help: "SLZB device WiFi connection status (1=connected, 0=disconnected)",
			},
			[]string{"device", "ssid", "ip_address", "mac_address", "gateway", "subnet_mask", "dns_server"},
		),

		SLZBHTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "slzb_http_requests_total",
				Help: "Total number of HTTP requests made by exporter to SLZB device API",
			},
			[]string{"device", "action", "status"},
		),

		SLZBHTTPErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "slzb_http_errors_total",
				Help: "Total number of HTTP errors when making requests to SLZB device API",
			},
			[]string{"device", "action", "error_type"},
		),

		SLZBDeviceReachable: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_device_reachable",
				Help: "SLZB device reachability status (1=reachable, 0=unreachable)",
			},
			[]string{"device"},
		),

		SLZBLastCollectionTime: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_last_collection_timestamp",
				Help: "Unix timestamp of the last successful collection from SLZB device",
			},
			[]string{"device"},
		),

		SLZBCollectionErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "slzb_collection_errors_total",
				Help: "Total number of collection errors for SLZB device by error type",
			},
			[]string{"device", "error_type"},
		),

		// NEW: Firmware Update Status
		SLZBFirmwareCurrentVersion: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_firmware_current_version",
				Help: "Current firmware version (always 1, used for joining with labels)",
			},
			[]string{"device", "version", "build_date"},
		),

		SLZBFirmwareUpdateAvailable: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_firmware_update_available",
				Help: "Firmware update availability (1=available, 0=not_available)",
			},
			[]string{"device", "available_version"},
		),

		SLZBFirmwareLastCheckTime: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_firmware_last_check_timestamp",
				Help: "Unix timestamp of last firmware check",
			},
			[]string{"device"},
		),

		// NEW: Configuration Management
		SLZBConfigBackupStatus: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_config_backup_status",
				Help: "Status of the last configuration backup (1=success, 0=failure)",
			},
			[]string{"device", "backup_type"},
		),

		SLZBConfigLastBackupTime: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_config_last_backup_timestamp",
				Help: "Unix timestamp of the last successful configuration backup",
			},
			[]string{"device", "backup_type"},
		),

		SLZBConfigFileCount: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_config_file_count",
				Help: "Number of configuration files on the device",
			},
			[]string{"device", "file_type"},
		),

		// NEW: Performance Benchmarks
		SLZBAPIResponseTimeSeconds: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "slzb_api_response_time_seconds",
				Help:    "Histogram of API response times in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10, 30, 60, 120, 300, 600, 1800, 3600},
			},
			[]string{"device", "action"},
		),

		SLZBAPITimeoutErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "slzb_api_timeout_errors_total",
				Help: "Total number of API timeout errors",
			},
			[]string{"device", "action"},
		),

		SLZBCollectionDurationSeconds: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "slzb_collection_duration_seconds",
				Help:    "Histogram of collection durations in seconds",
				Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60, 120, 300, 600, 1800, 3600},
			},
			[]string{"device"},
		),
	}

	// Add metric info for UI
	r.addMetricInfo("slzb_exporter_version_info", "Version information about the SLZB exporter", []string{"version", "commit", "build_date"})
	r.addMetricInfo("slzb_device_connected", "SLZB device connection status (1=connected, 0=disconnected)", []string{"device"})
	r.addMetricInfo("slzb_device_temperature_celsius", "SLZB device temperature in degrees Celsius", []string{"device"})
	r.addMetricInfo("slzb_device_uptime_seconds", "SLZB device uptime in seconds since last boot", []string{"device"})
	r.addMetricInfo("slzb_socket_uptime_seconds", "SLZB socket connection uptime in seconds since connection established", []string{"device"})
	r.addMetricInfo("slzb_socket_connected", "SLZB socket connection status (1=connected, 0=disconnected)", []string{"device", "connections"})
	r.addMetricInfo("slzb_device_operational_mode", "SLZB device operational mode (1=active, 0=inactive) with mode label", []string{"device", "mode"})
	r.addMetricInfo("slzb_device_heap_free_kb", "SLZB device free heap memory in kilobytes", []string{"device"})
	r.addMetricInfo("slzb_device_heap_size_kb", "SLZB device total heap memory in kilobytes", []string{"device"})
	r.addMetricInfo("slzb_device_heap_ratio", "SLZB device heap usage ratio as percentage (free heap / total heap * 100)", []string{"device"})
	r.addMetricInfo("slzb_device_ethernet_connected", "SLZB device Ethernet connection status (1=connected, 0=disconnected)", []string{"device", "ip_address", "mac_address", "gateway", "subnet_mask", "dns_server", "speed_mbps"})
	r.addMetricInfo("slzb_device_wifi_connected", "SLZB device WiFi connection status (1=connected, 0=disconnected)", []string{"device", "ssid", "ip_address", "mac_address", "gateway", "subnet_mask", "dns_server"})
	r.addMetricInfo("slzb_http_requests_total", "Total number of HTTP requests made by exporter to SLZB device API", []string{"device", "action", "status"})
	r.addMetricInfo("slzb_http_errors_total", "Total number of HTTP errors when making requests to SLZB device API", []string{"device", "action", "error_type"})
	r.addMetricInfo("slzb_device_reachable", "SLZB device reachability status (1=reachable, 0=unreachable)", []string{"device"})
	r.addMetricInfo("slzb_last_collection_timestamp", "Unix timestamp of the last successful collection from SLZB device", []string{"device"})
	r.addMetricInfo("slzb_collection_errors_total", "Total number of collection errors for SLZB device by error type", []string{"device", "error_type"})
	r.addMetricInfo("slzb_firmware_current_version", "Current firmware version (always 1, used for joining with labels)", []string{"device", "version", "build_date"})
	r.addMetricInfo("slzb_firmware_update_available", "Firmware update availability (1=available, 0=not_available)", []string{"device", "available_version"})
	r.addMetricInfo("slzb_firmware_last_check_timestamp", "Unix timestamp of last firmware check", []string{"device"})
	r.addMetricInfo("slzb_config_backup_status", "Status of the last configuration backup (1=success, 0=failure)", []string{"device", "backup_type"})
	r.addMetricInfo("slzb_config_last_backup_timestamp", "Unix timestamp of the last successful configuration backup", []string{"device", "backup_type"})
	r.addMetricInfo("slzb_config_file_count", "Number of configuration files on the device", []string{"device", "file_type"})
	r.addMetricInfo("slzb_api_response_time_seconds", "Histogram of API response times in seconds", []string{"device", "action"})
	r.addMetricInfo("slzb_api_timeout_errors_total", "Total number of API timeout errors", []string{"device", "action"})
	r.addMetricInfo("slzb_collection_duration_seconds", "Histogram of collection durations in seconds", []string{"device"})

	return r
}

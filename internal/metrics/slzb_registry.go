package metrics

import (
	promexporter_metrics "github.com/d0ugal/promexporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// SLZBRegistry wraps the promexporter registry with SLZB-specific metrics
type SLZBRegistry struct {
	*promexporter_metrics.Registry

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

	// Firmware Update Status
	SLZBFirmwareCurrentVersion  *prometheus.GaugeVec
	SLZBFirmwareUpdateAvailable *prometheus.GaugeVec
	SLZBFirmwareLastCheckTime   *prometheus.GaugeVec

	// Configuration Management
	SLZBConfigBackupStatus   *prometheus.GaugeVec
	SLZBConfigLastBackupTime *prometheus.GaugeVec
	SLZBConfigFileCount      *prometheus.GaugeVec

	// Performance Benchmarks
	SLZBAPIResponseTimeSeconds    *prometheus.HistogramVec
	SLZBAPITimeoutErrorsTotal     *prometheus.CounterVec
	SLZBCollectionDurationSeconds *prometheus.HistogramVec
}

// NewSLZBRegistry creates a new SLZB metrics registry
func NewSLZBRegistry(baseRegistry *promexporter_metrics.Registry) *SLZBRegistry {
	// Get the underlying Prometheus registry
	promRegistry := baseRegistry.GetRegistry()
	factory := promauto.With(promRegistry)

	slzb := &SLZBRegistry{
		Registry: baseRegistry,
	}

	// SLZB device metrics
	slzb.SLZBConnected = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_device_connected",
			Help: "SLZB device connection status (1=connected, 0=disconnected)",
		},
		[]string{"device"},
	)

	baseRegistry.AddMetricInfo("slzb_device_connected", "SLZB device connection status (1=connected, 0=disconnected)", []string{"device"})

	slzb.SLZBDeviceTemp = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_device_temperature_celsius",
			Help: "SLZB device temperature in degrees Celsius",
		},
		[]string{"device"},
	)

	baseRegistry.AddMetricInfo("slzb_device_temperature_celsius", "SLZB device temperature in degrees Celsius", []string{"device"})

	slzb.SLZBUptime = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_device_uptime_seconds",
			Help: "SLZB device uptime in seconds since last boot",
		},
		[]string{"device"},
	)

	baseRegistry.AddMetricInfo("slzb_device_uptime_seconds", "SLZB device uptime in seconds since last boot", []string{"device"})

	slzb.SLZBSocketUptime = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_socket_uptime_seconds",
			Help: "SLZB socket connection uptime in seconds since connection established",
		},
		[]string{"device"},
	)

	baseRegistry.AddMetricInfo("slzb_socket_uptime_seconds", "SLZB socket connection uptime in seconds since connection established", []string{"device"})

	slzb.SLZBSocketConnected = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_socket_connected",
			Help: "SLZB socket connection status (1=connected, 0=disconnected)",
		},
		[]string{"device", "connections"},
	)

	baseRegistry.AddMetricInfo("slzb_socket_connected", "SLZB socket connection status (1=connected, 0=disconnected)", []string{"device", "connections"})

	slzb.SLZBDeviceMode = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_device_operational_mode",
			Help: "SLZB device operational mode (1=active, 0=inactive) with mode label",
		},
		[]string{"device", "mode"},
	)

	baseRegistry.AddMetricInfo("slzb_device_operational_mode", "SLZB device operational mode (1=active, 0=inactive) with mode label", []string{"device", "mode"})

	slzb.SLZBHeapFree = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_device_heap_free_bytes",
			Help: "SLZB device free heap memory in bytes",
		},
		[]string{"device"},
	)

	baseRegistry.AddMetricInfo("slzb_device_heap_free_bytes", "SLZB device free heap memory in bytes", []string{"device"})

	slzb.SLZBHeapSize = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_device_heap_size_bytes",
			Help: "SLZB device total heap memory in bytes",
		},
		[]string{"device"},
	)

	baseRegistry.AddMetricInfo("slzb_device_heap_size_bytes", "SLZB device total heap memory in bytes", []string{"device"})

	slzb.SLZBHeapRatio = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_device_heap_ratio",
			Help: "SLZB device heap usage ratio as percentage (free heap / total heap * 100)",
		},
		[]string{"device"},
	)

	baseRegistry.AddMetricInfo("slzb_device_heap_ratio", "SLZB device heap usage ratio as percentage (free heap / total heap * 100)", []string{"device"})

	slzb.SLZBEthernetConnected = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_device_ethernet_connected",
			Help: "SLZB device Ethernet connection status (1=connected, 0=disconnected)",
		},
		[]string{"device", "ip_address", "mac_address", "gateway", "subnet_mask", "dns_server", "speed_mbps"},
	)

	baseRegistry.AddMetricInfo("slzb_device_ethernet_connected", "SLZB device Ethernet connection status (1=connected, 0=disconnected)", []string{"device", "ip_address", "mac_address", "gateway", "subnet_mask", "dns_server", "speed_mbps"})

	slzb.SLZBWifiConnected = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_device_wifi_connected",
			Help: "SLZB device WiFi connection status (1=connected, 0=disconnected)",
		},
		[]string{"device", "ssid", "ip_address", "mac_address", "gateway", "subnet_mask", "dns_server"},
	)

	baseRegistry.AddMetricInfo("slzb_device_wifi_connected", "SLZB device WiFi connection status (1=connected, 0=disconnected)", []string{"device", "ssid", "ip_address", "mac_address", "gateway", "subnet_mask", "dns_server"})

	slzb.SLZBHTTPRequestsTotal = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "slzb_http_requests_total",
			Help: "Total number of HTTP requests made by exporter to SLZB device API",
		},
		[]string{"device", "action", "status"},
	)

	baseRegistry.AddMetricInfo("slzb_http_requests_total", "Total number of HTTP requests made by exporter to SLZB device API", []string{"device", "action", "status"})

	slzb.SLZBHTTPErrorsTotal = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "slzb_http_errors_total",
			Help: "Total number of HTTP errors when making requests to SLZB device API",
		},
		[]string{"device", "action", "error_type"},
	)

	baseRegistry.AddMetricInfo("slzb_http_errors_total", "Total number of HTTP errors when making requests to SLZB device API", []string{"device", "action", "error_type"})

	slzb.SLZBDeviceReachable = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_device_reachable",
			Help: "SLZB device reachability status (1=reachable, 0=unreachable)",
		},
		[]string{"device"},
	)

	baseRegistry.AddMetricInfo("slzb_device_reachable", "SLZB device reachability status (1=reachable, 0=unreachable)", []string{"device"})

	slzb.SLZBLastCollectionTime = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_last_collection_timestamp",
			Help: "Unix timestamp of the last successful collection from SLZB device",
		},
		[]string{"device"},
	)

	baseRegistry.AddMetricInfo("slzb_last_collection_timestamp", "Unix timestamp of the last successful collection from SLZB device", []string{"device"})

	slzb.SLZBCollectionErrors = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "slzb_collection_errors_total",
			Help: "Total number of collection errors for SLZB device by error type",
		},
		[]string{"device", "error_type"},
	)

	baseRegistry.AddMetricInfo("slzb_collection_errors_total", "Total number of collection errors for SLZB device by error type", []string{"device", "error_type"})

	// Firmware Update Status
	slzb.SLZBFirmwareCurrentVersion = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_firmware_current_version",
			Help: "Current firmware version (always 1, used for joining with labels)",
		},
		[]string{"device", "version", "build_date"},
	)

	baseRegistry.AddMetricInfo("slzb_firmware_current_version", "Current firmware version (always 1, used for joining with labels)", []string{"device", "version", "build_date"})

	slzb.SLZBFirmwareUpdateAvailable = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_firmware_update_available",
			Help: "Firmware update availability (1=available, 0=not_available)",
		},
		[]string{"device", "available_version"},
	)

	baseRegistry.AddMetricInfo("slzb_firmware_update_available", "Firmware update availability (1=available, 0=not_available)", []string{"device", "available_version"})

	slzb.SLZBFirmwareLastCheckTime = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_firmware_last_check_timestamp",
			Help: "Unix timestamp of last firmware check",
		},
		[]string{"device"},
	)

	baseRegistry.AddMetricInfo("slzb_firmware_last_check_timestamp", "Unix timestamp of last firmware check", []string{"device"})

	// Configuration Management
	slzb.SLZBConfigBackupStatus = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_config_backup_status",
			Help: "Status of the last configuration backup (1=success, 0=failure)",
		},
		[]string{"device", "backup_type"},
	)

	baseRegistry.AddMetricInfo("slzb_config_backup_status", "Status of the last configuration backup (1=success, 0=failure)", []string{"device", "backup_type"})

	slzb.SLZBConfigLastBackupTime = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_config_last_backup_timestamp",
			Help: "Unix timestamp of the last successful configuration backup",
		},
		[]string{"device", "backup_type"},
	)

	baseRegistry.AddMetricInfo("slzb_config_last_backup_timestamp", "Unix timestamp of the last successful configuration backup", []string{"device", "backup_type"})

	slzb.SLZBConfigFileCount = factory.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "slzb_config_files",
			Help: "Number of configuration files on the device",
		},
		[]string{"device", "file_type"},
	)

	baseRegistry.AddMetricInfo("slzb_config_files", "Number of configuration files on the device", []string{"device", "file_type"})

	// Performance Benchmarks
	slzb.SLZBAPIResponseTimeSeconds = factory.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "slzb_api_response_time_seconds",
			Help:    "Histogram of API response times in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10, 30, 60, 120, 300, 600, 1800, 3600},
		},
		[]string{"device", "action"},
	)

	baseRegistry.AddMetricInfo("slzb_api_response_time_seconds", "Histogram of API response times in seconds", []string{"device", "action"})

	slzb.SLZBAPITimeoutErrorsTotal = factory.NewCounterVec(
		prometheus.CounterOpts{
			Name: "slzb_api_timeout_errors_total",
			Help: "Total number of API timeout errors when making requests to SLZB device",
		},
		[]string{"device", "action"},
	)

	baseRegistry.AddMetricInfo("slzb_api_timeout_errors_total", "Total number of API timeout errors when making requests to SLZB device", []string{"device", "action"})

	slzb.SLZBCollectionDurationSeconds = factory.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "slzb_collection_duration_seconds",
			Help:    "Histogram of collection duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10, 30, 60, 120, 300, 600, 1800, 3600},
		},
		[]string{"device"},
	)

	baseRegistry.AddMetricInfo("slzb_collection_duration_seconds", "Histogram of collection duration in seconds", []string{"device"})

	return slzb
}

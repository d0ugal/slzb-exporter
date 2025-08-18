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
}

// NewRegistry creates a new metrics registry
func NewRegistry() *Registry {
	return &Registry{
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
}

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

	// NEW: Zigbee Network Statistics
	SLZBZigbeePacketsReceived    *prometheus.CounterVec
	SLZBZigbeePacketsSent        *prometheus.CounterVec
	SLZBZigbeeErrorsTotal        *prometheus.CounterVec
	SLZBZigbeeNetworkDevices     *prometheus.GaugeVec
	SLZBZigbeeChannelUtilization *prometheus.GaugeVec
	SLZBZigbeeInterferenceLevel  *prometheus.GaugeVec

	// NEW: Firmware Update Status
	SLZBFirmwareCurrentVersion  *prometheus.GaugeVec
	SLZBFirmwareUpdateAvailable *prometheus.GaugeVec
	SLZBFirmwareLastCheckTime   *prometheus.GaugeVec

	// NEW: Configuration Management
	SLZBConfigBackupStatus   *prometheus.GaugeVec
	SLZBConfigLastBackupTime *prometheus.GaugeVec
	SLZBConfigFileCount      *prometheus.GaugeVec
	SLZBConfigTotalSizeBytes *prometheus.GaugeVec

	// NEW: Network Security Metrics
	SLZBSecurityKeyRotationTime *prometheus.GaugeVec
	SLZBEncryptionStatus        *prometheus.GaugeVec
	SLZBSecurityEventsTotal     *prometheus.CounterVec

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

		// NEW: Zigbee Network Statistics
		SLZBZigbeePacketsReceived: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "slzb_zigbee_packets_received_total",
				Help: "Total number of Zigbee packets received",
			},
			[]string{"device", "packet_type"},
		),

		SLZBZigbeePacketsSent: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "slzb_zigbee_packets_sent_total",
				Help: "Total number of Zigbee packets sent",
			},
			[]string{"device", "packet_type"},
		),

		SLZBZigbeeErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "slzb_zigbee_errors_total",
				Help: "Total number of Zigbee communication errors",
			},
			[]string{"device", "error_type"},
		),

		SLZBZigbeeNetworkDevices: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_zigbee_network_devices",
				Help: "Number of devices in the Zigbee network",
			},
			[]string{"device", "device_type"},
		),

		SLZBZigbeeChannelUtilization: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_zigbee_channel_utilization_percent",
				Help: "Zigbee channel utilization as percentage",
			},
			[]string{"device", "channel"},
		),

		SLZBZigbeeInterferenceLevel: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_zigbee_interference_level",
				Help: "Zigbee interference level (0-255, higher is more interference)",
			},
			[]string{"device", "channel"},
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
				Help: "Configuration backup status (1=success, 0=failed)",
			},
			[]string{"device", "backup_type"},
		),

		SLZBConfigLastBackupTime: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_config_last_backup_timestamp",
				Help: "Unix timestamp of last configuration backup",
			},
			[]string{"device", "backup_type"},
		),

		SLZBConfigFileCount: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_config_file_count",
				Help: "Number of configuration files",
			},
			[]string{"device", "file_type"},
		),

		SLZBConfigTotalSizeBytes: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_config_total_size_bytes",
				Help: "Total size of configuration files in bytes",
			},
			[]string{"device", "file_type"},
		),

		// NEW: Network Security Metrics
		SLZBSecurityKeyRotationTime: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_security_key_rotation_timestamp",
				Help: "Unix timestamp of last security key rotation",
			},
			[]string{"device", "key_type"},
		),

		SLZBEncryptionStatus: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "slzb_encryption_status",
				Help: "Encryption status (1=enabled, 0=disabled)",
			},
			[]string{"device", "encryption_type"},
		),

		SLZBSecurityEventsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "slzb_security_events_total",
				Help: "Total number of security events",
			},
			[]string{"device", "event_type", "severity"},
		),

		// NEW: Performance Benchmarks
		SLZBAPIResponseTimeSeconds: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "slzb_api_response_time_seconds",
				Help:    "API response time in seconds",
				Buckets: prometheus.DefBuckets,
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
				Help:    "Duration of collection cycles in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"device"},
		),
	}
}

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
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	metricsInfo := s.getMetricsInfo()

	// Generate metrics HTML dynamically
	metricsHTML := ""
	for _, metric := range metricsInfo {
		labelsStr := ""
		if len(metric.Labels) > 0 {
			var labelPairs []string
			for k, v := range metric.Labels {
				labelPairs = append(labelPairs, fmt.Sprintf(`%s="%s"`, k, v))
			}
			labelsStr = "{" + strings.Join(labelPairs, ", ") + "}"
		}

		metricsHTML += fmt.Sprintf(`
            <li><strong>%s%s:</strong> %s</li>`,
			metric.Name,
			labelsStr,
			metric.Help)
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    		<title>SLZB Exporter ` + versionInfo.Version + `</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 2rem;
            line-height: 1.6;
            color: #333;
        }
        h1 {
            color: #2c3e50;
            border-bottom: 2px solid #3498db;
            padding-bottom: 0.5rem;
        }
        h1 .version {
            font-size: 0.6em;
            color: #6c757d;
            font-weight: normal;
            margin-left: 0.5rem;
        }
        .endpoint {
            background: #f8f9fa;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 1rem;
            margin: 1rem 0;
        }
        .endpoint h3 {
            margin: 0 0 0.5rem 0;
            color: #495057;
        }
        .endpoint a {
            color: #007bff;
            text-decoration: none;
            font-weight: 500;
        }
        .endpoint a:hover {
            text-decoration: underline;
        }
        .description {
            color: #6c757d;
            font-size: 0.9rem;
        }
        .status {
            display: inline-block;
            padding: 0.25rem 0.5rem;
            border-radius: 4px;
            font-size: 0.8rem;
            font-weight: 500;
        }
        .status.healthy {
            background: #d4edda;
            color: #155724;
        }
        .status.metrics {
            background: #d1ecf1;
            color: #0c5460;
        }
        .status.ready {
            background: #d4edda;
            color: #155724;
        }
        .status.connected {
            background: #d4edda;
            color: #155724;
        }
        .status.disconnected {
            background: #f8d7da;
            color: #721c24;
        }
        .service-status {
            background: #e9ecef;
            border: 1px solid #dee2e6;
            border-radius: 8px;
            padding: 1rem;
            margin: 1rem 0;
        }
        .service-status h3 {
            margin: 0 0 0.5rem 0;
            color: #495057;
        }
        .service-status p {
            margin: 0.25rem 0;
            color: #6c757d;
        }
        .metrics-info {
            background: #e9ecef;
            border: 1px solid #dee2e6;
            border-radius: 8px;
            padding: 1rem;
            margin: 1rem 0;
        }
        .metrics-info h3 {
            margin: 0 0 0.5rem 0;
            color: #495057;
        }
        .metrics-info ul {
            margin: 0.5rem 0;
            padding-left: 1.5rem;
        }
        .metrics-info li {
            margin: 0.25rem 0;
            color: #6c757d;
        }
        .footer {
            margin-top: 2rem;
            padding-top: 1rem;
            border-top: 1px solid #dee2e6;
            text-align: center;
            color: #6c757d;
            font-size: 0.9rem;
        }
        .footer a {
            color: #007bff;
            text-decoration: none;
        }
        .footer a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    		<h1>SLZB Exporter<span class="version">` + versionInfo.Version + `</span></h1>
    
    <div class="endpoint">
        <h3><a href="/metrics">📊 Metrics</a></h3>
        <p class="description">Prometheus metrics endpoint</p>
        <span class="status metrics">Available</span>
    </div>

    <div class="endpoint">
        <h3><a href="/metrics-info">📋 Metrics Info</a></h3>
        <p class="description">Detailed metrics information with examples</p>
        <span class="status metrics">Available</span>
    </div>

    <div class="endpoint">
        <h3><a href="/health">❤️ Health Check</a></h3>
        <p class="description">Service health status</p>
        <span class="status healthy">Healthy</span>
    </div>

    <div class="service-status">
        <h3>Service Status</h3>
        <p><strong>Status:</strong> <span class="status ready">Ready</span></p>
        		<p><strong>SLZB Connection:</strong> <span class="status connected">Connected</span></p>
        <p><strong>Metrics Collection:</strong> <span class="status ready">Active</span></p>
    </div>

    <div class="metrics-info">
        <h3>Version Information</h3>
        <ul>
            <li><strong>Version:</strong> ` + versionInfo.Version + `</li>
            <li><strong>Commit:</strong> ` + versionInfo.Commit + `</li>
            <li><strong>Build Date:</strong> ` + versionInfo.BuildDate + `</li>
        </ul>
    </div>

    <div class="metrics-info">
        <h3>Configuration</h3>
        <ul>
            			<li><strong>SLZB API URL:</strong> ` + s.config.SLZB.APIURL + `</li>
            <li><strong>Collection Interval:</strong> ` + s.config.SLZB.Interval.String() + `</li>
        </ul>
    </div>

    <div class="metrics-info">
        <h3>Available Metrics</h3>
        <ul>` + metricsHTML + `
        </ul>
    </div>

    <div class="footer">
        <p>Copyright © 2025 Dougal Matthews. Licensed under <a href="https://opensource.org/licenses/MIT" target="_blank">MIT License</a>.</p>
        <p><a href="https://github.com/d0ugal/slzb-exporter" target="_blank">GitHub Repository</a> | <a href="https://github.com/d0ugal/slzb-exporter/issues" target="_blank">Report Issues</a></p>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")

	if _, err := w.Write([]byte(html)); err != nil {
		slog.Error("Failed to write HTML response", "error", err)
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

func (s *Server) handleMetricsInfo(w http.ResponseWriter, r *http.Request) {
	metricsInfo := s.getMetricsInfo()

	// Convert to JSON
	jsonResponse := fmt.Sprintf(`{
		"metrics": [
			%s
		],
		"total_count": %d,
		"generated_at": %d
	}`, s.generateMetricsJSON(metricsInfo), len(metricsInfo), time.Now().Unix())

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write([]byte(jsonResponse)); err != nil {
		slog.Error("Failed to write metrics info response", "error", err)
	}
}

func (s *Server) generateMetricsJSON(metricsInfo []MetricInfo) string {
	var jsonParts []string

	for _, metric := range metricsInfo {
		labelsJSON := "{"
		var labelPairs []string
		for k, v := range metric.Labels {
			labelPairs = append(labelPairs, fmt.Sprintf(`"%s": "%s"`, k, v))
		}
		labelsJSON += strings.Join(labelPairs, ", ") + "}"

		metricJSON := fmt.Sprintf(`{
			"name": "%s",
			"help": "%s",
			"example_value": "%s",
			"labels": %s
		}`, metric.Name, metric.Help, metric.ExampleValue, labelsJSON)

		jsonParts = append(jsonParts, metricJSON)
	}

	return strings.Join(jsonParts, ",\n			")
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/metrics-info", s.handleMetricsInfo)

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
	var metricsInfo []MetricInfo

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
			Help:         s.getMetricHelp(metric.field),
			ExampleValue: s.getExampleValue(metric.field),
			Labels:       s.getExampleLabels(metric.field),
		})
	}

	return metricsInfo
}

func (s *Server) getExampleLabels(metricName string) map[string]string {
	switch metricName {
	case "SLZBConnected", "SLZBDeviceTemp", "SLZBUptime", "SLZBHeapFree", "SLZBHeapSize", "SLZBHeapRatio":
		return map[string]string{"device": "slzb-01"}
	case "SLZBSocketConnected":
		return map[string]string{"device": "slzb-01", "connections": "5"}
	case "SLZBDeviceMode":
		return map[string]string{"device": "slzb-01", "mode": "coordinator"}
	case "SLZBEthernetConnected":
		return map[string]string{
			"device":      "slzb-01",
			"ip_address":  "192.168.1.100",
			"mac_address": "00:11:22:33:44:55",
			"gateway":     "192.168.1.1",
			"subnet_mask": "255.255.255.0",
			"dns_server":  "8.8.8.8",
			"speed_mbps":  "1000",
		}
	case "SLZBWifiConnected":
		return map[string]string{
			"device":      "slzb-01",
			"ssid":        "MyWiFi",
			"ip_address":  "192.168.1.101",
			"mac_address": "00:11:22:33:44:56",
			"gateway":     "192.168.1.1",
			"subnet_mask": "255.255.255.0",
			"dns_server":  "8.8.8.8",
		}
	case "SLZBHTTPRequestsTotal", "SLZBHTTPErrorsTotal":
		return map[string]string{"device": "slzb-01", "action": "get_status", "status": "200"}
	case "SLZBDeviceReachable", "SLZBLastCollectionTime":
		return map[string]string{"device": "slzb-01"}
	case "SLZBCollectionErrors":
		return map[string]string{"device": "slzb-01", "error_type": "timeout"}
	case "SLZBFirmwareCurrentVersion":
		return map[string]string{"device": "slzb-01", "version": "1.0.0", "build_date": "2024-01-01"}
	case "SLZBFirmwareUpdateAvailable":
		return map[string]string{"device": "slzb-01", "available_version": "1.0.1"}
	case "SLZBFirmwareLastCheckTime":
		return map[string]string{"device": "slzb-01"}
	case "SLZBConfigBackupStatus":
		return map[string]string{"device": "slzb-01", "backup_type": "full"}
	case "SLZBConfigLastBackupTime":
		return map[string]string{"device": "slzb-01", "backup_type": "full"}
	case "SLZBConfigFileCount":
		return map[string]string{"device": "slzb-01", "file_type": "configuration"}
	case "SLZBAPIResponseTimeSeconds":
		return map[string]string{"device": "slzb-01", "action": "get_status"}
	case "SLZBAPITimeoutErrorsTotal":
		return map[string]string{"device": "slzb-01", "action": "get_status"}
	case "SLZBCollectionDurationSeconds":
		return map[string]string{"device": "slzb-01"}
	default:
		return map[string]string{"device": "slzb-01"}
	}
}

func (s *Server) getExampleValue(metricName string) string {
	switch metricName {
	case "SLZBConnected", "SLZBSocketConnected", "SLZBDeviceMode", "SLZBEthernetConnected", "SLZBWifiConnected", "SLZBDeviceReachable", "SLZBFirmwareUpdateAvailable", "SLZBConfigBackupStatus":
		return "1"
	case "SLZBDeviceTemp":
		return "45.2"
	case "SLZBUptime", "SLZBSocketUptime":
		return "86400"
	case "SLZBHeapFree":
		return "512"
	case "SLZBHeapSize":
		return "1024"
	case "SLZBHeapRatio":
		return "50.0"
	case "SLZBHTTPRequestsTotal", "SLZBHTTPErrorsTotal", "SLZBCollectionErrors", "SLZBAPITimeoutErrorsTotal":
		return "42"
	case "SLZBLastCollectionTime", "SLZBFirmwareLastCheckTime", "SLZBConfigLastBackupTime":
		return "1704067200"
	case "SLZBFirmwareCurrentVersion":
		return "1"
	case "SLZBConfigFileCount":
		return "5"
	case "SLZBAPIResponseTimeSeconds", "SLZBCollectionDurationSeconds":
		return "0.125"
	default:
		return "0"
	}
}

func (s *Server) getMetricHelp(metricName string) string {
	switch metricName {
	case "SLZBConnected":
		return "SLZB device connection status (1=connected, 0=disconnected)"
	case "SLZBDeviceTemp":
		return "SLZB device temperature in degrees Celsius"
	case "SLZBUptime":
		return "SLZB device uptime in seconds since last boot"
	case "SLZBSocketUptime":
		return "SLZB socket connection uptime in seconds since connection established"
	case "SLZBSocketConnected":
		return "SLZB socket connection status (1=connected, 0=disconnected)"
	case "SLZBDeviceMode":
		return "SLZB device operational mode (1=active, 0=inactive) with mode label"
	case "SLZBHeapFree":
		return "SLZB device free heap memory in kilobytes"
	case "SLZBHeapSize":
		return "SLZB device total heap memory in kilobytes"
	case "SLZBHeapRatio":
		return "SLZB device heap usage ratio as percentage (free heap / total heap * 100)"
	case "SLZBEthernetConnected":
		return "SLZB device Ethernet connection status (1=connected, 0=disconnected)"
	case "SLZBWifiConnected":
		return "SLZB device WiFi connection status (1=connected, 0=disconnected)"
	case "SLZBHTTPRequestsTotal":
		return "Total number of HTTP requests made by exporter to SLZB device API"
	case "SLZBHTTPErrorsTotal":
		return "Total number of HTTP errors when making requests to SLZB device API"
	case "SLZBDeviceReachable":
		return "SLZB device reachability status (1=reachable, 0=unreachable)"
	case "SLZBLastCollectionTime":
		return "Unix timestamp of the last successful collection from SLZB device"
	case "SLZBCollectionErrors":
		return "Total number of collection errors for SLZB device by error type"
	case "SLZBFirmwareCurrentVersion":
		return "Current firmware version (always 1, used for joining with labels)"
	case "SLZBFirmwareUpdateAvailable":
		return "Firmware update availability (1=available, 0=not_available)"
	case "SLZBFirmwareLastCheckTime":
		return "Unix timestamp of last firmware check"
	case "SLZBConfigBackupStatus":
		return "Status of the last configuration backup (1=success, 0=failure)"
	case "SLZBConfigLastBackupTime":
		return "Unix timestamp of the last successful configuration backup"
	case "SLZBConfigFileCount":
		return "Number of configuration files on the device"
	case "SLZBAPIResponseTimeSeconds":
		return "Histogram of API response times in seconds"
	case "SLZBAPITimeoutErrorsTotal":
		return "Total number of API timeout errors"
	case "SLZBCollectionDurationSeconds":
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

package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
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
        <ul>
            			<li><strong>slzb_device_connected:</strong> SLZB device connection status</li>
            <li><strong>slzb_device_temperature_celsius:</strong> Device temperature</li>
            <li><strong>slzb_device_uptime_seconds:</strong> Device uptime</li>
            <li><strong>slzb_device_heap_free_kb:</strong> Free heap memory</li>
            <li><strong>slzb_device_heap_size_kb:</strong> Total heap memory</li>
            <li><strong>slzb_device_heap_ratio:</strong> Heap usage ratio as percentage</li>
            <li><strong>slzb_device_ethernet_speed_mbps:</strong> Ethernet connection speed</li>
            <li><strong>slzb_device_ethernet_connected:</strong> Ethernet connection status with network details</li>
            <li><strong>slzb_device_wifi_connected:</strong> WiFi connection status with network details</li>
            <li><strong>slzb_http_requests_total:</strong> Total HTTP requests</li>
            <li><strong>slzb_http_errors_total:</strong> HTTP error tracking</li>
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

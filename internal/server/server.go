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

type MetricInfo struct {
	Name         string
	Help         string
	ExampleValue string
	Labels       []string
}

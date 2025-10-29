package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/d0ugal/promexporter/app"
	"github.com/d0ugal/promexporter/logging"
	promexporter_metrics "github.com/d0ugal/promexporter/metrics"
	"github.com/d0ugal/slzb-exporter/internal/collectors"
	"github.com/d0ugal/slzb-exporter/internal/config"
	"github.com/d0ugal/slzb-exporter/internal/metrics"
	"github.com/d0ugal/slzb-exporter/internal/version"
)

func main() {
	// Parse command line flags
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information")
	flag.Parse()

	// Show version if requested
	if showVersion {
		fmt.Printf("slzb-exporter %s\n", version.Version)
		fmt.Printf("Commit: %s\n", version.Commit)
		fmt.Printf("Build Date: %s\n", version.BuildDate)
		os.Exit(0)
	}

	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Configure logging using promexporter
	logging.Configure(&logging.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
	})

	// Initialize metrics registry using promexporter
	metricsRegistry := promexporter_metrics.NewRegistry("slzb_exporter_info")

	// Add custom metrics to the registry
	slzbRegistry := metrics.NewSLZBRegistry(metricsRegistry)

	// Create and build application using promexporter
	application := app.New("SLZB Exporter").
		WithConfig(&cfg.BaseConfig).
		WithMetrics(metricsRegistry).
		WithVersionInfo(version.Version, version.Commit, version.BuildDate).
		Build()

	// Create collector with app reference for tracing
	slzbCollector := collectors.NewSLZBCollector(cfg, slzbRegistry, application)
	application.WithCollector(slzbCollector)

	if err := application.Run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

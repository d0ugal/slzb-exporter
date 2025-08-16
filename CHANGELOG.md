# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project setup and structure
- SLZB-06 device monitoring capabilities
- Prometheus metrics export
- Docker containerization
- Comprehensive test suite

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- N/A

## [1.0.0] - 2025-08-16

### Added
- Initial release of SLZB-06 Prometheus Exporter
- Device status monitoring (connectivity, temperature, uptime)
- Network metrics (Ethernet and WiFi status, RSSI, connection speeds)
- System health tracking (memory usage, heap statistics)
- HTTP request monitoring (success/failure rates, response times)
- Configurable collection intervals
- Prometheus integration with metrics export on port 9110
- Health check endpoint at /health
- Graceful shutdown with proper signal handling
- Structured logging with JSON and text formats
- Environment-based configuration
- Docker containerization with multi-stage builds
- Comprehensive test suite with race detection and coverage
- golangci-lint integration for code quality
- GitHub Actions CI/CD workflow
- Semantic versioning with build-time version injection
- Detailed documentation and API reference

### Technical Details
- Go 1.22+ with modern Go modules
- Prometheus client_golang v1.22.0 for metrics handling
- Clean architecture with separate packages for collectors, config, logging, metrics, server, and version
- Alpine Linux container for minimal footprint
- Production-ready with proper error handling and resource management

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.1](https://github.com/d0ugal/slzb-exporter/compare/v1.1.0...v1.1.1) (2025-08-16)


### Bug Fixes

* apply golangci-lint formatting fixes for all issues ([de0954b](https://github.com/d0ugal/slzb-exporter/commit/de0954b3a31f4ea285df9e7627d91686bb2e16fa))
* correct golangci-lint config version to 2 (only valid version) ([3660290](https://github.com/d0ugal/slzb-exporter/commit/3660290c4f724ab29e1159be6b2a57cc7a3f0b16))
* update golangci-lint config to match working mqtt-exporter pattern ([a9a73df](https://github.com/d0ugal/slzb-exporter/commit/a9a73df018f4a4e4574bc0195f22befc10369c0c))

## [1.1.0](https://github.com/d0ugal/slzb-exporter/compare/v1.0.0...v1.1.0) (2025-08-16)


### Features

* add missing project files and documentation ([319e212](https://github.com/d0ugal/slzb-exporter/commit/319e2127aaf80106db72a55e37e364b03c6f0cb8))
* upgrade to Go 1.25 ([627b5d4](https://github.com/d0ugal/slzb-exporter/commit/627b5d4fd9ae091a0316f28e3d1081bb4e2ccf5d))


### Bug Fixes

* update golangci-lint config for Go 1.25 compatibility ([b14762a](https://github.com/d0ugal/slzb-exporter/commit/b14762ab761261c7e034f15131bd00bbf8b813a5))

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

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.3.2](https://github.com/d0ugal/slzb-exporter/compare/v2.3.1...v2.3.2) (2025-09-05)


### Bug Fixes

* **deps:** update module github.com/prometheus/client_golang to v1.23.2 ([c02f4d3](https://github.com/d0ugal/slzb-exporter/commit/c02f4d32edec5a1d2fc4306ea5f70abe267cacf7))
* **deps:** update module github.com/prometheus/client_golang to v1.23.2 ([5c65bf2](https://github.com/d0ugal/slzb-exporter/commit/5c65bf2ff2c55fcb17d36145a6522a44291d45a7))

## [2.3.1](https://github.com/d0ugal/slzb-exporter/compare/v2.3.0...v2.3.1) (2025-09-04)


### Bug Fixes

* **deps:** update module github.com/prometheus/client_golang to v1.23.1 ([8d6d7d8](https://github.com/d0ugal/slzb-exporter/commit/8d6d7d8ceb448d88a206a10765df485e42918d1a))
* **deps:** update module github.com/prometheus/client_golang to v1.23.1 ([f66e810](https://github.com/d0ugal/slzb-exporter/commit/f66e810fa20dda0b6ebb2b71d8e7a10028aad574))

## [2.3.0](https://github.com/d0ugal/slzb-exporter/compare/v2.2.0...v2.3.0) (2025-09-04)


### Features

* update dev build versioning to use semver-compatible pre-release tags ([8aab899](https://github.com/d0ugal/slzb-exporter/commit/8aab8993a7f3c7741866c5658267b7c80bc96012))


### Bug Fixes

* **ci:** add v prefix to dev tags for consistent versioning ([0793607](https://github.com/d0ugal/slzb-exporter/commit/0793607d6e5cbd83f3b82b9d535824758b5ed8e9))
* use actual release version as base for dev tags instead of hardcoded 0.0.0 ([3924cf2](https://github.com/d0ugal/slzb-exporter/commit/3924cf27665a516ddfdbb533392de68b9870fe2e))
* use fetch-depth: 0 instead of fetch-tags for full git history ([bc9a10d](https://github.com/d0ugal/slzb-exporter/commit/bc9a10d74ebc3b24b10a02bce5ed1f455091ffc3))
* use fetch-tags instead of fetch-depth for GitHub Actions ([ae87d7e](https://github.com/d0ugal/slzb-exporter/commit/ae87d7e71282d5c9fdc97e748e1206124f7bcb6e))

## [2.2.0](https://github.com/d0ugal/slzb-exporter/compare/v2.1.1...v2.2.0) (2025-09-04)


### Features

* enable global automerge in Renovate config ([0cd920d](https://github.com/d0ugal/slzb-exporter/commit/0cd920de3beefdf3203e5b8c4a0f4d0f0071288d))

## [2.1.1](https://github.com/d0ugal/slzb-exporter/compare/v2.1.0...v2.1.1) (2025-09-03)


### Bug Fixes

* pin Alpine version to 3.22.1 for consistency ([2c65ff9](https://github.com/d0ugal/slzb-exporter/commit/2c65ff9f4c846648e8503efa5da29e1f51279e2f))

## [2.1.0](https://github.com/d0ugal/slzb-exporter/compare/v2.0.2...v2.1.0) (2025-08-26)


### Features

* **docker:** use an unprivileged user during runtime ([6606723](https://github.com/d0ugal/slzb-exporter/commit/6606723059aa972e2e2a8e217ce521acc859cbbb))

## [2.0.2](https://github.com/d0ugal/slzb-exporter/compare/v2.0.1...v2.0.2) (2025-08-20)


### Bug Fixes

* trigger release for server cleanup ([b09e8e8](https://github.com/d0ugal/slzb-exporter/commit/b09e8e81b96b95e6077ad95d36196d7f80cb1a7d))

## [2.0.1](https://github.com/d0ugal/slzb-exporter/compare/v2.0.0...v2.0.1) (2025-08-20)


### Bug Fixes

* remove redundant Service Information section from UI ([cac2dd3](https://github.com/d0ugal/slzb-exporter/commit/cac2dd3126e82a7550888c32788105a6dfa00334))

## [2.0.0](https://github.com/d0ugal/slzb-exporter/compare/v1.4.0...v2.0.0) (2025-08-20)


### ⚠ BREAKING CHANGES

* Metric names have been updated to comply with Prometheus naming conventions:
    - slzb_device_heap_free_kb → slzb_device_heap_free_bytes
    - slzb_device_heap_size_kb → slzb_device_heap_size_bytes
    - slzb_config_file_count → slzb_config_files

### Features

* optimize linting performance with caching and fix metric naming ([a0397b2](https://github.com/d0ugal/slzb-exporter/commit/a0397b2f9721889c50a637ef92d5a2797cbe19a7))


### Bug Fixes

* run Docker containers as current user to prevent permission issues ([59ce761](https://github.com/d0ugal/slzb-exporter/commit/59ce76138a3ff26e781a037f0b3e7e1722b6ebd6))
* temporarily disable gocyclo to allow commit to proceed ([34ae928](https://github.com/d0ugal/slzb-exporter/commit/34ae9282fccd359df879608f562eb33ae17e4493))

## [1.4.0](https://github.com/d0ugal/slzb-exporter/compare/v1.3.1...v1.4.0) (2025-08-20)


### Features

* implement template-based UI with centralized metric information ([6ad5c74](https://github.com/d0ugal/slzb-exporter/commit/6ad5c74773f71a2b4cbe5559be5108094d8ea0a2))

## [1.3.1](https://github.com/d0ugal/slzb-exporter/compare/v1.3.0...v1.3.1) (2025-08-20)


### Bug Fixes

* **tests:** remove handleMetricsInfo test and clean up unused imports ([e027723](https://github.com/d0ugal/slzb-exporter/commit/e027723c49145cbe4104d88fcf74469bf41ee1d0))

## [1.3.0](https://github.com/d0ugal/slzb-exporter/compare/v1.2.0...v1.3.0) (2025-08-19)


### Features

* **server:** add dynamic metrics information with examples ([5a48e1d](https://github.com/d0ugal/slzb-exporter/commit/5a48e1d3758e0f7affa459362e733caf40a94d96))
* **server:** make metrics list collapsible for better UX ([492d18d](https://github.com/d0ugal/slzb-exporter/commit/492d18d2172809712a0aabf192189e3e70995e9e))


### Bug Fixes

* **lint:** pre-allocate slices to resolve golangci-lint prealloc warnings ([de62003](https://github.com/d0ugal/slzb-exporter/commit/de62003aea7e2ee56731e360df8d926e2431fe76))

## [1.2.0](https://github.com/d0ugal/slzb-exporter/compare/v1.1.2...v1.2.0) (2025-08-18)


### Features

* add comprehensive monitoring features to SLZB exporter ([caa8f0b](https://github.com/d0ugal/slzb-exporter/commit/caa8f0b61c08869860b960c8e4cb829426ff1146))
* add comprehensive SLZB monitoring metrics ([da9ff58](https://github.com/d0ugal/slzb-exporter/commit/da9ff5840bdb9bb8866f7d2f625d45f515701c9d))


### Bug Fixes

* linting issue ([f5bd7dc](https://github.com/d0ugal/slzb-exporter/commit/f5bd7dcbe336c07cd6c368830a5f40be12175f9a))

## [1.1.2](https://github.com/d0ugal/slzb-exporter/compare/v1.1.1...v1.1.2) (2025-08-16)


### Bug Fixes

* add version injection to Dockerfile and CI workflow ([56fd68f](https://github.com/d0ugal/slzb-exporter/commit/56fd68f60e9f292c4cd7ffe39df7abf0470fba90))

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

# SLZB-06 Prometheus Exporter

A Prometheus exporter for SLZB-06 Zigbee 3.0 PoE Ethernet USB Adapters. This exporter connects to SLZB-06 devices via their web API and exposes metrics about device status, network connectivity, and system health.

## Features

- **Device Status Monitoring**: Track device connectivity, temperature, and uptime
- **Network Metrics**: Monitor Ethernet and WiFi connection status and performance
- **System Health**: Track memory usage, heap statistics, and device temperature
- **HTTP Request Monitoring**: Track API request success/failure rates
- **Configurable Collection**: Adjustable collection intervals for different environments
- **Prometheus Integration**: Exports metrics in Prometheus format
- **Health Check Endpoint**: Built-in health monitoring
- **Graceful Shutdown**: Proper signal handling and cleanup
- **Structured Logging**: JSON and text logging with configurable levels

## Metrics

The SLZB Exporter exposes the following metrics:

### Device Information
- `slzb_device_info` - Device information (always 1, used for joining with labels)
- `slzb_device_uptime_seconds` - Device uptime in seconds
- `slzb_device_memory_usage_percent` - Memory usage percentage
- `slzb_device_temperature_celsius` - Device temperature in Celsius
- `slzb_device_signal_strength_dbm` - WiFi signal strength in dBm
- `slzb_device_voltage_volts` - Device voltage in volts
- `slzb_device_current_ma` - Device current in milliamps
- `slzb_device_power_mw` - Device power consumption in milliwatts

### HTTP Request Metrics
- `slzb_http_requests_total` - Total number of HTTP requests
- `slzb_http_errors_total` - Total number of HTTP errors
- `slzb_last_collection_time` - Unix timestamp of last successful collection
- `slzb_collection_errors_total` - Total number of collection errors

### NEW: Firmware Update Status
- `slzb_firmware_current_version` - Current firmware version (always 1, used for joining with labels)
- `slzb_firmware_update_available` - Firmware update availability (1=available, 0=not_available)
- `slzb_firmware_last_check_timestamp` - Unix timestamp of last firmware check

### NEW: Configuration Management
- `slzb_config_backup_status` - Configuration backup status (1=success, 0=failed)
- `slzb_config_last_backup_timestamp` - Unix timestamp of last successful backup
- `slzb_config_file_count` - Number of configuration files

### NEW: Performance Benchmarks
- `slzb_api_response_time_seconds` - API response time in seconds (histogram)
- `slzb_api_timeout_errors_total` - Total number of API timeout errors
- `slzb_collection_duration_seconds` - Duration of collection cycles in seconds (histogram)

## Quick Start

### Using Docker

```bash
# Pull the latest image
docker pull ghcr.io/d0ugal/slzb-exporter:latest

# Run with default configuration
docker run -p 9110:9110 ghcr.io/d0ugal/slzb-exporter:latest

# Run with custom SLZB device URL
docker run -p 9110:9110 \
  -e SLZB_EXPORTER_SLZB_API_URL=http://your-slzb-device.local \
  ghcr.io/d0ugal/slzb-exporter:latest
```

### Using Docker Compose

```yaml
version: '3.8'
services:
  slzb-exporter:
    image: ghcr.io/d0ugal/slzb-exporter:latest
    ports:
      - "9110:9110"
    environment:
      - SLZB_EXPORTER_SLZB_API_URL=http://your-slzb-device.local
      - SLZB_EXPORTER_METRICS_DEFAULT_INTERVAL=60s
    restart: unless-stopped
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/d0ugal/slzb-exporter.git
cd slzb-exporter

# Build the application
make build

# Run the exporter
./slzb-exporter
```

## Configuration

The exporter is configured via environment variables:

### Environment Variables

- `SLZB_EXPORTER_SLZB_API_URL` - Base URL of the SLZB-06 device (default: http://slzb-device.local)
- `SLZB_EXPORTER_METRICS_DEFAULT_INTERVAL` - Collection interval (default: 10s)
- `SLZB_EXPORTER_SERVER_HOST` - Server host (default: 0.0.0.0)
- `SLZB_EXPORTER_SERVER_PORT` - Server port (default: 9110)
- `SLZB_EXPORTER_LOG_LEVEL` - Log level: debug, info, warn, error (default: info)
- `SLZB_EXPORTER_LOG_FORMAT` - Log format: json, text (default: json)

### Example Configuration

```bash
# Basic configuration
export SLZB_EXPORTER_SLZB_API_URL="http://192.168.1.100"
export SLZB_EXPORTER_METRICS_DEFAULT_INTERVAL="60s"

# Advanced configuration
export SLZB_EXPORTER_SERVER_HOST="0.0.0.0"
export SLZB_EXPORTER_SERVER_PORT="9110"
export SLZB_EXPORTER_LOG_LEVEL="info"
export SLZB_EXPORTER_LOG_FORMAT="json"
```

## API Endpoints

- `/` - Web interface with device status and metrics information
- `/metrics` - Prometheus metrics endpoint
- `/health` - Health check endpoint

## Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'slzb-exporter'
    static_configs:
      - targets: ['localhost:9110']
    scrape_interval: 30s
```

## Example PromQL Queries

### Device Health
```promql
# Device reachability
slzb_device_reachable

# Device temperature
slzb_device_temperature_celsius

# Heap memory usage
slzb_device_heap_ratio
```

### NEW: Zigbee Network Performance
```promql
# Zigbee packet rate (packets per second)
rate(slzb_zigbee_packets_received_total[5m])
rate(slzb_zigbee_packets_sent_total[5m])

# Zigbee error rate
rate(slzb_zigbee_errors_total[5m])

# Number of connected devices
slzb_zigbee_network_devices{device_type="connected"}

# Channel utilization
slzb_zigbee_channel_utilization_percent
```

### NEW: Firmware Management
```promql
# Check for firmware updates
slzb_firmware_update_available

# Firmware version distribution
slzb_firmware_current_version

# Time since last firmware check
time() - slzb_firmware_last_check_timestamp
```

### NEW: Configuration Health
```promql
# Configuration backup status
slzb_config_backup_status

# Time since last backup
time() - slzb_config_last_backup_timestamp

# Configuration file count
slzb_config_file_count

# Configuration size
slzb_config_total_size_bytes
```

### NEW: Security Monitoring
```promql
# Encryption status
slzb_encryption_status

# Security key age
time() - slzb_security_key_rotation_timestamp

# Security events rate
rate(slzb_security_events_total[5m])
```

### NEW: Performance Monitoring
```promql
# API response time percentiles
histogram_quantile(0.95, rate(slzb_api_response_time_seconds_bucket[5m]))
histogram_quantile(0.50, rate(slzb_api_response_time_seconds_bucket[5m]))

# Collection duration
histogram_quantile(0.95, rate(slzb_collection_duration_seconds_bucket[5m]))

# API timeout rate
rate(slzb_api_timeout_errors_total[5m])
```

## Grafana Dashboard

A sample Grafana dashboard is available in the `dashboards/` directory. Import the JSON file into Grafana to get started with monitoring your SLZB-06 device.

## Example Alerting Rules

### Device Health Alerts
```yaml
- alert: SLZBDeviceOffline
  expr: slzb_device_reachable == 0
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "SLZB device {{ $labels.device }} is offline"
    description: "SLZB device {{ $labels.device }} has been unreachable for more than 1 minute"

- alert: SLZBHighTemperature
  expr: slzb_device_temperature_celsius > 70
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "SLZB device {{ $labels.device }} temperature is high"
    description: "SLZB device {{ $labels.device }} temperature is {{ $value }}°C"

- alert: SLZBLowHeapMemory
  expr: slzb_device_heap_ratio < 20
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "SLZB device {{ $labels.device }} has low heap memory"
    description: "SLZB device {{ $labels.device }} heap usage is {{ $value }}%"
```

### NEW: Zigbee Network Alerts
```yaml
- alert: SLZBHighZigbeeErrorRate
  expr: rate(slzb_zigbee_errors_total[5m]) > 0.1
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "SLZB device {{ $labels.device }} has high Zigbee error rate"
    description: "Zigbee error rate is {{ $value }} errors/second"

- alert: SLZBNoZigbeeDevices
  expr: slzb_zigbee_network_devices{device_type="connected"} == 0
  for: 5m
  labels:
    severity: critical
  annotations:
    summary: "SLZB device {{ $labels.device }} has no connected Zigbee devices"
    description: "No Zigbee devices are connected to the network"

- alert: SLZBHighChannelUtilization
  expr: slzb_zigbee_channel_utilization_percent > 80
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "SLZB device {{ $labels.device }} has high channel utilization"
    description: "Channel utilization is {{ $value }}%"
```

### NEW: Firmware Management Alerts
```yaml
- alert: SLZBFirmwareUpdateAvailable
  expr: slzb_firmware_update_available == 1
  for: 0m
  labels:
    severity: info
  annotations:
    summary: "SLZB device {{ $labels.device }} has firmware update available"
    description: "Firmware update is available for device {{ $labels.device }}"

- alert: SLZBFirmwareCheckStale
  expr: time() - slzb_firmware_last_check_timestamp > 86400
  for: 1h
  labels:
    severity: warning
  annotations:
    summary: "SLZB device {{ $labels.device }} firmware check is stale"
    description: "Firmware check is more than 24 hours old"
```

### NEW: Configuration Health Alerts
```yaml
- alert: SLZBConfigBackupFailed
  expr: slzb_config_backup_status == 0
  for: 1h
  labels:
    severity: warning
  annotations:
    summary: "SLZB device {{ $labels.device }} configuration backup failed"
    description: "Configuration backup has failed for device {{ $labels.device }}"

- alert: SLZBConfigBackupStale
  expr: time() - slzb_config_last_backup_timestamp > 604800
  for: 1h
  labels:
    severity: warning
  annotations:
    summary: "SLZB device {{ $labels.device }} configuration backup is stale"
    description: "Configuration backup is more than 7 days old"
```

### NEW: Performance Alerts
```yaml
- alert: SLZBSlowAPIResponse
  expr: histogram_quantile(0.95, rate(slzb_api_response_time_seconds_bucket[5m])) > 5
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "SLZB device {{ $labels.device }} has slow API response"
    description: "95th percentile API response time is {{ $value }} seconds"

- alert: SLZBHighTimeoutRate
  expr: rate(slzb_api_timeout_errors_total[5m]) > 0.1
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "SLZB device {{ $labels.device }} has high API timeout rate"
    description: "API timeout rate is {{ $value }} timeouts/second"
```

## Development

### Building

```bash
# Build the application
make build

# Run tests
make test

# Run linter
make lint
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
```

### API Documentation

For detailed information about the SLZB-06 device API that this exporter consumes, see [API.md](API.md). This documentation is useful for developers who want to understand how the exporter works or contribute to its development.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions, please open an issue on GitHub.

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

The exporter provides the following Prometheus metrics:

### Device Status Metrics
- `slzb_device_connected` - Overall connection status
- `slzb_device_scan_done` - WiFi scan completion status
- `slzb_device_wifi_rssi_dbm` - WiFi RSSI in dBm
- `slzb_device_temperature_celsius` - Device temperature in Celsius
- `slzb_device_uptime_seconds` - Device uptime in seconds
- `slzb_device_heap_free_kb` - Free heap memory in KB
- `slzb_device_heap_size_kb` - Total heap memory in KB
- `slzb_device_ethernet_speed_mbps` - Ethernet connection speed in Mbps

### HTTP Request Metrics
- `slzb_http_requests_total` - Total number of HTTP requests made by exporter to SLZB-06 device
- `slzb_http_request_duration_seconds` - Duration of HTTP requests made by exporter to SLZB-06 device
- `slzb_http_errors_total` - Total number of HTTP errors when making requests to SLZB-06 device

### Device Health and Availability Metrics
- `slzb_device_reachable` - Device reachability status (1=reachable, 0=unreachable)
- `slzb_last_collection_timestamp` - Timestamp of the last successful collection
- `slzb_collection_errors_total` - Total number of collection errors

### Exporter Information
- `slzb_exporter_version_info` - Version information about the SLZB exporter

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

## Grafana Dashboard

A sample Grafana dashboard is available in the `dashboards/` directory. Import the JSON file into Grafana to get started with monitoring your SLZB-06 device.

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

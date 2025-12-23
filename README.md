# SLZB-06 Prometheus Exporter

A Prometheus exporter for SLZB-06 Zigbee 3.0 PoE Ethernet USB Adapters.

**Image**: `ghcr.io/d0ugal/slzb-exporter:v2.15.11`

## Metrics

### Device Information
- `slzb_device_connected` - Device connection status (1=connected, 0=disconnected)
- `slzb_device_temperature_celsius` - Device temperature in Celsius
- `slzb_device_uptime_seconds` - Device uptime in seconds
- `slzb_device_heap_free_kb` - Free heap memory in kilobytes
- `slzb_device_heap_size_kb` - Total heap memory in kilobytes
- `slzb_device_heap_ratio` - Heap usage ratio as percentage
- `slzb_device_ethernet_connected` - Ethernet connection status
- `slzb_device_wifi_connected` - WiFi connection status
- `slzb_device_operational_mode` - Device operational mode
- `slzb_device_reachable` - Device reachability status

### HTTP Request Metrics
- `slzb_http_requests_total` - Total number of HTTP requests
- `slzb_http_errors_total` - Total number of HTTP errors
- `slzb_last_collection_timestamp` - Unix timestamp of last successful collection
- `slzb_collection_errors_total` - Total number of collection errors

### Endpoints
- `GET /`: HTML dashboard with device status and metrics information
- `GET /metrics`: Prometheus metrics endpoint
- `GET /health`: Health check endpoint

## Quick Start

### Docker Compose

```yaml
version: '3.8'
services:
  slzb-exporter:
    image: ghcr.io/d0ugal/slzb-exporter:v2.15.11
    ports:
      - "9110:9110"
    environment:
      - SLZB_EXPORTER_SLZB_API_URL=http://your-slzb-device.local
    restart: unless-stopped
```

1. Update the SLZB device URL in the environment variables
2. Run: `docker-compose up -d`
3. Access metrics: `curl http://localhost:9110/metrics`

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

## Deployment

### Docker Compose (Environment Variables)

```yaml
version: '3.8'
services:
  slzb-exporter:
    image: ghcr.io/d0ugal/slzb-exporter:v2.15.11
    ports:
      - "9110:9110"
    environment:
      - SLZB_EXPORTER_SLZB_API_URL=http://192.168.1.100
      - SLZB_EXPORTER_METRICS_DEFAULT_INTERVAL=60s
    restart: unless-stopped
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: slzb-exporter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: slzb-exporter
  template:
    metadata:
      labels:
        app: slzb-exporter
    spec:
      containers:
      - name: slzb-exporter
        image: ghcr.io/d0ugal/slzb-exporter:v2.15.11
        ports:
        - containerPort: 9110
        env:
        - name: SLZB_EXPORTER_SLZB_API_URL
          value: "http://192.168.1.100"
        - name: SLZB_EXPORTER_METRICS_DEFAULT_INTERVAL
          value: "60s"
```

## Prometheus Integration

Add to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'slzb-exporter'
    static_configs:
      - targets: ['slzb-exporter:9110']
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

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Linting

```bash
make lint
```

## License

This project is licensed under the MIT License.
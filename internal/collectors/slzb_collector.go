package collectors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/d0ugal/promexporter/app"
	"github.com/d0ugal/promexporter/tracing"
	"github.com/d0ugal/slzb-exporter/internal/config"
	"github.com/d0ugal/slzb-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
)

// SLZBCollector collects metrics from SLZB devices
type SLZBCollector struct {
	config  *config.Config
	metrics *metrics.SLZBRegistry
	app     *app.App
	client  *http.Client

	deviceInfo map[string]string // Cache for device information
	deviceID   string            // Derived device identifier
}

// NewSLZBCollector creates a new SLZB collector
func NewSLZBCollector(cfg *config.Config, metricsRegistry *metrics.SLZBRegistry, app *app.App) *SLZBCollector {
	// Derive device ID from API URL
	deviceID := deriveDeviceID(cfg.SLZB.APIURL)

	return &SLZBCollector{
		config:  cfg,
		metrics: metricsRegistry,
		app:     app,
		client: &http.Client{
			Timeout: 15 * time.Second, // Longer timeout for low-power device
		},
		deviceInfo: make(map[string]string),
		deviceID:   deviceID,
	}
}

// deriveDeviceID creates a device identifier from the API URL
func deriveDeviceID(apiURL string) string {
	// Extract hostname/IP from URL
	if strings.HasPrefix(apiURL, "http://") {
		apiURL = strings.TrimPrefix(apiURL, "http://")
	} else if strings.HasPrefix(apiURL, "https://") {
		apiURL = strings.TrimPrefix(apiURL, "https://")
	}

	// Remove port if present
	if idx := strings.Index(apiURL, ":"); idx != -1 {
		apiURL = apiURL[:idx]
	}

	// Use the hostname/IP as device ID
	return apiURL
}

// Start begins collecting metrics
func (sc *SLZBCollector) Start(ctx context.Context) {
	go sc.collectLoop(ctx)
}

// collectLoop runs the main collection loop
func (sc *SLZBCollector) collectLoop(ctx context.Context) {
	ticker := time.NewTicker(sc.config.SLZB.Interval)
	defer ticker.Stop()

	// Collect immediately on start
	sc.collectMetrics(ctx)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping SLZB collector")
			return
		case <-ticker.C:
			sc.collectMetrics(ctx)
		}
	}
}

// collectMetrics collects all metrics from the SLZB device
func (sc *SLZBCollector) collectMetrics(ctx context.Context) {
	deviceID := sc.deviceID
	collectionStart := time.Now()
	successfulCollections := 0
	totalCollections := 0

	// Create span for collection cycle
	tracer := sc.app.GetTracer()

	var collectorSpan *tracing.CollectorSpan
	var spanCtx context.Context

	if tracer != nil && tracer.IsEnabled() {
		collectorSpan = tracer.NewCollectorSpan(ctx, "slzb-collector", "collect-metrics")

		collectorSpan.SetAttributes(
			attribute.String("device.id", deviceID),
			attribute.String("device.api_url", sc.config.SLZB.APIURL),
		)
		spanCtx = collectorSpan.Context()
		defer collectorSpan.End()
	} else {
		spanCtx = ctx
	}

	if collectorSpan != nil {
		collectorSpan.AddEvent("collection_started",
			attribute.String("device.id", deviceID),
		)
	}

	// Track collection success
	defer func() {
		// Update last collection timestamp if we had any successful collections
		if successfulCollections > 0 {
			sc.metrics.SLZBLastCollectionTime.With(prometheus.Labels{
				"device": deviceID,
			}).Set(float64(time.Now().Unix()))
		}

		// Record collection duration
		duration := time.Since(collectionStart).Seconds()
		sc.metrics.SLZBCollectionDurationSeconds.With(prometheus.Labels{
			"device": deviceID,
		}).Observe(duration)

		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("collection.duration_seconds", duration),
				attribute.Int("collection.successful", successfulCollections),
				attribute.Int("collection.total", totalCollections),
			)
			collectorSpan.AddEvent("collection_completed",
				attribute.Int("successful", successfulCollections),
				attribute.Int("total", totalCollections),
				attribute.Float64("duration_seconds", duration),
			)
		}

		// Log collection summary
		slog.Info("Collection cycle completed",
			"device", deviceID,
			"successful", successfulCollections,
			"total", totalCollections,
			"duration", duration)
	}()

	// Get device information and test reachability in one request
	deviceInfoStart := time.Now()
	deviceReachable := sc.collectDeviceInfo(spanCtx, deviceID)
	deviceInfoDuration := time.Since(deviceInfoStart).Seconds()

	if deviceReachable {
		sc.metrics.SLZBDeviceReachable.With(prometheus.Labels{
			"device": deviceID,
		}).Set(1)
		sc.metrics.SLZBConnected.With(prometheus.Labels{
			"device": deviceID,
		}).Set(1)

		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("device_info.duration_seconds", deviceInfoDuration),
				attribute.Bool("device.reachable", true),
			)
			collectorSpan.AddEvent("device_reachable",
				attribute.Float64("duration_seconds", deviceInfoDuration),
			)
		}

		successfulCollections++
	} else {
		sc.metrics.SLZBDeviceReachable.With(prometheus.Labels{
			"device": deviceID,
		}).Set(0)
		sc.metrics.SLZBConnected.With(prometheus.Labels{
			"device": deviceID,
		}).Set(0)
		sc.metrics.SLZBCollectionErrors.With(prometheus.Labels{
			"device":     deviceID,
			"error_type": "device_unreachable",
		}).Inc()

		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("device_info.duration_seconds", deviceInfoDuration),
				attribute.Bool("device.reachable", false),
			)
			collectorSpan.RecordError(fmt.Errorf("device unreachable"), attribute.String("device.id", deviceID))
			collectorSpan.AddEvent("device_unreachable",
				attribute.Float64("duration_seconds", deviceInfoDuration),
			)
		}

		slog.Error("Device unreachable", "device", deviceID)

		return
	}

	totalCollections++

	// Collect device information
	deviceInfo2Start := time.Now()
	if sc.collectDeviceInfo(spanCtx, deviceID) {
		deviceInfo2Duration := time.Since(deviceInfo2Start).Seconds()
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("device_info_2.duration_seconds", deviceInfo2Duration),
			)
			collectorSpan.AddEvent("device_info_collected",
				attribute.Float64("duration_seconds", deviceInfo2Duration),
			)
		}
		successfulCollections++
	} else {
		deviceInfo2Duration := time.Since(deviceInfo2Start).Seconds()
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("device_info_2.duration_seconds", deviceInfo2Duration),
			)
			collectorSpan.RecordError(fmt.Errorf("device info collection failed"), attribute.String("device.id", deviceID))
		}
	}

	totalCollections++

	// NEW: Collect firmware update status
	firmwareStart := time.Now()
	if sc.collectFirmwareStatus(spanCtx, deviceID) {
		firmwareDuration := time.Since(firmwareStart).Seconds()
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("firmware.duration_seconds", firmwareDuration),
			)
			collectorSpan.AddEvent("firmware_status_collected",
				attribute.Float64("duration_seconds", firmwareDuration),
			)
		}
		successfulCollections++
	} else {
		firmwareDuration := time.Since(firmwareStart).Seconds()
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("firmware.duration_seconds", firmwareDuration),
			)
			collectorSpan.RecordError(fmt.Errorf("firmware status collection failed"), attribute.String("device.id", deviceID))
		}
	}

	totalCollections++

	// NEW: Collect configuration management metrics
	configStart := time.Now()
	if sc.collectConfigurationMetrics(spanCtx, deviceID) {
		configDuration := time.Since(configStart).Seconds()
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("configuration.duration_seconds", configDuration),
			)
			collectorSpan.AddEvent("configuration_metrics_collected",
				attribute.Float64("duration_seconds", configDuration),
			)
		}
		successfulCollections++
	} else {
		configDuration := time.Since(configStart).Seconds()
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("configuration.duration_seconds", configDuration),
			)
			collectorSpan.RecordError(fmt.Errorf("configuration metrics collection failed"), attribute.String("device.id", deviceID))
		}
	}

	totalCollections++

	// Add a small delay between requests
	time.Sleep(500 * time.Millisecond)
}

// collectDeviceInfo collects device information and caches it
func (sc *SLZBCollector) collectDeviceInfo(ctx context.Context, deviceName string) bool {
	tracer := sc.app.GetTracer()

	var collectorSpan *tracing.CollectorSpan
	var spanCtx context.Context

	if tracer != nil && tracer.IsEnabled() {
		collectorSpan = tracer.NewCollectorSpan(ctx, "slzb-collector", "collect-device-info")
		collectorSpan.SetAttributes(
			attribute.String("device.id", deviceName),
			attribute.String("device.action", "0"),
		)
		spanCtx = collectorSpan.Context()
		defer collectorSpan.End()
	} else {
		spanCtx = ctx
	}

	startTime := time.Now()

	// Get device information from action 0
	apiStart := time.Now()
	resp, err := sc.client.Get(fmt.Sprintf("%s/api?action=0&page=0", sc.config.SLZB.APIURL))
	apiDuration := time.Since(apiStart).Seconds()

	if err != nil {
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("api_request.duration_seconds", apiDuration),
			)
			collectorSpan.RecordError(err, attribute.String("operation", "api-request"))
		}
		sc.handleDeviceInfoRequestError(deviceName, err)
		return false
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("Failed to close response body", "error", err)
		}
	}()

	// Record API response time and update metrics
	responseTime := time.Since(startTime).Seconds()
	sc.updateDeviceInfoMetrics(deviceName, resp, responseTime)

	if resp.StatusCode != http.StatusOK {
		sc.handleDeviceInfoHTTPError(deviceName, resp.StatusCode)
		return false
	}

	// Parse and process device data
	respValuesArr := resp.Header.Get("respValuesArr")
	if respValuesArr != "" {
		processStart := time.Now()
		result := sc.processDeviceData(spanCtx, deviceName, respValuesArr)
		processDuration := time.Since(processStart).Seconds()
		totalDuration := time.Since(startTime).Seconds()

		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("api_request.duration_seconds", apiDuration),
				attribute.Float64("process_data.duration_seconds", processDuration),
				attribute.Float64("total.duration_seconds", totalDuration),
				attribute.Bool("processing.success", result),
			)
			if result {
				collectorSpan.AddEvent("device_data_processed",
					attribute.Float64("total_duration_seconds", totalDuration),
				)
			} else {
				collectorSpan.AddEvent("device_data_processing_failed",
					attribute.Float64("total_duration_seconds", totalDuration),
				)
			}
		}

		return result
	}

	// Fallback to default values if header is not available
	sc.setDefaultDeviceInfo(deviceName)
	totalDuration := time.Since(startTime).Seconds()

	if collectorSpan != nil {
		collectorSpan.SetAttributes(
			attribute.Float64("api_request.duration_seconds", apiDuration),
			attribute.Float64("total.duration_seconds", totalDuration),
		)
		collectorSpan.AddEvent("device_info_defaults_set",
			attribute.Float64("total_duration_seconds", totalDuration),
		)
	}

	return true
}

// handleDeviceInfoRequestError handles request errors for device info collection
func (sc *SLZBCollector) handleDeviceInfoRequestError(deviceName string, err error) {
	sc.metrics.SLZBHTTPErrorsTotal.With(prometheus.Labels{
		"device":     deviceName,
		"action":     "0",
		"error_type": "request_error",
	}).Inc()
	sc.metrics.SLZBAPITimeoutErrorsTotal.With(prometheus.Labels{
		"device": deviceName,
		"action": "0",
	}).Inc()
	slog.Error("Failed to get device info", "error", err)
}

// updateDeviceInfoMetrics updates API metrics for device info collection
func (sc *SLZBCollector) updateDeviceInfoMetrics(deviceName string, resp *http.Response, responseTime float64) {
	sc.metrics.SLZBAPIResponseTimeSeconds.With(prometheus.Labels{
		"device": deviceName,
		"action": "0",
	}).Observe(responseTime)

	sc.metrics.SLZBHTTPRequestsTotal.With(prometheus.Labels{
		"device": deviceName,
		"action": "0",
		"status": strconv.Itoa(resp.StatusCode),
	}).Inc()
}

// handleDeviceInfoHTTPError handles HTTP errors for device info collection
func (sc *SLZBCollector) handleDeviceInfoHTTPError(deviceName string, statusCode int) {
	sc.metrics.SLZBHTTPErrorsTotal.With(prometheus.Labels{
		"device":     deviceName,
		"action":     "0",
		"error_type": "http_error",
	}).Inc()
	slog.Error("HTTP error getting device info", "status", statusCode)
}

// processDeviceData processes the device data from the response header
func (sc *SLZBCollector) processDeviceData(ctx context.Context, deviceName, respValuesArr string) bool {
	tracer := sc.app.GetTracer()

	var collectorSpan *tracing.CollectorSpan
	var spanCtx context.Context

	if tracer != nil && tracer.IsEnabled() {
		collectorSpan = tracer.NewCollectorSpan(ctx, "slzb-collector", "process-device-data")
		collectorSpan.SetAttributes(
			attribute.String("device.id", deviceName),
			attribute.Int("data.length", len(respValuesArr)),
		)
		spanCtx = collectorSpan.Context()
		defer collectorSpan.End()
	} else {
		spanCtx = ctx
	}

	unmarshalStart := time.Now()
	var deviceData map[string]string
	if err := json.Unmarshal([]byte(respValuesArr), &deviceData); err != nil {
		unmarshalDuration := time.Since(unmarshalStart).Seconds()
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("json_unmarshal.duration_seconds", unmarshalDuration),
			)
			collectorSpan.RecordError(err, attribute.String("operation", "json-unmarshal"))
		}
		sc.metrics.SLZBHTTPErrorsTotal.With(prometheus.Labels{
			"device":     deviceName,
			"action":     "0",
			"error_type": "json_error",
		}).Inc()
		slog.Error("Failed to parse respValuesArr header", "error", err)

		return false
	}
	unmarshalDuration := time.Since(unmarshalStart).Seconds()

	// Cache the device information
	sc.deviceInfo = deviceData

	// Update various device metrics
	metricsStart := time.Now()
	sc.updateDeviceBasicMetrics(deviceName, deviceData)
	sc.updateDeviceUptimeMetrics(deviceName, deviceData)
	sc.updateDeviceHeapMetrics(deviceName, deviceData)
	sc.updateDeviceNetworkMetrics(deviceName, deviceData)
	metricsDuration := time.Since(metricsStart).Seconds()
	totalDuration := time.Since(unmarshalStart).Seconds()

	if collectorSpan != nil {
		collectorSpan.SetAttributes(
			attribute.Float64("json_unmarshal.duration_seconds", unmarshalDuration),
			attribute.Float64("update_metrics.duration_seconds", metricsDuration),
			attribute.Float64("total.duration_seconds", totalDuration),
			attribute.Int("device_data.fields", len(deviceData)),
		)
		collectorSpan.AddEvent("device_data_processed",
			attribute.Int("fields", len(deviceData)),
			attribute.Float64("total_duration_seconds", totalDuration),
		)
	}

	slog.Info("Device info collected from respValuesArr", "device", deviceName, "info", deviceData)

	return true
}

// updateDeviceBasicMetrics updates basic device metrics (temperature, mode)
func (sc *SLZBCollector) updateDeviceBasicMetrics(deviceName string, deviceData map[string]string) {
	// Update device temperature
	if tempStr, ok := deviceData["deviceTemp"]; ok {
		if temp, err := strconv.ParseFloat(tempStr, 64); err == nil {
			sc.metrics.SLZBDeviceTemp.With(prometheus.Labels{
				"device": deviceName,
			}).Set(temp)
		}
	}

	// Extract device operational mode
	if operationalMode, ok := deviceData["operationalMode"]; ok {
		sc.metrics.SLZBDeviceMode.With(prometheus.Labels{
			"device": deviceName,
			"mode":   operationalMode,
		}).Set(1)
	}
}

// updateDeviceUptimeMetrics updates uptime-related metrics
func (sc *SLZBCollector) updateDeviceUptimeMetrics(deviceName string, deviceData map[string]string) {
	// Update device uptime
	if uptimeStr, ok := deviceData["uptime"]; ok {
		if uptimeSeconds := sc.parseUptime(uptimeStr); uptimeSeconds > 0 {
			sc.metrics.SLZBUptime.With(prometheus.Labels{
				"device": deviceName,
			}).Set(float64(uptimeSeconds))
		}
	}

	// Update socket uptime and connection status
	if socketUptimeStr, ok := deviceData["connectedSocket"]; ok {
		if socketUptimeSeconds := sc.parseUptime(socketUptimeStr); socketUptimeSeconds > 0 {
			sc.metrics.SLZBSocketUptime.With(prometheus.Labels{
				"device": deviceName,
			}).Set(float64(socketUptimeSeconds))
			sc.metrics.SLZBSocketConnected.With(prometheus.Labels{
				"device":      deviceName,
				"connections": "1",
			}).Set(1)
		}
	} else {
		sc.metrics.SLZBSocketConnected.With(prometheus.Labels{
			"device":      deviceName,
			"connections": "0",
		}).Set(0)
	}
}

// updateDeviceHeapMetrics updates heap-related metrics
func (sc *SLZBCollector) updateDeviceHeapMetrics(deviceName string, deviceData map[string]string) {
	var heapFree, heapSize float64

	heapFreeValid := false
	heapSizeValid := false

	// Parse heap free
	if heapFreeStr, ok := deviceData["espHeapFree"]; ok {
		if parsedHeapFree, err := strconv.ParseFloat(heapFreeStr, 64); err == nil {
			heapFree = parsedHeapFree
			sc.metrics.SLZBHeapFree.With(prometheus.Labels{
				"device": deviceName,
			}).Set(heapFree)

			heapFreeValid = true
		}
	}

	// Parse heap size
	if heapSizeStr, ok := deviceData["espHeapSize"]; ok {
		if parsedHeapSize, err := strconv.ParseFloat(heapSizeStr, 64); err == nil {
			heapSize = parsedHeapSize
			sc.metrics.SLZBHeapSize.With(prometheus.Labels{
				"device": deviceName,
			}).Set(heapSize)

			heapSizeValid = true
		}
	}

	// Calculate heap ratio if both values are valid
	if heapFreeValid && heapSizeValid && heapSize > 0 {
		heapRatio := (heapFree / heapSize) * 100.0 // Convert to percentage
		sc.metrics.SLZBHeapRatio.With(prometheus.Labels{
			"device": deviceName,
		}).Set(heapRatio)
		slog.Debug("Heap ratio calculated", "device", deviceName, "free", heapFree, "size", heapSize, "ratio", heapRatio)
	}
}

// updateDeviceNetworkMetrics updates network connection metrics
func (sc *SLZBCollector) updateDeviceNetworkMetrics(deviceName string, deviceData map[string]string) {
	ethConnected := false
	ipAddr := "unknown"
	macAddr := "unknown"
	gateway := "unknown"
	subnet := "unknown"
	dns := "unknown"
	speedMbps := "unknown"

	// Check ethernet connection status
	if ethConnection, ok := deviceData["ethConnection"]; ok && ethConnection == "Connected" {
		ethConnected = true

		// Get network details from device info
		if ip, ok := deviceData["ethIp"]; ok && ip != "" {
			ipAddr = ip
		}

		if mac, ok := deviceData["ethMac"]; ok && mac != "" {
			macAddr = mac
		}

		if gate, ok := deviceData["ethGate"]; ok && gate != "" {
			gateway = gate
		}

		if mask, ok := deviceData["etchMask"]; ok && mask != "" {
			subnet = mask
		}
		// DNS is not available in device info, so we'll use gateway as DNS
		if gate, ok := deviceData["ethGate"]; ok && gate != "" {
			dns = gate
		}
		// Get ethernet speed
		if ethSpeedStr, ok := deviceData["ethSpd"]; ok {
			if ethSpeed := sc.parseEthernetSpeed(ethSpeedStr); ethSpeed > 0 {
				speedMbps = fmt.Sprintf("%.0f", ethSpeed)
			}
		}
	}

	// Set ethernet connection metrics
	if ethConnected {
		sc.metrics.SLZBEthernetConnected.With(prometheus.Labels{
			"device":      deviceName,
			"ip_address":  ipAddr,
			"mac_address": macAddr,
			"gateway":     gateway,
			"subnet_mask": subnet,
			"dns_server":  dns,
			"speed_mbps":  speedMbps,
		}).Set(1)
		sc.metrics.SLZBWifiConnected.With(prometheus.Labels{
			"device":      deviceName,
			"ssid":        "none",
			"ip_address":  "none",
			"mac_address": "none",
			"gateway":     "none",
			"subnet_mask": "none",
			"dns_server":  "none",
		}).Set(0)
		slog.Info("Ethernet connected from device info", "device", deviceName, "ip", ipAddr, "mac", macAddr, "gateway", gateway, "subnet", subnet, "dns", dns, "speed", speedMbps)
	} else {
		sc.metrics.SLZBEthernetConnected.With(prometheus.Labels{
			"device":      deviceName,
			"ip_address":  "unknown",
			"mac_address": "unknown",
			"gateway":     "unknown",
			"subnet_mask": "unknown",
			"dns_server":  "unknown",
			"speed_mbps":  "unknown",
		}).Set(0)
		sc.metrics.SLZBWifiConnected.With(prometheus.Labels{
			"device":      deviceName,
			"ssid":        "unknown",
			"ip_address":  "unknown",
			"mac_address": "unknown",
			"gateway":     "unknown",
			"subnet_mask": "unknown",
			"dns_server":  "unknown",
		}).Set(0)
		slog.Info("Ethernet disconnected from device info", "device", deviceName)
	}
}

// setDefaultDeviceInfo sets default device information when header is not available
func (sc *SLZBCollector) setDefaultDeviceInfo(deviceName string) {
	sc.deviceInfo["name"] = "SLZB"
	sc.deviceInfo["model"] = "SLZB"
	sc.deviceInfo["firmware"] = "unknown"
	slog.Debug("Device info collected with defaults", "device", deviceName, "info", sc.deviceInfo)
}

// NEW: collectFirmwareStatus collects firmware version and update status
func (sc *SLZBCollector) collectFirmwareStatus(ctx context.Context, deviceName string) bool {
	tracer := sc.app.GetTracer()

	var collectorSpan *tracing.CollectorSpan
	var spanCtx context.Context

	if tracer != nil && tracer.IsEnabled() {
		collectorSpan = tracer.NewCollectorSpan(ctx, "slzb-collector", "collect-firmware-status")
		collectorSpan.SetAttributes(
			attribute.String("device.id", deviceName),
		)
		spanCtx = collectorSpan.Context()
		defer collectorSpan.End()
	} else {
		spanCtx = ctx
	}

	startTime := time.Now()

	// Get firmware information from device info (already collected)
	firmwareVersion := "unknown"
	if deviceInfo, ok := sc.deviceInfo["VERSION"]; ok {
		firmwareVersion = deviceInfo
		sc.metrics.SLZBFirmwareCurrentVersion.With(prometheus.Labels{
			"device":     deviceName,
			"version":    deviceInfo,
			"build_date": "unknown",
		}).Set(1)
	} else {
		if collectorSpan != nil {
			collectorSpan.AddEvent("firmware_version_not_found")
		}
	}

	// Check for firmware updates (this would require additional API calls)
	// For now, we'll set a default value
	sc.metrics.SLZBFirmwareUpdateAvailable.With(prometheus.Labels{
		"device":            deviceName,
		"available_version": "unknown",
	}).Set(0)
	sc.metrics.SLZBFirmwareLastCheckTime.With(prometheus.Labels{
		"device": deviceName,
	}).Set(float64(time.Now().Unix()))

	// Record API response time (using device info collection time)
	responseTime := time.Since(startTime).Seconds()
	sc.metrics.SLZBAPIResponseTimeSeconds.With(prometheus.Labels{
		"device": deviceName,
		"action": "firmware",
	}).Observe(responseTime)

	if collectorSpan != nil {
		collectorSpan.SetAttributes(
			attribute.Float64("collection.duration_seconds", responseTime),
			attribute.String("firmware.version", firmwareVersion),
		)
		collectorSpan.AddEvent("firmware_status_collected",
			attribute.String("version", firmwareVersion),
			attribute.Float64("duration_seconds", responseTime),
		)
	}

	slog.Debug("Firmware status collected", "device", deviceName, "response_time", responseTime)

	return true
}

// NEW: collectConfigurationMetrics collects configuration file metrics
func (sc *SLZBCollector) collectConfigurationMetrics(ctx context.Context, deviceName string) bool {
	tracer := sc.app.GetTracer()

	var collectorSpan *tracing.CollectorSpan

	if tracer != nil && tracer.IsEnabled() {
		collectorSpan = tracer.NewCollectorSpan(ctx, "slzb-collector", "collect-configuration-metrics")
		collectorSpan.SetAttributes(
			attribute.String("device.id", deviceName),
			attribute.String("device.action", "4"),
		)
		defer collectorSpan.End()
	}

	startTime := time.Now()

	// Get file list from action 4
	apiStart := time.Now()
	resp, err := sc.client.Get(fmt.Sprintf("%s/api?action=4&page=0", sc.config.SLZB.APIURL))
	apiDuration := time.Since(apiStart).Seconds()

	if err != nil {
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("api_request.duration_seconds", apiDuration),
			)
			collectorSpan.RecordError(err, attribute.String("operation", "api-request"))
		}
		sc.metrics.SLZBHTTPErrorsTotal.With(prometheus.Labels{
			"device":     deviceName,
			"action":     "4",
			"error_type": "request_error",
		}).Inc()
		sc.metrics.SLZBAPITimeoutErrorsTotal.With(prometheus.Labels{
			"device": deviceName,
			"action": "4",
		}).Inc()
		slog.Error("Failed to get configuration file list", "error", err)

		return false
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("Failed to close response body", "error", err)
		}
	}()

	// Record API response time
	responseTime := time.Since(startTime).Seconds()
	sc.metrics.SLZBAPIResponseTimeSeconds.With(prometheus.Labels{
		"device": deviceName,
		"action": "4",
	}).Observe(responseTime)

	sc.metrics.SLZBHTTPRequestsTotal.With(prometheus.Labels{
		"device": deviceName,
		"action": "4",
		"status": strconv.Itoa(resp.StatusCode),
	}).Inc()

	if resp.StatusCode != http.StatusOK {
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("api_request.duration_seconds", apiDuration),
				attribute.Int("http.status_code", resp.StatusCode),
			)
			collectorSpan.RecordError(fmt.Errorf("HTTP error: %d", resp.StatusCode), attribute.Int("status_code", resp.StatusCode))
		}
		sc.metrics.SLZBHTTPErrorsTotal.With(prometheus.Labels{
			"device":     deviceName,
			"action":     "4",
			"error_type": "http_error",
		}).Inc()
		slog.Error("HTTP error getting configuration file list", "status", resp.StatusCode)

		return false
	}

	// Parse file list response
	var fileList struct {
		Files []struct {
			Filename string `json:"filename"`
			Size     int    `json:"size"`
		} `json:"files"`
	}

	readStart := time.Now()
	body, err := io.ReadAll(resp.Body)
	readDuration := time.Since(readStart).Seconds()

	if err != nil {
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("read_body.duration_seconds", readDuration),
			)
			collectorSpan.RecordError(err, attribute.String("operation", "read-body"))
		}
		sc.metrics.SLZBHTTPErrorsTotal.With(prometheus.Labels{
			"device":     deviceName,
			"action":     "4",
			"error_type": "read_error",
		}).Inc()
		slog.Error("Failed to read configuration file list response", "error", err)

		return false
	}

	unmarshalStart := time.Now()
	if err := json.Unmarshal(body, &fileList); err != nil {
		unmarshalDuration := time.Since(unmarshalStart).Seconds()
		if collectorSpan != nil {
			collectorSpan.SetAttributes(
				attribute.Float64("json_unmarshal.duration_seconds", unmarshalDuration),
			)
			collectorSpan.RecordError(err, attribute.String("operation", "json-unmarshal"))
		}
		sc.metrics.SLZBHTTPErrorsTotal.With(prometheus.Labels{
			"device":     deviceName,
			"action":     "4",
			"error_type": "json_error",
		}).Inc()
		slog.Error("Failed to parse configuration file list", "error", err)

		return false
	}
	unmarshalDuration := time.Since(unmarshalStart).Seconds()

	// Process configuration files
	processStart := time.Now()
	totalSize := 0
	configFiles := 0
	backupFiles := 0

	for _, file := range fileList.Files {
		totalSize += file.Size

		if strings.Contains(file.Filename, "config") {
			configFiles++
		}

		if strings.Contains(file.Filename, "backup") {
			backupFiles++
		}
	}
	processDuration := time.Since(processStart).Seconds()

	// Update metrics
	sc.metrics.SLZBConfigFileCount.With(prometheus.Labels{
		"device":    deviceName,
		"file_type": "config",
	}).Set(float64(configFiles))
	sc.metrics.SLZBConfigFileCount.With(prometheus.Labels{
		"device":    deviceName,
		"file_type": "backup",
	}).Set(float64(backupFiles))

	// Set backup status (assuming success if backup files exist)
	if backupFiles > 0 {
		sc.metrics.SLZBConfigBackupStatus.With(prometheus.Labels{
			"device":      deviceName,
			"backup_type": "auto",
		}).Set(1)
		sc.metrics.SLZBConfigLastBackupTime.With(prometheus.Labels{
			"device":      deviceName,
			"backup_type": "auto",
		}).Set(float64(time.Now().Unix()))
	} else {
		sc.metrics.SLZBConfigBackupStatus.With(prometheus.Labels{
			"device":      deviceName,
			"backup_type": "auto",
		}).Set(0)
	}

	totalDuration := time.Since(startTime).Seconds()

	if collectorSpan != nil {
		collectorSpan.SetAttributes(
			attribute.Float64("api_request.duration_seconds", apiDuration),
			attribute.Float64("read_body.duration_seconds", readDuration),
			attribute.Float64("json_unmarshal.duration_seconds", unmarshalDuration),
			attribute.Float64("process_files.duration_seconds", processDuration),
			attribute.Float64("total.duration_seconds", totalDuration),
			attribute.Int("files.total", len(fileList.Files)),
			attribute.Int("files.config", configFiles),
			attribute.Int("files.backup", backupFiles),
			attribute.Int("files.total_size_bytes", totalSize),
		)
		collectorSpan.AddEvent("configuration_metrics_collected",
			attribute.Int("total_files", len(fileList.Files)),
			attribute.Int("config_files", configFiles),
			attribute.Int("backup_files", backupFiles),
		)
	}

	slog.Debug("Configuration metrics collected", "device", deviceName, "files", len(fileList.Files), "response_time", responseTime)

	return true
}

// parseUptime converts uptime string like "7 d 16:47:19" to seconds
func (sc *SLZBCollector) parseUptime(uptimeStr string) int64 {
	// Trim whitespace
	uptimeStr = strings.TrimSpace(uptimeStr)
	if uptimeStr == "" {
		slog.Warn("Empty uptime string")
		return 0
	}

	parts := strings.Fields(uptimeStr)
	if len(parts) != 3 {
		slog.Warn("Invalid uptime format - expected 3 parts", "uptime", uptimeStr, "parts", parts)
		return 0
	}

	// Validate format: "7 d 16:47:19"
	if parts[1] != "d" {
		slog.Warn("Invalid uptime format - expected 'd' separator", "uptime", uptimeStr, "separator", parts[1])
		return 0
	}

	// Parse days
	days, err := strconv.Atoi(parts[0])
	if err != nil {
		slog.Warn("Invalid days value", "uptime", uptimeStr, "days", parts[0], "error", err)
		return 0
	}

	if days < 0 {
		slog.Warn("Negative days value", "uptime", uptimeStr, "days", days)
		return 0
	}

	// Parse time part "16:47:19"
	timeParts := strings.Split(parts[2], ":")
	if len(timeParts) != 3 {
		slog.Warn("Invalid time format - expected HH:MM:SS", "uptime", uptimeStr, "time", parts[2])
		return 0
	}

	// Parse hours
	hours, err := strconv.Atoi(timeParts[0])
	if err != nil {
		slog.Warn("Invalid hours value", "uptime", uptimeStr, "hours", timeParts[0], "error", err)
		return 0
	}

	if hours < 0 || hours > 23 {
		slog.Warn("Invalid hours value - must be 0-23", "uptime", uptimeStr, "hours", hours)
		return 0
	}

	// Parse minutes
	minutes, err := strconv.Atoi(timeParts[1])
	if err != nil {
		slog.Warn("Invalid minutes value", "uptime", uptimeStr, "minutes", timeParts[1], "error", err)
		return 0
	}

	if minutes < 0 || minutes > 59 {
		slog.Warn("Invalid minutes value - must be 0-59", "uptime", uptimeStr, "minutes", minutes)
		return 0
	}

	// Parse seconds
	seconds, err := strconv.Atoi(timeParts[2])
	if err != nil {
		slog.Warn("Invalid seconds value", "uptime", uptimeStr, "seconds", timeParts[2], "error", err)
		return 0
	}

	if seconds < 0 || seconds > 59 {
		slog.Warn("Invalid seconds value - must be 0-59", "uptime", uptimeStr, "seconds", seconds)
		return 0
	}

	totalSeconds := int64(days*86400 + hours*3600 + minutes*60 + seconds)
	slog.Debug("Parsed uptime", "input", uptimeStr, "days", days, "hours", hours, "minutes", minutes, "seconds", seconds, "total_seconds", totalSeconds)

	return totalSeconds
}

// parseEthernetSpeed converts speed string like "100 Mbps" to numeric value
func (sc *SLZBCollector) parseEthernetSpeed(speedStr string) float64 {
	// Trim whitespace
	speedStr = strings.TrimSpace(speedStr)
	if speedStr == "" {
		slog.Warn("Empty ethernet speed string")
		return 100.0
	}

	parts := strings.Fields(speedStr)
	if len(parts) != 2 {
		slog.Warn("Invalid ethernet speed format - expected 2 parts", "speed", speedStr, "parts", parts)
		return 100.0
	}

	// Parse speed value
	speed, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		slog.Warn("Invalid speed value", "speed", speedStr, "value", parts[0], "error", err)
		return 100.0
	}

	if speed < 0 {
		slog.Warn("Negative speed value", "speed", speedStr, "value", speed)
		return 100.0
	}

	// Parse unit
	unit := strings.ToLower(parts[1])
	switch unit {
	case "mbps", "mbit/s":
		return speed
	case "gbps", "gbit/s":
		return speed * 1000
	case "kbps", "kbit/s":
		return speed / 1000
	default:
		slog.Warn("Unknown speed unit", "speed", speedStr, "unit", parts[1])
		return 100.0
	}
}

// Stop stops the collector
func (sc *SLZBCollector) Stop() {
	slog.Info("Stopping SLZB collector...")
	// No cleanup needed for HTTP client
}

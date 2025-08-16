package collectors

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/d0ugal/slzb-exporter/internal/config"
	"github.com/d0ugal/slzb-exporter/internal/metrics"
)

// SLZBCollector collects metrics from SLZB devices
type SLZBCollector struct {
	config  *config.Config
	metrics *metrics.Registry
	client  *http.Client

	deviceInfo map[string]string // Cache for device information
	deviceID   string            // Derived device identifier
}

// NewSLZBCollector creates a new SLZB collector
func NewSLZBCollector(cfg *config.Config, metricsRegistry *metrics.Registry) *SLZBCollector {
	// Derive device ID from API URL
	deviceID := deriveDeviceID(cfg.SLZB.APIURL)

	return &SLZBCollector{
		config:  cfg,
		metrics: metricsRegistry,
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
	sc.collectMetrics()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping SLZB collector")
			return
		case <-ticker.C:
			sc.collectMetrics()
		}
	}
}

// collectMetrics collects all metrics from the SLZB device
func (sc *SLZBCollector) collectMetrics() {
	deviceID := sc.deviceID
	collectionStart := time.Now()
	successfulCollections := 0
	totalCollections := 0

	// Track collection success
	defer func() {
		// Update last collection timestamp if we had any successful collections
		if successfulCollections > 0 {
			sc.metrics.SLZBLastCollectionTime.WithLabelValues(deviceID).Set(float64(time.Now().Unix()))
		}

		// Log collection summary
		slog.Info("Collection cycle completed",
			"device", deviceID,
			"successful", successfulCollections,
			"total", totalCollections,
			"duration", time.Since(collectionStart))
	}()

	// Get device information and test reachability in one request
	if sc.collectDeviceInfo(deviceID) {
		sc.metrics.SLZBDeviceReachable.WithLabelValues(deviceID).Set(1)
		sc.metrics.SLZBConnected.WithLabelValues(deviceID).Set(1)
		successfulCollections++
	} else {
		sc.metrics.SLZBDeviceReachable.WithLabelValues(deviceID).Set(0)
		sc.metrics.SLZBConnected.WithLabelValues(deviceID).Set(0)
		sc.metrics.SLZBCollectionErrors.WithLabelValues(deviceID, "device_unreachable").Inc()
		slog.Error("Device unreachable", "device", deviceID)
		return
	}
	totalCollections++

	// Collect device information
	if sc.collectDeviceInfo(deviceID) {
		successfulCollections++
	}
	totalCollections++

	// Add a small delay between requests
	time.Sleep(500 * time.Millisecond)

}

// collectDeviceInfo collects device information and caches it
func (sc *SLZBCollector) collectDeviceInfo(deviceName string) bool {
	// Get device information from action 0
	resp, err := sc.client.Get(fmt.Sprintf("%s/api?action=0&page=0", sc.config.SLZB.APIURL))
	if err != nil {
		sc.metrics.SLZBHTTPErrorsTotal.WithLabelValues(deviceName, "0", "request_error").Inc()
		slog.Error("Failed to get device info", "error", err)
		return false
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("Failed to close response body", "error", err)
		}
	}()

	sc.metrics.SLZBHTTPRequestsTotal.WithLabelValues(deviceName, "0", strconv.Itoa(resp.StatusCode)).Inc()

	if resp.StatusCode != http.StatusOK {
		sc.metrics.SLZBHTTPErrorsTotal.WithLabelValues(deviceName, "0", "http_error").Inc()
		slog.Error("HTTP error getting device info", "status", resp.StatusCode)
		return false
	}

	// Parse the respValuesArr header which contains device information
	respValuesArr := resp.Header.Get("respValuesArr")
	if respValuesArr != "" {
		var deviceData map[string]string
		if err := json.Unmarshal([]byte(respValuesArr), &deviceData); err != nil {
			sc.metrics.SLZBHTTPErrorsTotal.WithLabelValues(deviceName, "0", "json_error").Inc()
			slog.Error("Failed to parse respValuesArr header", "error", err)
		} else {
			// Cache the device information
			sc.deviceInfo = deviceData

			// Update metrics with real device data
			if tempStr, ok := deviceData["deviceTemp"]; ok {
				if temp, err := strconv.ParseFloat(tempStr, 64); err == nil {
					sc.metrics.SLZBDeviceTemp.WithLabelValues(deviceName).Set(temp)
				}
			}

			if uptimeStr, ok := deviceData["uptime"]; ok {
				// Parse uptime like "7 d 15:23:25" to seconds
				if uptimeSeconds := sc.parseUptime(uptimeStr); uptimeSeconds > 0 {
					sc.metrics.SLZBUptime.WithLabelValues(deviceName).Set(float64(uptimeSeconds))
				}
			}

			if socketUptimeStr, ok := deviceData["connectedSocket"]; ok {
				// Parse socket uptime like "7 d 18:28:12" to seconds
				if socketUptimeSeconds := sc.parseUptime(socketUptimeStr); socketUptimeSeconds > 0 {
					sc.metrics.SLZBSocketUptime.WithLabelValues(deviceName).Set(float64(socketUptimeSeconds))
					sc.metrics.SLZBSocketConnected.WithLabelValues(deviceName, "1").Set(1)
				}
			} else {
				sc.metrics.SLZBSocketConnected.WithLabelValues(deviceName, "0").Set(0)
			}

			// Extract device operational mode
			if operationalMode, ok := deviceData["operationalMode"]; ok {
				sc.metrics.SLZBDeviceMode.WithLabelValues(deviceName, operationalMode).Set(1)
			}

			var heapFree, heapSize float64
			heapFreeValid := false
			heapSizeValid := false

			if heapFreeStr, ok := deviceData["espHeapFree"]; ok {
				if parsedHeapFree, err := strconv.ParseFloat(heapFreeStr, 64); err == nil {
					heapFree = parsedHeapFree
					sc.metrics.SLZBHeapFree.WithLabelValues(deviceName).Set(heapFree)
					heapFreeValid = true
				}
			}

			if heapSizeStr, ok := deviceData["espHeapSize"]; ok {
				if parsedHeapSize, err := strconv.ParseFloat(heapSizeStr, 64); err == nil {
					heapSize = parsedHeapSize
					sc.metrics.SLZBHeapSize.WithLabelValues(deviceName).Set(heapSize)
					heapSizeValid = true
				}
			}

			// Calculate heap ratio if both values are valid
			if heapFreeValid && heapSizeValid && heapSize > 0 {
				heapRatio := (heapFree / heapSize) * 100.0 // Convert to percentage
				sc.metrics.SLZBHeapRatio.WithLabelValues(deviceName).Set(heapRatio)
				slog.Debug("Heap ratio calculated", "device", deviceName, "free", heapFree, "size", heapSize, "ratio", heapRatio)
			}

			// Extract ethernet connection status from device info
			ethConnected := false
			ipAddr := "unknown"
			macAddr := "unknown"
			gateway := "unknown"
			subnet := "unknown"
			dns := "unknown"
			speedMbps := "unknown"

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

			// Set ethernet connection metrics based on device info
			if ethConnected {
				sc.metrics.SLZBEthernetConnected.WithLabelValues(deviceName, ipAddr, macAddr, gateway, subnet, dns, speedMbps).Set(1)
				sc.metrics.SLZBWifiConnected.WithLabelValues(deviceName, "none", "none", "none", "none", "none", "none").Set(0)
				slog.Info("Ethernet connected from device info", "device", deviceName, "ip", ipAddr, "mac", macAddr, "gateway", gateway, "subnet", subnet, "dns", dns, "speed", speedMbps)
			} else {
				sc.metrics.SLZBEthernetConnected.WithLabelValues(deviceName, "unknown", "unknown", "unknown", "unknown", "unknown", "unknown", "unknown").Set(0)
				sc.metrics.SLZBWifiConnected.WithLabelValues(deviceName, "unknown", "unknown", "unknown", "unknown", "unknown", "unknown").Set(0)
				slog.Info("Ethernet disconnected from device info", "device", deviceName)
			}

			slog.Info("Device info collected from respValuesArr", "device", deviceName, "info", deviceData)
		}
	} else {
		// Fallback to default values if header is not available
		sc.deviceInfo["name"] = "SLZB"
		sc.deviceInfo["model"] = "SLZB"
		sc.deviceInfo["firmware"] = "unknown"
		slog.Debug("Device info collected with defaults", "device", deviceName, "info", sc.deviceInfo)
	}

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

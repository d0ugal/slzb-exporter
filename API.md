# SLZB-06 API Documentation

This document describes the SLZB-06 device API that this exporter consumes. This information is useful for developers who want to understand how the exporter works or contribute to its development.

## API Discovery

Based on exploration of the SLZB-06 device, the following API endpoints have been discovered:

## Available API Endpoints

Based on the firmware source code at [smlight-dev/slzb-06-firmware](https://github.com/smlight-dev/slzb-06-firmware/blob/main/src/web.cpp#L151), the API uses the following action enumeration:

```cpp
enum API_ACTION_t : uint8_t { 
    API_GET_PAGE,           // 0
    API_GET_PARAM,          // 1
    API_STARTWIFISCAN,      // 2
    API_WIFISCANSTATUS,     // 3
    API_GET_FILELIST,       // 4
    API_GET_FILE,           // 5
    API_SEND_HEX,           // 6
    API_WIFICONNECTSTAT,    // 7
    API_CMD,                // 8
    API_GET_LOG             // 9
};
```

All endpoints use the pattern: `GET /api?action=<action>&page=<page>`

### Action 0 (API_GET_PAGE): HTML Template
- **URL**: `/api?action=0&page=0`
- **Response**: Gzipped HTML template with data-replace placeholders
- **Content**: Device status page template with placeholders like `data-replace="operationalMode"`, `data-replace="connectedEther"`, etc.

### Action 1 (API_GET_PARAM): Get Parameters
- **URL**: `/api?action=1&param=<param_name>`
- **Response**: Text response with parameter value
- **Parameters**:
  - `refreshLogs` - Returns refresh logs setting
  - `coordMode` - Returns coordinator mode (1=wifi setup, otherwise coordinator_mode setting)

### Action 2 (API_STARTWIFISCAN): Start WiFi Scan
- **URL**: `/api?action=2&page=0`
- **Response**: "ok" if successful
- **Note**: Enables WiFi in STA mode if it was off

### Action 3 (API_WIFISCANSTATUS): WiFi Scan Results
- **URL**: `/api?action=3&page=0`
- **Response**: JSON with WiFi scan data
- **Example**:
```json
{
  "scanDone": true,
  "wifi": [
    {
      "ssid": "ctu",
      "rssi": -49,
      "channel": 6,
      "secure": 3
    }
  ]
}
```

### Action 4 (API_GET_FILELIST): File List
- **URL**: `/api?action=4&page=0`
- **Response**: JSON with configuration files
- **Example**:
```json
{
  "files": [
    {
      "filename": "configEther.json",
      "size": 36
    },
    {
      "filename": "configGeneral.json", 
      "size": 141
    }
  ]
}
```

### Action 5 (API_GET_FILE): Get File Content
- **URL**: `/api?action=5&filename=<filename>`
- **Response**: File content as text
- **Note**: Files are read from `/config/` directory

### Action 6 (API_SEND_HEX): Send Hex Data
- **URL**: `/api?action=6&hex=<hex_data>&size=<size>`
- **Response**: "ok" if successful
- **Note**: Sends hex data to Serial2 (Zigbee communication)

### Action 7 (API_WIFICONNECTSTAT): WiFi Connection Status
- **URL**: `/api?action=7&page=0`
- **Response**: JSON with WiFi connection status
- **Example**:
```json
{
  "connected": false
}
```
- **Note**: Returns `{"connected": true, "ip": "192.168.1.100"}` if connected

### Action 8 (API_CMD): Execute Commands
- **URL**: `/api?action=8&cmd=<command_id>`
- **Response**: "ok" if successful
- **Commands**:
  - `0` - CMD_ZB_ROUTER_RECON: Router reconnect
  - `1` - CMD_ZB_RST: Zigbee restart
  - `2` - CMD_ZB_BSL: Zigbee enable BSL
  - `3` - CMD_ESP_RES: ESP32 restart
  - `4` - CMD_ADAP_LAN: Adapter mode LAN
  - `5` - CMD_ADAP_USB: Adapter mode USB
  - `6` - CMD_LEDY_TOG: LED Yellow toggle
  - `7` - CMD_LEDB_TOG: LED Blue toggle
  - `8` - CMD_CLEAR_LOG: Clear log

### Action 9 (API_GET_LOG): Get Log
- **URL**: `/api?action=9&page=0`
- **Response**: Log content as text
- **Note**: Returns the device log output

## Data-Replace Placeholders Found

From the HTML template, these are the metrics available:

### Device Status
- `operationalMode` - Device operational mode
- `connectedEther` - Ethernet connected status
- `connectedSocketStatus` - Socket client connected status
- `wifiEnabled` - WiFi Client enabled status
- `wifiConnected` - WiFi Client connection status
- `wifiModeAP` - WiFi Access Point enabled status
- `wifiModeAPStatus` - WiFi Access Point status
- `uptime` - Device uptime
- `connectedSocket` - Socket uptime

### Device Information
- `hwRev` - Hardware revision/model
- `VERSION` - ESP32 Firmware version
- `espModel` - ESP32 version
- `deviceTemp` - ESP32 temperature (°C)
- `espCores` - Number of ESP32 cores
- `espFreq` - ESP32 frequency (MHz)
- `espFlashSize` - ESP32 flash size (MB)
- `espFlashType` - ESP32 flash type
- `espHeapFree` - ESP32 free heap (KB)
- `espHeapSize` - ESP32 total heap size (KB)

### Ethernet
- `ethConnection` - Connection status
- `ethDhcp` - DHCP status
- `ethIp` - IP Address
- `etchMask` - Subnet Mask
- `ethGate` - Default Gateway
- `ethSpd` - Connection speed
- `ethMac` - MAC address

### WiFi
- `wifiMode` - WiFi mode
- `wifiSsid` - SSID
- `wifiMac` - MAC Address
- `wifiIp` - IP Address
- `wifiSubnet` - Subnet Mask
- `wifiGate` - Default Gateway
- `wifiRssi` - RSSI
- `wifiDhcp` - DHCP status

## Notes

- All GET operations are safe and don't modify device state
- Some endpoints return gzipped content that needs decompression
- The device appears to use a simple action/page parameter system
- JSON responses are available for some endpoints
- Raw Zigbee data is available via action 9
- The web interface uses JavaScript to populate data-replace placeholders

## Future Exploration

To get the actual values for the data-replace placeholders, we need to find the endpoint that returns the JSON data with the actual values. This might be:
- A different action number
- A different parameter combination
- A separate API endpoint pattern
- Embedded in the JavaScript files

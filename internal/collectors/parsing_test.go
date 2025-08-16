package collectors

import (
	"testing"
)

func TestParseUptime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "Standard format with days and time",
			input:    "7 d 16:47:19",
			expected: 7*86400 + 16*3600 + 47*60 + 19, // 7 days, 16 hours, 47 minutes, 19 seconds
		},
		{
			name:     "Zero days with time",
			input:    "0 d 12:30:45",
			expected: 12*3600 + 30*60 + 45, // 12 hours, 30 minutes, 45 seconds
		},
		{
			name:     "Large number of days",
			input:    "365 d 23:59:59",
			expected: 365*86400 + 23*3600 + 59*60 + 59, // 365 days, 23 hours, 59 minutes, 59 seconds
		},
		{
			name:     "Single digit hours",
			input:    "1 d 5:30:15",
			expected: 1*86400 + 5*3600 + 30*60 + 15, // 1 day, 5 hours, 30 minutes, 15 seconds
		},
		{
			name:     "Zero minutes and seconds",
			input:    "2 d 10:00:00",
			expected: 2*86400 + 10*3600, // 2 days, 10 hours
		},
		{
			name:     "Just time without days",
			input:    "0 d 00:01:30",
			expected: 90, // 1 minute, 30 seconds
		},
		{
			name:     "Maximum time values",
			input:    "999 d 23:59:59",
			expected: 999*86400 + 23*3600 + 59*60 + 59,
		},
		{
			name:     "Single digit values",
			input:    "1 d 1:1:1",
			expected: 1*86400 + 1*3600 + 1*60 + 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &SLZBCollector{}
			result := collector.parseUptime(tt.input)
			if result != tt.expected {
				t.Errorf("parseUptime(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseUptimeInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Empty string",
			input: "",
		},
		{
			name:  "Invalid format - missing parts",
			input: "7 d",
		},
		{
			name:  "Invalid format - wrong separator",
			input: "7 days 16:47:19",
		},
		{
			name:  "Invalid format - no time",
			input: "7 d abc",
		},
		{
			name:  "Invalid format - negative time",
			input: "7 d -16:47:19",
		},
		{
			name:  "Invalid format - too many colons",
			input: "7 d 16:47:19:30",
		},
		{
			name:  "Invalid format - non-numeric days",
			input: "abc d 16:47:19",
		},
		{
			name:  "Invalid format - non-numeric time",
			input: "7 d abc:def:ghi",
		},
		{
			name:  "Invalid format - missing d",
			input: "7 16:47:19",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &SLZBCollector{}
			result := collector.parseUptime(tt.input)
			if result != 0 {
				t.Errorf("parseUptime(%q) = %d, want 0 for invalid input", tt.input, result)
			}
		})
	}
}

func TestParseEthernetSpeed(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "Mbps format",
			input:    "100 Mbps",
			expected: 100.0,
		},
		{
			name:     "Mbps with decimal",
			input:    "1000.5 Mbps",
			expected: 1000.5,
		},
		{
			name:     "Gbps format",
			input:    "1 Gbps",
			expected: 1000.0,
		},
		{
			name:     "Gbps with decimal",
			input:    "2.5 Gbps",
			expected: 2500.0,
		},
		{
			name:     "Kbps format",
			input:    "1000 Kbps",
			expected: 1.0,
		},
		{
			name:     "Mbit/s format",
			input:    "100 Mbit/s",
			expected: 100.0,
		},
		{
			name:     "Gbit/s format",
			input:    "1 Gbit/s",
			expected: 1000.0,
		},
		{
			name:     "Kbit/s format",
			input:    "1000 Kbit/s",
			expected: 1.0,
		},
		{
			name:     "Case insensitive Mbps",
			input:    "100 MBPS",
			expected: 100.0,
		},
		{
			name:     "Case insensitive Gbps",
			input:    "1 GBPS",
			expected: 1000.0,
		},
		{
			name:     "Zero speed",
			input:    "0 Mbps",
			expected: 0.0,
		},
		{
			name:     "Large speed value",
			input:    "10000 Mbps",
			expected: 10000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &SLZBCollector{}
			result := collector.parseEthernetSpeed(tt.input)
			if result != tt.expected {
				t.Errorf("parseEthernetSpeed(%q) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseEthernetSpeedInvalid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: 100.0, // Default value
		},
		{
			name:     "Invalid format - no unit",
			input:    "100",
			expected: 100.0, // Default value
		},
		{
			name:     "Invalid format - unknown unit",
			input:    "100 Tbps",
			expected: 100.0, // Default value
		},
		{
			name:     "Invalid format - non-numeric",
			input:    "abc Mbps",
			expected: 100.0, // Default value
		},
		{
			name:     "Invalid format - negative value",
			input:    "-100 Mbps",
			expected: 100.0, // Default value
		},
		{
			name:     "Invalid format - extra spaces",
			input:    "  100  Mbps  ",
			expected: 100.0, // Default value
		},
		{
			name:     "Invalid format - missing value",
			input:    "Mbps",
			expected: 100.0, // Default value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &SLZBCollector{}
			result := collector.parseEthernetSpeed(tt.input)
			if result != tt.expected {
				t.Errorf("parseEthernetSpeed(%q) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkParseUptime(b *testing.B) {
	collector := &SLZBCollector{}
	input := "7 d 16:47:19"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.parseUptime(input)
	}
}

func BenchmarkParseEthernetSpeed(b *testing.B) {
	collector := &SLZBCollector{}
	input := "100 Mbps"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.parseEthernetSpeed(input)
	}
}

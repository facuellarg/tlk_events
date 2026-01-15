package utils

import (
	"time"
)

// parseTimestamp parses a timestamp string in ISO 8601 format
// Supports formats: RFC3339 (2006-01-02T15:04:05Z07:00) and similar variations
func ParseTimestamp(timestamp string) (time.Time, error) {
	// Try parsing with different time formats
	formats := []string{
		time.RFC3339,                  // 2006-01-02T15:04:05Z07:00
		time.RFC3339Nano,              // 2006-01-02T15:04:05.999999999Z07:00
		"2006-01-02T15:04:05",         // Without timezone
		"2006-01-02 15:04:05",         // Space separator
		"2006-01-02T15:04:05.000Z",    // With milliseconds
		"2006-01-02T15:04:05.000000Z", // With microseconds
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, timestamp)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	return time.Time{}, lastErr
}

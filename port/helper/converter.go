package helper

import "time"

// GetStringValue Helper function to safely convert pointer types
func GetStringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// GetTimeValue Helper function to safely convert pointer types
func GetTimeValue(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}

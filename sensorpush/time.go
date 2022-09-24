package sensorpush

import (
	"time"
)

const (
	statusLayout = "Mon Jan 02 2006 15:04:05 MST"
)

func parseTimeStatus(s string) (time.Time, error) {
	return time.Parse(statusLayout, s)
}

func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

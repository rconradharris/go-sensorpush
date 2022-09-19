package sensorpush

import (
	"time"
)

const (
	timeLayout = "Mon Jan 02 2006 15:04:05 MST"
)

func parseTime(s string) (time.Time, error) {
	return time.Parse(timeLayout, s)
}

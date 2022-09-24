package sensorpush

import (
	"time"

	"github.com/rconradharris/go-sensorpush/units"
)

type Sample struct {
	Observed    time.Time
	Temperature *units.Temperature
}

func newSample(sr sampleResponse) (Sample, error) {
	s := Sample{}

	// Observed
	t, err := parseTime(sr.Observed)
	if err != nil {
		return s, err
	}
	s.Observed = t

	// Temperature
	if sr.Temperature != nil {
		temp := units.NewTemperatureF(*sr.Temperature)
		s.Temperature = &temp
	}

	return s, nil
}

type sampleResponse struct {
	Observed    string   `json:"observed"`
	Temperature *float32 `json:temperature"`
}

package sensorpush

import (
	"time"

	"github.com/rconradharris/go-sensorpush/units"
)

type Sample struct {
	Humidity    *units.Humidity
	Observed    time.Time
	Temperature *units.Temperature
}

type SampleSlice []*Sample

func (s SampleSlice) Len() int {
	return len(s)
}

func (s SampleSlice) Less(i, j int) bool {
	// Reverse chron
	//
	// The API appears to sort for us, but we sort ourselves to ensure
	// the sort is well-defined and stable
	return s[i].Observed.After(s[j].Observed)
}

func (s SampleSlice) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

func newSample(sr sampleResponse) (*Sample, error) {
	s := &Sample{}

	// Humidity
	if sr.Humidity != nil {
		hum := units.NewHumidity(*sr.Humidity)
		s.Humidity = &hum
	}

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
	Humidity    *float32 `json:humidity"`
	Observed    string   `json:"observed"`
	Temperature *float32 `json:temperature"`
}

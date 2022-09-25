package sensorpush

import (
	"sort"
	"time"
)

type Samples struct {
	LastTime time.Time

	Sensors SensorSamples

	Status       SampleStatus
	TotalSamples int
	TotalSensors int
	Truncated    bool
}

type SensorSamples map[string]SampleSlice

// IDs returns the sensor IDs, sorted
func (s SensorSamples) IDs() []string {
	ids := make([]string, 0, len(s))
	for id := range s {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

type samplesRequest struct {
	Active bool `json:"active"`
	//Bulk      bool     `json:"bulk"`
	//Format    string   `json:"format"`
	Limit     *int     `json:"limit,omitempty"`
	Measures  []string `json:"measures,omitempty"`
	Sensors   []string `json:"sensors,omitempty"`
	StartTime string   `json:"startTime,omitempty"`
	StopTime  string   `json:"stopTime,omitempty"`
	//Tags []string
}

type samplesResponse struct {
	LastTime string `json:"last_time"`

	Sensors sensorToSamplesResponse `json:"sensors"`

	Status       string `json:"status"`
	TotalSamples int    `json:"total_samples"` // swagger has this as a 'number' but a float32 doesn't make sense here
	TotalSensors int    `json:"total_sensors"`
	Truncated    bool   `json:"truncated"`
}

type sensorToSamplesResponse map[string][]sampleResponse

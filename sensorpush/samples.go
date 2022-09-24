package sensorpush

import (
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

type samplesRequest struct {
	Active bool `json:"active"`
	//Bulk      bool     `json:"bulk"`
	//Format    string   `json:"format"`
	Limit    *int     `json:"limit,omitempty"`
	Measures []string `json:"measures,omitempty"`
	//Sensors   []string `json:"sensors:`
	//StartTime string   `json:"startTime"`
	//StopTime  string   `json:"stopTime"`
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

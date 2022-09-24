package sensorpush

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type SampleService service

type SampleQueryFilter struct {
	//Active    bool
	//Bulk      bool
	//Format    SampleFormat
	Limit *int
	//Measures  []Measure
	//Sensors   []Sensor
	//StartTime time.Time
	//StopTime  time.Time
	//Tags []Tag
}

type samplesRequest struct {
	//Active    bool     `json:"active"`
	//Bulk      bool     `json:"bulk"`
	//Format    string   `json:"format"`
	Limit *int `json:"limit,omitempty"`
	//Measures  []string `json:"measures"`
	//Sensors   []string `json:"sensors:`
	//StartTime string   `json:"startTime"`
	//StopTime  string   `json:"stopTime"`
	//Tags []string
}

type samplesResponse struct {
	LastTime string `json:"last_time"`
	//TODO: Sensors
	Status       string `json:"status"`
	TotalSamples int    `json:"total_samples"` // swagger has this as a 'number' but a float32 doesn't make sense here
	TotalSensors int    `json:"total_sensors"`
	Truncated    bool   `json:"truncated"`
}

type Samples struct {
	LastTime     time.Time
	Status       SampleStatus
	TotalSamples int
	TotalSensors int
	Truncated    bool
}

// Query returns samples matching the criteria
func (s *SampleService) Query(ctx context.Context, f SampleQueryFilter) (Samples, error) {
	s0 := Samples{}
	sreq := samplesRequest{
		Limit: f.Limit,
	}

	req, err := s.c.NewRequest(ctx, http.MethodPost, "samples", sreq)
	if err != nil {
		return s0, err
	}

	ssresp := samplesResponse{}
	_, err = s.c.Do(req, &ssresp)
	if err != nil {
		return s0, err
	}

	ss := Samples{
		Status:       newSampleStatus(ssresp.Status),
		TotalSamples: ssresp.TotalSamples,
		TotalSensors: ssresp.TotalSensors,
		Truncated:    ssresp.Truncated,
	}

	// Last Time
	t, err := parseTime(ssresp.LastTime)
	if err != nil {
		return s0, err
	}
	ss.LastTime = t

	fmt.Printf("ssresp => %+v\n", ssresp)

	return ss, nil
}

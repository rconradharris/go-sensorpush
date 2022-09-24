package sensorpush

import (
	"context"
	"fmt"
	"net/http"
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
		Sensors:      make(SensorToSamples),
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

	// Samples
	for sensorID, sampResps := range ssresp.Sensors {
		samps := make([]Sample, 0, len(sampResps))
		for _, sr := range sampResps {
			s, err := newSample(sr)
			if err != nil {
				return s0, err
			}
			samps = append(samps, s)
		}
		ss.Sensors[sensorID] = samps
		//fmt.Printf("%s: %+v\n", k, v)
	}

	fmt.Printf("ssresp => %+v\n", ssresp)

	return ss, nil
}

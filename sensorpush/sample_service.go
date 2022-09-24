package sensorpush

import (
	"context"
	//"fmt"
	"net/http"
	"sort"
)

type SampleService service

type SampleQueryFilter struct {
	Active bool
	//Bulk      bool
	//Format    SampleFormat
	Limit    *int
	Measures []Measure
	//Sensors   []Sensor
	//StartTime time.Time
	//StopTime  time.Time
	//Tags []Tag
}

// Query returns samples matching the criteria
func (s *SampleService) Query(ctx context.Context, f SampleQueryFilter) (*Samples, error) {
	sreq := samplesRequest{
		Active: f.Active,
		Limit:  f.Limit,
	}

	if f.Measures != nil {
		ms := make([]string, 0, len(f.Measures))
		for _, m := range f.Measures {
			ms = append(ms, m.String())
		}
		sreq.Measures = ms
	}

	req, err := s.c.NewRequest(ctx, http.MethodPost, "samples", sreq)
	if err != nil {
		return nil, err
	}

	ssresp := samplesResponse{}
	_, err = s.c.Do(req, &ssresp)
	if err != nil {
		return nil, err
	}

	ss := &Samples{
		Sensors:      make(SensorSamples),
		Status:       newSampleStatus(ssresp.Status),
		TotalSamples: ssresp.TotalSamples,
		TotalSensors: ssresp.TotalSensors,
		Truncated:    ssresp.Truncated,
	}

	// Last Time
	t, err := parseTime(ssresp.LastTime)
	if err != nil {
		return nil, err
	}
	ss.LastTime = t

	// Samples
	for sensorID, sampResps := range ssresp.Sensors {
		samps := make(SampleSlice, 0, len(sampResps))
		for _, sr := range sampResps {
			s, err := newSample(sr)
			if err != nil {
				return nil, err
			}
			samps = append(samps, s)
		}

		sort.Sort(samps)
		ss.Sensors[sensorID] = samps
		//fmt.Printf("%s: %+v\n", k, v)
	}

	//fmt.Printf("ssresp => %+v\n", ssresp)

	return ss, nil
}

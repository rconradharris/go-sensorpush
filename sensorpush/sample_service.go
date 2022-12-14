package sensorpush

import (
	"context"
	//"fmt"
	"net/http"
	"sort"
	"time"
)

type SampleService service

type SampleQueryFilter struct {
	Active bool
	//Bulk      bool
	//Format    SampleFormat
	Limit     *int
	Measures  []Measure
	Sensors   []SensorID
	StartTime time.Time
	StopTime  time.Time
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

	if f.Sensors != nil {
		ss := make([]string, 0, len(f.Sensors))
		for _, m := range f.Sensors {
			ss = append(ss, m.String())
		}
		sreq.Sensors = ss
	}

	if !f.StartTime.IsZero() {
		sreq.StartTime = formatTime(f.StartTime)
	}

	if !f.StopTime.IsZero() {
		sreq.StopTime = formatTime(f.StopTime)
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
		Sensors:      make(SamplesMap),
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
	for idstr, sampResps := range ssresp.Sensors {
		samps := make(SampleSlice, 0, len(sampResps))
		for _, sr := range sampResps {
			s, err := newSample(sr)
			if err != nil {
				return nil, err
			}
			samps = append(samps, s)
		}

		sort.Sort(samps)
		id := NewSensorID(idstr)
		ss.Sensors[id] = samps
		//fmt.Printf("%s: %+v\n", k, v)
	}

	//fmt.Printf("ssresp => %+v\n", ssresp)

	return ss, nil
}

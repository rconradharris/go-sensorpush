package sensorpush

import (
	"context"
	"fmt"
	"net/http"
	"sort"
)

type SensorService service

// List returns the sensors matching the active criteria in alphabetical order
func (s *SensorService) List(ctx context.Context, active bool) (SensorSlice, error) {
	var s0 []*Sensor

	sreq := sensorsRequest{Active: &active}

	req, err := s.c.NewRequest(ctx, http.MethodPost, "devices/sensors", sreq)
	if err != nil {
		return s0, err
	}

	ssresp := sensorsResponse{}
	_, err = s.c.Do(req, &ssresp)
	if err != nil {
		return s0, err
	}

	sensors := make(SensorSlice, 0, len(ssresp))
	for id1, sresp := range ssresp {
		if id1 != sresp.ID {
			return s0, fmt.Errorf("ID mismatch %s != %s", id1, sresp.ID)
		}
		sensors = append(sensors, newSensor(sresp))
	}

	sort.Sort(sensors)

	return sensors, nil
}

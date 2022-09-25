package sensorpush

import (
	"context"
	"net/http"
)

type SensorService service

type SensorListFilter struct {
	Active *bool
}

// List returns a map of SensorIDs to Sensors
func (s *SensorService) List(ctx context.Context, f *SensorListFilter) (SensorMap, error) {
	sm := SensorMap{}

	sreq := sensorsRequest{}
	if f != nil {
		if f.Active != nil {
			sreq.Active = f.Active
		}
	}

	req, err := s.c.NewRequest(ctx, http.MethodPost, "devices/sensors", sreq)
	if err != nil {
		return sm, err
	}

	ssresp := sensorsResponse{}
	_, err = s.c.Do(req, &ssresp)
	if err != nil {
		return sm, err
	}

	for _, sresp := range ssresp {
		id := NewSensorID(sresp.ID)
		s := newSensor(sresp)
		sm[id] = s
	}

	return sm, nil
}

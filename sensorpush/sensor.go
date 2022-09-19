package sensorpush

import (
	"context"
	"fmt"
	"net/http"
)

type SensorType int

const (
	SensorTypeUnknown SensorType = iota
	SensorTypeHT1
	SensorTypeHTw
)

func (s SensorType) String() string {
	switch s {
	case SensorTypeHT1:
		return "HT1"
	case SensorTypeHTw:
		return "HT.w"
	}
	return "unknown-sensor-type"
}

func newSensorType(s string) SensorType {
	switch s {
	case "HT1":
		return SensorTypeHT1
	case "HT.w":
		return SensorTypeHTw
	}
	return SensorTypeUnknown
}

type SensorService service

type Alert struct {
	Enabled bool
	Max     float32
	Min     float32
}

type Alerts struct {
	Humidity    Alert
	Temperature Alert
}

type Calibration struct {
	Humidity    float32
	Temperature float32
}

type Sensor struct {
	Active         bool
	Address        string // MAC address
	Alerts         Alerts
	BatteryVoltage float32
	Calibration    Calibration
	DeviceID       string
	ID             string
	Name           string
	RSSI           float32 // Wireless signal strength in dB at last reading
	// TODO: tags
	Type SensorType
}

type sensorsRequest struct {
	Active bool `json:"active"`
}

type alertResponse struct {
	Enabled bool    `json:"enabled"`
	Max     float32 `json:"max"`
	Min     float32 `json:"min"`
}

type alertsResponse struct {
	Humidity    alertResponse `json:"humidity"`
	Temperature alertResponse `json:"temperature"`
}

type calibrationResponse struct {
	Humidity    float32 `json:"humidity"`
	Temperature float32 `json:"temperature"`
}

type sensorResponse struct {
	Active         bool                `json:"active"`
	Address        string              `json:"address"`
	Alerts         alertsResponse      `json:"alerts"`
	BatteryVoltage float32             `json:"battery_voltage"`
	Calibration    calibrationResponse `json:"calibration"`
	DeviceID       string              `json:"deviceId"`
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	RSSI           float32             `json:"rssi"`
	// TODO: tags
	Type string `json:"type"`
}

type sensorsResponse map[string]sensorResponse

func (s *SensorService) List(ctx context.Context, active bool) ([]*Sensor, error) {
	var s0 []*Sensor

	sreq := sensorsRequest{Active: active}

	req, err := s.c.NewRequest(ctx, http.MethodPost, "devices/sensors", sreq)
	if err != nil {
		return s0, err
	}

	ssresp := sensorsResponse{}
	_, err = s.c.Do(req, &ssresp)
	if err != nil {
		return s0, err
	}

	sensors := make([]*Sensor, 0, len(ssresp))
	for id1, sresp := range ssresp {
		if id1 != sresp.ID {
			return s0, fmt.Errorf("ID mismatch %s != %s", id1, sresp.ID)
		}

		a := sresp.Alerts
		c := sresp.Calibration
		s := &Sensor{
			Active:  sresp.Active,
			Address: sresp.Address,
			Alerts: Alerts{
				Humidity: Alert{
					Enabled: a.Humidity.Enabled,
					Max:     a.Humidity.Max,
					Min:     a.Humidity.Min,
				},
				Temperature: Alert{
					Enabled: a.Temperature.Enabled,
					Max:     a.Temperature.Max,
					Min:     a.Temperature.Min,
				},
			},
			BatteryVoltage: sresp.BatteryVoltage,
			Calibration: Calibration{
				Humidity:    c.Humidity,
				Temperature: c.Temperature,
			},
			DeviceID: sresp.DeviceID,
			ID:       sresp.ID,
			Name:     sresp.Name,
			RSSI:     sresp.RSSI,
			Type:     newSensorType(sresp.Type),
		}
		sensors = append(sensors, s)
	}

	return sensors, nil

}

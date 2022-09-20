package sensorpush

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/rconradharris/go-sensorpush/units"
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

type AlertHumidity struct {
	Enabled bool
	Max     units.Percentage
	Min     units.Percentage
}

type AlertTemperature struct {
	Enabled bool
	Max     units.Temperature
	Min     units.Temperature
}

type Alerts struct {
	Humidity    AlertHumidity
	Temperature AlertTemperature
}

type Calibration struct {
	Humidity    units.Percentage
	Temperature units.Temperature
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
	RSSI           int // Wireless signal strength in dB at last reading
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
	RSSI           int                 `json:"rssi"`
	// TODO: tags
	Type string `json:"type"`
}

type sensorsResponse map[string]sensorResponse

type SensorSlice []*Sensor

func (s SensorSlice) Len() int {
	return len(s)
}

func (s SensorSlice) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s SensorSlice) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

// List returns the sensors matching the active criteria in alphabetical order
func (s *SensorService) List(ctx context.Context, active bool) (SensorSlice, error) {
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

	sensors := make(SensorSlice, 0, len(ssresp))
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
				Humidity: AlertHumidity{
					Enabled: a.Humidity.Enabled,
					Max:     units.NewPercentageFloat32(a.Humidity.Max),
					Min:     units.NewPercentageFloat32(a.Humidity.Min),
				},
				Temperature: AlertTemperature{
					Enabled: a.Temperature.Enabled,
					Max:     units.NewTemperatureF(a.Temperature.Max),
					Min:     units.NewTemperatureF(a.Temperature.Min),
				},
			},
			BatteryVoltage: sresp.BatteryVoltage,
			Calibration: Calibration{
				Humidity:    units.NewPercentageFloat32(c.Humidity),
				Temperature: units.NewTemperatureF(c.Temperature),
			},
			DeviceID: sresp.DeviceID,
			ID:       sresp.ID,
			Name:     sresp.Name,
			RSSI:     sresp.RSSI,
			Type:     newSensorType(sresp.Type),
		}
		sensors = append(sensors, s)
	}

	sort.Sort(sensors)

	return sensors, nil
}

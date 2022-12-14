package sensorpush

import (
	"sort"

	"github.com/rconradharris/go-sensorpush/units"
)

type SensorID string

func NewSensorID(id string) SensorID {
	return SensorID(id)
}

func (s SensorID) String() string {
	return string(s)
}

type DeviceID string

func NewDeviceID(id string) DeviceID {
	return DeviceID(id)
}

func (id DeviceID) String() string {
	return string(id)
}

type Sensor struct {
	Active         bool
	Address        string // MAC address
	Alerts         Alerts
	BatteryVoltage *units.Voltage
	Calibration    Calibration
	DeviceID       DeviceID
	ID             SensorID
	Name           string
	RSSI           *units.SignalStrength // strength at last reading
	// TODO: tags
	Type SensorType
}

type SensorMap map[SensorID]*Sensor

// SensorsAlpha returns the Sensors sorted alphabetically by Name
func (sm SensorMap) SensorsAlpha() []*Sensor {
	ss := make(sensorSlice, 0, len(sm))
	for _, s := range sm {
		ss = append(ss, s)
	}
	sort.Sort(ss)
	return ss
}

type sensorSlice []*Sensor

func (s sensorSlice) Len() int {
	return len(s)
}

func (s sensorSlice) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s sensorSlice) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

type sensorsRequest struct {
	Active *bool `json:"active,omitempty"`
}

type sensorResponse struct {
	Active         bool                `json:"active"`
	Address        string              `json:"address"`
	Alerts         alertsResponse      `json:"alerts"`
	BatteryVoltage *float32            `json:"battery_voltage"`
	Calibration    calibrationResponse `json:"calibration"`
	DeviceID       string              `json:"deviceId"`
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	RSSI           *float32            `json:"rssi"`
	// TODO: tags
	Type string `json:"type"`
}

type sensorsResponse map[string]sensorResponse

func newSensor(sresp sensorResponse) *Sensor {
	a := sresp.Alerts
	c := sresp.Calibration
	s := &Sensor{
		Active:  sresp.Active,
		Address: sresp.Address,
		Alerts: Alerts{
			Humidity: AlertHumidity{
				Enabled: a.Humidity.Enabled,
			},
			Temperature: AlertTemperature{
				Enabled: a.Temperature.Enabled,
			},
		},
		Calibration: Calibration{
			HumidityDelta:    units.NewHumidityDelta(c.Humidity),
			TemperatureDelta: units.NewTemperatureDeltaF(c.Temperature),
		},
		DeviceID: NewDeviceID(sresp.DeviceID),
		ID:       NewSensorID(sresp.ID),
		Name:     sresp.Name,
		Type:     newSensorType(sresp.Type),
	}

	// Handle nullable fields

	// Humidity Alerts
	if v := a.Humidity.Max; v != nil {
		h := units.NewHumidity(*v)
		s.Alerts.Humidity.Max = &h
	}

	if v := a.Humidity.Min; v != nil {
		h := units.NewHumidity(*v)
		s.Alerts.Humidity.Min = &h
	}

	// Temperature Alerts
	if v := a.Temperature.Max; v != nil {
		t := units.NewTemperatureF(*v)
		s.Alerts.Temperature.Max = &t
	}

	if v := a.Temperature.Min; v != nil {
		t := units.NewTemperatureF(*v)
		s.Alerts.Temperature.Min = &t
	}

	// Battery Voltage
	if v := sresp.BatteryVoltage; v != nil {
		bv := units.NewVoltage(*v)
		s.BatteryVoltage = &bv
	}

	// RSSI
	if v := sresp.RSSI; v != nil {
		ss := units.NewSignalStrength(*v)
		s.RSSI = &ss
	}

	return s
}

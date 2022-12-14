package sensorpush

import (
	"time"

	"github.com/rconradharris/go-sensorpush/units"
)

type Sample struct {
	Altitude           *units.Distance
	BarometricPressure *units.Pressure
	DewPoint           *units.Temperature
	Humidity           *units.Humidity
	Observed           time.Time
	// TODO: Tags
	Temperature *units.Temperature
	VPD         *units.Pressure
}

type SampleSlice []*Sample

func (s SampleSlice) Len() int {
	return len(s)
}

func (s SampleSlice) Less(i, j int) bool {
	// Reverse chron
	//
	// The API appears to sort for us, but we sort ourselves to ensure
	// the sort is well-defined and stable
	return s[i].Observed.After(s[j].Observed)
}

func (s SampleSlice) Swap(i, j int) {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp
}

func newSample(sr sampleResponse) (*Sample, error) {
	s := &Sample{}

	// Altitude
	if sr.Altitude != nil {
		alt := units.NewDistanceFT(*sr.Altitude)
		s.Altitude = &alt
	}

	// Barometric Pressure
	if sr.BarometricPressure != nil {
		baro := units.NewPressureINHG(*sr.BarometricPressure)
		s.BarometricPressure = &baro
	}

	// Dew Point
	if sr.DewPoint != nil {
		dew := units.NewTemperatureF(*sr.DewPoint)
		s.DewPoint = &dew
	}

	// Humidity
	if sr.Humidity != nil {
		hum := units.NewHumidity(*sr.Humidity)
		s.Humidity = &hum
	}

	// Observed
	t, err := parseTime(sr.Observed)
	if err != nil {
		return s, err
	}
	s.Observed = t

	// Temperature
	if sr.Temperature != nil {
		temp := units.NewTemperatureF(*sr.Temperature)
		s.Temperature = &temp
	}

	// VPD
	if sr.VPD != nil {
		vpd := units.NewPressureKPA(*sr.VPD)
		s.VPD = &vpd
	}

	return s, nil
}

type sampleResponse struct {
	Altitude           *float32 `json:"altitude"`
	BarometricPressure *float32 `json:"barometric_pressure"`
	DewPoint           *float32 `json:"dewpoint"`
	Humidity           *float32 `json:"humidity"`
	Observed           string   `json:"observed"`
	Temperature        *float32 `json:"temperature"`
	VPD                *float32 `json:"vpd"`
}

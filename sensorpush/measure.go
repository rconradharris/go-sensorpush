package sensorpush

import (
	"fmt"
)

type Measure int

const (
	MeasureBarometricPressure Measure = iota
	MeasureDewPoint
	MeasureHumidity
	MeasureTemperature
	MeasureVPD
)

func (m Measure) String() string {
	switch m {
	case MeasureBarometricPressure:
		return "barometric_pressure"
	case MeasureDewPoint:
		return "dewpoint"
	case MeasureHumidity:
		return "humidity"
	case MeasureTemperature:
		return "temperature"
	case MeasureVPD:
		return "vpd"
	}
	return "unknown"
}

func ParseMeasure(s string) (Measure, error) {
	switch s {
	case "barometric_pressure":
		return MeasureBarometricPressure, nil
	case "dewpoint":
		return MeasureDewPoint, nil
	case "humidity":
		return MeasureHumidity, nil
	case "temperature":
		return MeasureTemperature, nil
	case "vpd":
		return MeasureVPD, nil
	}
	return MeasureBarometricPressure, fmt.Errorf("unknown measure '%s'", s)
}

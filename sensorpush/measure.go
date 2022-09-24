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
	case "baro":
		return MeasureBarometricPressure, nil
	case "dew":
		return MeasureDewPoint, nil
	case "hum":
		return MeasureHumidity, nil
	case "temp":
		return MeasureTemperature, nil
	case "vpd":
		return MeasureVPD, nil
	}
	return MeasureBarometricPressure, fmt.Errorf("unknown measure '%s'", s)
}

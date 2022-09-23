package sensorpush

import (
	"github.com/rconradharris/go-sensorpush/units"
)

type Calibration struct {
	HumidityDelta    units.HumidityDelta
	TemperatureDelta units.TemperatureDelta
}

type calibrationResponse struct {
	Humidity    float32 `json:"humidity"`
	Temperature float32 `json:"temperature"`
}

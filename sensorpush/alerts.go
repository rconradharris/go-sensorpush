package sensorpush

import (
	"github.com/rconradharris/go-sensorpush/units"
)

type AlertHumidity struct {
	Enabled bool
	Max     *units.Humidity
	Min     *units.Humidity
}

type AlertTemperature struct {
	Enabled bool
	Max     *units.Temperature
	Min     *units.Temperature
}

type Alerts struct {
	Humidity    AlertHumidity
	Temperature AlertTemperature
}

type alertResponse struct {
	Enabled bool     `json:"enabled"`
	Max     *float32 `json:"max"`
	Min     *float32 `json:"min"`
}

type alertsResponse struct {
	Humidity    alertResponse `json:"humidity"`
	Temperature alertResponse `json:"temperature"`
}

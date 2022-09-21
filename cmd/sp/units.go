package main

import (
	"fmt"

	"github.com/rconradharris/go-sensorpush/units"
)

type temperatureUnit int

func newTemperatureUnit(s string) (temperatureUnit, error) {
	switch s {
	case "f":
		return tempFahrenheit, nil
	case "c":
		return tempCelsius, nil
	}
	var t0 temperatureUnit
	return t0, fmt.Errorf("unknown temp unit specifier: %s", s)
}

const (
	tempFahrenheit temperatureUnit = iota
	tempCelsius
)

type unitsCfg struct {
	temperature temperatureUnit
}

type unitsFormatter struct {
	cfg unitsCfg
}

func newUnitsFormatter(cfg unitsCfg) *unitsFormatter {
	return &unitsFormatter{cfg: cfg}
}

func (f *unitsFormatter) Temperature(t units.Temperature) string {
	var v float32
	switch f.cfg.temperature {
	case tempFahrenheit:
		v = t.F()
	default:
		v = t.C()
	}
	return fmt.Sprintf("%.2f", v)
}

func (f *unitsFormatter) TemperatureDelta(t units.TemperatureDelta) string {
	var v float32
	switch f.cfg.temperature {
	case tempFahrenheit:
		v = t.F()
	default:
		v = t.C()
	}
	return fmt.Sprintf("%.2f", v)
}

func fmtHumidity(p units.Percentage) string {
	return fmt.Sprintf("%.2f", p.Norm())
}

func fmtSignalStrength(v int) string {
	return fmt.Sprintf("%d", v)
}

func fmtVoltage(v float32) string {
	return fmt.Sprintf("%.2f", v)
}

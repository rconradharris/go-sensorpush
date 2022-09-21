package main

import (
	"fmt"
	"strings"

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

func (u temperatureUnit) String() string {
	switch u {
	case tempCelsius:
		return "c"
	}
	return "f"
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
	case tempCelsius:
		v = t.C()
	default:
		v = t.F()
	}

	deg := strings.ToUpper(f.cfg.temperature.String())
	return fmt.Sprintf("%.1f °%s", v, deg)
}

func (f *unitsFormatter) TemperatureDelta(t units.TemperatureDelta) string {
	var v float32
	switch f.cfg.temperature {
	case tempFahrenheit:
		v = t.F()
	default:
		v = t.C()
	}

	sign := "+"
	if v < 0 {
		sign = "-"
	}
	deg := strings.ToUpper(f.cfg.temperature.String())
	return fmt.Sprintf("%s%.1f °%s", sign, v, deg)
}

func (f *unitsFormatter) Humidity(p units.Percentage) string {
	return fmt.Sprintf("%.1f", p.Norm())
}

func (f *unitsFormatter) HumidityDelta(p units.Percentage) string {
	sign := "+"
	if p.Norm() < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%.1f", sign, p.Norm())
}

func (f *unitsFormatter) SignalStrength(v int) string {
	return fmt.Sprintf("%d", v)
}

func (f *unitsFormatter) Voltage(v float32) string {
	return fmt.Sprintf("%.2f", v)
}

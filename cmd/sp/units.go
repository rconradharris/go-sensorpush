package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/rconradharris/go-sensorpush/units"
)

const (
	notAvail = "N/A"
)

type unitFlags struct {
	temp string
}

func addUnitFlags(fs *flag.FlagSet, f *unitFlags) {
	fs.StringVar(&f.temp, "temp", "f", "fahrenheit (\"f\") or celsius (\"c\")")
}

func newUnitsFormatter(uf *unitFlags) (*unitsFormatter, error) {
	cfg := unitsCfg{}

	if uf != nil {
		tempU, err := newTemperatureUnit(uf.temp)
		if err != nil {
			return nil, err
		}
		cfg.temperature = tempU
	}

	return &unitsFormatter{cfg: cfg}, nil
}

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

func (f *unitsFormatter) Temperature(t *units.Temperature) string {
	if t == nil {
		return notAvail
	}
	var v float32
	switch f.cfg.temperature {
	case tempCelsius:
		v = t.C()
	default:
		v = t.F()
	}

	deg := strings.ToUpper(f.cfg.temperature.String())
	return fmt.Sprintf("%.1f°%s", v, deg)
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
	return fmt.Sprintf("%s%.1f°%s", sign, v, deg)
}

func (f *unitsFormatter) Humidity(h *units.Humidity) string {
	if h == nil {
		return notAvail
	}
	return fmt.Sprintf("%.1f%%", h.Pct())
}

func (f *unitsFormatter) HumidityDelta(h units.HumidityDelta) string {
	v := h.Pct()
	sign := "+"
	if v < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%.1f%%", sign, v)
}

func (f *unitsFormatter) SignalStrength(v *units.SignalStrength) string {
	if v == nil {
		return notAvail
	}
	return fmt.Sprintf("%.0f dB", v.DB())
}

func (f *unitsFormatter) Time(t time.Time) string {
	return t.UTC().Format(time.RFC1123)
}

func (f *unitsFormatter) Voltage(v *units.Voltage) string {
	if v == nil {
		return notAvail
	}
	return fmt.Sprintf("%.1f V", v.V())
}

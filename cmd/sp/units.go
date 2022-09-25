package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/rconradharris/go-sensorpush/units"
)

const (
	notAvail = "N/A"
)

type unitFlags struct {
	baro string
	dist string
	temp string
	time string
}

func addUnitFlags(fs *flag.FlagSet, f *unitFlags) {
	fs.StringVar(&f.baro, "baro", "inhg", "inHg (\"inhg\") or millibars (\"mb\")")
	fs.StringVar(&f.dist, "dist", "ft", "feet (\"ft\") or meters (\"m\")")
	fs.StringVar(&f.temp, "temp", "f", "fahrenheit (\"f\") or celsius (\"c\")")
	fs.StringVar(&f.time, "time", "iso", "formats: human|iso")
}

type unitsFormatter struct {
	baroU units.PressureUnit
	distU units.DistanceUnit
	tempU units.TemperatureUnit
	timeF timeFormat
}

type timeFormat int

const (
	timeFormatHuman timeFormat = iota
	timeFormatISO
)

func parseTimeFormat(s string) (timeFormat, error) {
	var t0 timeFormat
	switch s {
	case "human":
		return timeFormatHuman, nil
	case "iso":
		return timeFormatISO, nil
	}
	return t0, fmt.Errorf("unknown time unit specifier: %s", s)
}

func (t timeFormat) String() string {
	switch t {
	case timeFormatHuman:
		return "human"
	case timeFormatISO:
		return "iso"
	}
	return "unknown"
}

func newUnitsFormatter(uf *unitFlags) (*unitsFormatter, error) {
	fmtU := &unitsFormatter{}

	if uf != nil {
		baroU, err := units.ParsePressureUnit(uf.baro)
		if err != nil {
			return nil, err
		}
		fmtU.baroU = baroU

		distU, err := units.ParseDistanceUnit(uf.dist)
		if err != nil {
			return nil, err
		}
		fmtU.distU = distU

		tempU, err := units.ParseTemperatureUnit(uf.temp)
		if err != nil {
			return nil, err
		}
		fmtU.tempU = tempU

		timeF, err := parseTimeFormat(uf.time)
		if err != nil {
			return nil, err
		}
		fmtU.timeF = timeF
	}

	return fmtU, nil
}

func (f *unitsFormatter) BarometricPressure(p *units.Pressure) string {
	if p == nil {
		return notAvail
	}

	var v float32
	switch f.baroU {
	case units.PressureUnitMB:
		v = p.MB()
	default:
		v = p.INHG()
	}

	unit := f.baroU.String()
	return fmt.Sprintf("%.1f%s", v, unit)
}

func (f *unitsFormatter) Distance(d *units.Distance) string {
	if d == nil {
		return notAvail
	}

	var v float32
	switch f.distU {
	case units.DistanceUnitM:
		v = d.M()
	default:
		v = d.FT()
	}

	unit := f.distU.String()
	return fmt.Sprintf("%.1f%s", v, unit)
}

func (f *unitsFormatter) Temperature(t *units.Temperature) string {
	if t == nil {
		return notAvail
	}

	var v float32
	switch f.tempU {
	case units.TemperatureUnitC:
		v = t.C()
	default:
		v = t.F()
	}

	unit := f.tempU.String()
	return fmt.Sprintf("%.1f°%s", v, unit)
}

func (f *unitsFormatter) TemperatureDelta(t units.TemperatureDelta) string {
	var v float32
	switch f.tempU {
	case units.TemperatureUnitC:
		v = t.C()
	default:
		v = t.F()
	}

	sign := "+"
	if v < 0 {
		sign = "-"
	}

	unit := f.tempU.String()
	return fmt.Sprintf("%s%.1f°%s", sign, v, unit)
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
	if t.IsZero() {
		return notAvail
	}

	switch f.timeF {
	case timeFormatISO:
		return f.ISOTime(t)
	}

	return f.HumanTime(t)
}

func (f *unitsFormatter) ISOTime(t time.Time) string {
	if t.IsZero() {
		return notAvail
	}
	return t.UTC().Format(time.RFC3339)
}

func (f *unitsFormatter) HumanTime(t time.Time) string {
	if t.IsZero() {
		return notAvail
	}
	return t.Local().Format(time.RFC822)
}

func (f *unitsFormatter) Voltage(v *units.Voltage) string {
	if v == nil {
		return notAvail
	}
	return fmt.Sprintf("%.1f V", v.V())
}

func (f *unitsFormatter) VPD(p *units.Pressure) string {
	if p == nil {
		return notAvail
	}

	return fmt.Sprintf("%.1f%s", p.KPA(), "kPa")
}

func (f *unitsFormatter) Int(v int) string {
	return fmt.Sprintf("%d", v)
}

func (f *unitsFormatter) Bool(v bool) string {
	if v {
		return "y"
	}
	return "n"
}

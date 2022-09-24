package units

import (
	"fmt"
)

type TemperatureUnit int

const (
	TemperatureUnitF TemperatureUnit = iota
	TemperatureUnitC
)

func ParseTemperatureUnit(s string) (TemperatureUnit, error) {
	var t0 TemperatureUnit
	switch s {
	case "f":
		return TemperatureUnitF, nil
	case "c":
		return TemperatureUnitC, nil
	}
	return t0, fmt.Errorf("unknown temperature unit specifier: %s", s)
}

func (u TemperatureUnit) String() string {
	switch u {
	case TemperatureUnitC:
		return "C"
	}
	return "F"
}

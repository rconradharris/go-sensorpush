package units

import (
	"fmt"
)

type PressureUnit int

const (
	PressureUnitINHG PressureUnit = iota
	PressureUnitMB
	PressureUnitKPA
)

func ParsePressureUnit(s string) (PressureUnit, error) {
	var p0 PressureUnit
	switch s {
	case "inhg":
		return PressureUnitINHG, nil
	case "kpa":
		return PressureUnitKPA, nil
	case "mb":
		return PressureUnitMB, nil
	}
	return p0, fmt.Errorf("unknown pressure unit specifier: %s", s)
}

func (u PressureUnit) String() string {
	switch u {
	case PressureUnitKPA:
		return "kPa"
	case PressureUnitMB:
		return "mb"
	}
	return "inHg"
}

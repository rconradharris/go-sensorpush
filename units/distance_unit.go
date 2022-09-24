package units

import (
	"fmt"
)

type DistanceUnit int

const (
	DistanceUnitFT DistanceUnit = iota
	DistanceUnitM
)

func ParseDistanceUnit(s string) (DistanceUnit, error) {
	var d0 DistanceUnit
	switch s {
	case "ft":
		return DistanceUnitFT, nil
	case "m":
		return DistanceUnitM, nil
	}
	return d0, fmt.Errorf("unknown distance unit specifier: %s", s)
}

func (u DistanceUnit) String() string {
	switch u {
	case DistanceUnitM:
		return "m"
	}
	return "ft"
}

package main

import (
	"fmt"
	"strings"

	sp "github.com/rconradharris/go-sensorpush/sensorpush"
)

const (
	fmtStrAttrVal   = "%s%-20s: %s\n"
	attrValRightPad = 20
	indentUnit      = "  "
)

func fmtAttrValHeading(b *strings.Builder, heading string, indent int) {
	indentS := strings.Repeat(indentUnit, indent)
	fmt.Fprintf(b, "%s%s\n", indentS, heading)
}

func fmtAttrVal(b *strings.Builder, attr, val string, indent int) {
	indentS := strings.Repeat(indentUnit, indent)
	padlen := attrValRightPad - len(indentS) // Make ':' line up even for indents
	paddedAttr := padRight(attr, " ", padlen)
	fmt.Fprintf(b, "%s%s: %s\n", indentS, paddedAttr, val)
}

func padRight(s, padUnit string, length int) string {
	padN := 1 + ((length - len(padUnit)) / len(padUnit))
	ps := s + strings.Repeat(padUnit, padN)
	return ps[:length]
}

// Returns sensor that matches:
//
// 1. Long ID exact match
// 2. Short ID exact match
// 3. Case-insensitive name
//
// Returns nil if no match is found
func findSensorByNameOrID(sm sp.SensorMap, nameOrID string) *sp.Sensor {
	id := sp.NewSensorID(nameOrID)
	if s, ok := sm[id]; ok {
		return s
	}

	lowerName := strings.ToLower(nameOrID)
	for _, s := range sm {
		if s.DeviceID.String() == nameOrID {
			return s
		}
		if strings.ToLower(s.Name) == lowerName {
			return s
		}
	}
	return nil
}

package main

import (
	"fmt"
	"strings"

	"github.com/rconradharris/go-sensorpush/units"
)

const (
	fmtStrAttrVal   = "%s%-20s: %s\n"
	attrValRightPad = 20
	indentUnit      = "  "
)

func fmtBool(b bool) string {
	if b {
		return "y"
	}
	return "n"
}

func fmtHumidity(v float32) string {
	return fmt.Sprintf("%.2f", v)
}

func fmtVoltage(v float32) string {
	return fmt.Sprintf("%.2f", v)
}

func fmtSignalStrength(v int) string {
	return fmt.Sprintf("%d", v)
}

func fmtTemperature(t units.Temperature) string {
	return fmt.Sprintf("%.2f", t.F())
}

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

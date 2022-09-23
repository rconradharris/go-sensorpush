package sensorpush

import (
	"fmt"
)

type SampleFormat int

const (
	SampleFormatJSON SampleFormat = iota
	SampleFormatCSV
)

func (f SampleFormat) String() string {
	switch f {
	case SampleFormatCSV:
		return "csv"
	}
	return "json"
}

func ParseSampleFormat(s string) (SampleFormat, error) {
	switch s {
	case "csv":
		return SampleFormatCSV, nil
	case "json":
		return SampleFormatJSON, nil
	}
	return SampleFormatJSON, fmt.Errorf("unknown sample format '%s'", s)
}

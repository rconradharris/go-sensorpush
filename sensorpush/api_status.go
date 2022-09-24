package sensorpush

type APIStatus int

const (
	APIStatusUnknown APIStatus = iota
	APIStatusOK
)

func (s APIStatus) String() string {
	switch s {
	case APIStatusOK:
		return "ok"
	}
	return "unknown"
}

func newAPIStatus(s string) APIStatus {
	switch s {
	case "ok":
		return APIStatusOK
	}
	return APIStatusUnknown
}

package sensorpush

type StatusEnum int

const (
	StatusUnknown StatusEnum = iota
	StatusOK
)

func (s StatusEnum) String() string {
	switch s {
	case StatusOK:
		return "ok"
	}
	return "unknown"
}

func newStatusEnum(s string) StatusEnum {
	switch s {
	case "ok":
		return StatusOK
	}
	return StatusUnknown
}

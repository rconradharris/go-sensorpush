package sensorpush

type SampleStatus int

const (
	SampleStatusUnknown SampleStatus = iota
	SampleStatusOK
)

func (s SampleStatus) String() string {
	switch s {
	case SampleStatusOK:
		return "ok"
	}
	return "unknown"
}

func newSampleStatus(s string) SampleStatus {
	switch s {
	case "OK":
		return SampleStatusOK
	}
	return SampleStatusUnknown
}

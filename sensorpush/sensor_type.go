package sensorpush

type SensorType int

const (
	SensorTypeUnknown SensorType = iota
	SensorTypeHT1
	SensorTypeHTw
	SensorTypeHTPxw
)

func (s SensorType) String() string {
	switch s {
	case SensorTypeHT1:
		return "HT1"
	case SensorTypeHTw:
		return "HT.w"
	case SensorTypeHTPxw:
		return "HTP.xw"
	}
	return "unknown-sensor-type"
}

func newSensorType(s string) SensorType {
	switch s {
	case "HT1":
		return SensorTypeHT1
	case "HT.w":
		return SensorTypeHTw
	case "HTP.xw":
		return SensorTypeHTPxw
	}
	return SensorTypeUnknown
}

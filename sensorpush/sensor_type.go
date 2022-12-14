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

func (s SensorType) Features() SensorFeatureSet {
	fs := newSensorFeatureSet(
		SensorFeatureHumidity,
		SensorFeatureTemperature,
	)
	switch s {
	case SensorTypeHT1:
	case SensorTypeHTw:
		fs.add(SensorFeatureDewPoint)
		fs.add(SensorFeatureVPD)
	case SensorTypeHTPxw:
		fs.add(SensorFeatureBarometricPressure)
		fs.add(SensorFeatureDewPoint)
		fs.add(SensorFeatureVPD)
	}
	return fs
}

package sensorpush

type SensorFeature int

const (
	SensorFeatureNull SensorFeature = iota
	SensorFeatureBarometricPressure
	SensorFeatureDewPoint
	SensorFeatureHumidity
	SensorFeatureTemperature
	SensorFeatureVPD
)

type SensorFeatureSet map[SensorFeature]struct{}

func newSensorFeatureSet(feats ...SensorFeature) SensorFeatureSet {
	fs := SensorFeatureSet{}
	fs[SensorFeatureNull] = struct{}{} // Null feature available everywhere
	for _, f := range feats {

		fs[f] = struct{}{}
	}
	return fs
}

func (fs SensorFeatureSet) add(f SensorFeature) {
	fs[f] = struct{}{}
}

func (fs SensorFeatureSet) Has(f SensorFeature) bool {
	_, ok := fs[f]
	return ok
}

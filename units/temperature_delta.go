package units

type TemperatureDelta float32 // stored natively as F

// NewTemperatureDeltaF creates a TemperatureDelta from a Fahrenheit value
func NewTemperatureDeltaF(f float32) TemperatureDelta {
	return TemperatureDelta(f)
}

// NewTemperatureDeltaC creates a TemperatureDelta from a Celsius value
func NewTemperatureDeltaC(f float32) TemperatureDelta {
	return TemperatureDelta(f * (9.0 / 5.0))
}

// F returns a temperature delta in Fahrenheit
func (d TemperatureDelta) F() float32 {
	return float32(d)
}

// C returns the temperature delta in Celsius
func (d TemperatureDelta) C() float32 {
	return float32(d) * (5.0 / 9.0)
}

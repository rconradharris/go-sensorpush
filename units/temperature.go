package units

type Temperature float32 // stored natively as F

// NewTemperatureF creates a Temperature from a Fahrenheit value
func NewTemperatureF(f float32) Temperature {
	return Temperature(f)
}

// F returns the temperature in Fahrenheit
func (t Temperature) F() float32 {
	return float32(t)
}

// NewTemperatureC creates a Temperature from a Celsius value
func NewTemperatureC(c float32) Temperature {
	return Temperature(cToF(c))
}

// C returns the temperature in Celsius
func (t Temperature) C() float32 {
	return fToC(float32(t))
}

func fToC(f float32) float32 {
	return (f - 32.0) * (5.0 / 9.0)
}

func cToF(c float32) float32 {
	return (9.0/5.0)*c + 32.0
}

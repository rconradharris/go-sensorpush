package units

type Humidity float32 // stored natively as a 0-100 %

// NewHumidity returns a new Humidity
func NewHumidity(h float32) Humidity {
	return Humidity(h)
}

// Pct returns a nominal range of [0, 100]
func (h Humidity) Pct() float32 {
	return float32(h)
}

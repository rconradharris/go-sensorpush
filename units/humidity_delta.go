package units

type HumidityDelta float32 // stored natively as a %

// NewHumidityDelta returns a new HumidityDelta
func NewHumidityDelta(h float32) HumidityDelta {
	return HumidityDelta(h)
}

// Pct returns a nominal range of [0, 100]
func (h HumidityDelta) Pct() float32 {
	return float32(h)
}

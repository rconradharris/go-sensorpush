package units

type SignalStrength float32 // stored natively in dB

// NewSignalStrength returns a new SignalStrength
func NewSignalStrength(v float32) SignalStrength {
	return SignalStrength(v)
}

// Returns the SignalStrength in dB
func (h SignalStrength) DB() float32 {
	return float32(h)
}

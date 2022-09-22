package units

type Voltage float32 // stored natively in Volts

// NewVoltage returns a new Voltage
func NewVoltage(v float32) Voltage {
	return Voltage(v)
}

// Returns the voltage in Volts
func (h Voltage) V() float32 {
	return float32(h)
}

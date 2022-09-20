package units

type Percentage float32

// NewPercentageInt returns a Percentage specified by an int in a nominal range [0, 100]
func NewPercentageInt(p int) Percentage {
	return Percentage(float32(p) / 100.0)
}

// NewPercentageFloat32 returns a Percentage specified by a float32 in a nominal range [0.0, 1.0]
func NewPercentageFloat32(p float32) Percentage {
	return Percentage(p)

}

// Norm returns nominal range of [0.0, 1.0]
func (p Percentage) Norm() float32 {
	return float32(p)
}

// Pct returns a nominal range of [0, 100]
func (p Percentage) Pct() int {
	return int(p * 100.0)
}

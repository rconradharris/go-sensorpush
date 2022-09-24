package units

type Distance float32 // stored natively as ft

// NewDistanceFT creates a Distance from a ft value
func NewDistanceFT(ft float32) Distance {
	return Distance(ft)
}

// FT returns the distance in feet
func (d Distance) FT() float32 {
	return float32(d)
}

// M returns the distance in meters
func (d Distance) M() float32 {
	return ftToM(float32(d))
}

func ftToM(ft float32) float32 {
	return ft / 3.281
}

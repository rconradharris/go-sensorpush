package units

import (
	"fmt"
)

type Pressure float32 // stored natively as inHg

// NewPressureINHG creates a Pressure from a inHg value
func NewPressureINHG(inHg float32) Pressure {
	return Pressure(inHg)
}

// NewPressureKPA creates a Pressure from a kPa value
func NewPressureKPA(kpa float32) Pressure {
	return Pressure(kpaToINHG(kpa))
}

// INHG returns the Pressure in inHg
func (p Pressure) INHG() float32 {
	return float32(p)
}

// KPA returns the Pressure in kilopascals
func (p Pressure) KPA() float32 {
	return inHgToKPA(float32(p))
}

// MB returns the Pressure in millibars
func (p Pressure) MB() float32 {
	return inHgToMB(float32(p))
}

func (p Pressure) String() string {
	return fmt.Sprintf("%.1finHg", p.INHG())
}

func inHgToKPA(inHg float32) float32 {
	return inHg * 3.386
}

func kpaToINHG(kpa float32) float32 {
	return kpa / 3.386
}

func inHgToMB(inHg float32) float32 {
	return inHgToKPA(inHg) * 10.0
}

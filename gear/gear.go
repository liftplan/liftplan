package gear

import (
	"errors"
	"fmt"
)

var (
	// ErrInputLessThanBar is an error message used to warn that a Rounding call requests
	// a value that is less than the weight of the empty bar, itself.
	ErrInputLessThanBar = errors.New("input weight is less than the minimum weight of the bar")
	ErrInvalidUnitGear  = errors.New("invalid unit: gear")
)

// Gear is a struct that represents the weight inputs for the Bar, Plates, and desired unites
// for the Gear to be measured in. Gear is primary used to calculate possible weight totals
// from a requested weight amount, and is meant to be used to convert a lift from a calculated
// percentage into the greatest possible incremental weight with the equipment provided.
type Gear struct {
	Bar    Bar    `json:"bar"`
	Plates Plates `json:"plates"`
	Unit   Unit   `json:"unit"`
}

// Default is a helper function to make it easy to setup default
// gear.
func Default(unit Unit) Gear {
	switch unit {
	case LBS:
		return Gear{
			Bar:    MensBarLBS,
			Plates: DefaultPlatesLBS,
			Unit:   LBS,
		}
	case KG:
		return Gear{
			Bar:    MensBarKG,
			Plates: DefaultPlatesKG,
			Unit:   KG,
		}
	default:
		return Gear{}
	}
}

// String outputs the string format for Gear.
func (g Gear) String() string {
	return fmt.Sprintf("Unit: %v, Bar: { %v }, Plates: { %v }", g.Unit, g.Bar, g.Plates)
}

// Min returns the minimum amount allowed for rounding. This is based
// on the bar weight converted to the Gear Units.
func (g Gear) Min() (float64, error) {
	return g.Bar.ConvertTo(g.Unit)
}

// Valid checks gear, plates, and bar for valid units and
// valid valudes for Plates and Bar
func (g Gear) Valid() error {
	if !g.Unit.Valid() {
		return ErrInvalidUnitGear
	}
	if err := g.Bar.Valid(); err != nil {
		return err
	}
	return g.Plates.Valid()
}

// Round takes a float64 number and returns the rounded total
// based on the bar and incremental plate weights, converted to
// the gear units. If the bar and incremental plates are in KG
// but the Gear units are LBS, the float64 will be returned in gear units
// of LBS.
func (g Gear) Round(weight float64) (float64, error) {
	if err := g.Valid(); err != nil {
		return 0, err
	}
	bar, plates, err := g.barFromWeight(weight)
	if err != nil {
		return 0, err
	}
	if weight == bar {
		return weight, nil
	}
	p, _ := ConvertFromTo(plates, g.Unit, g.Plates.Unit)
	pr, _ := g.Plates.Round(p)
	r, _ := ConvertFromTo(pr, g.Plates.Unit, g.Unit)
	return bar + r, nil
}

// barFromWeight takes a weight and returns the bar and plate weight
// in the units of Gear or returns an error.
func (g Gear) barFromWeight(weight float64) (b, p float64, err error) {
	bar, _ := g.Min()
	if weight-bar < -.0001 {
		return b, p, ErrInputLessThanBar
	}
	return bar, weight - bar, nil
}

// Recommend takes
func (g Gear) Recommend(weight float64) ([]float64, error) {
	if err := g.Valid(); err != nil {
		return nil, err
	}
	_, p, err := g.barFromWeight(weight)
	if err != nil {
		return nil, err
	}
	p, _ = ConvertFromTo(p, g.Unit, g.Plates.Unit)
	return Recommend(p, g.Plates.Weights)
}

// Equals checks for deep equality betten a comparable gear struct and itself
// and returns a boolean value.
func (g Gear) Equals(c Gear) bool {
	return g.Bar.Equals(c.Bar) && g.Plates.Equals(c.Plates) && g.Unit == c.Unit
}

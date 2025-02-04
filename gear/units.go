package gear

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	// ConversionFactorKGtoLBS is the weight of 1 KG in LBS
	ConversionFactorKGtoLBS float64 = 2.20462262185
)

var (
	// ErrInvalidUnit error are for any uint that is not represented by KG or LBS
	ErrInvalidUnit = errors.New("invalid unit")
)

// Unit is used to enumerate units for weights
type Unit uint

const (
	// KG represents the metric unit Kilograms
	KG Unit = iota
	// LBS represents the imperial unit for pounds
	LBS
	// TODO (NST): maybe add support for POOD?
)

// String prints the human readable value from the enum
func (u Unit) String() string {
	p := []string{"KG", "LBS"}
	if int(u) < len(p) {
		return p[u]
	}
	return ""
}

// StringToUnit is a simple map that converts from string to Unit type
var stringToUnit = map[string]Unit{
	"KG":  KG,
	"LBS": LBS,
}

// UnitFromString takes a string and returns a unit or error
func UnitFromString(s string) (Unit, error) {
	unit, ok := stringToUnit[s]
	if !ok {
		return 0, ErrInvalidUnit
	}
	return unit, nil
}

// MarshalJSON is used for human readable json.
func (u Unit) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, u)), nil
}

// UnmarshalJSON is used to convert human readable json to a Unit type
func (u *Unit) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	unit, err := UnitFromString(s)
	if err != nil {
		return err
	}
	*u = unit
	return nil
}

// Valid checks that a Unit is either KG or LBS and returns a boolean
func (u Unit) Valid() bool {
	return u == KG || u == LBS
}

// ConvertFromTo takes a float64 and a from and to unit and converts the float into the
// requested unit. For instance if you wanted to convert 55 lbs to kg,
// Convert(55.0, LBS, KG) would return 24.9476, nil.
func ConvertFromTo(w float64, from Unit, to Unit) (float64, error) {
	if !from.Valid() || !to.Valid() {
		return 0, ErrInvalidUnit
	}
	if from == to {
		return w, nil
	}
	if to == KG {
		return w / ConversionFactorKGtoLBS, nil
	}
	if to == LBS {
		return w * ConversionFactorKGtoLBS, nil
	}
	return 0, ErrInvalidUnit
}

package gear

import "fmt"

var (
	// MensBarKG represents a standard men's 20 KG barbell
	MensBarKG = Bar{
		Weight: 20.0,
		Unit:   KG,
	}
	// MensBarLBS represents a standard men's 45 LBS barbell
	MensBarLBS = Bar{
		Weight: 45.0,
		Unit:   LBS,
	}
	// WomensBarKG represents a standard women's 15 KG barbell
	WomensBarKG = Bar{
		Weight: 15.0,
		Unit:   KG,
	}
	// WomensBarLBS represents a standard women's 35 LBS barbell
	WomensBarLBS = Bar{
		Weight: 35.0,
		Unit:   LBS,
	}
)

// Bar contains the weight and units of a Barbell
type Bar struct {
	Weight float64 `json:"weight"`
	Unit   Unit    `json:"unit"`
}

// ConvertTo takes a Unit and returns the converted weight or an error.
// If the bar is already in requested unit, it simply returns the weight.
func (b Bar) ConvertTo(u Unit) (float64, error) {
	return ConvertFromTo(b.Weight, b.Unit, u)
}

// String prints the human readable string format of the bar data structure
func (b Bar) String() string {
	return fmt.Sprintf("Weight: %v, Unit: %v", b.Weight, b.Unit)
}

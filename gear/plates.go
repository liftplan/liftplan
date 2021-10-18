package gear

import (
	"errors"
	"fmt"
	"sort"
)

// Plates are are a set of plates and its corresponding unit.
type Plates struct {
	Weights []float64 `json:"weights"`
	Unit    Unit      `json:"unit"`
}

var (
	// ErrNoPlatesFound is the basic error for not finding plates.
	ErrNoPlatesFound = errors.New("no plates found")
)

// DefaultWeightsKB is the default set of weights in KB
var DefaultWeightsKB = []float64{
	1.25,
	2.5,
	5,
	10,
	15,
	20,
}

// DefaultWeightsLBS is the default set of weights in LBS
var DefaultWeightsLBS = []float64{
	2.5,
	5,
	10,
	25,
	35,
	45,
}

// DefaultPlatesLBS is the default set of plates in LBS
var DefaultPlatesLBS = Plates{
	Weights: DefaultWeightsLBS,
	Unit:    LBS,
}

// DefaultPlatesKG is the default set of plates in KB
var DefaultPlatesKG = Plates{
	Weights: DefaultWeightsKB,
	Unit:    KG,
}

func (p Plates) String() string {
	return fmt.Sprintf("Weights: %v, Unit: %v", p.Weights, p.Unit)
}

// Tidy cleans up the plates by removing duplicates and sorting them
func (p *Plates) Tidy() {
	p.Weights = tidy(p.Weights)
}

// Add takes a plate and adds it to the set of weights.
func (p *Plates) Add(plate float64) {
	p.Weights = addItem(p.Weights, plate)
}

// Remove takes a plate and removes it from the set of weights.
func (p *Plates) Remove(plate float64) {
	p.Weights = removeItem(p.Weights, plate)
}

// Min gets the smallest increment of plate in the Weights slice.
func (p Plates) Min() (float64, error) {
	if len(p.Weights) == 0 {
		return 0, ErrNoPlatesFound
	}
	p.Tidy()
	return p.Weights[0], nil
}

// Equals compares all values in Weights and Unit and returns true
// if all values are equal.
func (p Plates) Equals(c Plates) bool {
	return (p.Unit == c.Unit) && equal(p.Weights, c.Weights)
}

// Round takes a weight and uses its increment for rounding and doubles the increment
// to match the smallest increment per side of the bar. It returns the closest number
// below or even to the weight that is divisible by the doubled increment
func (p Plates) Round(weight float64) (float64, error) {
	m, err := p.Min()
	if err != nil {
		return 0, err
	}
	b := m * 2
	return float64(int(weight/b)) * b, nil
}

// equal compares to lists of floats and returns a boolean value on equality.
func equal(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// tidy takes a slice of float64, removes any duplicate values and sorts the output.
func tidy(input []float64) []float64 {
	m := make(map[float64]struct{})
	for _, x := range input {
		if x > 0 { // don't allow 0 or negative value plates
			m[x] = struct{}{}
		}
	}
	o := make([]float64, len(m))
	i := 0
	for k := range m {
		o[i] = k
		i++
	}
	sort.Float64s(o)
	return o
}

func addItem(slice []float64, item float64) []float64 {
	return tidy(append(slice, item))
}

func removeItem(slice []float64, item float64) (output []float64) {
	slice = tidy(slice)
	for _, x := range slice {
		if x != item {
			output = append(output, x)
		}
	}
	return
}

// Recommend takes a weight and a set of plates and returns
// a sorted recommendation of plates for one side of the bar
func Recommend(weight float64, plates []float64) ([]float64, error) {
	var rec []float64
	if len(plates) == 0 {
		return nil, ErrNoPlatesFound
	}
	sort.Float64s(plates)
	for i := len(plates) - 1; i >= 0; i-- {
		n := int(weight / (plates[i] * 2))
		for j := 1; j <= n; j++ {
			rec = append(rec, plates[i])
		}
		weight = weight - (float64(n) * (plates[i] * 2))
	}
	sort.Float64s(rec)
	return rec, nil
}

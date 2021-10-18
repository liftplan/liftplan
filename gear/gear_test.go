package gear

import (
	"testing"
)

func TestGear(t *testing.T) {
	t.Parallel()

	t.Run("Default", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			unit     Unit
			expected Gear
		}{
			{77, Gear{}},
			{KG, Gear{
				Bar:    MensBarKG,
				Plates: DefaultPlatesKG,
				Unit:   KG,
			}},
			{LBS, Gear{
				Bar:    MensBarLBS,
				Plates: DefaultPlatesLBS,
				Unit:   LBS,
			}},
		}
		for _, test := range tt {
			output := Default(test.unit)
			if !output.Equals(test.expected) {
				t.Errorf("expected: %v, output: %v", test.expected, output)
			}
		}
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    Gear
			expected string
		}{
			{Gear{
				Unit:   KG,
				Bar:    MensBarLBS,
				Plates: Plates{Weights: DefaultWeightsKB, Unit: KG},
			}, "Unit: KG, Bar: { Weight: 45, Unit: LBS }, Plates: { Weights: [1.25 2.5 5 10 15 20], Unit: KG }"},
		}
		for _, test := range tt {
			if test.input.String() != test.expected {
				t.Error("match failed for:", test.input.String(), test.expected)
			}
		}
	})
	t.Run("Min", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    Gear
			expected float64
			err      error
		}{
			{Gear{
				Unit:   LBS,
				Bar:    MensBarLBS,
				Plates: Plates{},
			}, MensBarLBS.Weight, nil},
			{Gear{
				Unit:   KG,
				Bar:    MensBarLBS,
				Plates: Plates{},
			}, 20.411656649988664, nil},
			{Gear{
				Unit:   Unit(5),
				Bar:    MensBarLBS,
				Plates: Plates{},
			}, 0.0, ErrInvalidUnit},
		}
		for _, test := range tt {
			o, err := test.input.Min()
			if err != test.err {
				t.Error("unexpected error:", err, test.err)
			}
			if o != test.expected {
				t.Error("unexpected result:", o, test.expected)
			}
		}
	})
	t.Run("Round", func(t *testing.T) {
		t.Parallel()
		g := Gear{
			Unit:   LBS,
			Bar:    MensBarLBS,
			Plates: Plates{Weights: DefaultWeightsKB, Unit: LBS},
		}
		tt := []struct {
			gear     Gear
			input    float64
			expected float64
			err      error
		}{
			{g, 89, 87.5, nil},
			{g, 46, 45, nil},
			{g, 45, 45, nil},
			{g, 44, 0, ErrInputLessThanBar},
			{Gear{
				Unit:   LBS,
				Bar:    MensBarLBS,
				Plates: Plates{Weights: []float64{}, Unit: LBS},
			}, 56, 0, ErrNoPlatesFound},
			{Gear{
				Unit:   Unit(5),
				Bar:    Bar{Weight: 45, Unit: Unit(5)},
				Plates: Plates{},
			}, 0, 0, ErrInvalidUnit},
			{Gear{
				Unit:   LBS,
				Bar:    Bar{Weight: 45, Unit: Unit(5)},
				Plates: Plates{},
			}, 0, 0, ErrInvalidUnit},
			{Gear{
				Unit:   LBS,
				Bar:    MensBarLBS,
				Plates: Plates{Weights: []float64{}, Unit: Unit(5)},
			}, 89, 0, ErrInvalidUnit},
			{Gear{
				Unit:   KG,
				Bar:    MensBarKG,
				Plates: Plates{Weights: []float64{}, Unit: Unit(5)},
			}, 89, 0, ErrInvalidUnit},
		}
		for _, test := range tt {
			o, err := test.gear.Round(test.input)
			if err != test.err {
				t.Error("unexpected error:", err, test.err)
			} else if o != test.expected {
				t.Error("unexpected result:", o, test.expected)
			}
		}
	})
	t.Run("Recommend", func(t *testing.T) {
		t.Parallel()
		w := DefaultWeightsLBS
		w = append(w, 1.25)
		g := Gear{
			Unit:   LBS,
			Bar:    MensBarLBS,
			Plates: Plates{Weights: w, Unit: LBS},
		}
		tt := []struct {
			gear     Gear
			weight   float64
			expected []float64
			err      error
		}{
			{g, 157.5, []float64{1.25, 10, 45}, nil},
			{g, 44, []float64{}, ErrInputLessThanBar},
		}
		for _, test := range tt {
			o, err := test.gear.Recommend(test.weight)
			if err != test.err {
				t.Error("unexpected error:", err, test.err)
			} else if !equal(o, test.expected) {
				t.Error("unexpected result:", o, test.expected)
			}
		}
	})
}

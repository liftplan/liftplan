package gear

import "testing"

func TestBar(t *testing.T) {
	t.Parallel()
	t.Run("String", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    Bar
			expected string
		}{
			{MensBarLBS, "Weight: 45, Unit: LBS"},
			{WomensBarKG, "Weight: 15, Unit: KG"},
		}
		for _, test := range tt {
			if test.input.String() != test.expected {
				t.Error("match failed for", test)
			}
		}
	})
	t.Run("ConvertTo", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			bar      Bar
			unit     Unit
			expected float64
			err      error
		}{
			{MensBarLBS, LBS, 45.0, nil},
			{MensBarLBS, KG, 20.411656649988664, nil},
			{WomensBarKG, LBS, 33.06933932775, nil},
			{WomensBarKG, KG, 15.0, nil},
			{WomensBarKG, Unit(5), 0.0, ErrInvalidUnit},
		}
		for _, test := range tt {
			o, err := test.bar.ConvertTo(test.unit)
			if err != test.err {
				t.Error("expected error mismatch:", err, test.err)
			}
			if o != test.expected {
				t.Error("match failed for", o, test.expected)
			}
		}
	})
	t.Run("Equals", func(t *testing.T) {
		t.Parallel()

		lbsBar20 := Bar{Weight: 20, Unit: LBS}
		kgBar20 := Bar{Weight: 20, Unit: KG}

		tt := []struct {
			bar      Bar
			comp     Bar
			expected bool
		}{
			{MensBarLBS, MensBarLBS, true},
			{MensBarLBS, WomensBarKG, false},
			{lbsBar20, kgBar20, false},
		}
		for _, test := range tt {
			if o := test.bar.Equals(test.comp); o != test.expected {
				t.Errorf("expected %v for %v and %v:", test.expected, test.bar, test.comp)
			}
		}
	})
}

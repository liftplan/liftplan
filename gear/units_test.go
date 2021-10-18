package gear

import (
	"testing"
)

func TestUnitFromString(t *testing.T) {
	t.Parallel()
	tt := []struct {
		input    string
		expected Unit
		err      error
	}{
		{"KG", KG, nil},
		{"LBS", LBS, nil},
		{"FOO", 0, ErrInvalidUnit},
	}
	for _, test := range tt {
		u, err := UnitFromString(test.input)
		if err != test.err {
			t.Errorf("unexpected error: %v, expected: %v", err, test.err)
		}
		if u != test.expected {
			t.Errorf("unexpected unit: %v, expected: %v", u, test.expected)
		}
	}
}

func TestUnit(t *testing.T) {
	t.Parallel()
	t.Run("String", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    Unit
			expected string
		}{
			{KG, "KG"},
			{LBS, "LBS"},
			{Unit(2), ""},
			{Unit(3), ""},
		}

		for _, test := range tt {
			if test.input.String() != test.expected {
				t.Error("match failed for", test)
			}
		}
	})
	t.Run("ConvertFromTo", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    float64
			from     Unit
			to       Unit
			expected float64
			err      error
		}{
			{0, LBS, KG, 0, nil},
			{55.0, LBS, KG, 24.947580349986147, nil},
			{55.0, KG, LBS, 121.25424420175, nil},
			{0, 5, LBS, 0, ErrInvalidUnit},
			{0, LBS, 5, 0, ErrInvalidUnit},
			{0, 5, KG, 0, ErrInvalidUnit},
		}

		for _, test := range tt {
			o, err := ConvertFromTo(test.input, test.from, test.to)
			if err != test.err {
				t.Error("expected error mismatch:", err, test.err)
			}
			if o != test.expected {
				t.Error("unexpected output:", o, test.expected)
			}
		}
	})
}

package fto

import (
	"bytes"
	"errors"
	"testing"

	"github.com/liftplan/liftplan"
	"github.com/liftplan/liftplan/gear"
)

func TestStrategy(t *testing.T) {
	t.Parallel()
	t.Run("Plan", func(t *testing.T) {
		m1 := Movement{
			Name:        "over-head press",
			TrainingMax: 175,
			Unit:        gear.LBS,
		}
		m2 := Movement{
			Name:        "squat",
			TrainingMax: 185,
			Unit:        gear.LBS,
		}
		s := Strategy{
			Movements: []Movement{m1, m2},
			Gear:      gear.Default(gear.LBS),
			Type:      FSLMULTI,
			Warmup:    true,
			JokerSets: true,
		}
		_, err := s.Plan(liftplan.JSON)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestSetType(t *testing.T) {
	t.Parallel()
	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    []byte
			expected SetType
			err      error
		}{
			{[]byte(`"Working"`), Working, nil},
			{[]byte(`"foo"`), Working, ErrInvalidSetType},
			{[]byte(`false`), Working, errors.New("json: cannot unmarshal bool into Go value of type string")},
		}
		for _, test := range tt {
			var s SetType
			if err := s.UnmarshalJSON(test.input); err != nil {
				if err.Error() != test.err.Error() {
					t.Error(err, test.err)
				}
			} else {
				if s != test.expected {
					t.Error(s, test.expected)
				}
			}
		}
	})
	t.Run("String", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    SetType
			expected string
		}{
			{Working, "Working"},
			{SetType(20), ""},
		}

		for _, test := range tt {
			if test.input.String() != test.expected {
				t.Error(test)
			}
		}

	})
}

func TestDeloadType(t *testing.T) {
	t.Parallel()
	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    []byte
			expected DeloadType
			err      error
		}{
			{[]byte(`"deload1"`), Deload1, nil},
			{[]byte(`"foo"`), Deload1, ErrInvalidDeloadType},
			{[]byte(`false`), Deload1, errors.New("json: cannot unmarshal bool into Go value of type string")},
		}
		for _, test := range tt {
			var d DeloadType
			if err := d.UnmarshalJSON(test.input); err != nil {
				if err.Error() != test.err.Error() {
					t.Error(err, test.err)
				}
			} else {
				if d != test.expected {
					t.Error(d, test.expected)
				}
			}
		}
	})
	t.Run("String", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    DeloadType
			expected string
		}{
			{Deload1, "deload1"},
			{DeloadType(20), "deload21"},
		}

		for _, test := range tt {
			if test.input.String() != test.expected {
				t.Error(test.input.String(), test.expected)
			}
		}
	})
	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    DeloadType
			expected []byte
		}{
			{Deload1, []byte(`"deload1"`)},
		}

		for _, test := range tt {
			b, _ := test.input.MarshalJSON()
			if !bytes.Equal(b, test.expected) {
				t.Error(test.input, test.expected)
			}
		}
	})
}

func TestStrategyType(t *testing.T) {
	t.Parallel()
	t.Run("String", func(t *testing.T) {
		t.Parallel()
		g := FSLMULTI
		b := StrategyType(55)
		if g.String() != "FSL Multiple Sets" {
			t.Error(g.String(), "!= FSL Multiple Sets")
		}
		if b.String() != "" {
			t.Error(b.String(), "!= ''")
		}
	})
	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		g := FSLMULTI
		b, _ := g.MarshalJSON()
		if !bytes.Equal(b, []byte(`"FSL Multiple Sets"`)) {
			t.Error(g, "!=", []byte(`"FSL Multiple Sets"`))
		}
	})
	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()

		goodBytes := []byte(`"FSL Multiple Sets"`)
		badBytes := []byte(`false`)

		tt := []struct {
			input    []byte
			expected StrategyType
			err      error
		}{
			{goodBytes, FSLMULTI, nil},
			{badBytes, StrategyType(0), errors.New("json: cannot unmarshal bool into Go value of type string")},
		}

		for _, test := range tt {
			var s StrategyType
			if err := s.UnmarshalJSON(test.input); err != nil {
				if err.Error() != test.err.Error() {
					t.Error(err, test.err)
				}
			} else {
				if s != test.expected {
					t.Error(s, test.expected)
				}
			}
		}
	})
}

func TestStrategyTypeFromString(t *testing.T) {
	t.Parallel()
	g, err := StrategyTypeFromString("FSL Multiple Sets")
	if err != nil {
		t.Error(err)
	}
	if g != StrategyType(0) {
		t.Error("mismatch", g, StrategyType(0), errors.New("json: cannot unmarshal bool into Go value of type string"))
	}
	{
		_, err := StrategyTypeFromString("blah")
		if err != ErrInvalidStrategyType {
			t.Error("error failed to return for 'blah'")
		}
	}
}

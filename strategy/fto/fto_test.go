package fto

import (
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
}

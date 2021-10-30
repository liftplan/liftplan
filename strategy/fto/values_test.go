package fto

import (
	"errors"
	"net/url"
	"testing"

	"github.com/liftplan/liftplan/gear"
)

func TestIsChecked(t *testing.T) {
	t.Parallel()
	v := url.Values{
		"foo": []string{"true"},
		"bar": []string{"false"},
		"bam": []string{},
	}
	tt := []struct {
		key      string
		vals     url.Values
		expected bool
	}{
		{"foo", v, true},
		{"bar", v, false},
		{"baz", v, false},
		{"bam", v, false},
	}

	for _, test := range tt {
		k, v, e := test.key, test.vals, test.expected
		if isChecked(k, v) != e {
			t.Errorf("key: %v, value: %v, expected: %v", k, v, e)
		}
	}
}

func TestFromValues(t *testing.T) {
	t.Parallel()

	m1 := Movement{
		Name:        "deadlift",
		TrainingMax: 175,
		Unit:        gear.LBS,
	}
	m2 := Movement{
		Name:        "bench press",
		TrainingMax: 400,
		Unit:        gear.LBS,
	}
	m3 := Movement{
		Name:        "overhead press",
		TrainingMax: 175,
		Unit:        gear.LBS,
	}
	m4 := Movement{
		Name:        "squat",
		TrainingMax: 400,
		Unit:        gear.LBS,
	}
	s1 := Strategy{
		Movements: []Movement{m1, m2, m3, m4},
		Gear:      gear.Default(gear.LBS),
		Type:      FSLMULTI,
		Warmup:    true,
		JokerSets: true,
	}

	goodVals, _ := s1.Values()
	missingStrat, _ := s1.Values()
	missingStrat.Del("fto.strategy")
	badStrat, _ := s1.Values()
	badStrat["fto.strategy"] = []string{"blah"}

	tt := []struct {
		input    url.Values
		expected Strategy
		err      error
	}{
		{url.Values{}, s1, gear.ErrMissingUnitQuery},
		{goodVals, s1, nil},
		{missingStrat, s1, errors.New("missing strategy in query")},
		{badStrat, s1, ErrInvalidStrategyType},
	}

	for _, test := range tt {
		_, err := FromValues(test.input)
		if err != nil {
			if err.Error() != test.err.Error() {
				t.Error(test.err, err)
			}
		}
	}
}

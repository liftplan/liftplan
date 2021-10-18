package gear

import (
	"testing"
)

func TestTidy(t *testing.T) {
	t.Parallel()
	tt := []struct {
		input    []float64
		expected []float64
	}{
		{[]float64{1, 1, 1, 1, 1}, []float64{1}},
		{[]float64{70, 50, 90, 90, 0}, []float64{50, 70, 90}},
	}
	for _, test := range tt {
		if o := tidy(test.input); !equal(o, test.expected) {
			t.Error("unexpected mismatch:", o, test.expected)
		}
	}
}

func TestAddItem(t *testing.T) {
	t.Parallel()
	tt := []struct {
		slice    []float64
		item     float64
		expected []float64
	}{
		{[]float64{}, 5, []float64{5}},
	}

	for _, test := range tt {
		if o := addItem(test.slice, test.item); !equal(o, test.expected) {
			t.Error("unexpected mismatch:", o, test.expected)
		}
	}
}

func TestRemoveItem(t *testing.T) {
	t.Parallel()
	tt := []struct {
		slice    []float64
		item     float64
		expected []float64
	}{
		{[]float64{}, 5, []float64{}},
		{[]float64{5, 5, 5}, 5, []float64{}},
		{[]float64{1, 2, 3, 4, 5}, 5, []float64{1, 2, 3, 4}},
	}

	for _, test := range tt {
		if o := removeItem(test.slice, test.item); !equal(o, test.expected) {
			t.Error("unexpected mismatch:", o, test.expected)
		}
	}
}

func TestRecommend(t *testing.T) {
	t.Parallel()
	tt := []struct {
		weight   float64
		plates   []float64
		expected []float64
		err      error
	}{
		{112.5, DefaultWeightsLBS, []float64{10, 45}, nil},
		{387.5, DefaultWeightsLBS, []float64{2.5, 10, 45, 45, 45, 45}, nil},
		{0, DefaultWeightsLBS, []float64{}, nil},
		{270, DefaultWeightsLBS, []float64{45, 45, 45}, nil},
		{271, DefaultWeightsLBS, []float64{45, 45, 45}, nil},
		{0, []float64{}, []float64{}, ErrNoPlatesFound},
	}
	for _, test := range tt {
		r, err := Recommend(test.weight, test.plates)
		if err != test.err {
			t.Error("unexpected error:", err, test.err)
		}
		if !equal(r, test.expected) {
			t.Errorf("not equal. result: %v, expected %v", r, test.expected)
		}
	}
}

func TestPlates(t *testing.T) {
	t.Parallel()
	t.Run("String", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    Plates
			expected string
		}{
			{Plates{Weights: DefaultWeightsKB, Unit: KG}, "Weights: [1.25 2.5 5 10 15 20], Unit: KG"},
			{Plates{Weights: DefaultWeightsLBS, Unit: LBS}, "Weights: [2.5 5 10 25 35 45], Unit: LBS"},
		}
		for _, test := range tt {
			if test.input.String() != test.expected {
				t.Error("match failed for", test.input.String(), test.expected)
			}
		}
	})
	t.Run("Tidy", func(t *testing.T) {
		t.Parallel()
		tt := []struct {
			input    Plates
			expected []float64
		}{
			{Plates{Weights: []float64{5, 5, 5, 5, 5}, Unit: KG}, []float64{5}},
			{Plates{Weights: []float64{5, 4, 3, 2, 1, 0, 0}, Unit: LBS}, []float64{1, 2, 3, 4, 5}},
		}
		for _, test := range tt {
			test.input.Tidy()
			if o := test.input.Weights; !equal(o, test.expected) {
				t.Error("match failed for", o, test.expected)
			}
		}
	})
	t.Run("Equals", func(t *testing.T) {
		t.Parallel()

		badUnit := DefaultPlatesKG
		badUnit.Unit = LBS
		badLen := DefaultPlatesKG
		badLen.Add(200)
		badWeights := DefaultPlatesKG
		badWeights.Add(200)
		badWeights.Remove(5)

		tt := []struct {
			plates   Plates
			comp     Plates
			expected bool
		}{
			{DefaultPlatesKG, DefaultPlatesKG, true},
			{DefaultPlatesKG, DefaultPlatesLBS, false},
			{DefaultPlatesKG, badUnit, false},
			{DefaultPlatesKG, badLen, false},
			{DefaultPlatesKG, badWeights, false},
		}
		for _, test := range tt {
			if o := test.plates.Equals(test.comp); o != test.expected {
				t.Errorf("expected %v for %v and %v:", test.expected, test.plates, test.comp)
			}
		}
	})
}

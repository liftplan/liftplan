package gear

import (
	"net/url"
	"testing"
)

func TestToValues(t *testing.T) {
	t.Parallel()

	goodGear := Default(LBS)
	badBar := Default(LBS)
	badBar.Bar.Unit = 5
	badPlates := Default(LBS)
	badPlates.Plates.Unit = 5

	tt := []struct {
		gear Gear
		vals url.Values
		err  error
	}{
		{goodGear, url.Values{
			"gear.bar.lbs":   []string{"45.00"},
			"gear.plate.lbs": []string{"2.50", "5.00", "10.00", "25.00", "35.00", "45.00"},
			"gear.unit":      []string{"lbs"},
		}, nil},
		{badBar, url.Values{}, ErrInvalidUnit},
		{badPlates, url.Values{}, ErrInvalidUnit},
	}

	for _, test := range tt {
		o, err := ToValues(test.gear)
		if err != test.err {
			t.Error("unexpected error:", err, test.err)
		}
		if test.err != nil {
			continue
		}
		for k, v := range o {
			r, ok := test.vals[k]
			if !ok {
				t.Error("missing:", k)
			}
			for i, val := range v {
				if r[i] != val {
					t.Error("value missmatch:", val, r[i])
				}
			}
		}
	}
}

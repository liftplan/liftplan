package gear

import (
	"encoding/xml"
	"errors"
	"io"
	"net/url"
	"strings"
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

func TestFromValues(t *testing.T) {
	t.Parallel()

	invalid := url.Values{}
	valid, _ := ToValues(Default(LBS))
	badPlates, _ := ToValues(Default(LBS))
	badPlates.Del("gear.plate.lbs")
	badPlatesVal, _ := ToValues(Default(LBS))
	badPlatesVal.Add("gear.plate.lbs", "foo")
	badBar, _ := ToValues(Default(LBS))
	badBar.Del("gear.bar.lbs")
	badBarVal, _ := ToValues(Default(LBS))
	badBarVal["gear.bar.lbs"] = []string{"foo"}
	badValErr := errors.New(`strconv.ParseFloat: parsing "foo": invalid syntax`)

	tt := []struct {
		values   url.Values
		expected Gear
		err      error
	}{
		{invalid, Gear{}, ErrMissingUnitQuery},
		{valid, Default(LBS), nil},
		{badPlates, Gear{}, ErrMissingPlatesQuery},
		{badPlatesVal, Gear{}, badValErr},
		{badBar, Gear{}, ErrMissingBarQuery},
		{badBarVal, Gear{}, badValErr},
	}

	for _, test := range tt {
		o, err := FromValues(test.values)
		if err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("expected error: %v, got: %v", test.err, err)
			}
		}

		if !o.Equals(test.expected) {
			t.Errorf("expected: %v, got: %v", test.expected, o)
		}
	}
}

func TestFormFields(t *testing.T) {
	t.Parallel()
	r := strings.NewReader(string(FormFields()))
	d := xml.NewDecoder(r)

	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity
	for {
		_, err := d.Token()
		if err == io.EOF {
			return
		}
		if err != nil {
			t.Error(err)
		}
	}

}

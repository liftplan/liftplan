package gear

import (
	"bytes"
	_ "embed" // for use with native embedding
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"strconv"
	"strings"
)

type values url.Values

const (
	namespace = "gear"
)

var (
	//go:embed templates/form.go.html
	formTemplate string
)

var (
	ErrMissingUnitQuery   = errors.New("missing unit in query")
	ErrMissingPlatesQuery = errors.New("missing plates in query")
	ErrMissingBarQuery    = errors.New("missing bar in query")
)

// options represent the input options for the gear form templates.
type options struct {
	Value   string
	Name    string
	Bars    []option
	Plates  []option
	Checked bool
}

// option represents a Value, Name, and Checked boolean for the configurable options.
type option struct {
	Value   string
	Name    string
	Checked bool
}

// ToValues takes a Gear and returns properly formatted url.Values to be used in a get URL.
func ToValues(g Gear) (url.Values, error) {
	vals := make(url.Values)
	unit := strings.ToLower(g.Unit.String())
	vals.Set(namespace+".unit", unit)
	b, err := g.Bar.ConvertTo(g.Unit)
	if err != nil {
		return vals, err
	}
	vals.Add(namespace+".bar."+unit, fmt.Sprintf("%.2f", b))
	for _, w := range g.Plates.Weights {
		p, err := ConvertFromTo(w, g.Plates.Unit, g.Unit)
		if err != nil {
			return vals, err
		}
		vals.Add(namespace+".plate."+unit, fmt.Sprintf("%.2f", p))
	}
	return vals, nil
}

// FromValues takes a set of values in `url.Values` format and returns gear and an error.
func FromValues(vals url.Values) (g Gear, err error) {
	v := values(vals)

	unit, err := v.unit()
	if err != nil {
		return g, err
	}
	plates, err := v.plates(unit)
	if err != nil {
		return g, err
	}
	bar, err := v.bar(unit)
	if err != nil {
		return g, err
	}
	g.Unit = unit
	g.Plates = plates
	g.Bar = bar
	return g, nil
}

func (v values) unit() (u Unit, err error) {
	units, ok := v[namespace+".unit"]
	if !ok {
		return u, ErrMissingUnitQuery
	}
	return UnitFromString(strings.ToUpper(units[0]))
}

func (v values) plates(unit Unit) (p Plates, err error) {
	plates, ok := v[namespace+".plate."+strings.ToLower(unit.String())]
	if !ok {
		return p, ErrMissingPlatesQuery
	}
	l := len(plates)
	pi := make([]float64, l)
	for i, plate := range plates {
		f, err := strconv.ParseFloat(plate, 64)
		if err != nil {
			return p, err
		}
		pi[i] = f
	}
	p.Weights = tidy(pi)
	p.Unit = unit
	return p, nil
}

func (v values) bar(unit Unit) (b Bar, err error) {

	w, ok := v[namespace+".bar."+strings.ToLower(unit.String())]
	if !ok {
		return b, ErrMissingBarQuery
	}
	weight, err := strconv.ParseFloat(w[0], 64)
	if err != nil {
		return b, err
	}
	b.Unit = unit
	b.Weight = weight
	return b, err
}

// FormFields returns an html snippet for choosing lifting gear for the submit form
func FormFields() (r template.HTML, err error) {
	lbs := options{
		Value: "lbs",
		Name:  "LBS",
		Bars: []option{
			{Value: "45", Name: "45 LBS"},
			{Value: "35", Name: "35 LBS"},
		},
		Plates: []option{
			{Value: "1.25"},
			{Value: "2.5", Checked: true},
			{Value: "5", Checked: true},
			{Value: "10", Checked: true},
			{Value: "15"},
			{Value: "25", Checked: true},
			{Value: "35"},
			{Value: "45", Checked: true},
			{Value: "100"},
		},
		Checked: true,
	}
	kg := options{
		Value: "kg",
		Name:  "KG",
		Bars: []option{
			{Value: "20", Name: "20 KG"},
			{Value: "15", Name: "15 KG"},
		},
		Plates: []option{
			{Value: "0.25"},
			{Value: "0.5"},
			{Value: "1.25", Checked: true},
			{Value: "2.5", Checked: true},
			{Value: "5", Checked: true},
			{Value: "10", Checked: true},
			{Value: "15", Checked: true},
			{Value: "20", Checked: true},
			{Value: "25", Checked: true},
			{Value: "50"},
		},
	}

	t, err := template.New(namespace).Parse(formTemplate)
	if err != nil {
		return r, err
	}
	var b bytes.Buffer
	if err := t.Execute(&b, []options{lbs, kg}); err != nil {
		return r, err
	}
	return template.HTML(b.String()), err
}

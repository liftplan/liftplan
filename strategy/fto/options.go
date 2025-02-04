package fto

import (
	"bytes"
	_ "embed" // used for embeding templates
	"html/template"

	"github.com/liftplan/liftplan"
)

var (
	//go:embed templates/form.go.html
	formTemplate string
)

type input struct {
	Template *template.Template
	Options  options
}

type options struct {
	Movements   []choice
	Selectables []choice
	Strategies  []choice
}

type choice struct {
	Name    string
	Value   string
	Checked bool
}

// FormFields returns a liftplan.FormFields
func FormFields() liftplan.FormFields {

	s := []choice{
		{Name: "Warmup", Value: "warmup", Checked: true},
		{Name: "Jokersets", Value: "jokersets", Checked: true},
		{Name: "Recommended Plates", Value: "recplates", Checked: false},
	}

	mo := []choice{
		{Name: "deadlift", Value: "0"},
		{Name: "bench press", Value: "1"},
		{Name: "overhead press", Value: "2"},
		{Name: "back squat", Value: "3"},
	}

	strats := []choice{
		{
			Name:    FSLMULTI.String(),
			Value:   FSLMULTI.String(),
			Checked: true,
		},
		{
			Name:    FSL.String(),
			Value:   FSL.String(),
			Checked: false,
		},
	}

	o := options{Selectables: s, Movements: mo, Strategies: strats}
	t, _ := template.New("fto").Parse(formTemplate)
	return input{Template: t, Options: o}
}

// Render returns the template.HTML for an input template
func (i input) Render() (template.HTML, error) {
	var b bytes.Buffer
	err := i.Template.Execute(&b, i.Options)
	return template.HTML(b.Bytes()), err
}

// Name returns a string with the name of the strategy methods
func (i input) Name() string { return "Beyond 5/3/1" }

// ShortCode is the code used for templating
func (i input) ShortCode() string { return "fto" }

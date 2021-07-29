package serve

import (
	"html/template"

	"github.com/liftplan/liftplan"
)

// Options represent HTML Gear and Method options for the webapp
type Options struct {
	Gear    template.HTML
	Methods []liftplan.FormFields
}

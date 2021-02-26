package serve

import (
	"html/template"

	"github.com/liftplan/liftplan"
)

type Options struct {
	Gear    template.HTML
	Methods []liftplan.FormFields
}

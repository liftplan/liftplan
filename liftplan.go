package liftplan

import (
	"html/template"
	"net/url"
)

// Format is used  as the datatype for Export formats.
type Format uint

const (
	// JSON format
	JSON Format = iota
	// HTML format
	HTML
)

// Liftplanner is an interface that wraps around 3 more basic interfaces
type Liftplanner interface {
	Planner
	Valuer
}

// FormFields is an interface for a strategy so that all of its formfields can be generated.
type FormFields interface {
	Renderer
	Namer
	ShortCoder
}

// Planner is an interface used to support methods that check and export plans in various formats.
type Planner interface {
	// Plan is used in combination with the Format type to choose an export format
	Plan(f Format) ([]byte, error)
}

// Renderer is used for rendering template.HTML data with a specific namespace.
type Renderer interface {
	Render() (template.HTML, error)
}

// Valuer exports the multipart.Form.Values or url.Values compatible values
type Valuer interface {
	Values() (url.Values, error)
}

// Namer is used to return the name as a string
type Namer interface {
	Name() string
}

// ShortCoder is the interface for a Liftplan's shortCode
type ShortCoder interface {
	ShortCode() string
}

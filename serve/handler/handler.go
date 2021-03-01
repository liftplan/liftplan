package handler

import (
	_ "embed" // used for embeding templates
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/liftplan/liftplan"
	"github.com/liftplan/liftplan/gear"
	"github.com/liftplan/liftplan/serve"
	"github.com/liftplan/liftplan/strategy/fto"
)

const (
	maxBytes = 100000
	maxAge   = 60 * 60 * 24 * 365 // cache for 1 year (in seconds)
)

var (
	//go:embed templates/plan.go.html
	planTemplate string
	//go:embed templates/root.go.html
	rootTemplate string
	//go:embed templates/footer.go.html
	footerTemplate string
	//go:embed templates/header.go.html
	headerTemplate string
)

func badRequestError(w http.ResponseWriter, err error) {
	w.Header().Del("Cache-Control")
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	log.Println(err)
}

func cacheControl(duration int, w http.ResponseWriter) {
	w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%v", duration))
}

// Root returns the root site for the form
func Root() http.HandlerFunc {
	t, err := pageTemplate(rootTemplate, "root")
	if err != nil {
		log.Fatal(err)
	}
	opts, err := getOptions()
	if err != nil {
		log.Fatal(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, opts); err != nil {
			badRequestError(w, err)
			return
		}
		cacheControl(maxAge, w)
		w.Header().Add("Content-Type", "text/html")
	}
}

// Plan returns a primary HTML plan
func Plan() http.HandlerFunc {
	t, err := pageTemplate(planTemplate, "plan")
	if err != nil {
		log.Fatal(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			formSubmit(w, r)
		case "GET":
			if wantsJSON(r) {
				w.Header().Add("Content-Type", "application/json")
				renderJSON(w, r)
			} else {
				w.Header().Add("Content-Type", "text/html")
				renderHTML(t, w, r)
			}
			cacheControl(maxAge, w)
		default:
			badRequestError(w, fmt.Errorf("invalid request method: %v", r.Method))
		}
	}
}

func formSubmit(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(maxBytes)
	p, err := plannerFromValues(r.Form)
	if err != nil {
		badRequestError(w, err)
		return
	}
	v, err := p.Values()
	if err != nil {
		badRequestError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("%v?%v", r.URL.Path, v.Encode()), 301)
}

func wantsJSON(r *http.Request) bool {
	if r.Header.Get("Accept") == "application/json" {
		return true
	}
	if r.URL.Query().Get("accept") == "application/json" {
		return true
	}
	return false
}

func renderJSON(w http.ResponseWriter, r *http.Request) {
	h, err := renderFromValues(r.URL.Query(), liftplan.JSON)
	if err != nil {
		badRequestError(w, err)
		return
	}
	w.Write(h)
}

func renderHTML(t *template.Template, w http.ResponseWriter, r *http.Request) {
	h, err := renderFromValues(r.URL.Query(), liftplan.HTML)
	if err != nil {
		badRequestError(w, err)
		return
	}
	if err := t.Execute(w, template.HTML(h)); err != nil {
		badRequestError(w, err)
		return
	}
}

func pageTemplate(core string, name string) (*template.Template, error) {
	t, err := template.New(name).Parse(core)
	if err != nil {
		return nil, err
	}
	if _, err := t.New("footer").Parse(footerTemplate); err != nil {
		return nil, err
	}
	if _, err := t.New("header").Parse(headerTemplate); err != nil {
		return nil, err
	}
	return t, nil
}

func getOptions() (serve.Options, error) {
	gf, err := gear.FormFields()
	if err != nil {
		return serve.Options{}, err
	}

	f, err := fto.FormFields()
	return serve.Options{
		Methods: []liftplan.FormFields{f},
		Gear:    gf,
	}, err
}

func plannerFromValues(vals url.Values) (liftplan.Liftplanner, error) {
	m, ok := vals["method"]
	if !ok {
		return nil, errors.New("missing method in query")
	}
	switch strings.ToLower(m[0]) {
	case "fto":
		return fto.FromValues(vals)
	default:
		return nil, errors.New("unrecognized method")
	}
}

func renderFromValues(vals url.Values, f liftplan.Format) ([]byte, error) {
	p, err := plannerFromValues(vals)
	if err != nil {
		return nil, err
	}
	return p.Plan(f)
}

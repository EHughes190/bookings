package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/EHughes190/bookings/pkg/config"
	"github.com/EHughes190/bookings/pkg/models"
)

//a map of functions which we can pass to templates from go.
var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

//Public func that renders html based on template files.
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	//gets template cache from the appConfigf
	tmplCache := app.TemplateCache

	//gets the correct template we want
	//ok is either true or false
	t, ok := tmplCache[tmpl]

	if !ok {
		log.Fatal("could not get template from template cache")
	}

	//we haven't read template from disk so we need to hold byte info of template from memory
	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	//write to template
	_, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println("error writing template to the browser", err)
	}
}

//Creates a cache of templates we can use to render pages
func CreateTemplateCache() (map[string]*template.Template, error) {
	//cache with string keys, and templates as values
	tmplCache := map[string]*template.Template{}

	// this gets a list of all files ending with page.tmpl, and stores
	// it in a slice of strings called pages
	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return tmplCache, err
	}

	for _, page := range pages {
		//e.g. "about.page.tmpl"
		name := filepath.Base(page)

		//template set of the name, including functions that we have in our function map
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return tmplCache, err
		}

		//return any files that match this file path (our layout files)
		matches, err := filepath.Glob("./templates/*.layout.tmpl")

		if err != nil {
			return tmplCache, err
		}

		// if the length of matches is > 0, we go through the slice
		// and parse all of the layouts available to us. We might not use
		// any of them in this iteration through the loop, but if the current
		// template we are working on (home.page.tmpl the first time through)
		// does use a layout, we need to have it available to us before we add it
		// to our template set
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return tmplCache, err
			}
		}

		tmplCache[name] = ts
	}
	return tmplCache, nil
}

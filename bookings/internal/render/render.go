package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/yahia-salah/learn-go/internal/config"
	"github.com/yahia-salah/learn-go/internal/models"
)

var functions = template.FuncMap{}
var app *config.AppConfig
var pathToTemplates = "./templates"

// Sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// Adds default data on all pages
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	td.CSRFToken = nosurf.Token(r)
	return td
}

// Renders the template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	// get the template cache from the app config
	var tc map[string]*template.Template
	var err error

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
		if err != nil {
			log.Fatal("Can't create cache", err)
		}
		app.TemplateCache = tc
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Can't find cached template:", tmpl)
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	err = t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

// Creates template cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from ./templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		}
		if err != nil {
			return myCache, err
		}

		myCache[name] = ts
	}

	return myCache, nil
}

package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/yahia-salah/learn-go/internal/config"
	"github.com/yahia-salah/learn-go/internal/models"
	"github.com/yahia-salah/learn-go/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func getRoutes() http.Handler {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	app.InProduction = false
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Can't create template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = true
	render.NewTemplates(&app)

	repo := NewRepo(&app)
	NewHandlers(repo)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Use(WriteToConsole)

	mux.Get("/", http.HandlerFunc(Repo.Home))
	mux.Get("/about", http.HandlerFunc(Repo.About))
	mux.Get("/rooms/generals-quarters", http.HandlerFunc(Repo.Generals))
	mux.Get("/rooms/majors-suite", http.HandlerFunc(Repo.Majors))
	mux.Get("/contact", http.HandlerFunc(Repo.Contact))
	mux.Get("/search-availability", http.HandlerFunc(Repo.SearchAvailability))
	mux.Post("/search-availability", http.HandlerFunc(Repo.PostAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(Repo.AvailabilityJSON))
	mux.Get("/make-reservation", http.HandlerFunc(Repo.Reservation))
	mux.Post("/make-reservation", http.HandlerFunc(Repo.PostReservation))
	mux.Get("/reservation-summary", http.HandlerFunc(Repo.ReservationSummary))

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// Writes something to the console on each request
func WriteToConsole(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hit the page:", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// Adds CFRS protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// Loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

// Creates template cache
func CreateTestTemplateCache() (map[string]*template.Template, error) {
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

package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/yahia-salah/learn-go/internal/config"
	"github.com/yahia-salah/learn-go/internal/driver"
	"github.com/yahia-salah/learn-go/internal/handlers"
	"github.com/yahia-salah/learn-go/internal/helpers"
	"github.com/yahia-salah/learn-go/internal/models"
	"github.com/yahia-salah/learn-go/internal/render"
)

const portNumber = ":8080"
const inPROD = false

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	log.Println("Starting server on port " + portNumber)

	db, err := run()

	if db != nil {
		defer db.SQL.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	log.Println("Starting server...")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Can't run server:", err)
	}
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})

	app.InProduction = inPROD

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	helpers.NewHelpers(&app)

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=p@ssw0rd")
	if err != nil {
		log.Fatal("Can't connect to database", err)
		return nil, err
	}

	log.Println("Creating template cache...")
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Can't create template cache", err)
		return db, err
	}
	app.TemplateCache = tc
	app.UseCache = false
	render.NewRenderer(&app)

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	return db, err
}

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yahia-salah/learn-go/pkg/config"
	"github.com/yahia-salah/learn-go/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Use(WriteToConsole)

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	mux.Get("/rooms/generals-quarters", http.HandlerFunc(handlers.Repo.Generals))
	mux.Get("/rooms/majors-suite", http.HandlerFunc(handlers.Repo.Majors))
	mux.Get("/contact", http.HandlerFunc(handlers.Repo.Contact))
	mux.Get("/reservation", http.HandlerFunc(handlers.Repo.Reservation))
	mux.Get("/make-reservation", http.HandlerFunc(handlers.Repo.MakeReservation))

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

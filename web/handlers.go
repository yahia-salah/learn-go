package main

import (
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.page.html")
	log.Println("Successfully handled /")
}

func About(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.page.html")
	log.Println("Successfully handled /about")
}

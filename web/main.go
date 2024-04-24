package main

import (
	"log"
	"net/http"
)

const portNumber = "8080"

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	log.Println("Listing to port: " + portNumber)
	_ = http.ListenAndServe(":"+portNumber, nil)
}

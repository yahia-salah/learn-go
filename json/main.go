package main

import (
	"encoding/json"
	"log"
)

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	HairColor string `json:"hair_color"`
	HasDog    bool   `json:"has_dog"`
}

func main() {
	someJson := `
	[
		{
			"first_name":"Clark",
			"last_name":"Kent",
			"hair_color":"black",
			"has_dog":true
		},
		{
			"first_name":"Bruce",
			"last_name":"Wayne",
			"hair_color":"black",
			"has_dog":false
		}
	]
	`

	var unmarshalled []Person

	err := json.Unmarshal([]byte(someJson), &unmarshalled)

	if err != nil {
		log.Println("Error unmarshalling json", err)
	}

	log.Println(unmarshalled)

	unmarshalled = append(unmarshalled, Person{"Yahia", "Salah", "brown", false})

	marshalled, err := json.MarshalIndent(unmarshalled, "", "")

	if err != nil {
		log.Println("Error marshalling", err)
	}

	log.Println(string(marshalled))
}

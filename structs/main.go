package main

import "fmt"

type contactInfo struct {
	email   string
	zipCode int
}

type person struct {
	firstName string
	lastName  string
	contactInfo
}

func main() {
	p := person{firstName: "Yahia", lastName: "Salah"}

	p.contactInfo.email = "yahiatnt@gmail.com"
	p.contactInfo.zipCode = 11321

	p.updateFirstName("Yunus")

	p.print()
}

func (p person) print() {
	fmt.Printf("%+v", p)
}

func (p *person) updateFirstName(newFirstName string) {
	p.firstName = newFirstName
}

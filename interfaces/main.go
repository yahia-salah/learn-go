package main

import "fmt"

type bot interface {
	getGreerting() string
}
type englishBot struct{}
type spanishBot struct{}

func (englishBot) getGreerting() string {
	return "Hello!"
}

func (spanishBot) getGreerting() string {
	return "Hola!"
}

func printGreeting(b bot) {
	fmt.Println(b.getGreerting())
}

func main() {
	eb := englishBot{}
	sb := spanishBot{}

	printGreeting(eb)
	printGreeting(sb)
}

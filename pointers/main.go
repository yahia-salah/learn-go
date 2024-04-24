package main

import "log"

func main() {
	myString := "Green"

	log.Println("myString=", myString)

	changeUsingPointer(&myString, "Red")

	log.Println("myString=", myString)
}

func changeUsingPointer(s *string, newValue string) {
	*s = newValue
}

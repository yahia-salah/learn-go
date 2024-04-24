package main

import "fmt"

func main() {
	// var colors map[string]string
	// colors:= make(map[string]string)
	colors := map[string]string{
		"red":   "#ff0000",
		"green": "#00ff00",
		"blue":  "#0000ff",
	}

	colors["yellow"] = "#ffff00"

	//delete(colors, "green")

	printMap(colors)
}

func printMap(c map[string]string) {
	for key, value := range c {
		fmt.Printf("key:%v, value:%v\n", key, value)
	}
}

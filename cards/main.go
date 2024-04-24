package main

func main() {
	cards := newDeck()

	// hand, remainingDeck := deal(cards, 5)

	// fmt.Println("Hand Deck:")
	// hand.print()
	// fmt.Println("Remaining Deck")
	// remainingDeck.print()

	// fmt.Println(cards.toString())
	// cards.saveToFile("my_cards.txt")

	// cards2 := newDeckFromFile("my_cards.txt")
	// fmt.Println(cards2.toString())

	cards.shuffle()
	cards.print()
}

func newCard() string {
	return "Five of Diamonds"
}

package main

import "fmt"

type Hand struct {
	cards  []*Card
	coordY int
}

func newHand() *Hand {
	hand := new(Hand)
	// cooordY is 7/8 of the screen height, on the bottom of the screen
	hand.coordY = 7 * (screenHeight / 8)
	hand.cards = make([]*Card, 0)
	return hand

}

func (h *Hand) addCard(newCard *Card) {
	if len(h.cards) == 0 {
		// If no cards in the list, place the new card at the initial position
		newCard.X = screenWidth / 2
		newCard.Y = h.coordY
		h.cards = append(h.cards, newCard)
	} else {
		for _, card := range h.cards {
			card.X -= card.W / 2
			card.Y = h.coordY
		}
		// Place the new card to the right of the last card
		lastCard := h.cards[len(h.cards)-1]
		newCard.X += lastCard.X + lastCard.W
		newCard.Y = h.coordY
		h.cards = append(h.cards, newCard)
		fmt.Println("Cards in hand:", len(h.cards))
	}
}

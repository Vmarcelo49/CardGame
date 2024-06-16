package main

type Hand struct {
	cards  []*Card
	coordY int
}

const handCoordY = 7 * (screenHeight / 8)

func newHand() *Hand {
	return &Hand{
		coordY: handCoordY,
		cards:  nil,
	}
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
		// Place the new card to the side of the last card
		lastCardX := h.cards[len(h.cards)-1].X
		lastCardW := h.cards[len(h.cards)-1].W
		newCard.X += lastCardX + lastCardW
		newCard.Y = h.coordY
		h.cards = append(h.cards, newCard)
		// fmt.Println("Cards in hand:", len(h.cards))
	}
}

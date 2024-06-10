package main

type Duel struct {
	//player1Health, player2Health, turnCount int
	//featuredCard                            *ebiten.Image
	p1Hand, p2Hand *Hand
	p1Deck, p2Deck *Deck
	//p1GY, p2GY                              *GY
	field *Field // Precisa ser refeito para os dois jogadores

}

func newDuel() *Duel {
	placeHolderDeck := "Deck/testDeck.txt"

	duel := new(Duel)

	duel.p1Deck = newDeck(placeHolderDeck, 1)
	duel.p2Deck = newDeck(placeHolderDeck, 2)
	duel.field = newField()

	// TODO: implementar p2 e o rps
	duel.p1Hand = newHand()
	duel.p2Hand = newHand()

	return duel
}

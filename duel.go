package main

type Duel struct {
	//player1Health, player2Health, turnCount int
	//featuredCard                            *ebiten.Image
	p1Hand, p2Hand *Hand
	p1Deck, p2Deck *Deck
	//p1GY, p2GY                              *GY
	field *Field // Precisa ser refeito para os dois jogadores

}

const (
	placeHolderDeck = "Deck/testDeck.txt"
)

func (g *Game) newDuel() *Duel {

	duel := new(Duel)

	duel.p1Deck = newDeck(placeHolderDeck, 1, g.texMap[-1].Bounds().Dx(), g.texMap[-1].Bounds().Dy())
	duel.p2Deck = newDeck(placeHolderDeck, 2, g.texMap[-1].Bounds().Dx(), g.texMap[-1].Bounds().Dy())
	duel.field = newField()

	// TODO: implementar p2 e o rps
	duel.p1Hand = newHand()
	duel.p2Hand = newHand()

	return duel
}

package main

import "fmt"

type whichTurnPlayer bool

const (
	player   whichTurnPlayer = true
	opponent whichTurnPlayer = false
)

type Gamestate struct {
	P1, P2 *Player
	Field  Field

	TurnCount         int
	CurrentPlayerTurn whichTurnPlayer
}

func (g *Gamestate) update() {
	//help
}

// equals compares two gamestates, if they are different, it returns false
func (gs *Gamestate) equals(other *Gamestate) bool {
	// Comparar HP dos jogadores
	if gs.P1.HP != other.P1.HP || gs.P2.HP != other.P2.HP {
		return false
	}

	// Comparar Decks
	if len(gs.P1.Deck) != len(other.P1.Deck) || len(gs.P2.Deck) != len(other.P2.Deck) {
		return false
	}
	for i := range gs.P1.Deck {
		if gs.P1.Deck[i].ID != other.P1.Deck[i].ID {
			return false
		}
	}
	for i := range gs.P2.Deck {
		if gs.P2.Deck[i].ID != other.P2.Deck[i].ID {
			return false
		}
	}

	// Comparar Mãos
	if len(gs.P1.Hand) != len(other.P1.Hand) || len(gs.P2.Hand) != len(other.P2.Hand) {
		return false
	}
	for i := range gs.P1.Hand {
		if gs.P1.Hand[i].ID != other.P1.Hand[i].ID {
			return false
		}
	}
	for i := range gs.P2.Hand {
		if gs.P2.Hand[i].ID != other.P2.Hand[i].ID {
			return false
		}
	}

	// Comparar Cemitérios
	if len(gs.P1.GY) != len(other.P1.GY) || len(gs.P2.GY) != len(other.P2.GY) {
		return false
	}
	for i := range gs.P1.GY {
		if gs.P1.GY[i].ID != other.P1.GY[i].ID {
			return false
		}
	}
	for i := range gs.P2.GY {
		if gs.P2.GY[i].ID != other.P2.GY[i].ID {
			return false
		}
	}

	// Comparar Campos
	if len(gs.Field.P1) != len(other.Field.P1) || len(gs.Field.P2) != len(other.Field.P2) {
		return false
	}
	for i := range gs.Field.P1 {
		if gs.Field.P1[i].ID != other.Field.P1[i].ID {
			return false
		}
	}
	for i := range gs.Field.P2 {
		if gs.Field.P2[i].ID != other.Field.P2[i].ID {
			return false
		}
	}

	// Comparar Turno
	if gs.TurnCount != other.TurnCount {
		return false
	}
	if gs.CurrentPlayerTurn != other.CurrentPlayerTurn {
		return false
	}

	return true
}

func copyGamestate(gs *Gamestate) *Gamestate {
	newGs := &Gamestate{
		P1:                copyPlayer(gs.P1),
		P2:                copyPlayer(gs.P2),
		Field:             copyField(gs.Field),
		TurnCount:         gs.TurnCount,
		CurrentPlayerTurn: gs.CurrentPlayerTurn,
	}
	return newGs
}

func copyPlayer(p *Player) *Player {
	newP := &Player{
		HP:   p.HP,
		Hand: make([]*Card, len(p.Hand)),
		GY:   make([]*Card, len(p.GY)),
		Deck: make([]*Card, len(p.Deck)),
	}
	for i, c := range p.Hand {
		newP.Hand[i] = copyCard(c) // Copia profunda da carta
	}
	for i, c := range p.GY {
		newP.GY[i] = copyCard(c) // Copia profunda da carta
	}
	for i, c := range p.Deck {
		newP.Deck[i] = copyCard(c) // Copia profunda da carta
	}
	return newP
}

func copyCard(c *Card) *Card {
	newCard := newCardFromID(c.ID)
	newCard.X = c.X
	newCard.Y = c.Y
	newCard.Selected = c.Selected
	return newCard
}

func copyField(f Field) Field {
	newF := Field{
		P1: make([]*Card, len(f.P1)),
		P2: make([]*Card, len(f.P2)),
	}
	for i, c := range f.P1 {
		newF.P1[i] = copyCard(c) // Copia profunda da carta
	}
	for i, c := range f.P2 {
		newF.P2[i] = copyCard(c) // Copia profunda da carta
	}
	return newF
}

// Creates a new gamestate. Returns errors if any deck file is invalid.
func newGameState(deckPathP1, deckPathP2 string) (*Gamestate, error) {
	gamestate := &Gamestate{}
	deck1, err := newDeck(deckPathP1)
	if err != nil {
		return nil, err
	}
	gamestate.P1 = &Player{
		HP:   100,
		Deck: deck1,
	}
	deck2, err := newDeck(deckPathP2)
	if err != nil {
		return nil, err
	}
	gamestate.P2 = &Player{
		HP:   100,
		Deck: deck2,
	}
	return gamestate, nil
}

type Player struct {
	HP   int
	Hand []*Card
	GY   []*Card
	Deck []*Card
}

func (p *Player) drawCard() {
	if len(p.Deck) == 0 {
		fmt.Println("Deck is empty")
		return
	}
	p.Hand = append(p.Hand, p.Deck[0])
	p.Deck = p.Deck[1:]
}

type Field struct {
	P1, P2 []*Card
}

// Returns all cards from both players fields
func (f *Field) AllCards() []*Card {
	all := make([]*Card, 0, len(f.P1)+len(f.P2))
	all = append(all, f.P1...) // Adiciona todas as cartas de P1
	all = append(all, f.P2...) // Adiciona todas as cartas de P2
	return all
}

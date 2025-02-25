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
	if !gs.P1.equals(other.P1) || !gs.P2.equals(other.P2) {
		return false
	}

	if !gs.Field.equals(&other.Field) {
		return false
	}

	if gs.TurnCount != other.TurnCount || gs.CurrentPlayerTurn != other.CurrentPlayerTurn {
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

func (p *Player) modifyHP(amount int) {
	originalValue := p.HP
	p.HP += amount
	fmt.Printf("Player HP value modified was %d, now is %d", originalValue, p.HP)
}

func (p *Player) drawCard() {
	if len(p.Deck) == 0 {
		fmt.Println("Deck is empty")
		return
	}
	p.Hand = append(p.Hand, p.Deck[0])
	p.Deck = p.Deck[1:]
}

func (p *Player) equals(other *Player) bool {
	if p.HP != other.HP {
		return false
	}
	if !compareCardSlices(p.Deck, other.Deck) || !compareCardSlices(p.Hand, other.Hand) || !compareCardSlices(p.GY, other.GY) {
		return false
	}
	return true
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

func (f *Field) equals(other *Field) bool {
	if !compareCardSlices(f.P1, other.P1) || !compareCardSlices(f.P2, other.P2) {
		return false
	}
	return true
}

func compareCardSlices(a, b []*Card) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].ID != b[i].ID {
			return false
		}
	}
	return true
}

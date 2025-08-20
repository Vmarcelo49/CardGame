package state

import (
	"fmt"

	"github.com/Vmarcelo49/CardGame/pkg/card"
)

type WhichTurnPlayer bool

const (
	Player   WhichTurnPlayer = true
	Opponent WhichTurnPlayer = false
)

// Gamestate holds the current game state
type Gamestate struct {
	P1, P2            *Player
	Field             Field
	TurnCount         int
	CurrentPlayerTurn WhichTurnPlayer
}

// Update updates the game state logic
func (g *Gamestate) Update() {
	// Game logic updates go here
}

// Equals compares two gamestates
func (gs *Gamestate) Equals(other *Gamestate) bool {
	if !gs.P1.Equals(other.P1) || !gs.P2.Equals(other.P2) {
		return false
	}

	if !gs.Field.Equals(&other.Field) {
		return false
	}

	if gs.TurnCount != other.TurnCount || gs.CurrentPlayerTurn != other.CurrentPlayerTurn {
		return false
	}

	return true
}

// GetP1 returns player 1
func (gs *Gamestate) GetP1() *Player {
	return gs.P1
}

// GetP2 returns player 2
func (gs *Gamestate) GetP2() *Player {
	return gs.P2
}

// GetField returns the game field
func (gs *Gamestate) GetField() *Field {
	return &gs.Field
}

// CopyGamestate creates a deep copy of the gamestate
func CopyGamestate(gs *Gamestate) *Gamestate {
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
		Hand: make([]*card.Card, len(p.Hand)),
		GY:   make([]*card.Card, len(p.GY)),
		Deck: make([]*card.Card, len(p.Deck)),
	}
	for i, c := range p.Hand {
		newP.Hand[i] = copyCard(c)
	}
	for i, c := range p.GY {
		newP.GY[i] = copyCard(c)
	}
	for i, c := range p.Deck {
		newP.Deck[i] = copyCard(c)
	}
	return newP
}

func copyCard(c *card.Card) *card.Card {
	newCard := card.NewCardFromID(c.ID)
	newCard.X = c.X
	newCard.Y = c.Y
	newCard.Selected = c.Selected
	return newCard
}

func copyField(f Field) Field {
	newF := Field{
		P1: make([]*card.Card, len(f.P1)),
		P2: make([]*card.Card, len(f.P2)),
	}
	for i, c := range f.P1 {
		newF.P1[i] = copyCard(c)
	}
	for i, c := range f.P2 {
		newF.P2[i] = copyCard(c)
	}
	return newF
}

// NewGameState creates a new gamestate with decks from files
func NewGameState(deckPathP1, deckPathP2 string, newDeckFunc func(string) ([]*card.Card, error)) (*Gamestate, error) {
	gamestate := &Gamestate{}
	deck1, err := newDeckFunc(deckPathP1)
	if err != nil {
		return nil, err
	}
	gamestate.P1 = &Player{
		HP:   100,
		Deck: deck1,
	}
	deck2, err := newDeckFunc(deckPathP2)
	if err != nil {
		return nil, err
	}
	gamestate.P2 = &Player{
		HP:   100,
		Deck: deck2,
	}
	return gamestate, nil
}

// Player represents a game player
type Player struct {
	HP   int
	Hand []*card.Card
	GY   []*card.Card
	Deck []*card.Card
}

// ModifyHP modifies the player's HP
func (p *Player) ModifyHP(amount int) {
	originalValue := p.HP
	p.HP += amount
	fmt.Printf("Player HP value modified was %d, now is %d", originalValue, p.HP)
}

// DrawCard draws a card from deck to hand
func (p *Player) DrawCard() {
	if len(p.Deck) == 0 {
		fmt.Println("Deck is empty")
		return
	}
	p.Hand = append(p.Hand, p.Deck[0])
	p.Deck = p.Deck[1:]
}

// GetHand returns the player's hand
func (p *Player) GetHand() []*card.Card {
	return p.Hand
}

// GetDeck returns the player's deck
func (p *Player) GetDeck() []*card.Card {
	return p.Deck
}

// GetGY returns the player's graveyard
func (p *Player) GetGY() []*card.Card {
	return p.GY
}

// Equals compares two players
func (p *Player) Equals(other *Player) bool {
	if p.HP != other.HP {
		return false
	}
	if !compareCardSlices(p.Deck, other.Deck) || !compareCardSlices(p.Hand, other.Hand) || !compareCardSlices(p.GY, other.GY) {
		return false
	}
	return true
}

// Field represents the game field
type Field struct {
	P1, P2 []*card.Card
}

// AllCards returns all cards from both players' fields
func (f *Field) AllCards() []*card.Card {
	all := make([]*card.Card, 0, len(f.P1)+len(f.P2))
	all = append(all, f.P1...)
	all = append(all, f.P2...)
	return all
}

// GetP1 returns player 1's field cards
func (f *Field) GetP1() []*card.Card {
	return f.P1
}

// GetP2 returns player 2's field cards
func (f *Field) GetP2() []*card.Card {
	return f.P2
}

// Equals compares two fields
func (f *Field) Equals(other *Field) bool {
	if !compareCardSlices(f.P1, other.P1) || !compareCardSlices(f.P2, other.P2) {
		return false
	}
	return true
}

func compareCardSlices(a, b []*card.Card) bool {
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

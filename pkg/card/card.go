package card

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// CardType represents the type of a card
type CardType int

const (
	Creature CardType = iota
	Spell
	Permanent
	Dominion
)

// Keyword represents card keywords using bit flags
type Keyword uint8

const (
	None     Keyword = 0
	Attacker         = 1 << iota
	Piercer
	Blocker
)

// cardFlags represents flags that control card behavior
type cardFlags uint

const (
	CannotBeUsed        cardFlags = 0
	CanBeNormalSummoned cardFlags = 1 << iota
	CanBeUsedSpeed1
	CanBeUsedSpeed2
)

// Stats represents attack and life points
type Stats struct {
	Attack int
	Life   int
}

// Card represents a game card
type Card struct {
	Name     string
	ID       int
	CType    CardType
	SubType  string
	Text     string
	Effect   func() error
	Stats    Stats
	Keywords Keyword
	Flags    cardFlags
	X, Y     int
	Selected bool
}

// Damageable interface for entities that can take damage
type Damageable interface {
	ModifyHP(amount int)
}

// ModifyHP modifies the card's life points
func (c *Card) ModifyHP(amount int) {
	originalValue := c.Stats.Life
	c.Stats.Life += amount
	fmt.Printf("Card HP value modified was %d, now is %d", originalValue, c.Stats.Life)
}

// GiveKeyword adds a keyword to a creature card
func GiveKeyword(key Keyword, target *Card) error {
	if target.CType != Creature {
		return fmt.Errorf("card is not a creature")
	}
	target.Keywords |= key
	return nil
}

// GetLocation returns the location of the card in the game state
// This function will be moved to the logic package to avoid circular dependencies
func (c *Card) GetLocation(playerHands, playerDecks, playerGYs, fieldP1, fieldP2 [][]*Card) string {
	// For now, return unknown - this will be implemented in the logic package
	return "UNKNOWN"
}

// CheckCardClicked checks if the mouse position is over the card
func CheckCardClicked(c *Card, mouseX, mouseY int) bool {
	// Card dimensions (these should be configurable)
	cardWidth, cardHeight := 100, 140

	return mouseX >= c.X && mouseX <= c.X+cardWidth &&
		mouseY >= c.Y && mouseY <= c.Y+cardHeight
}

// CreateCardImage creates a visual representation of a card
func CreateCardImage(cardFrameImage *ebiten.Image, c *Card) (*ebiten.Image, error) {
	if cardFrameImage == nil {
		return nil, fmt.Errorf("card frame image is nil")
	}

	// For now, return the frame image
	// In a full implementation, this would add text and other card details
	return cardFrameImage, nil
}

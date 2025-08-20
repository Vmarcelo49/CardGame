package logic

import (
	"fmt"
	"log"

	"github.com/Vmarcelo49/CardGame/pkg/card"
	"github.com/Vmarcelo49/CardGame/pkg/state"
	"github.com/Vmarcelo49/CardGame/pkg/ui"
)

// GameLogic handles the main game logic updates
type GameLogic struct {
	inputHandler *ui.InputHandler
}

// NewGameLogic creates a new game logic handler
func NewGameLogic() *GameLogic {
	return &GameLogic{
		inputHandler: ui.NewInputHandler(),
	}
}

// UpdateSelectCard handles card selection logic
func UpdateSelectCard(gamestate *state.Gamestate, mouse *ui.Mouse) {
	var selectedCard *card.Card

	for _, c := range gamestate.P1.Hand {
		if c.Selected {
			selectedCard = c
			break
		}
	}

	// Check if a new card was clicked
	for _, c := range gamestate.P1.Hand {
		if card.CheckCardClicked(c, mouse.X, mouse.Y) && mouse.LeftPressed {
			if selectedCard == nil {
				fmt.Println("Card selected:", c.Name)
				c.Selected = true
				NewCardClickedFunc(c, gamestate)
			} else if selectedCard != c {
				fmt.Println("Card deselected:", selectedCard.Name)
				selectedCard.Selected = false
				fmt.Println("Card selected:", c.Name)
				c.Selected = true
			}
			return
		}
	}

	// If clicked elsewhere, deselect the current card
	if mouse.LeftPressed && selectedCard != nil && !card.CheckCardClicked(selectedCard, mouse.X, mouse.Y) {
		fmt.Println("Card deselected:", selectedCard.Name)
		selectedCard.Selected = false
	}
}

// UpdateGameLogic performs the main game logic update cycle
func UpdateGameLogic(gamestate, previousGamestate *state.Gamestate, mouse *ui.Mouse,
	buttons []*ui.Button, inputHandler *ui.InputHandler,
	drawCardFunc func(string, uint), exitDuelFunc func()) (*state.Gamestate, error) {

	// Check input
	if err := ui.CheckInput(mouse, buttons, inputHandler, drawCardFunc, exitDuelFunc); err != nil {
		log.Println(err)
		return nil, err
	}

	// Update card selection
	UpdateSelectCard(gamestate, mouse)

	// Update game state
	gamestate.Update()

	// Update renderer if game state changed
	if previousGamestate != nil && !gamestate.Equals(previousGamestate) {
		log.Println("Gamestate changed")
	}

	// Return a copy of the current game state
	return state.CopyGamestate(gamestate), nil
}

// GetCardLocation returns the location of a card in the game state
func GetCardLocation(c *card.Card, gs *state.Gamestate) string {
	// Check P1 locations
	for _, card := range gs.P1.Hand {
		if c == card {
			return "P1HAND"
		}
	}
	for _, card := range gs.P1.Deck {
		if c == card {
			return "P1DECK"
		}
	}
	for _, card := range gs.P1.GY {
		if c == card {
			return "P1GY"
		}
	}
	for _, card := range gs.Field.P1 {
		if c == card {
			return "P1FIELD"
		}
	}

	// Check P2 locations
	for _, card := range gs.P2.Hand {
		if c == card {
			return "P2HAND"
		}
	}
	for _, card := range gs.P2.Deck {
		if c == card {
			return "P2DECK"
		}
	}
	for _, card := range gs.P2.GY {
		if c == card {
			return "P2GY"
		}
	}
	for _, card := range gs.Field.P2 {
		if c == card {
			return "P2FIELD"
		}
	}

	return "UNKNOWN"
}

// NewCardClickedFunc handles logic when a card is clicked
func NewCardClickedFunc(c *card.Card, gs *state.Gamestate) {
	location := GetCardLocation(c, gs)

	if c.CType == card.Creature && location == "P1HAND" {
		fmt.Printf("Creature card %s clicked in hand\n", c.Name)
	}
	if c.CType == card.Spell && location == "P1HAND" {
		fmt.Printf("Spell card %s clicked in hand\n", c.Name)
	}
}

// EncapsulateButtonFunc positions a button relative to a card
func EncapsulateButtonFunc(c *card.Card, button *ui.Button) {
	button.X = c.X
	button.Y = c.Y - button.H

	originalFunc := button.Function
	button.Function = func() error {
		button.AlreadyClicked = true
		button.X = -5000 // Move off screen
		if originalFunc != nil {
			return originalFunc()
		}
		return nil
	}
}

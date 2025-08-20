package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// Mouse represents the mouse state
type Mouse struct {
	X, Y          int
	LeftPressed   bool
	RightPressed  bool
	MiddlePressed bool
}

// UpdateMouseState updates the mouse state with current input
func (m *Mouse) UpdateMouseState() {
	m.X, m.Y = ebiten.CursorPosition()
	m.LeftPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	m.RightPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	m.MiddlePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
}

// InputHandler handles keyboard input
type InputHandler struct {
	keyStates map[ebiten.Key]bool
}

// NewInputHandler creates a new input handler
func NewInputHandler() *InputHandler {
	return &InputHandler{
		keyStates: make(map[ebiten.Key]bool),
	}
}

// HandleKeyPress handles a key press with the given action
func (ih *InputHandler) HandleKeyPress(key ebiten.Key, action func()) {
	if ebiten.IsKeyPressed(key) {
		if !ih.keyStates[key] {
			action()
		}
		ih.keyStates[key] = true
	} else {
		ih.keyStates[key] = false
	}
}

// CheckInput processes input for duel buttons and keyboard
func CheckInput(mouse *Mouse, buttons []*Button, inputHandler *InputHandler,
	drawCardFunc func(string, uint),
	exitDuelFunc func()) error {

	for _, button := range buttons {
		if err := button.CheckClicked(mouse); err != nil {
			return err
		}
	}

	inputHandler.HandleKeyPress(ebiten.KeySpace, func() {
		drawCardFunc("player", 1)
	})

	inputHandler.HandleKeyPress(ebiten.KeyQ, func() {
		drawCardFunc("opp", 1)
		fmt.Println("Drew card for opponent")
	})

	inputHandler.HandleKeyPress(ebiten.KeyEscape, func() {
		exitDuelFunc()
	})

	return nil
}

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Mouse struct {
	X, Y          int
	LeftPressed   bool
	RightPressed  bool
	MiddlePressed bool
}

func (m *Mouse) UpdateMouseState() {
	m.X, m.Y = ebiten.CursorPosition()
	m.LeftPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	m.RightPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	m.MiddlePressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
}

func (g *Game) keyboardInput() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.duel.p1Deck.drawCard(g.duel.p1Hand)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.loadMainMenu()
	}
	return nil
}

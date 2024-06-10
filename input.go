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

func keyboardInput(deck *Deck, hand *Hand) error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		deck.drawCard(hand)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		//maluquice do caralho, pq retornar um erro pra fechar o jogo?
		// TODO: Mudar para voltar ao menu depois
		return ebiten.Termination
	}
	return nil
}

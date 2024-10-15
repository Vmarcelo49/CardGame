package main

import (
	"fmt"

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

func (g *Game) checkInput() error {
	for _, button := range g.duelButtons {
		if err := button.checkClicked(g.mouse); err != nil {
			return err
		}
	}
	g.checkKey(ebiten.KeySpace, func() {
		drawCard(g.gamestate, "player", 1)
	})
	g.checkKey(ebiten.KeyQ, func() {
		drawCard(g.gamestate, "opp", 1)
		fmt.Println(len(g.gamestate.P2.Hand))
	})

	g.checkKey(ebiten.KeyEscape, func() {
		g.exitingDuel = true
		g.currentScene = MainMenu
	})
	return nil
}

func (g *Game) checkKey(key ebiten.Key, action func()) {
	if ebiten.IsKeyPressed(key) {
		if !g.keyStates[key] {
			action() // Executa o comando uma vez
		}
		g.keyStates[key] = true
	} else {
		g.keyStates[key] = false
	}
}

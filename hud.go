package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type HUDDuel struct {
	player1Health, player2Health int
	featuredCard                 *ebiten.Image
}

type Button struct {
	x, y, width, height int
	image               *ebiten.Image
	text                string
	function            func()
}

func (b *Button) in(x, y int) bool {
	return x > b.x && x < b.x+b.width && y > b.y && y < b.y+b.height
}

func createButtonSlice() []*Button {
	var buttonSlice []*Button
	for i := 0; i < 4; i++ {
		buttonSlice = append(buttonSlice, &Button{x: screenWidth/2 - 100, y: screenHeight/2 + 50*i, width: 200, height: 50, text: "Button " + string(i)})
	}
	return buttonSlice
}

type HUDMainMenu struct {
	button []*Button
}

func (g *Game) drawMainMenu(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Main Menu")
	for _, b := range g.hudMainMenu.button {

	}
}

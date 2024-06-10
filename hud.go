package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type HUDDuel struct {
	player1Health, player2Health int
	featuredCard                 *ebiten.Image
}

type Button struct {
	x, y, w, h    int
	image         *ebiten.Image
	function      func() error
	alreadClicked bool
}

func newButton(x, y int, texto string, function func() error) *Button {
	newImage := ebiten.NewImage(screenWidth/8, screenHeight/8)
	newImage.Fill(color.White)

	// draw text on the image
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(10, 10)
	textOp.ColorScale.ScaleWithColor(color.Black)

	text.Draw(newImage, texto, &text.GoTextFace{
		Source: font,
		Size:   20.0,
	}, textOp)

	return &Button{x, y, screenWidth / 8, screenHeight / 8, newImage, function, false}
}

func (b *Button) checkClicked(m *Mouse) error {
	if m.X > b.x && m.X < b.x+b.w && m.Y > b.y && m.Y < b.y+b.h && m.LeftPressed {
		b.alreadClicked = true
		fmt.Println("Button clicked")
		return b.function()
	}
	return nil
}

func (b *Button) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(b.image, op)
}

func addButton(buttonSlice []*Button, text string, function func() error) []*Button {
	buttonX := screenWidth / 8
	buttonY := screenHeight / 8
	x := (screenWidth - buttonX) / 2
	y := (screenHeight-buttonY)/2 + len(buttonSlice)*(buttonY+10)
	buttonSlice = append(buttonSlice, newButton(x, y, text, function))
	return buttonSlice
}

// Cria os botões do menu principal
func (g *Game) createButtons() ([]*Button, error) {
	var buttons []*Button
	buttons = addButton(buttons, "Duel", func() error {
		g.loadDuelMode()
		g.mainMenuButtons = nil // Go doesn clear the buttons, so we need to do it manually

		return nil
	})
	buttons = addButton(buttons, "Deck Editor", func() error {
		fmt.Println("Soon...")
		return nil
	})
	buttons = addButton(buttons, "Exit", func() error {
		return ebiten.Termination
	})

	return buttons, nil

}

func (g *Game) DrawMainMenu(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	ebitenutil.DebugPrint(screen, "Main Menu")
	for _, b := range g.mainMenuButtons {
		b.draw(screen)
	}
}

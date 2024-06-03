package main

import (
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
	x, y     int
	image    *ebiten.Image
	function func()
}

func newButton(x, y int, texto string, function func()) *Button {
	newImage := ebiten.NewImage(screenWidth/8, screenHeight/8)
	newImage.Fill(color.White)

	// draw text on the image
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(10, 10)
	textOp.ColorScale.ScaleWithColor(color.Black)

	text.Draw(newImage, texto, &text.GoTextFace{
		Source: font,
		Size:   fontSize,
	}, textOp)

	return &Button{x, y, newImage, function}
}

func (b *Button) checkClicked(m *Mouse) {
	if m.X > b.x && m.X < b.x+b.image.Bounds().Dy() && m.Y > b.y && m.Y < b.y+b.image.Bounds().Dx() && m.LeftPressed == true {
		b.function()
		return
	}
	return
}

func (b *Button) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(b.image, op)
}

func addButton(buttonSlice []*Button, image *ebiten.Image, text string, function func()) []*Button {
	x := (screenWidth - image.Bounds().Dx()) / 2
	y := (screenHeight-image.Bounds().Dy())/2 + len(buttonSlice)*(image.Bounds().Dy()+10)
	buttonSlice = append(buttonSlice, newButton(x, y, text, function))
	return buttonSlice
}

func (g *Game) DrawMainMenu(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	ebitenutil.DebugPrint(screen, "Main Menu")
	for _, b := range g.mainMenuButtons {
		b.draw(screen)
	}
}

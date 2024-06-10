package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var contador int

type Field struct {
	player1Field []*Card
	player2Field []*Card

	player1Graveyard []*Card
	player2Graveyard []*Card

	middlescreen float32
}

func newField() *Field {
	field := new(Field)

	field.middlescreen = (screenHeight / 2) - (screenHeight / 18)

	return field
}

func (f *Field) draw(screen *ebiten.Image, textureMap map[int]*ebiten.Image) {
	if len(f.player1Field) > 0 {
		contador = f.player1Field[0].W
	}
	for i, card := range f.player1Field {
		card.X = (screenHeight / 2) + (contador * i)
		card.Y = int(f.middlescreen)
		card.draw(screen, textureMap[card.ID])
	}

	for _, card := range f.player2Field {
		// but where lol TODO: fix player 2
		card.draw(screen, textureMap[card.ID])
	}

	vector.DrawFilledRect(screen, 0, f.middlescreen, screenWidth, 2, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
}

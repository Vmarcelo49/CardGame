package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var valor int

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

func (f *Field) draw(screen *ebiten.Image) {
	if len(f.player1Field) > 0 {
		valor = f.player1Field[0].W
	}
	for i, card := range f.player1Field {
		card.X = (screenHeight / 2) + (valor * i)
		card.Y = int(f.middlescreen)
		card.draw(screen)
	}

	for _, card := range f.player2Field {
		card.draw(screen)
	}

	vector.DrawFilledRect(screen, 0, f.middlescreen, screenWidth, 2, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)
}

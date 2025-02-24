package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//CardRes = 1500x2100
	//CardArt = 1400x1050
	//scalingFactor    float64
	border           = 20.0
	margemHorizontal = 50.0
	margemVertical   = 250.0
	fontSize         = 20.0 * 10 // fits 14 characters
	maxTextHeight    = 125.0
)

//type Effect string

type CardType int

const (
	Creature CardType = iota
	Spell
	Permanent
	Dominion
)

type Stats struct {
	Attack int
	Life   int
}

type Card struct {
	// things that are parsed from the yaml file
	Name     string
	ID       int
	CType    CardType
	SubType  string
	Text     string
	Effect   func()
	Stats    Stats //can be nil
	Keywords Keyword
	Flags    cardFlags
	// Render related
	X, Y int
	//W, H     int
	Selected bool
}

// Checa se X e Y estão dentro da carta
func (g *Game) checkCardClicked(c *Card) bool {
	return g.mouse.X > c.X && g.mouse.X < c.X+g.duelRenderer.cardSizeW && g.mouse.Y > c.Y && g.mouse.Y < c.Y+g.duelRenderer.cardSizeH
}

// Cria uma imagem de uma carta, com o background com o cardframe fornecido, em seu tamanho normal.
func createCardImage(cardFrameIm *ebiten.Image, card *Card) (*ebiten.Image, error) {
	pathCardArt := fmt.Sprint("./assets/image/cardArt/", card.ID, ".png")
	cardArt, err := newImageFromPath(pathCardArt)
	if err != nil {
		return nil, err
	}
	// Imagem que será combinada com o frame da carta
	image := ebiten.NewImage(cardFrameIm.Bounds().Dx(), cardFrameIm.Bounds().Dy())

	// Draw no card frame primeiro
	op := &ebiten.DrawImageOptions{}
	image.DrawImage(cardFrameIm, op)
	op.GeoM.Translate(margemHorizontal, margemVertical)
	image.DrawImage(cardArt, op)
	op.GeoM.Reset()

	// Draw no nome da carta
	op.GeoM.Translate(border, border)
	fontSize = float64(cardFrameIm.Bounds().Dx()) / 10 // this image is huge, so the font size is also huge
	textImg := newTextImageMultiline(card.Name, color.White, fontSize, int(cardFrameIm.Bounds().Dx())-int(border), int(maxTextHeight))
	image.DrawImage(textImg, op)
	op.GeoM.Reset()

	// Draw no texto da carta
	op.GeoM.Translate(border, float64(cardArt.Bounds().Dy()+int(margemVertical))+border)
	effImg := newTextImageMultiline(card.Text, color.White, fontSize, int(cardFrameIm.Bounds().Dx())-int(border), int(maxTextHeight))
	image.DrawImage(effImg, op)
	op.GeoM.Reset()

	return image, nil
}

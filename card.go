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
	Effect   func() error
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
	textImg := newTextImageMultiline(card.Name, color.White, fontSize, int(cardFrameIm.Bounds().Dx())-int(border))
	image.DrawImage(textImg, op)
	op.GeoM.Reset()

	// Draw no texto da carta
	op.GeoM.Translate(border, float64(cardArt.Bounds().Dy()+int(margemVertical))+border)
	effImg := newTextImageMultiline(card.Text, color.White, fontSize, int(cardFrameIm.Bounds().Dx())-int(border))
	image.DrawImage(effImg, op)
	op.GeoM.Reset()

	return image, nil
}

func (c *Card)modifyHP(amount int){
	originalValue := c.Stats.Life
	c.Stats.Life += amount
	fmt.Printf("Card HP value modified was %d, now is %d"originalValue,c.Stats.Life)
}

func (c *Card) getLocation(gs Gamestate) string {
	// Check P1 locations
	for _, card := range gs.P1.Hand {
		if c == card {
			return "P1HAND"
		}
	}
	for _, card := range gs.P1.Deck {
		if c == card {
			return "P1DECK"
		}
	}
	for _, card := range gs.P1.GY {
		if c == card {
			return "P1GY"
		}
	}
	for _, card := range gs.Field.P1 {
		if c == card {
			return "P1FIELD"
		}
	}

	// Check P2 locations
	for _, card := range gs.P2.Hand {
		if c == card {
			return "P2HAND"
		}
	}
	for _, card := range gs.P2.Deck {
		if c == card {
			return "P2DECK"
		}
	}
	for _, card := range gs.P2.GY {
		if c == card {
			return "P2GY"
		}
	}
	for _, card := range gs.Field.P2 {
		if c == card {
			return "P2FIELD"
		}
	}

	// Default return if card is not found
	return "UNKNOWN"
}
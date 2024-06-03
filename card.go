package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"gopkg.in/yaml.v3"
)

var (
	//CardRes = 1500x2100
	//CardArt = 1400x1050
	scalingFactor    float64
	border           = 20.0
	margemHorizontal = 50.0
	margemVertical   = 250.0
	fontSize         = 20.0 * 10 // fits 14 characters
)

type Effect string

type CardType int

const (
	Creature CardType = iota
	Spell
	Dominion
)

type Stats struct {
	Attack              int  `yaml:"attack"`
	Life                int  `yaml:"life"`
	CanBeNormalSummoned bool `yaml:"CanBeNormalSummoned"`
}

type Card struct {
	// things that are parsed from the yaml file
	Name    string   `yaml:"Name"`
	ID      int      `yaml:"ID"`
	CType   CardType `yaml:"Type"`
	SubType string   `yaml:"Subtype"`
	Text    string   `yaml:"Text"`
	Effects []Effect `yaml:"Effects"`
	Stats   Stats    `yaml:"Stats"` //can be nil
	// Card effect related
	Keywords []Keyword
	// more gerenal stuff
	X, Y           int
	W, H           int
	ScaleX, ScaleY float64
	Selected       bool
	// ebiten stuff
	Image *ebiten.Image //should be removed later, check Draw
}

func (c *Card) scale(x, y float64) {
	c.ScaleX = x
	c.ScaleY = y
}

func (c *Card) moveTo(x, y int) {
	c.X = x
	c.Y = y

}

func (c *Card) in(x, y int) bool {
	return x > c.X && x < c.X+int(float64(c.Image.Bounds().Dx())*c.ScaleX) && y > c.Y && y < c.Y+int(float64(c.Image.Bounds().Dy())*c.ScaleY)
}

func (c *Card) resetScale() {
	c.ScaleX = scalingFactor
	c.ScaleY = scalingFactor
}

// function to draw the card, first translate to the position of the card, then draw the image
func (c *Card) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.ScaleX, c.ScaleY) //always scale first
	op.GeoM.Translate(float64(c.X), float64(c.Y))
	screen.DrawImage(c.Image, op)
}

func createCard(card *Card, cardFrameIm *ebiten.Image, pathCardArt string) *Card {
	card.ScaleX = scalingFactor
	card.ScaleY = scalingFactor

	card.W = int(float64(cardFrameIm.Bounds().Dx()) * card.ScaleX)
	card.H = int(float64(cardFrameIm.Bounds().Dy()) * card.ScaleY)

	CardArt, _, err := ebitenutil.NewImageFromFile(pathCardArt)
	if err != nil {
		log.Panic(err)
	}

	// Create a new image to hold the combined result
	card.Image = ebiten.NewImage(cardFrameIm.Bounds().Dx(), cardFrameIm.Bounds().Dy())

	// Draw no card frame primeiro
	op := &ebiten.DrawImageOptions{}
	card.Image.DrawImage(cardFrameIm, op)
	op.GeoM.Translate(margemHorizontal, margemVertical)
	card.Image.DrawImage(CardArt, op)

	// text
	// todo: justify text
	if len(card.Name) > 14 {
		fontSize = 15.0 * 10 // fits 22 characters
	}

	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(border, border)
	textOp.ColorScale.ScaleWithColor(color.White)

	text.Draw(card.Image, card.Name, &text.GoTextFace{
		Source: font,
		Size:   fontSize,
	}, textOp)

	// reset font size
	fontSize = 20.0 * 10

	return card
}

func newCardFromID(id int) (*Card, error) {
	card, err := parseCard(id)
	cardFrameIm := new(ebiten.Image)
	if card.CType == 0 { // will be changed to handle other types of cards
		cardFrameIm, _, err = ebitenutil.NewImageFromFile("Image/CardFrame/CardFrame.png")
		if err != nil {
			return nil, err
		}
	}
	cardArtPath := fmt.Sprint("Image/CardArt/", card.ID, ".png")

	card = createCard(card, cardFrameIm, cardArtPath)

	if err != nil {
		return nil, err
	}
	return card, nil

}

func parseCard(id int) (*Card, error) {
	filename := fmt.Sprint("Cards/", id, ".yaml")
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Panic(err)
	}
	card := new(Card)

	err = yaml.Unmarshal(file, card) // card is not being filled
	fmt.Println(card.Name, card.ID, card.CType, card.SubType, card.Text, card.Effects, card.Stats)
	return card, err
}

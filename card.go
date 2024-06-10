package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
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

//type Effect string

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
	Effects Effect   `yaml:"Effects"` // will be remade in the YAML file
	Stats   Stats    `yaml:"Stats"`   //can be nil
	// Card effect related
	Keywords []Keyword
	// more gerenal stuff
	X, Y           int
	W, H           int
	ScaleX, ScaleY float64
	Selected       bool
	// ebiten stuff
	//Image *ebiten.Image //should be removed later, check Draw
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
	return x > c.X && x < c.X+int(float64(c.W)*c.ScaleX) && y > c.Y && y < c.Y+int(float64(c.H)*c.ScaleY)
}

func (c *Card) resetScale() {
	c.ScaleX = scalingFactor
	c.ScaleY = scalingFactor
}

// function to draw the card, first translate to the position of the card, then draw the image
func (c *Card) draw(screen *ebiten.Image, cardImage *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.ScaleX, c.ScaleY) //always scale first
	op.GeoM.Translate(float64(c.X), float64(c.Y))
	screen.DrawImage(cardImage, op)
}

// Cria uma carta com texto e o frame da carta
func createCardImage(cardFrameIm *ebiten.Image, cardName, pathCardArt string) (*ebiten.Image, error) {
	CardArt, err := createImageFromPath(pathCardArt)
	if err != nil {
		return nil, err
	}
	// Imagem que será combinada com o frame da carta
	image := ebiten.NewImage(cardFrameIm.Bounds().Dx(), cardFrameIm.Bounds().Dy())

	// Draw no card frame primeiro
	op := &ebiten.DrawImageOptions{}
	image.DrawImage(cardFrameIm, op)
	op.GeoM.Translate(margemHorizontal, margemVertical)
	image.DrawImage(CardArt, op)

	// text
	// todo: justify text
	if len(cardName) > 14 {
		fontSize = 15.0 * 10 // fits 22 characters
	}

	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(border, border)
	textOp.ColorScale.ScaleWithColor(color.White)

	text.Draw(image, cardName, &text.GoTextFace{
		Source: font,
		Size:   fontSize,
	}, textOp)

	// reset font size
	fontSize = 20.0 * 10

	return image, nil
}

func newCardFromID(id int) (*Card, error) {
	card, err := parseCard(id)
	cardFrameIm := new(ebiten.Image)
	if card.CType == 0 { // will be changed to handle other types of cards
		cardFrameIm, err = createImageFromPath("Image/CardFrame/CardFrame.png")
		if err != nil {
			return nil, err
		}
	}

	card.ScaleX = scalingFactor
	card.ScaleY = scalingFactor

	card.W = int(float64(cardFrameIm.Bounds().Dx()) * scalingFactor)
	card.H = int(float64(cardFrameIm.Bounds().Dy()) * scalingFactor)
	if err != nil {
		return nil, err
	}
	cardFrameIm.Deallocate()
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
	// fmt.Println(card.Name, card.ID, card.CType, card.SubType, card.Text, card.Effects, card.Stats)
	return card, err
}

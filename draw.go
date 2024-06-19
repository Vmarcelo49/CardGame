package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// A imbutida do ebitenutil.NewImageFromFile não é recomendada para uso em produção, então não irei usar.
func createImageFromPath(path string) (*ebiten.Image, error) {
	imgBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	Rawimg, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}
	img := ebiten.NewImageFromImage(Rawimg)
	return img, nil
}

// Carrega as imagens da mão, campo e cemitério
// Checa se a imagem já foi carregada, se não carregada, carrega a imagem
func (g *Game) loadImages() {
	for _, card := range g.duel.p1Hand.cards {
		if _, ok := g.texMap[card.ID]; !ok {
			fmt.Println("Carregando imagem de ID:", card.ID)
			img, err := createCardImage(g.texMap[-1], card.Name, card.Text, fmt.Sprint("Image/CardArt/", card.ID, ".png"))
			if err != nil {
				log.Panic(err)
			}
			g.texMap[card.ID] = img
		}
	}
	// TODO: arrumar isso pros dois players
	for _, card := range g.duel.field.player1Field {
		if _, ok := g.texMap[card.ID]; !ok {
			img, err := createCardImage(g.texMap[0], card.Name, card.Text, fmt.Sprint("Image/CardArt/", card.ID, ".png"))
			if err != nil {
				log.Panic(err)
			}
			g.texMap[card.ID] = img
		}
	}
}

func (g *Game) freeImages() {
	for key := range g.texMap {
		// frame e facedown não são imagens de cartas, então não precisam ser deletadas
		if key == 0 || key == -1 {
			continue
		}
		delete(g.texMap, key)
	}

}

func (c *Card) draw(screen *ebiten.Image, cardImage *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.ScaleX, c.ScaleY) //always scale first
	op.GeoM.Translate(float64(c.X), float64(c.Y))
	screen.DrawImage(cardImage, op)
}

func (d *Deck) draw(screen *ebiten.Image, cardBack *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scalingFactor, scalingFactor)
	op.GeoM.Translate(float64(d.X), float64(d.Y))
	screen.DrawImage(cardBack, op)
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

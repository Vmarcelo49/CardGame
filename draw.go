package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
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
			img, err := createCardImage(g.texMap[0], card.Name, fmt.Sprint("Image/CardArt/", card.ID, ".png"))
			if err != nil {
				log.Panic(err)
			}
			g.texMap[card.ID] = img
		}
	}
	// TODO: arrumar isso pros dois players
	for _, card := range g.duel.field.player1Field {
		if _, ok := g.texMap[card.ID]; !ok {
			img, err := createCardImage(g.texMap[0], card.Name, fmt.Sprint("Image/CardArt/", card.ID, ".png"))
			if err != nil {
				log.Panic(err)
			}
			g.texMap[card.ID] = img
		}
	}
}

func (g *Game) freeImages() {
	for key := range g.texMap {
		// 0 é a imagem padrão, não precisa ser deletada
		if key == 0 {
			continue
		}
		delete(g.texMap, key)
	}
}

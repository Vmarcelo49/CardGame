package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	_ "net/http/pprof"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	backgroundColor = color.RGBA{R: 31, G: 31, B: 31, A: 255}
	font            *text.GoTextFaceSource
)

const (
	screenWidth  = 1280
	screenHeight = 720
	//
	MainMenu Scene = iota
	RockPaperScissors
	DuelScene
)

type Scene uint8

type Game struct {
	//Main Menu
	scene           Scene
	mainMenuButtons []*Button
	mouse           *Mouse
	//Duel
	duel *Duel

	texMap map[int]*ebiten.Image
}

func init() {
	fontBytes, err := os.ReadFile("Font/Ubuntu-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatal(err)
	}
	font = s
	scalingFactor = getScalingFactor(screenWidth, screenHeight)

}

func getScalingFactor(currentWidth, currentHeight int) float64 {
	baseWidth, baseHeight := 1280, 720
	baseScalingFactor := 0.1

	// Calculate the scaling factor based on the current resolution
	widthRatio := float64(currentWidth) / float64(baseWidth)
	heightRatio := float64(currentHeight) / float64(baseHeight)

	// Use the smaller ratio to ensure the scaled image fits both dimensions
	if widthRatio < heightRatio {
		return baseScalingFactor * widthRatio
	} else {
		return baseScalingFactor * heightRatio
	}
}

func newGame() *Game {
	back, err := createImageFromPath("Image/CardFrame/CardBackside.png")
	frame, err := createImageFromPath("Image/CardFrame/CardFrame.png")
	if err != nil {
		log.Panic(err)
	}
	texMap := map[int]*ebiten.Image{
		-1: frame,
		0:  back,
	}
	return &Game{
		mouse:  &Mouse{},
		scene:  MainMenu,
		texMap: texMap,
	}

}

func (g *Game) loadDuelMode() {
	g.duel = newDuel()
}

func (g *Game) Update() error {
	g.mouse.UpdateMouseState()
	var exit error
	switch g.scene {
	case DuelScene:
		exit = logic(g)
		g.loadImages()
	case MainMenu:
		if g.mainMenuButtons == nil {
			g.mainMenuButtons, _ = g.createButtons()
		}
		for _, b := range g.mainMenuButtons {
			if !b.alreadClicked {
				exit = b.checkClicked(g.mouse)
			}
		}

	case RockPaperScissors:
		//rock paper scissors logic
	}

	// Can return the error to end the game.

	return exit
}

func (g *Game) DrawDuel(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	g.duel.p1Deck.draw(screen, g.texMap[0])
	g.duel.p2Deck.draw(screen, g.texMap[0])

	for _, card := range g.duel.p1Hand.cards {
		if g.texMap[card.ID] == nil {
			log.Panic(fmt.Sprintln("Card ID: ", card.ID, " is nil"))
			card.draw(screen, g.texMap[0])
		}
		card.draw(screen, g.texMap[card.ID])
	}
	for i := range g.duel.p2Hand.cards {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(g.duel.p2Hand.cards[i].X), float64(g.duel.p2Hand.cards[i].Y))

		screen.DrawImage(g.texMap[0], op)

	}
	g.duel.field.draw(screen, g.texMap)
	// g.duel.p1GY.draw(screen, lastCardSentToP1GY.ID, g.texMap)

}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.scene {
	case DuelScene:
		g.DrawDuel(screen)
	case MainMenu:
		g.DrawMainMenu(screen)
	case RockPaperScissors:
		fmt.Println("RockPaperScissors")
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("MyGame")
	game := newGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

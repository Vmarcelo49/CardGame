package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	backgroundColor = color.RGBA{R: 31, G: 31, B: 31, A: 255}
	cardBack        *ebiten.Image
	font            *text.GoTextFaceSource
)

const (
	screenWidth  = 1280
	screenHeight = 720
	//
	MainMenu Scene = iota
	RockPaperScissors
	Duel
)

type Scene uint8

type Game struct {
	mouse       *Mouse
	hand        *Hand
	deck        *Deck
	field       *Field
	turnCount   int
	scene       Scene
	HUDMainMenu *HUDMainMenu
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
	return &Game{
		mouse: &Mouse{},
		hand:  &Hand{},
		deck:  &Deck{},
		field: &Field{},
	}

}

func (g *Game) loadDuelMode() {
	//currently only for p1
	var err error
	cardBack, _, err = ebitenutil.NewImageFromFile("Image/CardFrame/CardBackside.png")
	if err != nil {
		log.Panic(err)
	}
	// Todo: make the card images separate from the card object

	g.hand = newHand()

	cardIDS, err := getCardIDs("Deck/testDeck.txt")
	if err != nil {
		log.Panic(err)
	}
	g.deck = newDeck()

	for _, id := range cardIDS {
		card, err := newCardFromID(id)
		if err != nil {
			log.Panic(err)
		}
		g.deck.cards = append(g.deck.cards, card)
	}

	g.field = newField()
}

func (g *Game) Update() error {
	g.mouse.UpdateMouseState()
	var exit error
	switch g.scene {
	case Duel:
		exit = logic(g)
	case MainMenu:
		//menu logic

		// then
		g.scene = Duel
		g.loadDuelMode()
	case RockPaperScissors:
		//rock paper scissors logic
	}

	// Can return the error to end the game.

	return exit
}

func (g *Game) DrawDuel(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	debugText := fmt.Sprintf("Cards in Hand: %d, Cards on Field %d", len(g.hand.cards), len(g.field.player1Field))
	ebitenutil.DebugPrint(screen, debugText)
	g.deck.draw(screen)
	g.field.draw(screen)
	for _, card := range g.hand.cards {
		card.draw(screen)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.scene {
	case Duel:
		g.DrawDuel(screen)
	case MainMenu:
		//menu logic
	case RockPaperScissors:
		//rock paper scissors logic
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

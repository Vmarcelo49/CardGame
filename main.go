package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	_ "github.com/silbinarywolf/preferdiscretegpu" // This is needed for Windows to prefer the discrete GPU
)

var (
	backgroundColor = color.RGBA{R: 31, G: 31, B: 31, A: 255}
	font            *text.GoTextFaceSource
	exitFlag        error
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
	keyStates map[ebiten.Key]bool
	//Main Menu
	currentScene    Scene
	mainMenuButtons []*Button
	mouse           *Mouse
	//Duel
	duelRenderer      *DuelRenderer
	gamestate         *Gamestate
	previousGamestate *Gamestate
	exitingDuel       bool
	//label        *Label
	otherImgs []*Label
}

func (g *Game) Update() error {
	g.mouse.UpdateMouseState()

	switch g.currentScene {
	case DuelScene:
		g.updateGameLogic()
	case MainMenu:
		if g.exitingDuel {
			g.freeImages()
			g.duelRenderer = nil
			runtime.GC()
			g.exitingDuel = false
		}
		if g.mainMenuButtons == nil {
			g.mainMenuButtons = g.newButtons()
		}
		for _, b := range g.mainMenuButtons {
			if !b.alreadyClicked { // evita chamar a função de criar o duelo mais de uma vez.
				exitFlag = b.checkClicked(g.mouse)
			}
		}

	case RockPaperScissors:
		// TODO
	}

	// Can return the error to end the game.

	return exitFlag
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.currentScene {
	case DuelScene:
		g.DrawDuel(screen)
	case MainMenu:
		g.DrawMainMenu(screen)
	case RockPaperScissors:
		fmt.Println("It will be done someday...")
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func newGame() *Game {
	game := &Game{}
	game.mouse = &Mouse{}
	game.currentScene = MainMenu
	game.keyStates = make(map[ebiten.Key]bool)

	return game
}

func init() {
	loadFont()
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("MyGame")

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}

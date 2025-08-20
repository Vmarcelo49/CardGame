package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/silbinarywolf/preferdiscretegpu"

	"github.com/Vmarcelo49/CardGame/pkg/game"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("CardGame")

	gameInstance, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	if err := ebiten.RunGame(gameInstance); err != nil {
		log.Fatal(err)
	}
}

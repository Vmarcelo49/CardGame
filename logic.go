package main

import (
	"errors"
	"fmt"
	"log"
)

var (
	SelectedCardIndex = -1
)

// Envia carta apartir do index para o slice de destino, removendo a carta do slice de origem.
func sendCardTo(destination []*Card, source []*Card, index int) ([]*Card, []*Card, error) {
	if len(source) == 0 {
		return destination, source, errors.New("error: Source is empty, unable to send card")
	}
	if index < 0 || index >= len(source) {
		return destination, source, errors.New("error: Index out of bounds, unable to send card")
	}

	destination = append(destination, source[index])
	source = append(source[:index], source[index+1:]...)
	return destination, source, nil
}

func selectCard(g *Game) {
	for i, card := range g.hand.cards {
		if card.in(g.mouse.X, g.mouse.Y) {
			card.ScaleX = 0.11 //trocar isso por um highlight no momento causa problemas em resoluções diferentes TODO
			card.ScaleY = 0.11
			if g.mouse.LeftPressed {
				card.Selected = true // Maybe not needed
				SelectedCardIndex = i
				fmt.Println("Card Selected")
				break
			}
		} else {
			card.resetScale()
		}
	}
}

func deselectOrMoveCard(g *Game) {
	if g.mouse.LeftPressed {
		if !g.hand.cards[SelectedCardIndex].in(g.mouse.X, g.mouse.Y) && g.mouse.Y > g.hand.coordY {
			g.hand.cards[SelectedCardIndex].resetScale()
			SelectedCardIndex = -1
		}
		if g.mouse.Y < g.hand.coordY {
			var err error
			g.hand.cards[SelectedCardIndex].resetScale()
			g.field.player1Field, g.hand.cards, err = sendCardTo(g.field.player1Field, g.hand.cards, SelectedCardIndex)
			if err != nil {
				log.Println(err) // Probably should never return the error
			}

			//g.hand.cards[SelectedCardIndex].moveTo(screenWidth/2, int(g.field.middlescreen)) // will be replaced

			SelectedCardIndex = -1
		}
	}
}

func logic(g *Game) error {
	if SelectedCardIndex == -1 {
		selectCard(g)
	}
	if SelectedCardIndex > -1 {
		deselectOrMoveCard(g)
	}

	// Who goes first
	// will be implemented later

	// if 0, its beginning of the game, decide who goes first, then draw cards

	return keyboardInput(g.deck, g.hand)
}
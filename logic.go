package main

import (
	"fmt"
	"log"
)

// puts the specified button on top of the card
func (g *Game) updateButton(card *Card, button *Button) {
	button.x = card.X
	button.y = card.Y - button.h

	// Atualiza a função do botão para usar o efeito da carta selecionada
	button.function = func() error {
		button.alreadyClicked = true
		button.x = -5000
		err := card.normalSummon(g.gamestate)
		if err != nil {
			return err
		}
		return nil
	}
}

func (g *Game) updateSelectCard() {
	var selectedCard *Card

	for _, card := range g.gamestate.P1.Hand {
		if card.Selected {
			selectedCard = card
			break
		}
	}

	// Verifique se uma nova carta foi clicada
	for _, card := range g.gamestate.P1.Hand {
		if g.checkCardClicked(card) && g.mouse.LeftPressed {
			if selectedCard == nil {
				// Se nenhuma carta estiver selecionada, selecione esta
				fmt.Println("Card selected:", card.Name)
				card.Selected = true
				g.newCardClickedFunc(card, g.gamestate)

			} else if selectedCard != card {
				// Se uma carta diferente estiver selecionada, deselecione a atual e selecione a nova
				fmt.Println("Card deselected:", selectedCard.Name)
				selectedCard.Selected = false

				fmt.Println("Card selected:", card.Name)
				card.Selected = true

			}
			return
		}
	}

	// Se o clique for em qualquer outro lugar (e não em uma carta), deselecione a carta atual
	if g.mouse.LeftPressed && selectedCard != nil && !g.checkCardClicked(selectedCard) {
		fmt.Println("Card deselected:", selectedCard.Name)
		selectedCard.Selected = false

	}
}

func (g *Game) updateGameLogic() {
	// Verifica entradas de teclado e cliques em botões da Hud do duelo
	if err := g.checkInput(); err != nil {
		log.Println(err)
	}

	// Verifica se o mouse foi clicado em alguma carta
	g.updateSelectCard()

	// Atualiza o estado do jogo
	g.gamestate.update()

	// Atualiza o conteúdo visível se houver mudança no estado do jogo
	if g.previousGamestate != nil && !g.gamestate.equals(g.previousGamestate) {
		log.Println("Gamestate changed")
		g.updateDuelRenderer()
	}
	// Salva o estado atual do jogo
	g.previousGamestate = copyGamestate(g.gamestate)
}

func (g *Game) newCardClickedFunc(card *Card, gs *Gamestate) {
	if card.CType == 0 && card.getLocation(gs) == "P1HAND" {
		g.updateButton(card, g.duelButtons[0])
		fmt.Println("NormalSummon")
		fmt.Println(len(g.duelButtons))
	}
	if card.CType == 1 && card.getLocation(gs) == "P1HAND" {
		g.duelButtons[1].function = card.Effect // Reminder: Effect is a `func() error`
		g.duelButtons[1].x = card.X
		g.duelButtons[1].y = card.Y - g.duelButtons[1].h
		fmt.Println("eff")
	}
}

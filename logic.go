package main

import (
	"fmt"
	"log"
	"time"
)

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
				g.otherImgs[0].x = float64(card.X)
				g.otherImgs[0].y = float64(card.Y) - 31
				time.AfterFunc(2*time.Second, func() {
					card.Selected = false
					g.otherImgs[0].x = 5000
				})
			} else if selectedCard != card {
				// Se uma carta diferente estiver selecionada, deselecione a atual e selecione a nova
				fmt.Println("Card deselected:", selectedCard.Name)
				selectedCard.Selected = false

				fmt.Println("Card selected:", card.Name)
				card.Selected = true
				g.otherImgs[0].x = float64(card.X)
				g.otherImgs[0].y = float64(card.Y) - 31
				time.AfterFunc(2*time.Second, func() {
					card.Selected = false
					g.otherImgs[0].x = 5000
				})
			}
			// Não faça nada se a mesma carta for clicada novamente
			return
		}
	}

	// Se o clique for em qualquer outro lugar (e não em uma carta), deselecione a carta atual
	if g.mouse.LeftPressed && selectedCard != nil && !g.checkCardClicked(selectedCard) {
		fmt.Println("Card deselected:", selectedCard.Name)
		selectedCard.Selected = false
		g.otherImgs[0].x = 5000
	}
}

func (g *Game) updateGameLogic() {
	// Verifica entradas de teclado (ex: tecla ESC para sair)
	g.checkInput()

	// Verifica se o mouse foi clicado em alguma carta
	g.updateSelectCard()

	// Atualiza o estado do jogo
	g.gamestate.update()

	// Atualiza o conteúdo visível se houver mudança no estado do jogo
	if g.previousGamestate != nil && !g.gamestate.equals(g.previousGamestate) {
		log.Println(" Gamestate changed")
		g.updateDuelRenderer()
	}
	// Salva o estado atual do jogo
	g.previousGamestate = copyGamestate(g.gamestate)
}

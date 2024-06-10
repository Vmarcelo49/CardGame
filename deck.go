package main

import (
	"errors"
	"log"
)

type Deck struct {
	cards []*Card
	X, Y  int
}

// Nenhuma carta deve ser renderizada apartir do deck, exceto a imagem de carta virada para baixo.
// O Deck é o local onde as cartas começam no jogo.

// Como a mão, o deck pode receber cartas, porém em vez de adicionar uma carta por vez, o deck recebe um slice de cartas inteiro.
// No futuro, o meio de adicionar cartas ao deck será alterado para outro meio, como um arquivo de texto ou um banco de dados.

func newDeck(deckFilePath string, player int, cardWidth, cardHeight int) *Deck {
	cardIDS, err := getCardIDs(deckFilePath)
	if err != nil {
		log.Panic(err)
	}
	deck := new(Deck)

	for _, id := range cardIDS {
		card, err := newCardFromID(id, cardWidth, cardHeight)
		if err != nil {
			log.Panic(err)
		}
		deck.cards = append(deck.cards, card)

	}
	if player == 1 {
		deck.X = 5 * (screenWidth / 6)
		deck.Y = 7 * (screenHeight / 8)
	}
	if player == 2 {
		deck.X = 1 * (screenWidth / 6)
		deck.Y = -1 * (screenHeight / 8)
	}
	return deck
}

// drawCard, coloca a carta 0 do slice em uma variavel temporaria, remove a carta do slice do deck e coloca a carta na mão recebida como argumento.
// Com essa função, dá para se editar diretamente o slice de cartas da mão.

func (d *Deck) drawCard(hand *Hand) error {
	if len(d.cards) == 0 {
		return errors.New("error: Deck is empty, unable to draw card")
	}
	tempCard := d.cards[0]
	d.cards = d.cards[1:]
	hand.addCard(tempCard)
	return nil
}

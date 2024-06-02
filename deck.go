package main

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"
)

type Deck struct {
	cards []*Card
	X, Y  int
}

// Nenhuma carta deve ser renderizada apartir do deck, exceto a imagem de carta virada para baixo.
// O Deck é o local onde as cartas começam no jogo.

// Como a mão, o deck pode receber cartas, porém em vez de adicionar uma carta por vez, o deck recebe um slice de cartas inteiro.
// No futuro, o meio de adicionar cartas ao deck será alterado para outro meio, como um arquivo de texto ou um banco de dados.

func newDeck() *Deck {
	deck := new(Deck)
	// O X do deck é 5/6 da tela, no lado direito da tela.
	deck.X = 5 * (screenWidth / 6)

	// O valor da Altura do deck é 7/8 da tela, no topo da tela, como na mão.
	deck.Y = 7 * (screenHeight / 8)

	return deck
}

// drawDeck desenha a imagem de carta virada para baixo no deck.

func (d *Deck) draw(screen *ebiten.Image) {
	// A carta virada para baixo é uma imagem de escopo global e não é algo diretamente relacionado com o deck.

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scalingFactor, scalingFactor) //segue o padrão de escala, sempre escalar primeiro
	op.GeoM.Translate(float64(d.X), float64(d.Y))
	screen.DrawImage(cardBack, op)
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

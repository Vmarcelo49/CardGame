package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type DuelRenderer struct {
	cardImgMap                       map[int]*ebiten.Image
	cardSizeW, cardSizeH             int
	scaling                          float64
	P1DeckLocationX, P1DeckLocationY float64
	P2DeckLocationX, P2DeckLocationY float64

	// Field locations
	P1FieldLocationY float64
	P2FieldLocationY float64
}

func getScalingFactor(currentWidth, currentHeight int) float64 {
	baseWidth, baseHeight := 1280, 720
	baseScalingFactor := 0.1

	// Calculate the scaling factor based on the current resolution
	widthRatio := float64(currentWidth) / float64(baseWidth)
	heightRatio := float64(currentHeight) / float64(baseHeight)

	// Use the smaller ratio to ensure the scaled image fits both dimensions
	var scaling float64
	if widthRatio < heightRatio {
		scaling = baseScalingFactor * widthRatio
	} else {
		scaling = baseScalingFactor * heightRatio
	}

	return scaling
}

func (g *Game) loadDuelRenderer() error {
	g.duelRenderer = &DuelRenderer{}
	cardBackImg, err := newImageFromPath("assets/image/cardFrame/cardBackside.png")
	if err != nil {
		return err
	}
	cardFrameImg, err := newImageFromPath("assets/image/cardFrame/cardFrame.png")
	if err != nil {
		return err
	}
	g.duelRenderer.cardImgMap = map[int]*ebiten.Image{
		-1: cardFrameImg,
		0:  cardBackImg,
	}

	g.duelRenderer.scaling = getScalingFactor(screenWidth, screenHeight)
	g.duelRenderer.cardSizeW = int(float64(cardFrameImg.Bounds().Dx()) * g.duelRenderer.scaling) // around 150 pixels
	g.duelRenderer.cardSizeH = int(float64(cardFrameImg.Bounds().Dy()) * g.duelRenderer.scaling)

	g.duelRenderer.P1DeckLocationX = float64(screenWidth - g.duelRenderer.cardSizeW)
	g.duelRenderer.P1DeckLocationY = float64(screenHeight - g.duelRenderer.cardSizeH)

	g.duelRenderer.P2DeckLocationX = 0
	g.duelRenderer.P2DeckLocationY = 0

	// Field locations
	g.duelRenderer.P1FieldLocationY = screenHeight/2 + 5
	g.duelRenderer.P2FieldLocationY = screenHeight/2 - 5

	return nil
}

func (g *Game) updateVisibleCards() {
	for _, card := range g.gamestate.P1.Hand {
		if _, ok := g.duelRenderer.cardImgMap[card.ID]; !ok { // se a carta não existe no map...
			img, err := createCardImage(g.duelRenderer.cardImgMap[-1], card)
			if err != nil {
				log.Panic(err)
			}
			g.duelRenderer.cardImgMap[card.ID] = img
		}
	}

	for _, card := range g.gamestate.Field.AllCards() {
		if _, ok := g.duelRenderer.cardImgMap[card.ID]; !ok {
			img, err := createCardImage(g.duelRenderer.cardImgMap[0], card)
			if err != nil {
				log.Panic(err)
			}
			g.duelRenderer.cardImgMap[card.ID] = img
		}
	}
}

func (g *Game) updateCardLocations() {
	// Atualiza a posição das cartas na mão do jogador 1
	if len(g.gamestate.P1.Hand) == 1 {
		// Coloca a primeira carta na posição inicial
		g.gamestate.P1.Hand[0].X = screenWidth/2 - g.duelRenderer.cardSizeW/2
		g.gamestate.P1.Hand[0].Y = int(g.duelRenderer.P1DeckLocationY)
	}
	if len(g.gamestate.P1.Hand) > 1 {
		startX := screenWidth/2 - (g.duelRenderer.cardSizeW * (len(g.gamestate.P1.Hand) - 1) / 2)
		for i, card := range g.gamestate.P1.Hand {
			card.X = startX + i*g.duelRenderer.cardSizeW
			card.Y = int(g.duelRenderer.P1DeckLocationY)
		}
	}

	// Atualiza a posição das cartas no campo do jogador 1
	if len(g.gamestate.Field.P1) > 0 {
		// Coloca a primeira carta na posição inicial
		g.gamestate.Field.P1[0].X = screenWidth / 2
		g.gamestate.Field.P1[0].Y = int(g.duelRenderer.P1FieldLocationY)

		// Posiciona as demais cartas ao lado da última carta adicionada
		for i := 1; i < len(g.gamestate.Field.P1); i++ {
			prevCard := g.gamestate.Field.P1[i-1]
			card := g.gamestate.Field.P1[i]
			card.X = prevCard.X + g.duelRenderer.cardSizeW
			card.Y = int(g.duelRenderer.P1FieldLocationY)
		}
	}

	// Atualiza a posição das cartas na mão do jogador 2
	if len(g.gamestate.P2.Hand) > 0 {
		// Coloca a primeira carta na posição inicial
		g.gamestate.P2.Hand[0].X = screenWidth / 2
		g.gamestate.P2.Hand[0].Y = int(g.duelRenderer.P2DeckLocationY)

		// Posiciona as demais cartas ao lado da última carta adicionada
		if len(g.gamestate.P2.Hand) > 1 {
			startX := screenWidth/2 - (g.duelRenderer.cardSizeW * (len(g.gamestate.P2.Hand) - 1) / 2)
			for i, card := range g.gamestate.P2.Hand {
				card.X = startX + i*g.duelRenderer.cardSizeW
				card.Y = int(g.duelRenderer.P2DeckLocationY)
			}
		}

	}

	// Atualiza a posição das cartas no campo do jogador 2
	if len(g.gamestate.Field.P2) > 0 {
		// Coloca a primeira carta na posição inicial
		g.gamestate.Field.P2[0].X = screenWidth / 2
		g.gamestate.Field.P2[0].Y = int(g.duelRenderer.P2FieldLocationY)

		// Posiciona as demais cartas ao lado da última carta adicionada
		for i := 1; i < len(g.gamestate.Field.P2); i++ {
			prevCard := g.gamestate.Field.P2[i-1]
			card := g.gamestate.Field.P2[i]
			card.X = prevCard.X + g.duelRenderer.cardSizeW
			card.Y = int(g.duelRenderer.P2FieldLocationY)
		}
	}
}

func (g *Game) freeImages() {
	for key := range g.duelRenderer.cardImgMap {
		// frame e facedown não são imagens de cartas, então não precisam ser deletadas
		if key == 0 || key == -1 {
			continue
		}
		delete(g.duelRenderer.cardImgMap, key)
	}

}

func (g *Game) DrawDuel(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	ebitenutil.DebugPrint(screen, "Duel mode")

	// decks
	if len(g.gamestate.P1.Deck) > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.duelRenderer.scaling, g.duelRenderer.scaling)
		op.GeoM.Translate(g.duelRenderer.P1DeckLocationX, g.duelRenderer.P1DeckLocationY)
		screen.DrawImage(g.duelRenderer.cardImgMap[0], op)

	}
	if len(g.gamestate.P2.Deck) > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.duelRenderer.scaling, g.duelRenderer.scaling)
		op.GeoM.Translate(g.duelRenderer.P2DeckLocationX, g.duelRenderer.P2DeckLocationY)
		screen.DrawImage(g.duelRenderer.cardImgMap[0], op)
	}

	// field
	for _, card := range g.gamestate.Field.P1 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.duelRenderer.scaling, g.duelRenderer.scaling)
		op.GeoM.Translate(float64(card.X), g.duelRenderer.P1FieldLocationY)
		screen.DrawImage(g.duelRenderer.cardImgMap[card.ID], op)

	}
	for _, card := range g.gamestate.Field.P2 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.duelRenderer.scaling, g.duelRenderer.scaling)
		op.GeoM.Translate(float64(card.X), g.duelRenderer.P2FieldLocationY)

		// rotate card vertically
		op.GeoM.Scale(1, -1)

		screen.DrawImage(g.duelRenderer.cardImgMap[card.ID], op)
	}
	vector.DrawFilledRect(screen, 0, screenHeight/2, screenWidth, 2, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false) // divisória dos campos

	// hands
	for _, card := range g.gamestate.P1.Hand {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.duelRenderer.scaling, g.duelRenderer.scaling)               //always scale first
		op.GeoM.Translate(float64(card.X), float64(g.duelRenderer.P1DeckLocationY)) // Mesma altura do deck
		if g.duelRenderer.cardImgMap[card.ID] == nil {
			screen.DrawImage(g.duelRenderer.cardImgMap[0], op)
		} else {
			screen.DrawImage(g.duelRenderer.cardImgMap[card.ID], op)
		}
	}
	for _, card := range g.gamestate.P2.Hand {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.duelRenderer.scaling, g.duelRenderer.scaling)
		op.GeoM.Translate(float64(card.X), float64(g.duelRenderer.P2DeckLocationY))
		screen.DrawImage(g.duelRenderer.cardImgMap[0], op) // facedown

	}

	// UI

	for _, im := range g.otherImgs {
		// if x and y are inside the screen, draw the image
		if im.x > 0 && im.y > 0 {
			im.draw(screen)
		}
	}
	for _, b := range g.duelButtons {
		// if x and y are inside the screen, draw the button
		if b.x > 0 && b.y > 0 {
			b.draw(screen)
		}
	}

}

func (g *Game) DrawMainMenu(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	ebitenutil.DebugPrint(screen, "Main Menu")
	for _, b := range g.mainMenuButtons {
		b.draw(screen)
	}
}

func (g *Game) updateDuelRenderer() {
	g.updateVisibleCards()
	g.updateCardLocations()
	//g.updateUi()

}

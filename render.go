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

// getScalingFactor calculates the scaling factor based on the current screen dimensions.
func getScalingFactor(currentWidth, currentHeight int) float64 {
	baseWidth, baseHeight := 1280, 720
	baseScalingFactor := 0.1

	widthRatio := float64(currentWidth) / float64(baseWidth)
	heightRatio := float64(currentHeight) / float64(baseHeight)

	if widthRatio < heightRatio {
		return baseScalingFactor * widthRatio
	}
	return baseScalingFactor * heightRatio
}

// loadDuelRenderer initializes the DuelRenderer with card images and locations.
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
	g.duelRenderer.cardSizeW = int(float64(cardFrameImg.Bounds().Dx()) * g.duelRenderer.scaling)
	g.duelRenderer.cardSizeH = int(float64(cardFrameImg.Bounds().Dy()) * g.duelRenderer.scaling)

	g.duelRenderer.P1DeckLocationX = float64(screenWidth - g.duelRenderer.cardSizeW)
	g.duelRenderer.P1DeckLocationY = float64(screenHeight - g.duelRenderer.cardSizeH)

	g.duelRenderer.P2DeckLocationX = 0
	g.duelRenderer.P2DeckLocationY = 0

	g.duelRenderer.P1FieldLocationY = screenHeight/2 + 5
	g.duelRenderer.P2FieldLocationY = screenHeight/2 - 5

	return nil
}

// updateVisibleCards updates the images of the cards that are currently visible.
func (g *Game) updateVisibleCards() {
	g.updateCardImages(g.gamestate.P1.Hand, -1)
	g.updateCardImages(g.gamestate.Field.AllCards(), 0)
}

// updateCardImages updates the images of the given cards if they are not already in the cardImgMap.
func (g *Game) updateCardImages(cards []*Card, defaultImgKey int) {
	for _, card := range cards {
		if _, ok := g.duelRenderer.cardImgMap[card.ID]; !ok {
			img, err := createCardImage(g.duelRenderer.cardImgMap[defaultImgKey], card)
			if err != nil {
				log.Panic(err)
			}
			g.duelRenderer.cardImgMap[card.ID] = img
		}
	}
}

// updateCardLocations updates the locations of the cards in the players' hands and fields.
func (g *Game) updateCardLocations() {
	g.updateHandLocations(g.gamestate.P1.Hand, g.duelRenderer.P1DeckLocationY)
	g.updateFieldLocations(g.gamestate.Field.P1, g.duelRenderer.P1FieldLocationY)
	g.updateHandLocations(g.gamestate.P2.Hand, g.duelRenderer.P2DeckLocationY)
	g.updateFieldLocations(g.gamestate.Field.P2, g.duelRenderer.P2FieldLocationY)
}

// updateHandLocations updates the locations of the cards in a player's hand.
func (g *Game) updateHandLocations(hand []*Card, deckLocationY float64) {
	if len(hand) == 1 {
		hand[0].X = screenWidth/2 - g.duelRenderer.cardSizeW/2
		hand[0].Y = int(deckLocationY)
	} else if len(hand) > 1 {
		startX := screenWidth/2 - (g.duelRenderer.cardSizeW * (len(hand) - 1) / 2)
		for i, card := range hand {
			card.X = startX + i*g.duelRenderer.cardSizeW
			card.Y = int(deckLocationY)
		}
	}
}

// updateFieldLocations updates the locations of the cards on a player's field.
func (g *Game) updateFieldLocations(field []*Card, fieldLocationY float64) {
	if len(field) > 0 {
		field[0].X = screenWidth / 2
		field[0].Y = int(fieldLocationY)
		for i := 1; i < len(field); i++ {
			prevCard := field[i-1]
			field[i].X = prevCard.X + g.duelRenderer.cardSizeW
			field[i].Y = int(fieldLocationY)
		}
	}
}

// freeImages removes all card images from the cardImgMap except for the default images.
func (g *Game) freeImages() {
	for key := range g.duelRenderer.cardImgMap {
		if key == 0 || key == -1 {
			continue
		}
		delete(g.duelRenderer.cardImgMap, key)
	}
}

// drawDeck draws the deck of cards at the specified location.
func (g *Game) drawDeck(screen *ebiten.Image, deck []*Card, deckLocationX, deckLocationY float64) {
	if len(deck) > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.duelRenderer.scaling, g.duelRenderer.scaling)
		op.GeoM.Translate(deckLocationX, deckLocationY)
		screen.DrawImage(g.duelRenderer.cardImgMap[0], op)
	}
}

// drawField draws the cards on the field at the specified location.
func (g *Game) drawField(screen *ebiten.Image, field []*Card, fieldLocationY float64, rotate bool) {
	for _, card := range field {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.duelRenderer.scaling, g.duelRenderer.scaling)
		op.GeoM.Translate(float64(card.X), fieldLocationY)
		if rotate {
			op.GeoM.Scale(1, -1)
		}
		screen.DrawImage(g.duelRenderer.cardImgMap[card.ID], op)
	}
}

// drawHand draws the cards in a player's hand at the specified location.
func (g *Game) drawHand(screen *ebiten.Image, hand []*Card, deckLocationY float64) {
	for _, card := range hand {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.duelRenderer.scaling, g.duelRenderer.scaling)
		op.GeoM.Translate(float64(card.X), deckLocationY)
		screen.DrawImage(g.duelRenderer.cardImgMap[card.ID], op)
	}
}

// Draws only the facedown img
func (g *Game) drawHandP2(screen *ebiten.Image, hand []*Card, deckLocationY float64) {
	for _, card := range hand {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(g.duelRenderer.scaling, g.duelRenderer.scaling)
		op.GeoM.Translate(float64(card.X), deckLocationY)
		screen.DrawImage(g.duelRenderer.cardImgMap[0], op)
	}
}

// DrawDuel draws the duel screen, including decks, fields, and hands.
func (g *Game) DrawDuel(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	ebitenutil.DebugPrint(screen, "Duel mode")

	g.drawDeck(screen, g.gamestate.P1.Deck, g.duelRenderer.P1DeckLocationX, g.duelRenderer.P1DeckLocationY)
	g.drawDeck(screen, g.gamestate.P2.Deck, g.duelRenderer.P2DeckLocationX, g.duelRenderer.P2DeckLocationY)

	g.drawField(screen, g.gamestate.Field.P1, g.duelRenderer.P1FieldLocationY, false)
	g.drawField(screen, g.gamestate.Field.P2, g.duelRenderer.P2FieldLocationY, true)

	vector.DrawFilledRect(screen, 0, screenHeight/2, screenWidth, 2, color.RGBA{R: 255, G: 255, B: 255, A: 255}, false)

	g.drawHand(screen, g.gamestate.P1.Hand, g.duelRenderer.P1DeckLocationY)
	g.drawHandP2(screen, g.gamestate.P2.Hand, g.duelRenderer.P2DeckLocationY)

	for _, im := range g.otherImgs {
		if im.x > 0 && im.y > 0 {
			im.draw(screen)
		}
	}
	for _, b := range g.duelButtons {
		if b.x > 0 && b.y > 0 {
			b.draw(screen)
		}
	}
}

// DrawMainMenu draws the main menu screen.
func (g *Game) DrawMainMenu(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	ebitenutil.DebugPrint(screen, "Main Menu")
	for _, b := range g.mainMenuButtons {
		b.draw(screen)
	}
}

// updateDuelRenderer updates the duel renderer by updating visible cards and their locations.
func (g *Game) updateDuelRenderer() {
	g.updateVisibleCards()
	g.updateCardLocations()
}

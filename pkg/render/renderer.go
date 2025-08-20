package render

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/Vmarcelo49/CardGame/pkg/assets"
	"github.com/Vmarcelo49/CardGame/pkg/card"
	"github.com/Vmarcelo49/CardGame/pkg/ui"
)

// DuelRenderer handles rendering of the duel scene
type DuelRenderer struct {
	CardImgMap                       map[int]*ebiten.Image
	CardSizeW, CardSizeH             int
	Scaling                          float64
	P1DeckLocationX, P1DeckLocationY float64
	P2DeckLocationX, P2DeckLocationY float64
	P1FieldLocationY                 float64
	P2FieldLocationY                 float64
}

// GetScalingFactor calculates the scaling factor based on screen dimensions
func GetScalingFactor(currentWidth, currentHeight int) float64 {
	baseWidth, baseHeight := 1280, 720
	baseScalingFactor := 0.1

	widthRatio := float64(currentWidth) / float64(baseWidth)
	heightRatio := float64(currentHeight) / float64(baseHeight)

	if widthRatio < heightRatio {
		return baseScalingFactor * widthRatio
	}
	return baseScalingFactor * heightRatio
}

// LoadDuelRenderer initializes the DuelRenderer with card images and locations
func LoadDuelRenderer() (*DuelRenderer, error) {
	renderer := &DuelRenderer{}
	cardBackImg, err := assets.NewImageFromPath("assets/image/cardFrame/cardBackside.png")
	if err != nil {
		return nil, err
	}
	cardFrameImg, err := assets.NewImageFromPath("assets/image/cardFrame/cardFrame.png")
	if err != nil {
		return nil, err
	}
	renderer.CardImgMap = map[int]*ebiten.Image{
		-1: cardFrameImg,
		0:  cardBackImg,
	}

	renderer.Scaling = GetScalingFactor(1280, 720)
	renderer.CardSizeW = int(225 * renderer.Scaling)
	renderer.CardSizeH = int(315 * renderer.Scaling)

	renderer.P1DeckLocationX = 50
	renderer.P1DeckLocationY = 520
	renderer.P2DeckLocationX = 50
	renderer.P2DeckLocationY = 50

	renderer.P1FieldLocationY = 375
	renderer.P2FieldLocationY = 225

	return renderer, nil
}

// UpdateVisibleCards updates the images of cards that are currently visible
func (dr *DuelRenderer) UpdateVisibleCards(p1Hand, p2Hand, fieldCards []*card.Card) {
	dr.updateCardImages(p1Hand, -1)
	dr.updateCardImages(fieldCards, 0)
}

// updateCardImages updates the images of the given cards if they are not already in the cardImgMap
func (dr *DuelRenderer) updateCardImages(cards []*card.Card, defaultImgKey int) {
	for _, c := range cards {
		if _, ok := dr.CardImgMap[c.ID]; !ok {
			img, err := card.CreateCardImage(dr.CardImgMap[defaultImgKey], c)
			if err != nil {
				log.Panic(err)
			}
			dr.CardImgMap[c.ID] = img
		}
	}
}

// UpdateCardLocations updates the locations of cards in players' hands and fields
func (dr *DuelRenderer) UpdateCardLocations(p1Hand, p1Field, p2Hand, p2Field []*card.Card) {
	dr.updateHandLocations(p1Hand, dr.P1DeckLocationY)
	dr.updateFieldLocations(p1Field, dr.P1FieldLocationY)
	dr.updateHandLocations(p2Hand, dr.P2DeckLocationY)
	dr.updateFieldLocations(p2Field, dr.P2FieldLocationY)
}

// updateHandLocations updates the locations of cards in a player's hand
func (dr *DuelRenderer) updateHandLocations(hand []*card.Card, deckLocationY float64) {
	cardSpacing := 120.0
	startX := 150.0

	for i, c := range hand {
		c.X = int(startX + float64(i)*cardSpacing)
		c.Y = int(deckLocationY)
	}
}

// updateFieldLocations updates the locations of cards on a player's field
func (dr *DuelRenderer) updateFieldLocations(field []*card.Card, fieldLocationY float64) {
	cardSpacing := 120.0
	startX := 150.0

	for i, c := range field {
		c.X = int(startX + float64(i)*cardSpacing)
		c.Y = int(fieldLocationY)
	}
}

// FreeImages removes all card images from the cardImgMap except for the default images
func (dr *DuelRenderer) FreeImages() {
	for key := range dr.CardImgMap {
		if key != -1 && key != 0 {
			delete(dr.CardImgMap, key)
		}
	}
}

// DrawDuel draws the duel scene
func (dr *DuelRenderer) DrawDuel(screen *ebiten.Image, p1Hand, p1Deck, p1Field, p2Hand, p2Deck, p2Field []*card.Card,
	p1HP, p2HP int, buttons []*ui.Button, labels []*ui.Label) {
	screen.Fill(color.RGBA{R: 31, G: 31, B: 31, A: 255})

	// Draw field divider
	vector.DrawFilledRect(screen, 0, 360, 1280, 5, color.White, false)

	// Draw decks
	dr.drawDeck(screen, p1Deck, dr.P1DeckLocationX, dr.P1DeckLocationY)
	dr.drawDeck(screen, p2Deck, dr.P2DeckLocationX, dr.P2DeckLocationY)

	// Draw fields
	dr.drawField(screen, p1Field, dr.P1FieldLocationY, false)
	dr.drawField(screen, p2Field, dr.P2FieldLocationY, true)

	// Draw hands
	dr.drawHand(screen, p1Hand, dr.P1DeckLocationY)
	dr.drawHandP2(screen, p2Hand, dr.P2DeckLocationY)

	// Draw UI elements
	for _, button := range buttons {
		button.Draw(screen)
	}
	for _, label := range labels {
		label.Draw(screen)
	}
}

// DrawMainMenu draws the main menu
func DrawMainMenu(screen *ebiten.Image, buttons []*ui.Button) {
	screen.Fill(color.RGBA{R: 31, G: 31, B: 31, A: 255})
	ebitenutil.DebugPrint(screen, "Main Menu")
	for _, b := range buttons {
		b.Draw(screen)
	}
}

// drawDeck draws the deck of cards at the specified location
func (dr *DuelRenderer) drawDeck(screen *ebiten.Image, deck []*card.Card, deckLocationX, deckLocationY float64) {
	if len(deck) > 0 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(dr.Scaling, dr.Scaling)
		op.GeoM.Translate(deckLocationX, deckLocationY)
		screen.DrawImage(dr.CardImgMap[0], op)
	}
}

// drawField draws the cards on the field at the specified location
func (dr *DuelRenderer) drawField(screen *ebiten.Image, field []*card.Card, fieldLocationY float64, rotate bool) {
	for _, c := range field {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(dr.Scaling, dr.Scaling)
		if rotate {
			op.GeoM.Rotate(3.14159) // 180 degrees
			op.GeoM.Translate(float64(c.X)+float64(dr.CardSizeW), float64(c.Y)+float64(dr.CardSizeH))
		} else {
			op.GeoM.Translate(float64(c.X), float64(c.Y))
		}
		screen.DrawImage(dr.CardImgMap[c.ID], op)
	}
}

// drawHand draws the cards in a player's hand at the specified location
func (dr *DuelRenderer) drawHand(screen *ebiten.Image, hand []*card.Card, deckLocationY float64) {
	for _, c := range hand {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(dr.Scaling, dr.Scaling)
		op.GeoM.Translate(float64(c.X), float64(c.Y))
		if c.Selected {
			op.ColorScale.ScaleWithColor(color.RGBA{R: 255, G: 255, B: 0, A: 255})
		}
		screen.DrawImage(dr.CardImgMap[c.ID], op)
	}
}

// drawHandP2 draws only the face-down cards for player 2
func (dr *DuelRenderer) drawHandP2(screen *ebiten.Image, hand []*card.Card, deckLocationY float64) {
	for _, c := range hand {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(dr.Scaling, dr.Scaling)
		op.GeoM.Rotate(3.14159)
		op.GeoM.Translate(float64(c.X)+float64(dr.CardSizeW), float64(c.Y)+float64(dr.CardSizeH))
		screen.DrawImage(dr.CardImgMap[0], op)
	}
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func loadFont() {
	fontBytes, err := os.ReadFile("assets/font/Ubuntu-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	textFaceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatal(err)
	}
	font = textFaceSource
}

// Creates a new image from path
func newImageFromPath(path string) (*ebiten.Image, error) {
	imgBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	decodedImg, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image from %s: %w", path, err)
	}
	img := ebiten.NewImageFromImage(decodedImg)
	return img, nil
}

// deck

func newDeck(deckFilePath string) ([]*Card, error) {
	cardIDs, err := getCardIDs(deckFilePath)
	if err != nil {
		return nil, err
	}
	deck := []*Card{}
	for _, cardID := range cardIDs {
		deck = append(deck, newCardFromID(cardID))
	}
	return deck, nil
}

// getCardIDs reads a .txt file with card IDs and returns a slice of integers.
func getCardIDs(filename string) ([]int, error) {
	var cardIDs []int
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line) // Trim whitespace
		if line == "" {
			continue // Skip empty lines
		}
		cardID, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Failed to convert '%s' to integer: %v\n", line, err)
			continue // Skip this line and continue with the next
		}
		cardIDs = append(cardIDs, cardID)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return cardIDs, nil
}

func (g *Game) loadDuelMode() {
	g.currentScene = DuelScene
	g.loadDuelRenderer()
	var err error

	deck := "./deck/testDeck.txt"
	g.gamestate, err = newGameState(deck, deck)
	if err != nil {
		log.Panic("erro criando gamestate: ", err)
	}

	// Hud Items

	g.otherImgs = make([]*Label, 0)

	player1HPImg := ebiten.NewImage(100, 50)
	player1HPImg.DrawImage(newTextImageMultiline("P1: 100", color.White, 20, 100), &ebiten.DrawImageOptions{})
	g.otherImgs = append(g.otherImgs, &Label{5, screenHeight/2 + 5, player1HPImg, 0})

	player2HPImg := ebiten.NewImage(100, 50)
	player2HPImg.DrawImage(newTextImageMultiline("P2: 100", color.White, 20, 150), &ebiten.DrawImageOptions{})
	g.otherImgs = append(g.otherImgs, &Label{5, screenHeight/2 - 25, player2HPImg, 0})

	turnButtonImg := ebiten.NewImage(75, 75)
	turnButtonImg.Fill(color.RGBA{0, 255, 0, 255})
	g.otherImgs = append(g.otherImgs, &Label{screenWidth - 15 - 75, screenHeight/2 - (75 / 2), turnButtonImg, 0})

	turnCountImg := ebiten.NewImage(100, 50)
	turnCountImg.DrawImage(newTextImageMultiline("Turn: 1", color.White, 20, 100), &ebiten.DrawImageOptions{})
	g.otherImgs = append(g.otherImgs, &Label{screenWidth - 15 - 75, screenHeight/2 - (75 / 2) - 20, turnCountImg, 0})

	ButtonNormalSummon := newButton(g.duelRenderer.cardSizeW, g.duelRenderer.cardSizeH/10, 5000, 5000, "Normal Summon", func() error { // remember to change the function
		fmt.Println("Original button function called")
		return nil
	})
	g.duelButtons = []*Button{ButtonNormalSummon}

}

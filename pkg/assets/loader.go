package assets

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/Vmarcelo49/CardGame/pkg/card"
)

const (
	PlayerHPWidth  = 100
	PlayerHPHeight = 50
	TurnButtonSize = 75
)

// LoadFont loads the game font from assets
func LoadFont() (*text.GoTextFaceSource, error) {
	fontBytes, err := os.ReadFile("assets/font/Ubuntu-Regular.ttf")
	if err != nil {
		return nil, err
	}

	textFaceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		return nil, err
	}
	return textFaceSource, nil
}

// NewImageFromPath creates a new image from a file path
func NewImageFromPath(path string) (*ebiten.Image, error) {
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

// NewDeck creates a new deck from a file path
func NewDeck(deckFilePath string) ([]*card.Card, error) {
	cardIDs, err := getCardIDs(deckFilePath)
	if err != nil {
		return nil, err
	}
	deck := []*card.Card{}
	for _, cardID := range cardIDs {
		deck = append(deck, card.NewCardFromID(cardID))
	}
	return deck, nil
}

// getCardIDs reads a .txt file with card IDs and returns a slice of integers
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
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cardID, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Failed to convert '%s' to integer: %v\n", line, err)
			continue
		}
		cardIDs = append(cardIDs, cardID)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return cardIDs, nil
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
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

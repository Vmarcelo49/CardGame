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

const (
	playerHPWidth  = 100
	playerHPHeight = 50
	turnButtonSize = 75
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

func (g *Game) loadDuelMode() error {
	g.currentScene = DuelScene
	if err := g.loadDuelRenderer(); err != nil {
		return fmt.Errorf("failed to load duel renderer: %w", err)
	}

	deck := "./deck/testDeck.txt"
	gameState, err := newGameState(deck, deck)
	if err != nil {
		return fmt.Errorf("failed to create game state: %w", err)
	}
	g.gamestate = gameState

	if err := g.setupDuelUI(); err != nil {
		return fmt.Errorf("failed to setup duel UI: %w", err)
	}

	return nil
}


func (g *Game) setupDuelUI() error {
	g.otherImgs = make([]*Label, 0)

	// Player 1 HP
	player1HPImg, err := g.createHPImage("P1: 100", screenHeight/2+5)
	if err != nil {
		return err
	}
	g.otherImgs = append(g.otherImgs, player1HPImg)

	// Player 2 HP
	player2HPImg, err := g.createHPImage("P2: 100", screenHeight/2-25)
	if err != nil {
		return err
	}
	g.otherImgs = append(g.otherImgs, player2HPImg)

	// Turn Button
	turnButtonImg := ebiten.NewImage(turnButtonSize, turnButtonSize)
	turnButtonImg.Fill(color.RGBA{0, 255, 0, 255})
	g.otherImgs = append(g.otherImgs, &Label{screenWidth - 15 - turnButtonSize, screenHeight/2 - (turnButtonSize / 2), turnButtonImg, 0})

	// Turn Count
	turnCountImg, err := g.createTextImage("Turn: 1", screenWidth-15-turnButtonSize, screenHeight/2-(turnButtonSize/2)-20)
	if err != nil {
		return err
	}
	g.otherImgs = append(g.otherImgs, turnCountImg)

	// Normal Summon Button
	ButtonNormalSummon := newButton(g.duelRenderer.cardSizeW, g.duelRenderer.cardSizeH/10, 5000, 5000, "Normal Summon", func() error {
		fmt.Println("Normal Summon button clicked")
		return nil
	})
	g.duelButtons = []*Button{ButtonNormalSummon}

	return nil
}

func (g *Game) createHPImage(text string, yPos int) (*Label, error) {
	img := ebiten.NewImage(playerHPWidth, playerHPHeight)
	textImg := newTextImageMultiline(text, color.White, 20, playerHPWidth)
	img.DrawImage(textImg, &ebiten.DrawImageOptions{})
	return &Label{5, yPos, img, 0}, nil
}

func (g *Game) createTextImage(text string, xPos, yPos int) (*Label, error) {
	img := ebiten.NewImage(playerHPWidth, playerHPHeight)
	textImg := newTextImageMultiline(text, color.White, 20, playerHPWidth)
	img.DrawImage(textImg, &ebiten.DrawImageOptions{})
	return &Label{xPos, yPos, img, 0}, nil
}
package main

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Retorna uma lista de strings que cabem dentro de uma largura máxima
func breakTextIntoLines(txt string, fontsize float64, maxWidth int) []string {
	var lines []string

	face := &text.GoTextFace{
		Source: font,
		Size:   fontsize,
	}
	words := strings.Fields(txt)
	if len(words) == 0 {
		return lines
	}

	currentLine := words[0]
	var textW float64
	for _, word := range words[1:] {
		testLine := currentLine + " " + word
		textW, _ = text.Measure(testLine, face, 3)
		if int(textW) > maxWidth {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}
	}
	lines = append(lines, currentLine)
	return lines
}

func createTextImage(texto string, cor color.Color, fontsize float64) (*ebiten.Image, float64, float64) {
	// Função auxiliar para criar e medir uma imagem de texto
	face := &text.GoTextFace{
		Source: font,
		Size:   fontsize,
	}
	textSizeW, textSizeH := text.Measure(texto, face, 0)
	textImage := ebiten.NewImage(int(textSizeW), int(textSizeH))

	textOp := &text.DrawOptions{}
	textOp.ColorScale.ScaleWithColor(cor)
	text.Draw(textImage, texto, face, textOp)

	return textImage, textSizeW, textSizeH
}

func newTextImage(texto string, cor color.Color, fontsize float64) (*ebiten.Image, float64, float64) {
	// Determina o tamanho da fonte se não for especificado
	if fontsize <= 0 {
		fontsize = screenWidth / 60 // 25.6 em uma tela 1280x720
	}
	return createTextImage(texto, cor, fontsize)
}

// Cria uma imagem de texto quebrada em várias linhas, ainda não considera altura maxima
func newTextImageMultiline(texto string, cor color.Color, fontsize float64, maxWidth int) *ebiten.Image {
	face := &text.GoTextFace{
		Source: font,
		Size:   fontsize,
	}
	textSizeW, textSizeH := text.Measure(texto, face, 0)
	if int(textSizeW) <= maxWidth {
		textImage, _, _ := createTextImage(texto, cor, fontsize)
		return textImage
	} else {
		lines := breakTextIntoLines(texto, fontsize, maxWidth)
		hSizeOfLines := len(lines) * int(textSizeH)
		textOp := &text.DrawOptions{}
		textOp.ColorScale.ScaleWithColor(cor)
		textImage := ebiten.NewImage(int(textSizeW), hSizeOfLines)

		for i, line := range lines {
			lineY := float64(i) * textSizeH
			textOp.GeoM.Reset()
			textOp.GeoM.Translate(0, lineY)
			text.Draw(textImage, line, face, textOp)
		}
		return textImage
	}
}

type Button struct {
	x, y, w, h     int
	image          *ebiten.Image
	function       func() error
	alreadyClicked bool
}

func newButton(x, y int, texto string, function func() error) *Button {
	newImage := ebiten.NewImage(screenWidth/8, screenHeight/8)
	newImage.Fill(color.White)

	// draw text on the image
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(10, 10)
	textOp.ColorScale.ScaleWithColor(color.Black)

	text.Draw(newImage, texto, &text.GoTextFace{
		Source: font,
		Size:   20.0,
	}, textOp)

	return &Button{x, y, screenWidth / 8, screenHeight / 8, newImage, function, false}
}

func (b *Button) checkClicked(m *Mouse) error {
	if m.X > b.x && m.X < b.x+b.w && m.Y > b.y && m.Y < b.y+b.h && m.LeftPressed {
		b.alreadyClicked = true
		return b.function()
	} else {
		return nil
	}

}

func (b *Button) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	screen.DrawImage(b.image, op)
}

func addButton(buttonSlice []*Button, text string, function func() error) []*Button {
	buttonX := screenWidth / 8
	buttonY := screenHeight / 8
	x := (screenWidth - buttonX) / 2
	y := (screenHeight-buttonY)/2 + len(buttonSlice)*(buttonY+10)
	buttonSlice = append(buttonSlice, newButton(x, y, text, function))
	return buttonSlice
}

// Cria os botões do menu principal
func (g *Game) createButtons() ([]*Button, error) {
	var buttons []*Button
	buttons = addButton(buttons, "Duel", func() error {
		g.loadDuelMode()
		g.mainMenuButtons = nil // Go doesn clear the buttons, so we need to do it manually

		return nil
	})
	buttons = addButton(buttons, "Deck Editor", func() error {
		fmt.Println("Soon...")
		return nil
	})
	buttons = addButton(buttons, "Exit", func() error {
		return ebiten.Termination
	})

	return buttons, nil

}

func (g *Game) DrawMainMenu(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	ebitenutil.DebugPrint(screen, "Main Menu")
	for _, b := range g.mainMenuButtons {
		b.draw(screen)
	}
}

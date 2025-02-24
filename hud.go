package main

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Retorna uma lista de strings que cabem dentro de uma largura e altura máximas
func breakTextIntoLines(txt string, fontsize float64, maxWidth, maxHeight int) []string {
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
	var textW, textH float64
	for _, word := range words[1:] {
		testLine := currentLine + " " + word
		textW, textH = text.Measure(testLine, face, 3)
		if int(textW) > maxWidth || len(lines)*int(textH) >= maxHeight {
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

// Cria uma imagem de texto quebrada em várias linhas, considerando largura e altura máximas
func newTextImageMultiline(texto string, cor color.Color, fontsize float64, maxWidth, maxHeight int) *ebiten.Image {
	face := &text.GoTextFace{
		Source: font,
		Size:   fontsize,
	}
	textSizeW, textSizeH := text.Measure(texto, face, 0)
	if int(textSizeW) <= maxWidth && int(textSizeH) <= maxHeight {
		textImage, _, _ := createTextImage(texto, cor, fontsize)
		return textImage
	} else {
		lines := breakTextIntoLines(texto, fontsize, maxWidth, maxHeight)
		hSizeOfLines := len(lines) * int(textSizeH)
		textOp := &text.DrawOptions{}
		textOp.ColorScale.ScaleWithColor(cor)
		textImage := ebiten.NewImage(maxWidth, hSizeOfLines)

		for i, line := range lines {
			lineY := float64(i) * textSizeH
			textOp.GeoM.Reset()
			textOp.GeoM.Translate(0, lineY)
			text.Draw(textImage, line, face, textOp)
		}
		return textImage
	}
}

type Label struct {
	x, y     float64
	image    *ebiten.Image
	duration float64
}

func newTextLabel(text string, x, y float64) *Label {
	labelImage := ebiten.NewImage(145, 30)
	labelImage.Fill(color.White)
	op := &ebiten.DrawImageOptions{}
	labelImage.DrawImage(newTextImageMultiline(text, color.Black, 20, 200, 30), op)
	return &Label{x, y, labelImage, 0}
}

func (l *Label) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(l.x), float64(l.y))
	screen.DrawImage(l.image, op)
}

type Button struct {
	x, y, w, h     int
	image          *ebiten.Image
	clickedImage   *ebiten.Image
	function       func() error
	alreadyClicked bool
}

func (b *Button) checkClicked(m *Mouse) error {
	if m.X > b.x && m.X < b.x+b.w && m.Y > b.y && m.Y < b.y+b.h && m.LeftPressed {
		b.alreadyClicked = true
		return b.function()
	} else {
		b.alreadyClicked = false
		return nil
	}
}

func (b *Button) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(b.y))
	if b.alreadyClicked {
		screen.DrawImage(b.clickedImage, op)
	} else {
		screen.DrawImage(b.image, op)
	}
}

func newButton(w, h, x, y int, innerText string, function func() error) *Button {
	newImage := ebiten.NewImage(w, h)
	newImage.Fill(color.White)

	clickedImage := ebiten.NewImage(w, h)
	clickedImage.Fill(color.Gray{Y: 0x80}) // Different color for clicked state

	// draw text on the image
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(10, 10) // pequena margem
	textOp.ColorScale.ScaleWithColor(color.Black)

	text.Draw(newImage, innerText, &text.GoTextFace{
		Source: font,
		Size:   15.0,
	}, textOp)

	text.Draw(clickedImage, innerText, &text.GoTextFace{
		Source: font,
		Size:   15.0,
	}, textOp)

	return &Button{x, y, w, h, newImage, clickedImage, function, false}
}

// Cria os botões do menu principal.
func (g *Game) newMainMenuButtons() []*Button {
	buttonW := screenWidth / 8
	buttonH := screenHeight / 8
	x := (screenWidth - buttonW) / 2

	buttonDuel := newButton(buttonW, buttonH, x, screenHeight/2, "Duel", func() error {
		g.loadDuelMode()
		g.mainMenuButtons = nil // reset buttons, avoid being clicked again in other scenes
		return nil
	})
	buttonDeckEditor := newButton(buttonW, buttonH, x, screenHeight/2+buttonH+10, "Deck Editor", func() error {
		fmt.Println("Soon...")
		return nil
	})
	buttonExit := newButton(buttonW, buttonH, x, screenHeight/2+buttonH*2+20, "Exit", func() error {
		return ebiten.Termination
	})

	return []*Button{buttonDuel, buttonDeckEditor, buttonExit}
}

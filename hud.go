package main

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
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

type Label struct {
	x, y     float64
	image    *ebiten.Image
	duration float64
}

func newTextLabel(text string, x, y float64) *Label {
	labelImage := ebiten.NewImage(145, 30)
	labelImage.Fill(color.White)
	op := &ebiten.DrawImageOptions{}
	labelImage.DrawImage(newTextImageMultiline(text, color.Black, 20, 200), op)
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
	function       func() error
	alreadyClicked bool
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

func newButton(buttonSlice []*Button, innerText string, function func() error) []*Button {
	buttonW := screenWidth / 8
	buttonH := screenHeight / 8
	x := (screenWidth - buttonW) / 2
	y := (screenHeight-buttonH)/2 + len(buttonSlice)*(buttonH+10)

	newImage := ebiten.NewImage(buttonW, buttonH)
	newImage.Fill(color.White)

	// draw text on the image
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(10, 10) // pequena margem
	textOp.ColorScale.ScaleWithColor(color.Black)

	text.Draw(newImage, innerText, &text.GoTextFace{
		Source: font,
		Size:   20.0,
	}, textOp)

	button := &Button{x, y, screenWidth / 8, screenHeight / 8, newImage, function, false}

	buttonSlice = append(buttonSlice, button)
	return buttonSlice
}

// Cria os botões do menu principal.
func (g *Game) newButtons() []*Button {
	var buttons []*Button
	buttons = newButton(buttons, "Duel", func() error {
		g.loadDuelMode()
		g.loadDuelRenderer()
		g.mainMenuButtons = nil // Seems like this doesnt unload by itself when its not used

		return nil
	})
	buttons = newButton(buttons, "Deck Editor", func() error {
		fmt.Println("Soon...")
		return nil
	})
	buttons = newButton(buttons, "Exit", func() error {
		return ebiten.Termination
	})

	return buttons
}

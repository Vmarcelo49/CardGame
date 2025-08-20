package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Label represents a text label
type Label struct {
	X, Y  float64
	Image *ebiten.Image
	Alpha float32
}

// Draw renders the label to the screen
func (l *Label) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(l.X, l.Y)
	screen.DrawImage(l.Image, op)
}

// NewTextLabel creates a new text label
func NewTextLabel(text string, x, y float64, font *text.GoTextFaceSource) *Label {
	labelImage := ebiten.NewImage(145, 30)
	labelImage.Fill(color.White)
	op := &ebiten.DrawImageOptions{}
	labelImage.DrawImage(newTextImageMultiline(text, color.Black, 20, 200, 30, font), op)
	return &Label{x, y, labelImage, 0}
}

// Button represents a clickable button
type Button struct {
	X, Y, W, H     int
	Image          *ebiten.Image
	ClickedImage   *ebiten.Image
	Function       func() error
	AlreadyClicked bool
}

// CheckClicked checks if the button was clicked
func (b *Button) CheckClicked(m *Mouse) error {
	if m.X > b.X && m.X < b.X+b.W && m.Y > b.Y && m.Y < b.Y+b.H && m.LeftPressed {
		b.AlreadyClicked = true
		return b.Function()
	} else {
		b.AlreadyClicked = false
		return nil
	}
}

// Draw renders the button to the screen
func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.X), float64(b.Y))
	if b.AlreadyClicked {
		screen.DrawImage(b.ClickedImage, op)
	} else {
		screen.DrawImage(b.Image, op)
	}
}

// NewButton creates a new button
func NewButton(w, h, x, y int, innerText string, function func() error, font *text.GoTextFaceSource) *Button {
	newImage := ebiten.NewImage(w, h)
	newImage.Fill(color.White)

	clickedImage := ebiten.NewImage(w, h)
	clickedImage.Fill(color.Gray{Y: 0x80})

	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(10, 10)
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

// CreateMainMenuButtons creates the main menu buttons
func CreateMainMenuButtons(screenWidth, screenHeight int, font *text.GoTextFaceSource,
	duelFunc, deckEditorFunc, exitFunc func() error) []*Button {
	buttonW := screenWidth / 8
	buttonH := screenHeight / 8
	x := (screenWidth - buttonW) / 2

	buttonDuel := NewButton(buttonW, buttonH, x, screenHeight/2, "Duel", duelFunc, font)
	buttonDeckEditor := NewButton(buttonW, buttonH, x, screenHeight/2+buttonH+10, "Deck Editor", deckEditorFunc, font)
	buttonExit := NewButton(buttonW, buttonH, x, screenHeight/2+buttonH*2+20, "Exit", exitFunc, font)

	return []*Button{buttonDuel, buttonDeckEditor, buttonExit}
}

// Text utility functions

func BreakTextIntoLines(txt string, fontsize float64, maxWidth, maxHeight int, font *text.GoTextFaceSource) []string {
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

func CreateTextImage(texto string, cor color.Color, fontsize float64, font *text.GoTextFaceSource) (*ebiten.Image, float64, float64) {
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

func newTextImageMultiline(txt string, cor color.Color, fontsize float64, maxWidth, maxHeight int, font *text.GoTextFaceSource) *ebiten.Image {
	lines := BreakTextIntoLines(txt, fontsize, maxWidth, maxHeight)

	face := &text.GoTextFace{
		Source: font,
		Size:   fontsize,
	}

	_, lineHeight := text.Measure("A", face, 0)
	totalHeight := float64(len(lines)) * lineHeight

	textImage := ebiten.NewImage(maxWidth, int(totalHeight))

	for i, line := range lines {
		textOp := &text.DrawOptions{}
		textOp.GeoM.Translate(0, float64(i)*lineHeight)
		textOp.ColorScale.ScaleWithColor(cor)
		text.Draw(textImage, line, face, textOp)
	}

	return textImage
}

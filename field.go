package main

var contador int

type Field struct {
	player1Field []*Card
	player2Field []*Card

	player1Graveyard []*Card
	player2Graveyard []*Card

	middlescreen float32
}

func newField() *Field {
	field := new(Field)

	field.middlescreen = (screenHeight / 2) - (screenHeight / 18)

	return field
}

package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Card struct {
	Name    string `yaml:"cardName"`
	ID      int    `yaml:"cardID"`
	CType   int    `yaml:"cardType"`
	Artpath string `yaml:"artPath"`
}

func main() {
	card := Card{
		Name:    "Card Name",
		ID:      1,
		CType:   2,
		Artpath: "path/to/art"}

	out, err := yaml.Marshal(card)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	cardName := fmt.Sprint(card.ID, ".yaml")
	file, err := os.Create(cardName)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer file.Close()

	file.Write(out)

	fmt.Println(string(out))
}

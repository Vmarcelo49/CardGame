package main

import (
	"fmt"
)

type Keyword uint8

const (
	None     Keyword = 0
	Attacker         = 1 << iota
	Piercer
	Blocker
)

// TODO: add when this effect should happen, on summon on the case of first card.
func giveKeyword(args ...interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 arguments, got %d", len(args))
	}

	key, ok1 := args[0].(Keyword)
	target, ok2 := args[1].(*Card)
	if !ok1 || !ok2 {
		return fmt.Errorf("invalid argument types")
	}

	if target.CType != Creature {
		return fmt.Errorf("Card is not a creature")
	}

	// Adiciona o keyword ao campo de keywords usando uma operação bitwise OR
	target.Keywords |= key

	return nil
}

type cardFlags uint

const (
	cannotBeUsed        cardFlags = 0         // 0000 0000
	canBeNormalSummoned cardFlags = 1 << iota // 0000 0001
	canBeUsedSpeed1                           // 0000 0010
	canBeUsedSpeed2                           // 0000 0011
)

type Damageable interface {
    modifyHP(amount int)
}

func inflictDamage(args ...interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 arguments, got %d", len(args))
	}

	amount, ok1 := args[0].(int)
	if !ok1 {
		return fmt.Errorf("first argument must be an int")
	}

	switch target := args[1].(type) {
	case *Card:
		target.Stats.Life -= amount
		fmt.Println("Damage inflicted on Card:", target.Stats.Life)
	case *Player:
		target.HP -= amount
		fmt.Println("Damage inflicted on Player:", target.HP)
	default:
		return fmt.Errorf("invalid target type")
	}
	return nil
}

func inflictDamage2(amount int, target Damageable) {
	target.modifyHP(-amount)
}

func tapinhaEff() error{
	amount := 5
	target,err := getTarget("player,card")
	if err != nil{
		return errors.New("Unable to get a target for a card effect")
	}
	return inflictDamage(amount,target)
}

func getTarget(strValues string) interface{}{
	//split the str with ,
	// check for each valid target and return a pointer to it
}

func drawCard(gamestate *Gamestate, target string, amount uint) {
	if target == "player" {
		for i := 0; i < int(amount); i++ {
			gamestate.P1.drawCard()
		}
	} else {
		for i := 0; i < int(amount); i++ {
			gamestate.P2.drawCard()
		}

	}

}

func (c *Card) getItself(gamestate *Gamestate) bool {
	for _, card := range gamestate.P1.Hand {
		if card == c {
			return true
		}
	}
	return false
}

func removeCard(place []*Card, card *Card) []*Card {
	for i, c := range place {
		if c == card {
			return append(place[:i], place[i+1:]...)
		}
	}
	return place
}

func (c *Card) normalSummon(gamestate *Gamestate) error {
	if c.CType != Creature {
		return fmt.Errorf("This card is not a creature")
	}
	if c.Flags&canBeNormalSummoned == 0 {
		return fmt.Errorf("This card cannot be normal summoned")
	}

	gamestate.Field.P1 = append(gamestate.Field.P1, c)
	gamestate.P1.Hand = removeCard(gamestate.P1.Hand, c)

	fmt.Println("Card normal summoned:", c.Name)
	return nil
}

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
	canBeUsedSpeed2                           // 0000 0100
)

func inflictDamage(args ...interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 2 arguments, got %d", len(args))
	}

	amount, ok1 := args[0].(int)
	target, ok2 := args[1].(*Card)
	//target, ok2 := args[1].(*playerhp)
	if !ok1 || !ok2 {
		return fmt.Errorf("invalid argument types")
	}

	target.Stats.Life -= amount

	return nil
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

package main

import (
	"fmt"
)

type Effect struct {
	name     string // maybe not needed
	function func(args ...interface{}) error
	args     []interface{}
}

// Runs all functions in the Funcs slice and returns a slice of errors
func (f *Effect) Run() error {
	err := f.function(f.args...)
	return err
}

func newCardEffect(name string, fun func(args ...interface{}) error, args ...interface{}) Effect {
	return Effect{
		name:     name,
		function: fun,
		args:     args,
	}
}

func (c *Card) firstCardEffExample() {
	effect := newCardEffect("Gives a keyword to itself", giveKeyword, Attacker, c)
	c.Effects = effect
}

var funcMap = map[string]func(args ...interface{}) error{
	"giveKeyword": giveKeyword,
}

type Keyword uint8

const (
	None Keyword = iota
	Attacker
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
	// Doesnt start a Chain
	if target.CType != Creature {
		return fmt.Errorf("Card is not a creature")
	}
	switch key {
	case Attacker:
		target.Keywords = append(target.Keywords, Attacker)
	case Piercer:
		target.Keywords = append(target.Keywords, Piercer)
	case Blocker:
		target.Keywords = append(target.Keywords, Blocker)
	}
	return nil
}

package main

import (
	"errors"
	"fmt"
)

type Effect struct {
	functions []func(args ...interface{}) error
	args      []interface{}
	text      string
}

func newEffect(text string, fun func(args ...interface{}) error, args ...interface{}) Effect {
	return Effect{
		functions: []func(args ...interface{}) error{fun},
		args:      args,
		text:      text}
}

type Card struct {
	x, y int
	eff  Effect
}

func (c *Card) run() error {
	return c.eff.functions[0](c.eff.args...)
}

func say(args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("no arguments provided")
	}
	for _, myvalue := range args {
		fmt.Println(myvalue)
	}
	return nil
}

func main() {
	effect := newEffect("Prints stuff", say, "Hello", "World", 10, 1)
	myCard := &Card{10, 10, effect}
	if err := myCard.run(); err != nil {
		fmt.Println("Error:", err)
	} // Call the actual function using the reference and handle error
}

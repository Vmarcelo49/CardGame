package main

import (
	"errors"
	"fmt"
)

var funcMap = map[string]func() error{
	"giveKeyword": giveKeyword,
}

type Funcs struct {
	functions []func() error
}

// Run executa todas as funções armazenadas em Funcs e retorna um slice de erros
func (f *Funcs) Run() []error {
	var errors []error
	for _, fun := range f.functions {
		err := fun()
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (f *Funcs) AddFuncByName(name string) error {
	if fn, ok := funcMap[name]; ok {
		f.functions = append(f.functions, fn)
		return nil
	}
	return errors.New(fmt.Sprint("Function ", name, " not found"))
}

type Keyword uint8

const (
	None Keyword = iota
	Attacker
	Piercer
	Blocker
)

func giveKeyword(key Keyword, target *Card) {
	// Doesnt start a Chain
	if target.CType != Creature {
		fmt.Println("Card is not a creature")
		return
	}
	switch key {
	case Attacker:
		target.Keywords = append(target.Keywords, Attacker)
	case Piercer:
		target.Keywords = append(target.Keywords, Piercer)
	case Blocker:
		target.Keywords = append(target.Keywords, Blocker)
	}

}

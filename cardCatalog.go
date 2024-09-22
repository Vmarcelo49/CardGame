package main

// Cria uma carta apartir de um ID, se o ID não existir, retorna uma carta padrão indicando erro, não dá os valores de X,Y,W,H e Selected
func newCardFromID(cardID int) *Card {
	if card, found := cardCatalog[cardID]; found {
		// Retorna uma cópia da carta para evitar alterações no catálogo original
		newCard := *card
		return &newCard
	}

	// Se o ID da carta não for encontrado, retorna uma carta padrão indicando erro
	return &Card{
		Name:  "Carta Fora do Range de IDs",
		ID:    0,
		CType: 0,
		Stats: Stats{
			Attack: 0,
			Life:   1,
		},
		SubType:  "Nil",
		Text:     "ID dessa carta não existe no catálogo, arrume o deck.",
		Flags:    0,
		Keywords: 0,
	}
}

var cardCatalog = map[int]*Card{
	1: {
		Name:    "Ricardo, o Cavalo",
		ID:      1,
		SubType: "Animal",
		Text:    "A violência e ignorância resolve boa parte dos problemas. - Um cavalo humanoide qualquer",
		CType:   Creature,
		Stats: Stats{
			Attack: 15,
			Life:   5,
		},
		Flags:    canBeNormalSummoned,
		Keywords: Attacker,
	},
	2: {
		Name:    "Tapinha",
		ID:      2,
		CType:   Spell,
		SubType: "fast?",
		Text:    "Cause 5 de dano a qualquer coisa. - Um tapinha não dói",
		Flags:   canBeUsedSpeed2,
	},
	3: {
		Name:    "Totem da drenagem vital",
		ID:      3,
		CType:   Permanent,
		SubType: "Permanent Token",
		Text:    "Para ativar é necessário ter mais de 10 de vida. Receba 15 de dano cada começo de turno. Criaturas que você controla ganham ataque igual a vida. ",
		Flags:   cannotBeUsed,
	},
}

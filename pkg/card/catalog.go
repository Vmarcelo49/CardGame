package card

// NewCardFromID creates a card from an ID, returns default error card if ID not found
func NewCardFromID(cardID int) *Card {
	if card, found := cardCatalog[cardID]; found {
		newCard := *card
		return &newCard
	}

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
		Flags:    CanBeNormalSummoned,
		Keywords: Attacker,
	},
	2: {
		Name:    "Tapinha",
		ID:      2,
		CType:   Spell,
		SubType: "fast?",
		Text:    "Cause 5 de dano a qualquer coisa. - Um tapinha não dói",
		Flags:   CanBeUsedSpeed2,
	},
	3: {
		Name:    "Totem da drenagem vital",
		ID:      3,
		CType:   Permanent,
		SubType: "Permanent Token",
		Text:    "Para ativar é necessário ter mais de 10 de vida. Receba 15 de dano cada começo de turno. Criaturas que você controla ganham ataque igual a vida. ",
		Flags:   CannotBeUsed,
	},
}

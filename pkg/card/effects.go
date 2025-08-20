package card

// DrawCard will be moved to logic package to avoid circular dependencies
// RemoveCard removes a card from a slice of cards
func RemoveCard(place []*Card, card *Card) []*Card {
	for i, c := range place {
		if c == card {
			return append(place[:i], place[i+1:]...)
		}
	}
	return place
}

// InflictDamage deals damage to a target
func InflictDamage(amount int, target Damageable) {
	target.ModifyHP(-amount)
}

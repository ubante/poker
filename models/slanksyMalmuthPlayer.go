package models

// This player will go all-in if their hold cards are at a certain
// Sklansky-Malmuth score.
// https://en.wikipedia.org/wiki/Texas_hold_%27em_starting_hands#Sklansky_hand_groups
type SlanksyMalmuthPlayer struct {
	GenericPlayer
	threshold int
}

func NewSklanskyMalmuthPlayer(name string, level int) SlanksyMalmuthPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(SlanksyMalmuthPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.threshold = level

	return *newPlayer
}

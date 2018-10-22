package models

import "fmt"

// This Player is pure aggro.  He will always shove.
type AllInAlwaysPlayer struct {
	GenericPlayer
}

func NewAllInAlwaysPlayer(name string) AllInAlwaysPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(AllInAlwaysPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack

	return *newPlayer
}

func (ap *AllInAlwaysPlayer) chooseAction(t *Table) {
	fmt.Println("I am always all-in always.")
	ap.allIn()
}
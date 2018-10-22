package models

import "fmt"

// This Player likes to see a specific street.  But then check/folds.
type StreetTestPlayer struct {
	GenericPlayer
	foldingStreet string
}

func NewStreetTestPlayer(name string, street string) StreetTestPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(StreetTestPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.foldingStreet = street
	newPlayer.typeName = "StreetTestPlayer"

	return *newPlayer
}

func (stp *StreetTestPlayer) chooseAction(t *Table) {
	if stp.hasFolded {
		fmt.Println("I have already folded so no action.  How did this codepath happen btw?.")
		return
	}

	if t.bettingRound == stp.foldingStreet {
		fmt.Printf("(%s) @%s: I've seen my street so check-folding\n", stp.getName(), stp.foldingStreet)
		stp.checkOrFold(t)
	} else {
		fmt.Printf("(%s) @%s: check-calling\n", stp.getName(), stp.foldingStreet)
		stp.checkOrCall(t)
	}
}

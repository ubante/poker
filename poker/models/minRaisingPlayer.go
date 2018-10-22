package models

import "fmt"

type MinRaisingPlayer struct {
	GenericPlayer
}

// This Player makes sure that no street goes unbet.
func NewMinRaisingPlayer(name string) MinRaisingPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(MinRaisingPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.typeName = "MinRaisingPlayer"

	return *newPlayer
}

func (mrp *MinRaisingPlayer) chooseAction(t *Table) {
	currentTableBet := t.getMaxBet()
	if currentTableBet > mrp.bet {
		fmt.Printf("My bet ($%d) is below the table's bet ($%d) so calling.\n", mrp.bet, currentTableBet)
		mrp.call(t)
	} else {
		if mrp.bet == 0 {
			if t.bigBlindValue > mrp.stack {
				fmt.Println("I'm down to less than a BB so just checking.")
				mrp.checkOrCall(t)
				return
			}
			fmt.Printf("To keep the bluffers nervous, min-raising $%d\n", t.bigBlindValue)
			mrp.raise(t.bigBlindValue)
		} else {
			fmt.Println("I have bet but I'm now behind the bet.  Calling")
			mrp.call(t)
		}
	}
}
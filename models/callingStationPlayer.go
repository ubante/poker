package models

import "fmt"

// This Player will always call.
type CallingStationPlayer struct {
	GenericPlayer
}

// Repeating the constructor is kinda lame.
func NewCallingStationPlayer(name string) CallingStationPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(CallingStationPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack

	return *newPlayer
}

func (csp *CallingStationPlayer) chooseAction(t *Table) {
	fmt.Println("This is CSP overriding chooseAction")
	currentTableBet := t.getMaxBet()

	if currentTableBet > csp.bet {
		fmt.Printf("My bet ($%d) is below the table's bet ($%d) so calling.\n", csp.bet, currentTableBet)
		csp.call(t)
	} else {
		fmt.Printf("My bet ($%d) is >= the table bet ($%d) so just checking.\n", csp.bet, currentTableBet)
		csp.check(t)
	}
}
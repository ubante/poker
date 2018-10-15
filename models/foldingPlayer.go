package models

import "fmt"

// This Player will check if able, otherwise will fold.  This models an
// online player who has disconnected.  If the present players are
// aggressive enough, it is common to see a FoldingPlayer make it to the
// money.
type FoldingPlayer struct {
	GenericPlayer
}

func NewFoldingPlayer(name string) FoldingPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(FoldingPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.typeName = "FoldingPlayer"

	return *newPlayer
}

func (fp *FoldingPlayer) chooseAction(t *Table) {
	fmt.Println("This is FP overriding chooseAction")
	currentTableBet := t.getMaxBet()

	if currentTableBet > fp.bet {
		fmt.Printf("My bet ($%d) is below the table's bet ($%d) so folding like a champ.\n", fp.bet, currentTableBet)
		fp.fold()
	} else {
		fmt.Printf("My bet ($%d) is >= the table bet ($%d) so I'm still in it.\n", fp.bet, currentTableBet)
		fp.check(t)
	}
}
package models

import "fmt"

// This Player is a calling station unless the current bet is >=5x his
// current non-zero bet.  In that case, he will fold.
type CallToFivePlayer struct {
	GenericPlayer
}

func NewCallToFivePlayer(name string) CallToFivePlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(CallToFivePlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.typeName = "CallToFivePlayer"

	return *newPlayer
}

func (ct5p *CallToFivePlayer) chooseAction(t *Table) {
	fmt.Println("This is CT5 overriding chooseAction")
	currentTableBet := t.getMaxBet()

	// First check or call if I haven't bet yet.
	if ct5p.bet == 0 {
		ct5p.checkOrCall(t)
	}

	// Then call up to 5x my bet.
	if currentTableBet >= 5*ct5p.bet {
		fmt.Printf("My bet ($%d) is less than a fifth of the table's bet ($%d) so folding.\n",
			ct5p.bet, currentTableBet)
		ct5p.fold()
	} else if currentTableBet > ct5p.bet {
		fmt.Printf("My bet ($%d) is below the table's bet ($%d) so calling.\n", ct5p.bet, currentTableBet)
		ct5p.call(t)
	} else {
		fmt.Printf("My bet ($%d) is >= the table bet ($%d) so just checking.\n", ct5p.bet, currentTableBet)
		ct5p.check(t)
	}
}
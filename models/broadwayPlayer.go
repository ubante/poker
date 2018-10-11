package models

import (
	"fmt"
	"os"
)

// This player will bet X BB if his holecards are 10 or better.
// Otherwise, check-fold.  For the streets, let's keep it simple and
// compute present odds.  If it's greater than 50%, then check-call,
// otherwise, check-fold.
type BroadwayPlayer struct {
	GenericPlayer
	preFlopBetMultiplier int
}

func NewBroadwayPlayer(name string, preFlopBet int) BroadwayPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(BroadwayPlayer)
	newPlayer.name = name 
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.preFlopBetMultiplier = preFlopBet

	return *newPlayer
}

func (bp *BroadwayPlayer) chooseAction(t *Table) {
	if bp.hasFolded {
		fmt.Println("I have already folded so no action.  How did this codepath happen btw?.")
		os.Exit(4)
		return
	}

	fmt.Println(bp.String())
	if t.bettingRound == "PREFLOP" {
		if bp.holeCards.cardSet.cards[0].NumericalRank >= 10 && bp.holeCards.cardSet.cards[1].NumericalRank >= 10 {
			myHappyPlace := t.bigBlindValue*bp.preFlopBetMultiplier
			if t.getMaxBet() < myHappyPlace {
				fmt.Printf("(%s) BroadwayPlayer has a good hand - raising from %d to %d.\n",
					bp.getName(), t.getMaxBet(), myHappyPlace)
				bp.raiseUpTo(myHappyPlace)
			} else {
				bp.call(t)
			}
		} else {
			bp.checkOrFold(t)
		}

		return
	}

	bp.computePresentOdds(*t, bp.getHoleCardsCardSet())
	if bp.winningOdds > 0.50 {
		fmt.Printf("We have present winning odds of %3.2f so check-calling anything.\n", bp.winningOdds)
		bp.checkOrCall(t)
	} else {
		fmt.Printf("We have present winning odds of %3.2f which isn't great so check-folding.\n",
			bp.winningOdds)
		bp.checkOrFold(t)
	}

}
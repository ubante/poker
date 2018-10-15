package models

import (
	"goven/poker/matrix"
	"fmt"
)

// This player is based on OddsComputingPlayer.  It will compute the
// current odds at each street and bet accordingly.

type VariableOddsPlayer struct {
	OddsComputingPlayer
}

func NewVariableOddsPlayer(name string, smLevel int, computedOddsPercentageLevel float64) VariableOddsPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(VariableOddsPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.smThreshold = smLevel
	newPlayer.oddsThreshold = computedOddsPercentageLevel
	newPlayer.typeName = "VariableOddsPlayer"

	// These should be overwritten by the caller.
	newPlayer.preFlopRaise = -1
	newPlayer.postFlopRaise = -1
	newPlayer.turnRaise = -1
	newPlayer.riverRaise = -1

	return *newPlayer
}

func (vop *VariableOddsPlayer) chooseAction(t *Table) {
	if t.bettingRound == "PREFLOP" {
		// Check the SM score of hole cards to decide what to do.
		smScore := matrix.GetSMModifiedScore()
		hcScore := smScore.GetScoreOfHoleCardStrings(vop.holeCards.cardSet.cards[0].ToString(),
			vop.holeCards.cardSet.cards[1].ToString())
		if hcScore > vop.smThreshold {
			fmt.Println("Holecards have a score of", hcScore, "so just checkfolding.")
			vop.checkOrFold(t)
			return
		}

		proposedRaise := t.bigBlindValue*vop.preFlopRaise
		fmt.Println("Holecards have a score of", hcScore, "so raising.")
		fmt.Println("I have a stack of:", vop.stack)
		if t.getMaxBet() >= proposedRaise {
			fmt.Println("Holecards have a score of", hcScore, "and bet is", t.getMaxBet(),
				"which is greater than", vop.preFlopRaise, "so just calling.")
			vop.call(t)
		} else {
			fmt.Println("Holecards have a score of", hcScore, "and bet is", t.getMaxBet(),
				"which is less than", vop.preFlopRaise, "times the BB so raising.")
			betDiff := proposedRaise - vop.bet
			vop.raiseUpTo(betDiff)
			fmt.Println("I have raised up", betDiff, "to", proposedRaise)
		}

		return
	}

	vop.computePresentOdds(*t, vop.getHoleCardsCardSet())
	// TODO need logic for the other streets
}
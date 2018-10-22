package models

import (
	"fmt"
	"os"
	"goven/poker/matrix"
)

// This player will go all-in if their hold cards are at a certain
// Sklansky-Malmuth score.
// https://en.wikipedia.org/wiki/Texas_hold_%27em_starting_hands#Sklansky_hand_groups
type SlanksyMalmuthModifiedPlayer struct {
	GenericPlayer
	threshold int  // Lesser values are better.
}

func NewSklanskyMalmuthModifiedPlayer(name string, level int) SlanksyMalmuthModifiedPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(SlanksyMalmuthModifiedPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.threshold = level
	newPlayer.typeName = "SklanskyMalmuthModifiedPlayer"

	return *newPlayer
}

func (smp *SlanksyMalmuthModifiedPlayer) chooseAction(t *Table) {
	if smp.hasFolded {
		fmt.Println("I have already folded so no action.  How did this codepath happen btw?.")
		os.Exit(4)
		return
	}

	if t.bettingRound == "PREFLOP" {
		// Check the SM score of hole cards to decide what to do.
		smScore := matrix.GetSMModifiedScore()
		hcScore := smScore.GetScoreOfHoleCardStrings(smp.holeCards.cardSet.cards[0].ToString(),
			smp.holeCards.cardSet.cards[1].ToString())
		if hcScore > smp.threshold {
			fmt.Println("Holecards have a score of", hcScore, "so just checkfolding.")
			smp.checkOrFold(t)
		} else {
			fmt.Println("Holecards have a score of", hcScore, "so going all-in.")
			smp.allIn()
		}

	} else {
		fmt.Printf("(%s) SM%d: check-calling\n", smp.getName(), smp.threshold)
		smp.checkOrCall(t)
	}
}
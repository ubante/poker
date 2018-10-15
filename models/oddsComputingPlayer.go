package models

import (
	"fmt"
	"os"
	"goven/poker/matrix"
	"math"
)

// This player will go all-in if their hold cards are at a certain
// Sklansky-Malmuth score.

// Preflop, this player will go 5BB if its SM-modified score is above a
// given level.  Otherwise, this player will check-fold.  If reraised,
// this player will check entire stack.
//
// On the streets, this player will compute the odds of winning.  If
// odds are greater than the given threshold, raise to full pot or
// check to full stack.  If odds are lesser than the given threshold,
// fold to raises that don't pay out given the computed odds.

type OddsComputingPlayer struct {
	GenericPlayer
	smThreshold   int     // Lesser values means play tighter.
	oddsThreshold float64 // Percentage of winning odds that will trigger a raise.
	preFlopRaise  int     // Multiplier of big blind
	postFlopRaise float64 // Multiplier of pot
	turnRaise     float64 // Multiplier of pot
	riverRaise    float64 // Multiplier of pot
	winningOdds   float64 // Computed odds of having the best hand if evaluated right now
}

func NewOddsComputingPlayer(name string, smLevel int, computedOddsPercentageLevel float64) OddsComputingPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(OddsComputingPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.smThreshold = smLevel
	newPlayer.oddsThreshold = computedOddsPercentageLevel
	newPlayer.typeName = "OddsComputingPlayer"

	// These should be overwritten by the caller.
	newPlayer.preFlopRaise = -1
	newPlayer.postFlopRaise = -1
	newPlayer.turnRaise = -1
	newPlayer.riverRaise = -1

	return *newPlayer
}

// TODO need a better toString() so we can see all the raising fields.

// TODO If this player is the last player, it will still try to bet
// which doesn't hurt but is wasteful and could hurt with a future
// player. The tournament class should know this and just finish the
// hand.
func (ocp *OddsComputingPlayer) chooseAction(t *Table) {
	if ocp.hasFolded {
		fmt.Println("I have already folded so no action.  How did this codepath happen btw?.")
		os.Exit(110)
		return
	}

	// Make sure that some fields are set.  This should be done soon
	// after the player object is instantiated.
	if ocp.preFlopRaise == -1 {
		fmt.Printf("FATAL: %s does not have all its ___Raise fields properly set.", ocp.name)
		os.Exit(4)
	}

	if t.bettingRound == "PREFLOP" {
		// Check the SM score of hole cards to decide what to do.
		smScore := matrix.GetSMModifiedScore()
		hcScore := smScore.GetScoreOfHoleCardStrings(ocp.holeCards.cardSet.cards[0].ToString(),
			ocp.holeCards.cardSet.cards[1].ToString())
		if hcScore > ocp.smThreshold {
			fmt.Println("Holecards have a score of", hcScore, "so just checkfolding.")
			ocp.checkOrFold(t)
			return
		}

		proposedRaise := t.bigBlindValue*ocp.preFlopRaise
		fmt.Println("Holecards have a score of", hcScore, "so raising.")
		fmt.Println("I have a stack of:", ocp.stack)
		if t.getMaxBet() >= proposedRaise {
			fmt.Println("Holecards have a score of", hcScore, "and bet is", t.getMaxBet(),
				"which is greater than", ocp.preFlopRaise, "so just calling.")
			ocp.call(t)
		} else {
			fmt.Println("Holecards have a score of", hcScore, "and bet is", t.getMaxBet(),
				"which is less than", ocp.preFlopRaise, "times the BB so raising.")
			betDiff := proposedRaise - ocp.bet
			ocp.raiseUpTo(betDiff)
			fmt.Println("I have raised up", betDiff, "to", proposedRaise)
		}

		return
	}

	// Compute the odds.
	ocp.computePresentOdds(*t, ocp.getHoleCardsCardSet())
	var streetRaise float64
	switch t.bettingRound {
	case "FLOP":
		streetRaise = ocp.postFlopRaise
	case "TURN":
		streetRaise = ocp.turnRaise
	case "RIVER":
		streetRaise = ocp.riverRaise
	default:
		fmt.Println("FATAL: The bettingRound was incorrectly defined to:", t.bettingRound)
		os.Exit(22)
	}

	if ocp.winningOdds > ocp.oddsThreshold {
		fmt.Printf("There is a %3.2f chance of winning which is greater than the threshold of %3.2f " +
			"so raising.\n", ocp.winningOdds, ocp.oddsThreshold)
		raiseTarget := float64(t.pot.getValue()) * streetRaise
		raiseTargetInt := int(math.Round(raiseTarget))
		fmt.Printf("The pot is $%d and we have a multiplier of %3.2f so raising up to $%d.\n",
			t.pot.getValue(), streetRaise, raiseTargetInt)
		ocp.raiseUpTo(raiseTargetInt)
		return
	} else {
		fmt.Printf("The odds of %s winning, %3.2f%%, did not meet the odds threshold of %3.2f%% so just " +
			"checkfolding.\n", ocp.name, ocp.winningOdds, ocp.oddsThreshold)
		ocp.checkOrFold(t)
		return
	}
}

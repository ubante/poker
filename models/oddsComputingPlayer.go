package models

import (
	"fmt"
	"os"
	"goven/poker/matrix"
	"sort"
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

	return *newPlayer
}

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
	ocp.computeOdds(*t, ocp.getHoleCardsCardSet())
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
		fmt.Printf("The odds of winning, %3.2f%%, did not meet the odds threshold of %3.2f%% so just " +
			"checkfolding.\n", ocp.winningOdds, ocp.oddsThreshold)
		ocp.checkOrFold(t)
		return
	}
}

func (ocp *OddsComputingPlayer) computeOdds(table Table, heroHCCS CardSet) {
	fmt.Println("Computing odds....")
	fmt.Printf(table.GetStatus())
	heroCombinedCardSet := heroHCCS.Combine(*table.community.cards)
	fmt.Println("Hero's combined cards:", heroCombinedCardSet)

	heroCombinedCardSet.FindBestHand()
	fmt.Println("Hero's best eval is:", heroCombinedCardSet.bestEval)

	// Make a new deck and remove the hero's hole cards and the
	// post-flop community cards.
	nonCheatingDeck := NewDeck()  // Note that this deck is not shuffled.
	for _, card := range heroCombinedCardSet.cards {
		// Despite the method name, this will just remove the card from
		// the deck.
		nonCheatingDeck.getCardOfValue(card.ToString())
	}

	// Brute force the villian's hands.
	deckLength := nonCheatingDeck.length()
	fmt.Println("\nThere are", deckLength, "cards left in the deck.")
	comboCounter := 0
	heroLoses := 0
	heroTies := 0
	strongestVillainNewStreetHand := NewCardSet()
	winningVillainHandMap := make(map[int]int)

	for i := 0; i < deckLength-1; i++ {
		for j := i+1; j < deckLength; j++ {
			comboCounter++

			villainCardSet := NewCardSet()
			villainCardSet.Add(*nonCheatingDeck.cardSet.cards[i])
			villainCardSet.Add(*nonCheatingDeck.cardSet.cards[j])
			villainCombinedCardSet := villainCardSet.Combine(*table.community.cards)
			villainCombinedCardSet.FindBestHand()

			// A higher score is better here.
			if villainCombinedCardSet.bestEval.flattenedScore > heroCombinedCardSet.bestEval.flattenedScore {
				heroLoses++

				if strongestVillainNewStreetHand.isEmpty() {
					strongestVillainNewStreetHand = villainCombinedCardSet
				} else if villainCombinedCardSet.bestEval.flattenedScore > strongestVillainNewStreetHand.bestEval.flattenedScore {
					strongestVillainNewStreetHand = villainCombinedCardSet
				}

				// See Evaluation() for the full list.  Higher is better
				// where 9 is a straight flush and 1 is a high card.
				primaryRank := villainCombinedCardSet.bestEval.primaryRank
				if _, ok := winningVillainHandMap[primaryRank]; ok {
					winningVillainHandMap[primaryRank]++
				} else {
					winningVillainHandMap[primaryRank] = 1
				}

				continue
			}

			if villainCombinedCardSet.bestEval.flattenedScore == heroCombinedCardSet.bestEval.flattenedScore {
				heroTies++

				continue
			}
		}
	}

	fmt.Println("Just to repeat, Hero's best eval is:", heroCombinedCardSet.bestEval)

	heroWins := comboCounter - heroLoses - heroTies
	ocp.winningOdds = 100*float64(heroWins)/float64(comboCounter)

	fmt.Printf("Of the %d possibilities,\n %d (%4.1f%%) result in loss for the hero,\n %d (%4.1f%%) " +
		"result in ties,\n and %d (%4.1f%%) result in wins.",
		comboCounter, heroLoses, 100*float64(heroLoses)/float64(comboCounter), heroTies,
		100*float64(heroTies)/float64(comboCounter), heroWins, ocp.winningOdds)

	// Break down the hands where the villian wins by hand rank.
	var sortedRanks []int
	for rank := range winningVillainHandMap {
		sortedRanks = append(sortedRanks, rank)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedRanks)))

	fmt.Println("\nHere's the breakdown of hands beat the hero's hand:")
	for _, primaryRank := range sortedRanks {
		fmt.Printf("%16s: %4.1f%% (%d) \n", decodeEvaluationPrimaryRank(primaryRank),
			100*float64(winningVillainHandMap[primaryRank])/float64(comboCounter),
			winningVillainHandMap[primaryRank])
	}

	fmt.Println("\nThe strongest possible villain hand is:\n", strongestVillainNewStreetHand.bestEval)
	fmt.Println("... done computing odds.")

}
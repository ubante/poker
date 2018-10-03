package models

import (
	"fmt"
	"os"
	"goven/poker/matrix"
	"sort"
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
	smThreshold int  // Lesser values means play tighter.
	oddsThreshold int  // Percentage of winning odds tha trigger a raise.
}

func NewOddsComputingPlayer(name string, smLevel int, computedOddsPercentageLevel int) OddsComputingPlayer {
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
		} else {
			fmt.Println("Holecards have a score of", hcScore, "so going all-in.")
			ocp.allIn()
		}

		return
	}

	// Compute the odds.

	// Below block is just until I finish the logic to compute the odds.
	fmt.Printf("(%s) SM%d: check-calling\n", ocp.getName(), ocp.smThreshold)
	ocp.checkOrCall(t)
}

// This should be a method to avoid namespace collisions.
func computeOdds(table Table, heroHCCS CardSet) {
	fmt.Printf(table.GetStatus())
	heroCombinedCardSet := heroHCCS.Combine(*table.community.cards)
	fmt.Println("Hero's combined cards:", heroCombinedCardSet)

	heroCombinedCardSet.FindBestHand()
	fmt.Println("Hero's best eval is:", heroCombinedCardSet.bestEval)

	// Brute force the villian's hands.
	deckLength := len(table.deck.cardSet.cards)
	fmt.Println("\nThere are", deckLength, "cards left in the deck.")
	comboCounter := 0
	heroLoses := 0
	heroTies := 0
	strongestVillainNewStreetHand := NewCardSet()
	winningVillianHandMap := make(map[int]int)

	for i := 0; i < deckLength-1; i++ {
		for j := i+1; j < deckLength; j++ {
			comboCounter++
			//fmt.Printf("%2d: %s %s\n", comboCounter, table.deck.cardSet.cards[i], table.deck.cardSet.cards[j])

			villainCardSet := NewCardSet()
			villainCardSet.Add(*table.deck.cardSet.cards[i])
			villainCardSet.Add(*table.deck.cardSet.cards[j])
			villainCombinedCardSet := villainCardSet.Combine(*table.community.cards)
			villainCombinedCardSet.FindBestHand()

			// A higher score is better here.
			if villainCombinedCardSet.bestEval.flattenedScore > heroCombinedCardSet.bestEval.flattenedScore {
				heroLoses++
				//fmt.Println("Hero LOSES to Villain:", villainCombinedCardSet.bestEval)

				if strongestVillainNewStreetHand.isEmpty() {
					strongestVillainNewStreetHand = villainCombinedCardSet
				} else if villainCombinedCardSet.bestEval.flattenedScore > strongestVillainNewStreetHand.bestEval.flattenedScore {
					strongestVillainNewStreetHand = villainCombinedCardSet
				}

				// See Evaluation() for the full list.  Higher is better
				// where 9 is a straight flush and 1 is a high card.
				primaryRank := villainCombinedCardSet.bestEval.primaryRank
				if _, ok := winningVillianHandMap[primaryRank]; ok {
					winningVillianHandMap[primaryRank]++
				} else {
					winningVillianHandMap[primaryRank] = 1
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
	fmt.Printf("Of the %d possibilities,\n %d (%4.1f%%) result in loss for the hero,\n %d (%4.1f%%) result in ties,\n and %d (%4.1f%%) result in wins.",
		comboCounter, heroLoses, 100*float64(heroLoses)/float64(comboCounter), heroTies,
		100*float64(heroTies)/float64(comboCounter), heroWins, 100*float64(heroWins)/float64(comboCounter))

	// Break down the hands where the villian wins by hand rank.
	var sortedRanks []int
	for rank := range winningVillianHandMap {
		sortedRanks = append(sortedRanks, rank)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedRanks)))

	fmt.Println("\nHere's the breakdown of hands beat the hero's hand:")
	for _, primaryRank := range sortedRanks {
		fmt.Printf("%16s: %4.1f%% (%d) \n", decodeEvaluationPrimaryRank(primaryRank),
			100*float64(winningVillianHandMap[primaryRank])/float64(comboCounter),
			winningVillianHandMap[primaryRank])
	}


	fmt.Println("\nThe strongest possible villain hand is:\n", strongestVillainNewStreetHand.bestEval)
}
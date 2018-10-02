package models

import (
	"fmt"
	"sort"
)

// This is just to cast a GenericPlayer to Player.
//func castToPlayer(player Player) Player {
//	return player
//}

func compute(table Table, heroHCCS CardSet) {
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
				//fmt.Println("Hero TIES with Villain:", villainCombinedCardSet.bestEval)

				continue
			}

			//fmt.Println("  hero WINS")
		}
	}

	fmt.Println("Just to repeat, Hero's best eval is:", heroCombinedCardSet.bestEval)

	heroWins := comboCounter - heroLoses - heroTies
	fmt.Printf("Of the %d possibilities,\n %d (%4.1f%%) result in loss for the hero,\n %d (%4.1f%%) result in ties,\n and %d (%4.1f%%) result in wins.",
		comboCounter, heroLoses, 100*float64(heroLoses)/float64(comboCounter), heroTies,
		100*float64(heroTies)/float64(comboCounter), heroWins, 100*float64(heroWins)/float64(comboCounter))
	//fmt.Println()

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

func VillainGuesser() {
	var table Table
	table.Initialize()
	villain := NewGenericPlayer("Victor")
	table.AddPlayer(&villain)
	hero := NewGenericPlayer("Hari")
	table.AddPlayer(&hero)

	fmt.Printf(table.GetStatus())

	table.Preset()
	hero.addHoleCard(*table.deck.getCard())
	hero.addHoleCard(*table.deck.getCard())

	// Use the below to test hole cards.  Comment out the above
	// getCard() lines.
	//hero.addHoleCard(*table.deck.getCardOfValue("HA"))
	//hero.addHoleCard(*table.deck.getCardOfValue("SA"))

	fmt.Printf(table.GetStatus())
	heroHCCS := hero.getHoleCardsCardSet()

	fmt.Println("\n================= Flop =================")
	table.DealFlop()
	compute(table, heroHCCS)

	fmt.Println("\n================= Turn =================")
	table.dealTurn()
	compute(table, heroHCCS)

	fmt.Println("\n================= River =================")
	table.dealRiver()
	compute(table, heroHCCS)

}

func ComputeFlop(heroHandString, communityString string) {
	var table Table
	table.Initialize()
	hero := NewGenericPlayer("Hank")
	table.AddPlayer(&hero)
	table.Preset()

	heroCard1 := string([]rune(heroHandString)[0:2])
	heroCard2 := string([]rune(heroHandString)[2:4])
	hero.addHoleCard(*table.deck.getCardOfValue(heroCard1))
	hero.addHoleCard(*table.deck.getCardOfValue(heroCard2))
	fmt.Println(table.GetStatus())

	flopCard1 := string([]rune(communityString)[0:2])
	flopCard2 := string([]rune(communityString)[2:4])
	flopCard3 := string([]rune(communityString)[4:6])
	table.community.add(*table.deck.getCardOfValue(flopCard1))
	table.community.add(*table.deck.getCardOfValue(flopCard2))
	table.community.add(*table.deck.getCardOfValue(flopCard3))
	fmt.Println(table.GetStatus())

	heroHCCS := hero.getHoleCardsCardSet()

	fmt.Println("\n================= Given Flop =================")
	compute(table, heroHCCS)

	fmt.Println("\n================= Random Turn =================")
	table.dealTurn()
	compute(table, heroHCCS)

	fmt.Println("\n================= Random River =================")
	table.dealRiver()
	compute(table, heroHCCS)

}

 /*
Starting to guess....
=====================
Adding Hari to the front.
------
 -- 2 players left (game #0)
Hari: [] $0/$1000
Victor: [] $0/$1000
Pot: 0
Community:
Bet totals: 0
Stack totals: 2000
------
------
 -- 2 players left (game #1)
Hari: [ST H3] $0/$1000
Victor: [] $0/$1000
Pot: 0
Community:
Bet totals: 0
Stack totals: 2000
------

================= Flop =================
------
FLOP -- 2 players left (game #1)
Hari: [ST H3] $0/$1000
Victor: [] $0/$1000
Pot: 0
Community: S4 C6 CT
Bet totals: 0
Stack totals: 2000
------
Hero's combined cards: ST H3 S4 C6 CT
Hero's best eval is: ST H3 S4 C6 CT: pair with ranks: [2 10 6 4 3 0]

There are 47 cards left in the deck.
Just to repeat, Hero's best eval is: ST H3 S4 C6 CT: pair with ranks: [2 10 6 4 3 0]
Of the 1081 possibilities,
 116 (10.7%) result in loss for the hero,
 6 ( 0.6%) result in ties,
 and 959 (88.7%) result in wins.
Here's the breakdown of hands beat the hero's hand:
 three of a kind:  0.6% (7)
        two pair:  1.9% (21)
        one pair:  8.1% (88)

The strongest possible villain hand is:
 HT DT S4 C6 CT: three of a kind with ranks: [4 10 6 0 4 0]

================= Turn =================
------
TURN -- 2 players left (game #1)
Hari: [ST H3] $0/$1000
Victor: [] $0/$1000
Pot: 0
Community: S4 C6 CT D7
Bet totals: 0
Stack totals: 2000
------
Hero's combined cards: ST H3 S4 C6 CT D7
Hero's best eval is: ST S4 C6 CT D7: pair with ranks: [2 10 7 6 4 0]

There are 46 cards left in the deck.
Just to repeat, Hero's best eval is: ST S4 C6 CT D7: pair with ranks: [2 10 7 6 4 0]
Of the 1035 possibilities,
 179 (17.3%) result in loss for the hero,
 14 ( 1.4%) result in ties,
 and 842 (81.4%) result in wins.
Here's the breakdown of hands beat the hero's hand:
        straight:  4.3% (44)
 three of a kind:  1.0% (10)
        two pair:  4.3% (45)
        one pair:  7.7% (80)

The strongest possible villain hand is:
 C9 D8 C6 CT D7: straight with ranks: [5 10 0 0 0 0]

================= River =================
------
RIVER -- 2 players left (game #1)
Hari: [ST H3] $0/$1000
Victor: [] $0/$1000
Pot: 0
Community: S4 C6 CT D7 H6
Bet totals: 0
Stack totals: 2000
------
Hero's combined cards: ST H3 S4 C6 CT D7 H6
Hero's best eval is: ST C6 CT D7 H6: two pairs with ranks: [3 10 6 7 0 0]

There are 45 cards left in the deck.
Just to repeat, Hero's best eval is: ST C6 CT D7 H6: two pairs with ranks: [3 10 6 7 0 0]
Of the 990 possibilities,
 216 (21.8%) result in loss for the hero,
 28 ( 2.8%) result in ties,
 and 746 (75.4%) result in wins.
Here's the breakdown of hands beat the hero's hand:
  four of a kind:  0.1% (1)
      full house:  2.3% (23)
        straight:  4.4% (44)
 three of a kind:  7.1% (70)
        two pair:  7.9% (78)

The strongest possible villain hand is:
 D6 S6 C6 CT H6: quads with ranks: [8 6 10 0 0 0]

Process finished with exit code 0

  */
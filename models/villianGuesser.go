package models

import (
	"fmt"
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
		comboCounter, heroLoses, 100*float64(heroLoses)/float64(comboCounter), heroTies, 100*float64(heroTies)/float64(comboCounter),
		heroWins, 100*float64(heroWins)/float64(comboCounter))

	fmt.Println("\nThe strongest possible villain hand is:\n", strongestVillainNewStreetHand.bestEval)
}

func Guess() {
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

 /*
Starting to guess....
=====================
Adding Hari to the end.
------
 -- 2 players left (game #0)
Victor: [] $0/$1000
Hari: [] $0/$1000
Pot: 0
Community:
Bet totals: 0
Stack totals: 2000
------
------
 -- 2 players left (game #1)
Victor: [] $0/$1000
Hari: [HQ SQ] $0/$1000
Pot: 0
Community:
Bet totals: 0
Stack totals: 2000
------

================= Flop =================
------
FLOP -- 2 players left (game #1)
Victor: [] $0/$1000
Hari: [HQ SQ] $0/$1000
Pot: 0
Community: H4 DJ SK
Bet totals: 0
Stack totals: 2000
------
Hero's combined cards: HQ SQ H4 DJ SK
Hero's best eval is: HQ SQ H4 DJ SK: pair with ranks: [2 12 13 11 4 0]

There are 47 cards left in the deck.
Just to repeat, Hero's best eval is: HQ SQ H4 DJ SK: pair with ranks: [2 12 13 11 4 0]
Of the 1081 possibilities,
 156 (14.4%) result in loss for the hero,
 1 ( 0.1%) result in ties,
 and 924 (85.5%) result in wins.
The strongest possible villain hand is:
 HK DK H4 DJ SK: three of a kind with ranks: [4 13 11 0 4 0]

================= Turn =================
------
TURN -- 2 players left (game #1)
Victor: [] $0/$1000
Hari: [HQ SQ] $0/$1000
Pot: 0
Community: H4 DJ SK CT
Bet totals: 0
Stack totals: 2000
------
Hero's combined cards: HQ SQ H4 DJ SK CT
Hero's best eval is: HQ SQ DJ SK CT: pair with ranks: [2 12 13 11 10 0]

There are 46 cards left in the deck.
Just to repeat, Hero's best eval is: HQ SQ DJ SK CT: pair with ranks: [2 12 13 11 10 0]
Of the 1035 possibilities,
 190 (18.4%) result in loss for the hero,
 1 ( 0.1%) result in ties,
 and 844 (81.5%) result in wins.
The strongest possible villain hand is:
 DA CQ DJ SK CT: straight with ranks: [5 14 0 0 0 0]

================= River =================
------
RIVER -- 2 players left (game #1)
Victor: [] $0/$1000
Hari: [HQ SQ] $0/$1000
Pot: 0
Community: H4 DJ SK CT D2
Bet totals: 0
Stack totals: 2000
------
Hero's combined cards: HQ SQ H4 DJ SK CT D2
Hero's best eval is: HQ SQ DJ SK CT: pair with ranks: [2 12 13 11 10 0]

There are 45 cards left in the deck.
Just to repeat, Hero's best eval is: HQ SQ DJ SK CT: pair with ranks: [2 12 13 11 10 0]
Of the 990 possibilities,
 217 (21.9%) result in loss for the hero,
 1 ( 0.1%) result in ties,
 and 772 (78.0%) result in wins.
The strongest possible villain hand is:
 DA CQ DJ SK CT: straight with ranks: [5 14 0 0 0 0]

Process finished with exit code 0

 */
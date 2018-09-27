package models

import (
	"fmt"
)

// This is just to cast a GenericPlayer to Player.
//func castToPlayer(player Player) Player {
//	return player
//}

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
	//table.DealHoleCards()
	fmt.Printf(table.GetStatus())

	fmt.Println("\n================= Flop =================")
	table.DealFlop()
	fmt.Printf(table.GetStatus())

	// At this point, players have their two hole cards and the flop has
	// been dealt.  Pretend we don't know the villain's hands and find
	// the odds of our hero having a better hand.
	heroHCCS := hero.getHoleCardsCardSet()
	heroCombinedCardSet := heroHCCS.Combine(*table.community.cards)
	//heroHCCS.
	fmt.Println("Hero's combined cards:", heroCombinedCardSet)

	heroCombinedCardSet.FindBestHand()
	fmt.Println("Hero's best eval is:", heroCombinedCardSet.bestEval)

	// Brute force the villian's hands.
	deckLength := len(table.deck.cardSet.cards)
	fmt.Println("\nThere are", deckLength, "cards left in the deck.")
	comboCounter := 0
	heroLoses := 0
	heroTies := 0
	strongestVillainFlopHand := NewCardSet()

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

				if strongestVillainFlopHand.isEmpty() {
					strongestVillainFlopHand = villainCombinedCardSet
				} else if villainCombinedCardSet.bestEval.flattenedScore > strongestVillainFlopHand.bestEval.flattenedScore {
					strongestVillainFlopHand = villainCombinedCardSet
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

	fmt.Println("\nThe strongest possible villain hand is:\n", strongestVillainFlopHand.bestEval)
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
FLOP -- 2 players left (game #1)
Hari: [SJ S6] $0/$1000
Victor: [ST DQ] $0/$1000
Pot: 0
Community: D9 CT S5
Bet totals: 0
Stack totals: 2000
------
Hero's combined cards: SJ S6 D9 CT S5
Hero's best eval is: SJ S6 D9 CT S5: high card with ranks: [1 11 10 9 6 5]

<snip>

Hero LOSES to Villain's best eval is: D5 DK SQ HK DA: pair with ranks: [2 13 14 12 5 0]
Hero LOSES to Villain's best eval is: ST HJ SQ HK DA: straight with ranks: [5 14 0 0 0 0]
Hero LOSES to Villain's best eval is: ST DK SQ HK DA: pair with ranks: [2 13 14 12 10 0]
Hero LOSES to Villain's best eval is: C6 DK SQ HK DA: pair with ranks: [2 13 14 12 6 0]
Hero LOSES to Villain's best eval is: HJ DK SQ HK DA: pair with ranks: [2 13 14 12 11 0]
Just to repeat, Hero's best eval is: D4 C4 SQ HK DA: pair with ranks: [2 4 14 13 12 0]
Of the 990 possibilities,
 375 (37.9%) result in loss for the hero,
 1 ( 0.1%) result in ties,
 and 614 (62.0%) result in wins.



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
Hari: [CT D8] $0/$1000
Victor: [] $0/$1000
Pot: 0
Community:
Bet totals: 0
Stack totals: 2000
------

================= Flop =================
------
FLOP -- 2 players left (game #1)
Hari: [CT D8] $0/$1000
Victor: [] $0/$1000
Pot: 0
Community: C2 HT C4
Bet totals: 0
Stack totals: 2000
------
Hero's combined cards: CT D8 C2 HT C4
Hero's best eval is: CT D8 C2 HT C4: pair with ranks: [2 10 8 4 2 0]

There are 47 cards left in the deck.
Just to repeat, Hero's best eval is: CT D8 C2 HT C4: pair with ranks: [2 10 8 4 2 0]
Of the 1081 possibilities,
 92 ( 8.5%) result in loss for the hero,
 6 ( 0.6%) result in ties,
 and 983 (90.9%) result in wins.
The strongest possible villain hand is:
 DT ST C2 HT C4: three of a kind with ranks: [4 10 4 0 2 0]

 */
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
	table.DealHoleCards()
	//fmt.Printf(table.GetStatus())

	table.DealFlop()
	fmt.Printf(table.GetStatus())

	// At this point, players have their two hole cards and the flop has
	// been dealt.  Pretend we don't know the villain's hands and find
	// the odds of our hero having a better hand.
	heroHoleCards := hero.getHoleCards()
	heroCombinedCardSet := heroHoleCards.Combine(*table.community.cards)
	fmt.Println("Hero's combined cards:", heroCombinedCardSet)

	heroCombinedCardSet.FindBestHand()
	fmt.Println("Hero's best eval is:", heroCombinedCardSet.bestEval)

	// Brute force the villian's hands.
	deckLength := len(table.deck.cardSet.cards)
	fmt.Println("\nThere are", deckLength, "cards left in the deck.")
	comboCounter := 0
	betterHandCounter := 0
	for i := 0; i < deckLength-1; i++ {
		for j := i+1; j < deckLength; j++ {
			comboCounter++
			fmt.Printf("%2d: %s %s\n", comboCounter, table.deck.cardSet.cards[i], table.deck.cardSet.cards[j])
			//var vhc HoleCards
		}
	}


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

Process finished with exit code 0

 */
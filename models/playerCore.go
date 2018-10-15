package models

import (
	"fmt"
	"os"
	"sort"
)

// This is most-likely an anti-pattern.
type Player interface {
	fold()
	allIn()
	check(t *Table)
	preset()
	setNextPlayer(p *Player)
	getNextPlayer() *Player
	setPreviousPlayer(p *Player)
	getPreviousPlayer() *Player
	setName(n string)
	getName() string
	raise(raiseAmount int)
	setBet(newBet int)
	getBet() int
	addToStack(payout int)
	getStack() int
	//addHoleCard(c models.Card)
	addHoleCard(c Card)
	getHoleCardsCardSet() CardSet
	payBlind(blindAmount int)
	checkHasFolded() bool
	checkIsAllIn() bool
	chooseAction(t *Table)
}

type GenericPlayer struct {
	name           string
	bankRoll	   int
	typeName	   string

	// The below gets set before each tournament.
	nextPlayer     *Player
	previousPlayer *Player

	// The below get preset before each game.
	holeCards HoleCards
	stack     int
	bet       int
	hasFolded bool
	isAllIn   bool
	winningOdds	float64 // Computed odds of having the best hand if evaluated right now
}

// GenericPlayer constructor
// http://www.golangpatterns.info/object-oriented/constructors
// Maybe this could be replaced with new() and some helper lines.
func NewGenericPlayer(name string) GenericPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars
	newPlayer := GenericPlayer{name, 0,"GenericPlayer", nil, nil,
		hc, initialStack, 0,false, false, 0.0}
	return newPlayer
}

func (gp *GenericPlayer) getStatus() string {
	status := ""
	status = fmt.Sprintf("%s: [%s] stack=%d bet=%d", gp.name, gp.holeCards.ToString(), gp.stack, gp.bet)

	return status
}

func (gp GenericPlayer) String() string {
	return fmt.Sprintf("%s: [%s] $%d/$%d", gp.name, gp.holeCards.ToString(), gp.bet, gp.stack)
}

/*
Preset before each game.

Maybe use "prepare()" instead of "Preset()" because the latter implies
something you do afterwards.
*/
func (gp *GenericPlayer) preset() {
	// Maybe NewGenericPlayer can call this?
	ecs := NewCardSet()
	gp.holeCards = HoleCards{cardSet: &ecs}
	gp.bet = 0
	gp.hasFolded = false
	gp.isAllIn = false
}

func (gp *GenericPlayer) setNextPlayer(p *Player) {
	gp.nextPlayer = p
}

func (gp *GenericPlayer) getNextPlayer() *Player {
	return gp.nextPlayer
}

func (gp *GenericPlayer) setPreviousPlayer(p *Player) {
	gp.previousPlayer = p
}

func (gp *GenericPlayer) getPreviousPlayer() *Player {
	return gp.previousPlayer
}

func (gp *GenericPlayer) setName(n string) {
	gp.name = n
}

func (gp *GenericPlayer) getName() string {
	return gp.name
}

func (gp *GenericPlayer) setBet(newBet int) {
	gp.bet = newBet
}

func (gp *GenericPlayer) getBet() int {
	return gp.bet
}

func (gp *GenericPlayer) addToStack(payout int) {
	gp.stack += payout
}

func (gp *GenericPlayer) getStack() int {
	return gp.stack
}

//func (gp *GenericPlayer) addHoleCard(c models.Card) {
func (gp *GenericPlayer) addHoleCard(c Card) {
	gp.holeCards.Add(c)
}

func (gp *GenericPlayer) getHoleCardsCardSet() CardSet {
	return *gp.holeCards.cardSet
}

func (gp *GenericPlayer) payBlind(blindAmount int) {
	if gp.stack > blindAmount {
		gp.stack -= blindAmount
		gp.bet = blindAmount
	} else {
		gp.allIn()
	}
}

func (gp *GenericPlayer) check(t *Table) {
	// This may one day do something.

	return
}

func (gp *GenericPlayer) call(t *Table) {
	maxBet := t.getMaxBet()
	increase := maxBet - gp.bet
	stackBefore := gp.stack
	betBefore := gp.bet

	if increase >= gp.stack {
		gp.allIn()
	} else {
		gp.bet = maxBet
		gp.stack -= increase
	}

	fmt.Printf("calling - max bet:%d, increase:%d, bet before/after:%d/%d, stack before/after:%d/%d\n",
		maxBet, increase, betBefore, gp.bet, stackBefore, gp.stack)
}

func (gp *GenericPlayer) fold() {
	gp.hasFolded = true
	gp.holeCards.toss()
}

// I know there is a difference between bet and raise but that seems
// like irrelevant semantics here.
func (gp *GenericPlayer) raise(raiseAmount int) {
	if gp.stack < raiseAmount {
		fmt.Println(gp.getName(),"tried to raise more than its stack.  This is a fatal error.")
		fmt.Printf("Stack: $%d, Raise: $%d\n", gp.stack, raiseAmount)
		//fmt.Println(gp)
		os.Exit(9)
	}

	// The below block may be redundant if this method was called from
	// allIn() but it's best to have it in case raise() is called from
	// a method that doesn't check if the player is all-in or not.
	if raiseAmount == gp.stack {
		gp.isAllIn = true
	}

	gp.bet += raiseAmount
	gp.stack -= raiseAmount
}

// This is the safer alternative to gp.raise().
func (gp *GenericPlayer) raiseUpTo(raiseAmount int) {
	if raiseAmount >= gp.stack {
		gp.allIn()
		return
	}

	gp.raise(raiseAmount)
}

func (gp *GenericPlayer) allIn() {
	gp.raise(gp.stack)
	gp.isAllIn = true
}

func (gp *GenericPlayer) checkOrCall(t *Table) {
	if t.getMaxBet() == 0 {
		gp.check(t)
		return
	}

	gp.call(t)

	return
}

func (gp *GenericPlayer) checkOrFold(t *Table) {
	if t.getMaxBet() > gp.getBet() {
		fmt.Printf("%s thinks the pot ($%d) is too rich - folding.\n", gp.getName(), t.getMaxBet())
		gp.fold()
		return
	}

	gp.check(t)
	return
}

func (gp *GenericPlayer) chooseActionPreflop(t *Table) {
	fmt.Println("Using the default preflop action of check-calling")
	gp.checkOrCall(t)

	return
}

func (gp *GenericPlayer) chooseActionFlop(t *Table) {
	fmt.Println("Using the default flop action of check-folding")
	gp.checkOrFold(t)

	return
}

func (gp *GenericPlayer) chooseActionTurn(t *Table) {
	fmt.Println("Using the default turn action of check-folding")
	gp.checkOrFold(t)

	return
}

func (gp *GenericPlayer) chooseActionRiver(t *Table) {
	fmt.Println("Using the default river action of check-folding")
	gp.checkOrFold(t)

	return
}

func (gp *GenericPlayer) chooseAction(t *Table) {
	// This is handled by Table() but be redundant for clearness.
	if gp.hasFolded {
		fmt.Println("I,", gp.name, "have already checkHasFolded so no action.")
		return
	}
	if gp.isAllIn {
		fmt.Println("I,", gp.name, "am all in so no action.")
		return
	}

	fmt.Println(gp.getStatus())

	switch t.bettingRound {
	case "PREFLOP":
		gp.chooseActionPreflop(t)
	case "FLOP":
		gp.chooseActionFlop(t)
	case "TURN":
		gp.chooseActionTurn(t)
	case "RIVER":
		gp.chooseActionRiver(t)
	}
	return
}

func (gp *GenericPlayer) checkHasFolded() bool {
	return gp.hasFolded
}

func (gp *GenericPlayer) checkIsAllIn() bool {
	return gp.isAllIn
}

func (gp *GenericPlayer) computePresentOdds(table Table, heroHCCS CardSet) {
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
	gp.winningOdds = 100*float64(heroWins)/float64(comboCounter)

	fmt.Printf("Of the %d possibilities,\n %d (%4.1f%%) result in loss for the hero,\n %d (%4.1f%%) " +
		"result in ties,\n and %d (%4.1f%%) result in wins.",
		comboCounter, heroLoses, 100*float64(heroLoses)/float64(comboCounter), heroTies,
		100*float64(heroTies)/float64(comboCounter), heroWins, gp.winningOdds)

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

package models

import (
	"fmt"
	"os"
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
	getHoleCards() CardSet
	payBlind(blindAmount int)
	checkHasFolded() bool
	checkIsAllIn() bool
	chooseAction(t *Table)
}

type GenericPlayer struct {
	name           string
	nextPlayer     *Player
	previousPlayer *Player

	// The below get preset after each game.
	holeCards HoleCards
	stack     int
	bet       int
	hasFolded bool
	isAllIn   bool
}

// GenericPlayer constructor
// http://www.golangpatterns.info/object-oriented/constructors
// Maybe this could be replaced with new() and some helper lines.
func NewGenericPlayer(name string) GenericPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars
	newPlayer := GenericPlayer{name, nil, nil, hc, initialStack, 0,
		false, false}
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

Maybe use "prepare()" instead of "preset()" because the latter implies
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

func (gp *GenericPlayer) getHoleCards() CardSet {
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
		fmt.Println("Player tried to raise more than its stack.  This is a fatal error.")
		fmt.Println(gp)
		os.Exit(9)
	}

	gp.bet += raiseAmount
	gp.stack -= raiseAmount
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

// This Player will always call.
type CallingStationPlayer struct {
	GenericPlayer
}

// Repeating the constructor is kinda lame.
func NewCallingStationPlayer(name string) CallingStationPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(CallingStationPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack

	return *newPlayer
}

func (csp *CallingStationPlayer) chooseAction(t *Table) {
	fmt.Println("This is CSP overriding chooseAction")
	currentTableBet := t.getMaxBet()

	if currentTableBet > csp.bet {
		fmt.Printf("My bet ($%d) is below the table's bet ($%d) so calling.\n", csp.bet, currentTableBet)
		csp.call(t)
	} else {
		fmt.Printf("My bet ($%d) is >= the table bet ($%d) so just checking.\n", csp.bet, currentTableBet)
		csp.check(t)
	}
}

// This Player likes to see a specific street.  But then check/folds.
type StreetTestPlayer struct {
	GenericPlayer
	foldingStreet string
}

func NewStreetTestPlayer(name string, street string) StreetTestPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(StreetTestPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.foldingStreet = street

	return *newPlayer
}

func (stp *StreetTestPlayer) chooseAction(t *Table) {
	if stp.hasFolded {
		fmt.Println("I have already folded so no action.  How did this codepath happen btw?.")
		return
	}

	if t.bettingRound == stp.foldingStreet {
		fmt.Printf("(%s) @%s: I've seen my street so check-folding\n", stp.getName(), stp.foldingStreet)
		stp.checkOrFold(t)
	} else {
		fmt.Printf("(%s) @%s: check-calling\n", stp.getName(), stp.foldingStreet)
		stp.checkOrCall(t)
	}
}

type MinRaisingPlayer struct {
	GenericPlayer
}

// This Player makes sure that no street goes unbet.
func NewMinRaisingPlayer(name string) MinRaisingPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(MinRaisingPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack

	return *newPlayer
}

func (mrp *MinRaisingPlayer) chooseAction(t *Table) {
	currentTableBet := t.getMaxBet()
	if currentTableBet > mrp.bet {
		fmt.Printf("My bet ($%d) is below the table's bet ($%d) so calling.\n", mrp.bet, currentTableBet)
		mrp.call(t)
	} else {
		if mrp.bet == 0 {
			if t.bigBlindValue > mrp.stack {
				fmt.Println("I'm down to less than a BB so just checking.")
				mrp.checkOrCall(t)
				return
			}
			fmt.Printf("To keep the bluffers nervous, min-raising $%d\n", t.bigBlindValue)
			mrp.raise(t.bigBlindValue)
		} else {
			fmt.Println("I have bet but I'm now behind the bet.  Calling")
			mrp.call(t)
		}
	}
}

// This Player will check if able, otherwise will fold.  This models an
// online player who has disconnected.  If the present players are
// aggressive enough, it is common to see a FoldingPlayer make it to the
// money.
type FoldingPlayer struct {
	GenericPlayer
}

func NewFoldingPlayer(name string) FoldingPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(FoldingPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack

	return *newPlayer
}

func (fp *FoldingPlayer) chooseAction(t *Table) {
	fmt.Println("This is FP overriding chooseAction")
	currentTableBet := t.getMaxBet()

	if currentTableBet > fp.bet {
		fmt.Printf("My bet ($%d) is below the table's bet ($%d) so folding like a champ.\n", fp.bet, currentTableBet)
		fp.fold()
	} else {
		fmt.Printf("My bet ($%d) is >= the table bet ($%d) so I'm still in it.\n", fp.bet, currentTableBet)
		fp.check(t)
	}
}

// This Player is pure aggro.  He will always shove.
type AllInAlwaysPlayer struct {
	GenericPlayer
}

func NewAllInAlwaysPlayer(name string) AllInAlwaysPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(AllInAlwaysPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack

	return *newPlayer
}

func (ap *AllInAlwaysPlayer) chooseAction(t *Table) {
	fmt.Println("I am always all-in always.")
	ap.allIn()
}

// This Player is a calling station unless the current bet is >=5x his
// current non-zero bet.  In that case, he will fold.
type CallToFivePlayer struct {
	GenericPlayer
}

func NewCallToFivePlayer(name string) CallToFivePlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(CallToFivePlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack

	return *newPlayer
}

func (ct5p *CallToFivePlayer) chooseAction(t *Table) {
	fmt.Println("This is CT5 overriding chooseAction")
	currentTableBet := t.getMaxBet()

	// First check or call if I haven't bet yet.
	if ct5p.bet == 0 {
		ct5p.checkOrCall(t)
	}

	// Then call up to 5x my bet.
	if currentTableBet >= 5*ct5p.bet {
		fmt.Printf("My bet ($%d) is less than a fifth of the table's bet ($%d) so folding.\n",
			ct5p.bet, currentTableBet)
		ct5p.fold()
	} else if currentTableBet > ct5p.bet {
		fmt.Printf("My bet ($%d) is below the table's bet ($%d) so calling.\n", ct5p.bet, currentTableBet)
		ct5p.call(t)
	} else {
		fmt.Printf("My bet ($%d) is >= the table bet ($%d) so just checking.\n", ct5p.bet, currentTableBet)
		ct5p.check(t)
	}

}

// These are future players.

// This player will go all-in if their hold cards are at a certain
// Sklansky-Malmuth score.
// https://en.wikipedia.org/wiki/Texas_hold_%27em_starting_hands#Sklansky_hand_groups
type SlanksyMalmuthPlayer struct {
	GenericPlayer
	threshold int
}

func NewSklanskyMalmuthPlayer(name string, level int) SlanksyMalmuthPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(SlanksyMalmuthPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.threshold = level

	return *newPlayer
}

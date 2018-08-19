package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"gopkg.in/inconshreveable/log15.v2"
	"time"
	"os"
)

type Card struct {
	suit         string
	numericaRank int
	rank         string
}

func NewCard(s string, nr int) Card {
	var c Card
	c.suit = s
	c.numericaRank = nr

	if nr == 14 {
		c.rank = "A" // Aces are aces.
	}

	switch nr {
	case 14:
		c.rank = "A" // Aces are aces.
	case 13:
		c.rank = "K"
	case 12:
		c.rank = "Q"
	case 11:
		c.rank = "J"
	case 10:
		c.rank = "T"
	default:
		c.rank = strconv.Itoa(nr)
	}

	return c
}

func (c *Card) toString() string {
	return fmt.Sprintf("%s%s", c.suit, c.rank)
}

func (c Card) String() string {
	return c.toString()
}

func (c Card) isSuited(c2 Card) bool {
	if c.suit == c2.suit {
		return true
	}

	return false
}

func (c Card) isPaired(c2 Card) bool {
	if c.numericaRank == c2.numericaRank {
		return true
	}

	return false
}

type CardSet struct {
	cards []*Card
	bestHand []*CardSet
	possibleHands []*CardSet
}

func NewCardSet(cards ...Card) CardSet {
	var cs CardSet

	for _, card := range cards {
		cs.add(card)
	}

	return cs
}

// I should rename toString to String and explicitly call
// SomeType.String() when I want a string outside of Println()
func (cs CardSet) toString() string {
	var toString string

	for _, c := range cs.cards {
		if toString == "" {
			toString = c.toString()
			continue
		}

		toString += fmt.Sprintf(" %s", c.toString())
	}

	return toString
}

func (cs CardSet) String() string {
	return cs.toString()
}

func (cs *CardSet) add(c Card) {
	cs.cards = append(cs.cards, &c)
}

/*
This will accept a cardset and return the union of it with this.
*/
func (cs *CardSet) combine(cs2 CardSet) CardSet {
	var combined CardSet

	for _, card := range cs.cards {
		combined.add(*card)
	}
	for _, card := range cs2.cards {
		combined.add(*card)
	}

	return combined
}

// Find all the possible combinations.
// I like it gross.
func (cs *CardSet) setPossibleHands() {
	cards := cs.cards
	for a := 0; a < len(cards)-4; a++ {
		for b := a+1; b < len(cards)-3; b++ {
			for c := b+1; c < len(cards)-2; c++ {
				for d := c+1; d < len(cards)-1; d++ {
					for e := d+1; e < len(cards); e++ {
						possibleHand := NewCardSet(*cards[a], *cards[b], *cards[c], *cards[d], *cards[e])
						cs.possibleHands = append(cs.possibleHands, &possibleHand)
					}
				}
			}
		}
	}
}

func (cs *CardSet) findBestHand() {
	cs.setPossibleHands()
	for i, ph := range cs.possibleHands {
		fmt.Println(i, "evalling:", ph)
		eval := NewEvaluation(*ph)
		fmt.Println(eval)
	}
}

/*
I couldn't find a pop function - weird.
https://groups.google.com/forum/#!topic/Golang-nuts/obZI4uyZTe0
*/
func (cs *CardSet) pop() *Card {
	card := cs.cards[0]
	copy(cs.cards, cs.cards[1:])
	cs.cards = cs.cards[:len(cs.cards)-1]

	return card
}

func (cs *CardSet) shuffle() {
	rand.Shuffle(cs.length(), func(i, j int) {
		cs.cards[i], cs.cards[j] = cs.cards[j], cs.cards[i]
	})
}

func (cs CardSet) length() int {
	return len(cs.cards)
}

/*
Primary ranks are:
  9: Straight flush
  8: Four of a kind
  7: Full house
  6: Flush
  5: Straight
  4: Three of a kind
  3: Two pair
  2: One pair
  1: High card

Secondary ranks vary with the primary rank.  For example,  for
a straight, the secondary rank will be the highest card.  For
a full house, the secondary rank will be the numeric value of
the three of a kind.

Tertiary rank also vary's with the primary rank.  For a flush,
this would be the second highest card.  For a full house, this
would be the numeric value of the pair.

This just goes on and on.  For a flush and high card, it is
possible to have a rank for each of the five cards.  For the
other types of hands, these will remain at zero.

It is possible for different hands to have the same evaluation.
*/
type Evaluation struct {
	cardSet   *CardSet
	humanEval string
	allRanks  []string  // do I still need this?
	primaryRank int
	secondaryRank int
	tertiaryRank int
	quaternaryRank int
	quinaryRank int
	senaryRank int
}

func NewEvaluation(cardSet CardSet) *Evaluation {
	var eval Evaluation
	eval.cardSet = &cardSet
	eval.evaluate()
	return &eval
}

func (e Evaluation) String() string {
	toString := fmt.Sprintf("%s: %s with ranks: [%d %d %d %d %d %d]", e.cardSet, e.humanEval,
		e.primaryRank, e.secondaryRank, e.tertiaryRank, e.quaternaryRank, e.quinaryRank, e.senaryRank)

	return toString
}

func (e *Evaluation) evaluate() {
	e.humanEval = "TBDeval"

}

func (e Evaluation) compare(otherEval Evaluation) bool {
	// TODO flesh this out
	return false
}

type HoleCards struct {
	cardSet *CardSet
}

func (hc *HoleCards) toString() string {
	if hc.cardSet == nil {
		return ""
	}

	return hc.cardSet.toString()
}

func (hc HoleCards) String() string {
	return hc.toString()
}

func (hc *HoleCards) add(c Card) {
	hc.cardSet.add(c)
}

func (hc *HoleCards) empty() {
	ecs := NewCardSet()
	hc.cardSet = &ecs
}

func (hc *HoleCards) toss() {
	hc.empty()
}

type Deck struct {
	cardset *CardSet // Deck is not a CardSet.
}

func NewDeck() *Deck {
	var d Deck
	ecs := NewCardSet()
	d.cardset = &ecs

	for _, suit := range []string{"S", "H", "D", "C"} {
		// https://stackoverflow.com/questions/21950244/is-there-a-way-to-iterate-over-a-range-of-integers-in-golang
		for numericRank := range [13]int{} {
			newCard := NewCard(suit, numericRank+2)
			d.cardset.add(newCard)
		}
	}

	return &d
}

func (d *Deck) getStatus() string {
	status := fmt.Sprintf("Deck has %d cards:\n", len(d.cardset.cards))
	for i, card := range d.cardset.cards {
		if i != 0 && i%13 == 0 {
			status += "\n"
		}
		status += card.toString()
		status += " "
	}

	return status
}

func (d Deck) length() int {
	return d.cardset.length()
}

func (d *Deck) shuffle() {
	d.cardset.shuffle()
}

func (d *Deck) getCard() *Card {
	return d.cardset.pop()
}

// Maybe I can get away with just using a CardSet instead of this
// struct.
type Community struct {
	cards *CardSet
}

func NewCommunity() Community {
	ecs := NewCardSet()

	var c Community
	c.cards = &ecs

	return c
}

// My first String().  Will eventually replace all the getStatus()
// methods.
func (c Community) String() string {
	return c.cards.String()
	//return c.cards.getStatus()
}

func (c *Community) add(card Card) {
	c.cards.add(card)
}

/*
In live poker, there can be multiple pots when someone goes in but
there's still action on the table.

Here, we will have one pot, the main pot.  But we'll keep track of
each players equity.  Once the game is complete, we'll segment the
pot.  Using the equity, we know how much of the pot is split across
all players.  This is the first segment.  And how much of the pot is
split across all-1 players.  This is the second segment.  This will
continue until all equity is accounted for.

For each segment, we'll find take the involved players and find who
has the best hand and that player will get that segment of the pot.
Folded players will have dead hands so their equity remains in the
pot but their hand cannot win.

If all players have equal equity, then there will be one segment
and the player with the best hand gets the whole pot.

For example, if one player goes all in and everyone else calls
without going all in themselves, and there is no further betting,
then each player has equal equity and we still have one segment.

For another example, if PlayerA goes all on the flop and everyone
folds except PlayerB and PlayerC.  Those two players call.  Then
on the turn, they check-check.  Then on the river, they raise-call,
then we have two segments.  The first segment are for exactly the
three remaining players.  The second segment is just for PlayerB
and PlayerC.

Note that bets do not enter a pot until all betting is complete for
the round.
*/
type Pot struct {
	value  int // This could be gotten from summing equity.
	equity map[*Player]int
}

func NewPot() Pot {
	var pot Pot

	pot.value = 0
	pot.equity = make(map[*Player]int)

	return pot
}

func (p Pot) String() string {
	toString := fmt.Sprintf("POT is $%d\n", p.value)

	for p, value := range p.equity {
		player := *p
		toString += fmt.Sprintf("%s has equity: $%d\n", player.getName(), value)
	}

	return toString
}

func (p *Pot) addEquity(playerBet int, player *Player) {
	p.value += playerBet
	p.equity[player] += playerBet
}

func (p *Pot) getSegments() map[int][]*Player {

	// First invert the map.
	invertedMap := make(map[int][]*Player)
	for p, equity := range p.equity {
		player := *p
		invertedMap[equity] = append(invertedMap[equity], &player)
	}

	// Then sort the values of equity.  Note that this also removes
	// duplicate values.
	var sortedEquity []int
	for eq := range invertedMap {
		sortedEquity = append(sortedEquity, eq)
	}

	// Loop through the equities and create the reverse map.
	segments := make(map[int][]*Player)
	for _, equity := range sortedEquity {
		fmt.Printf("$%d: \n", equity)
		for _, p := range invertedMap[equity] {
			player := *p
			fmt.Printf("    %s\n", player.getName())
			segments[equity] = append(segments[equity], &player)
		}
		fmt.Println()
	}

	fmt.Println("Returning from getSegments().")
	return segments
}

// This is probably an anti-pattern.
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
	setBet(newBet int)
	getBet() int
	getStack() int
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
	status = fmt.Sprintf("%s: [%s] stack=%d bet=%d", gp.name, gp.holeCards.toString(), gp.stack, gp.bet)

	return status
}

func (gp GenericPlayer) String() string {
	return fmt.Sprintf("%s: [%s] $%d/$%d", gp.name, gp.holeCards.toString(), gp.bet, gp.stack)
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

func (gp *GenericPlayer) getStack() int {
	return gp.stack
}

func (gp *GenericPlayer) addHoleCard(c Card) {
	gp.holeCards.add(c)
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

func (gp *GenericPlayer) allIn() {
	gp.bet += gp.stack
	gp.stack = 0
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
// current bet.  In that case, he will fold.
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


/**
This breaks my brain.
*/
type Table struct {
	players []*Player
	gameCtr int

	// The below get preset before each game.
	community        Community
	button           Player
	smallBlindValue  int
	smallBlindPlayer *Player
	bigBlindValue    int
	bigBlindPlayer   *Player
	bettingRound     string
	deck             *Deck
	pot              Pot
}

// getStatus is more verbose than toString.
func (t *Table) getStatus() string {
	status := "------\n"
	status += fmt.Sprintf("%s -- %d players\n", t.bettingRound, len(t.players))

	for _, player := range t.players {
		status += fmt.Sprintf("%s\n", *player)
	}

	status += fmt.Sprintf("Pot: %d\n", t.pot.value)
	status += fmt.Sprintf("Community: %s\n", t.community)
	status += "------\n"
	return status
}

/*
This happens at the start of tournaments.
*/
func (t *Table) initialize() {
	t.gameCtr = 0

	// https://stackoverflow.com/questions/33994677/pick-a-random-value-from-a-go-slice
	rand.Seed(time.Now().Unix())
}

/*
This happens at the start of games.
*/
func (t *Table) preset() {
	t.gameCtr++

	t.deck = NewDeck()
	t.deck.shuffle()
	t.pot = NewPot()
	t.community = NewCommunity()

	for _, p := range t.players {
		player := *p
		player.preset()
	}
}

func (t *Table) addPlayerPointerVersion(player *Player) {
	if len(t.players) == 0 {
		t.players = append(t.players, player)
		return
	}

	playerDerefd := *player

	initialPlayer := *t.players[0]
	lastPlayerPtr := t.players[len(t.players)-1]
	lastPlayer := *lastPlayerPtr

	lastPlayer.setNextPlayer(player)
	playerDerefd.setPreviousPlayer(&lastPlayer)
	playerDerefd.setNextPlayer(t.players[0])
	initialPlayer.setPreviousPlayer(player)

	t.players = append(t.players, player)

	fmt.Printf("CURRENT PLAYER: %s > %s > %s\n",
		lastPlayer.getName(), playerDerefd.getName(), initialPlayer.getName())
	return
}

func (t *Table) addPlayer(player Player) {
	if len(t.players) == 0 {
		//fmt.Println(t)
		//fmt.Println("The table is empty so adding", player.getName())
		t.players = append(t.players, &player)
		//fmt.Println("The table now has length:", len(t.players))
		//fmt.Println(t)  // &{[0xc042054320] 0 {<nil>} <nil> 0 <nil> 0 <nil>  <nil> {0 map[]}}

		//fmt.Printf("1b -> ")
		//t.printPlayerList()

		return
	}

	//t.printPlayerList()
	initialPlayer := *t.players[0]
	//fmt.Println("initial player is:", initialPlayer.getName())
	lastPlayerPtr := t.players[len(t.players)-1]
	lastPlayer := *lastPlayerPtr
	//fmt.Println("last player is:", lastPlayer.getName())

	//var peter Player
	//zubin := NewGenericPlayer("zubin")
	//converterSlice := []*Player(&zubin)
	//converterSlice = append(converterSlice, &zubin)

	//fmt.Println("table.players is of type:", reflect.TypeOf(t.players))      // []*main.Player
	//fmt.Println("table.gameCtr is of type:", reflect.TypeOf(t.gameCtr))      // int
	//fmt.Println("table.Pot is of type:    ", reflect.TypeOf(t.pot))          // main.Pot
	//fmt.Println("lastPlayerPtr is of type:", reflect.TypeOf(lastPlayerPtr))  // *main.Player
	//fmt.Println("lastPlayer is of type:   ", reflect.TypeOf(lastPlayer))     // *main.GenericPlayer
	//fmt.Println("&lastPlayer is of type:  ", reflect.TypeOf(&lastPlayer))    // *main.Player
	//fmt.Println("zubin is of type:        ", reflect.TypeOf(zubin))          // main.GenericPlayer
	//fmt.Println("&zubin is of type:       ", reflect.TypeOf(&zubin))         // *main.GenericPlayer
	//fmt.Println("peter is of type:        ", reflect.TypeOf(peter))          // <nil>
	//fmt.Println("&peter is of type:       ", reflect.TypeOf(&peter))         // *main.Player

	lastPlayer.setNextPlayer(&player)
	//player.setPreviousPlayer(&lastPlayer)
	player.setPreviousPlayer(lastPlayerPtr)
	player.setNextPlayer(t.players[0])
	initialPlayer.setPreviousPlayer(&player)

	t.players = append(t.players, &player)

	// below error: Can't use *GenericPlayer as type *Player
	//player.setPreviousPlayer(&peter)

	//player.setName("fakeo")
	fmt.Printf("CURRENT PLAYER: %s > %s > %s\n",
		lastPlayer.getName(), player.getName(), initialPlayer.getName())
	return
}

func (t Table) printPlayerList() {
	if len(t.players) == 1 {
		lonePlayer := *t.players[0]
		fmt.Println("There is only one player at the table:", lonePlayer.getName())
		return
	}
	//fmt.Println("Here's the table: ", t)

	fmt.Println("Players: ")
	for _, p := range t.players {
		player := *p // Needed because t.players is a slice of *Player.
		np := player.getNextPlayer()
		nextPlayer := *np
		nnp := nextPlayer.getNextPlayer()
		nextNextPlayer := *nnp
		fmt.Println(player.getName(), nextPlayer.getName(), nextNextPlayer.getName())
	}
}

func (t Table) printLinkList(reverse bool, p *Player) {
	// the zero value of a bool is false

	if p == nil {
		p = t.players[0]
	}
	player := *p

	if reverse == false {
		if player.getNextPlayer() == t.players[0] {
			fmt.Println(player.getName())
			return
		}
	} else {
		if player.getPreviousPlayer() == t.players[0] {
			fmt.Println(player.getName())
			return
		}
	}

	fmt.Printf("%s -> ", player.getName())
	time.Sleep(1000)
	if reverse == false {
		t.printLinkList(reverse, player.getNextPlayer()) // error (fixed)
	} else {
		t.printLinkList(reverse, player.getPreviousPlayer())
	}
	return
}

func (t *Table) assignInitialButtonAndBlinds() {
	n := rand.Int() % len(t.players)
	t.button = *t.players[n]

	fmt.Println("Assigning the button to:", t.button.getName())
	t.smallBlindPlayer = t.button.getNextPlayer()
	smallBlindDerefd := *t.smallBlindPlayer
	fmt.Println("Assigning SB to:", smallBlindDerefd.getName())
	t.bigBlindPlayer = smallBlindDerefd.getNextPlayer()
	bigBlindPlayerDerefd := *t.bigBlindPlayer
	fmt.Println("Assigning BB to:", bigBlindPlayerDerefd.getName())
}

func (t *Table) defineBlinds(sb int) {
	t.smallBlindValue = sb
	t.bigBlindValue = sb * 2
}

func (t *Table) postBlinds() (table Table) {
	bbp := *t.bigBlindPlayer
	bbp.payBlind(t.bigBlindValue)
	fmt.Println(bbp.getName(), "just paid the blind of $", t.bigBlindValue, "and has $",
		bbp.getStack(), "left.")
	sbp := *t.smallBlindPlayer
	sbp.payBlind(t.smallBlindValue)
	fmt.Println(sbp.getName(), "just paid the blind of $", t.smallBlindValue, "and has $",
		sbp.getStack(), "left.")

	return
}

func (t *Table) dealHoleCards() {

	for _, p := range t.players {
		player := *p
		player.addHoleCard(*t.deck.getCard())
		player.addHoleCard(*t.deck.getCard())
	}
}

func (t *Table) getMaxBet() int {
	var maxBet int
	if t.bettingRound == "PREFLOP" {
		maxBet = t.bigBlindValue
	} else {
		maxBet = 0
	}

	for _, p := range t.players {
		player := *p
		if maxBet < player.getBet() {
			maxBet = player.getBet()
		}
	}

	return maxBet
}

func (t *Table) getPlayerAction(playerPtr *Player) {
	//if player.checkHasFolded() {
	player := *playerPtr
	if player.checkHasFolded() {
		fmt.Println(player.getName(), "has folded so no action.")
		return
	}

	if player.checkIsAllIn() {
		fmt.Println(player.getName(), "is all-in so no action.")
		return
	}

	fmt.Println(player.getName(), "has action - finding it.")
	player.chooseAction(t)

	return
}

/*
Return false unless all non folded players either have the same bet or
are all-in.
*/
func (t *Table) checkBetParity() bool {
	maxBet := t.getMaxBet()
	for _, p := range t.players {
		player := *p
		if player.checkHasFolded() || player.checkIsAllIn() {
			continue
		}

		if player.getBet() != maxBet {
			fmt.Println(player.getName(), "needs to take action.  Current bet is $", player.getBet(),
				"which is less than the max bet of $", maxBet)
			return false
		}
	}

	return true
}

func (t *Table) genericBet(firstBetter *Player) {
	firstBetterDerefd := *firstBetter
	fmt.Println("The first better is", firstBetterDerefd.getName())
	log15.Info("The first better is", firstBetterDerefd.getName())

	better := firstBetterDerefd
	t.getPlayerAction(&better)
	better = *better.getNextPlayer()

	// First go around the table.
	for better != firstBetterDerefd {
		fmt.Println(better.getName(), "is the better.")

		if t.checkForOnePlayer() {
			fmt.Println("We are down to one player.")
			return
		}

		t.getPlayerAction(&better)
		better = *better.getNextPlayer()

	}

	fmt.Println("After going around the table once, we have:")
	fmt.Println(t.getStatus())

	// There may be raises and re-raises so handle that.
	for {
		if t.checkForOnePlayer() {
			fmt.Println("There is only one player left with action.")
			break
		}
		time.Sleep(1000)

		if t.checkBetParity() {
			fmt.Println("Everyone had a chance to bet and everyone is all-in, has checkHasFolded or has called.")
			break
		}

		// These players have no action.
		if better.checkHasFolded() || better.checkIsAllIn() {
			fmt.Println(better.getName(), "has no action.")
			better = *better.getNextPlayer()
			continue
		}

		t.getPlayerAction(&better)
		better = *better.getNextPlayer()
	}
}

func (t *Table) preFlopBet() {
	bigBlindPlayer := *t.bigBlindPlayer
	firstBetter := bigBlindPlayer.getNextPlayer()
	t.genericBet(firstBetter)
}

func (t *Table) postPreFlopBet() {
	t.genericBet(t.smallBlindPlayer)
}

func (t *Table) countFoldedPlayers() {

}

func (t *Table) checkForOnePlayer() bool {
	remaingPlayerCount := len(t.players)
	if remaingPlayerCount > 1 {
		return false
	}

	return true
}

func (t *Table) moveBetsToPot() {
	fmt.Println("Moving bets to pot.")

	for _, p := range t.players {
		player := *p
		t.pot.addEquity(player.getBet(), p)
		player.setBet(0)
	}
}

func (t *Table) dealFlop() {
	t.bettingRound = "FLOP"
	for i := 1; i <= 3; i++ {
		card := t.deck.getCard()
		t.community.add(*card)
	}
}

func (t *Table) dealTurn() {
	t.bettingRound = "TURN"
	card := t.deck.getCard()
	t.community.add(*card)
}

func (t *Table) dealRiver() {
	t.bettingRound = "RIVER"
	card := t.deck.getCard()
	t.community.add(*card)
}

func (t *Table) payWinners() {
	fmt.Println("The pot:", t.pot)

	// To properly test this loop, we need different player types first.
	for segmentAmount, segmentPlayers := range t.pot.getSegments() {
		fmt.Printf("$%d: ", segmentAmount)
		for _, p := range segmentPlayers {
			player := *p
			fmt.Printf("%s, ", player.getName())
		}
		fmt.Println()
		fmt.Println(segmentAmount, "->", segmentPlayers)
		t.payWinnersForSegment(segmentAmount, segmentPlayers)
	}
}

func (t *Table) payWinnersForSegment(segmentValue int, players []*Player) {
	fmt.Println("Finding the winner of the", segmentValue, "dollar segment.")

	var activePlayers []Player
	for _, p := range players {
		player := *p

		// You can't win if you don't play.
		if player.checkHasFolded() {
			continue
		}
		activePlayers = append(activePlayers, player)
	}

	fmt.Println("There are", len(activePlayers), "active players in this segment.")
	for _, ap := range activePlayers {
		fmt.Println("-", ap.getName())

		// I will pay the cost of reevaluating these hands so I don't
		// have to add more methods to the Player interface.
		aphc := ap.getHoleCards()
		combinedCardset := aphc.combine(*t.community.cards)
		fmt.Println("Combined cardset:", combinedCardset)
		combinedCardset.findBestHand()
	}
}


func runTournament() {
	var table Table
	table.initialize()

	temp := NewAllInAlwaysPlayer("Adam")
	table.addPlayer(&temp)
	temp2 := NewGenericPlayer("Bert")
	table.addPlayer(&temp2)
	tempCSP := NewCallingStationPlayer("Cali")
	table.addPlayer(&tempCSP)
	temp4 := NewGenericPlayer("Dale")
	table.addPlayer(&temp4)
	temp5 := NewGenericPlayer("Eyor")
	table.addPlayer(&temp5)
	temp6 := NewCallToFivePlayer("Fred")
	table.addPlayer(&temp6)
	temp7 := NewFoldingPlayer("Greg")
	table.addPlayer(&temp7)
	temp8 := NewGenericPlayer("Hill")
	table.addPlayer(&temp8)
	temp9 := NewGenericPlayer("Igor")
	table.addPlayer(&temp9)

	table.printLinkList(false, nil)
	table.printLinkList(true, nil)
	fmt.Print("\n\n")

	// Set an initial small blind value.
	table.defineBlinds(25)

	for i := 1; i <= 2; i++ {
		fmt.Println("============================")
		table.preset()
		fmt.Printf("This is game #%d.\n", table.gameCtr)
		table.assignInitialButtonAndBlinds()
		table.bettingRound = "PREFLOP"

		if i == 2 {
			fmt.Println("Figure out how to pay winners then we'll continue.")
			fmt.Println(table.getStatus())
			os.Exit(3)
		}

		table.postBlinds()
		fmt.Println(table.getStatus())
		table.dealHoleCards()
		table.preFlopBet()
		table.moveBetsToPot()
		fmt.Println(table.getStatus())

		table.bettingRound = "FLOP"
		fmt.Println("Dealing the flop.")
		table.dealFlop()
		table.postPreFlopBet()
		fmt.Println(table.getStatus())

		table.bettingRound = "TURN"
		fmt.Println("Dealing the turn.")
		table.dealTurn()
		table.postPreFlopBet()
		fmt.Println(table.getStatus())

		table.bettingRound = "RIVER"
		fmt.Println("Dealing the river.")
		table.dealRiver()
		table.postPreFlopBet()
		fmt.Println(table.getStatus())

		fmt.Println("Finding and paying the winners.")
		table.payWinners()
	}

}

/*
We will run multiple tournaments to find the best player type.
Each tournament has multiple poker _games_.
Each game has multiple betting _rounds_.
Each winner of each game has the best _hand_ - multiple winners are possible.
*/

func main() {
	//card := NewCard("H", 12)
	//cardSet := NewCardSet(card)
	//fmt.Println(cardSet)
	//card2 := NewCard("S", 4)
	//cardSet.add(card2)
	//fmt.Println(cardSet)
	//fmt.Println("cardset1 =", cardSet)
	//
	////ecs := getEmptyCardSet()
	//ecs := NewCardSet()
	//fmt.Println("ECS:", ecs)
	//
	//cs2 := NewCardSet()
	//cs2.add(NewCard("C", 13))
	//cs2.add(NewCard("C", 12))
	//fmt.Println("cardset2 =", cs2)
	//
	//cs3 := cs2.combine(cardSet)
	//fmt.Println("cardset3 =", cs3)
	//
	//cs4 := NewCardSet(card, card2)
	//fmt.Println("cardset4 =", cs4)
	//
	//cs5 := NewCardSet(card, card2, NewCard("C", 4))
	//fmt.Println("using type:", reflect.TypeOf(card2))
	//fmt.Println("cardset5 =", cs5)
	//for i, c := range cs5.cards {
	//	fmt.Println(i, c)
	//}
	//
	//os.Exit(3)

	//
	////holeCards := new(HoleCards)
	//holeCards := HoleCards{cardset: &ecs}
	//fmt.Println("Holecards:", holeCards)
	//holeCards.add(card)
	//fmt.Println("Holecards:", holeCards)
	//holeCards.add(card2)
	//fmt.Println("Holecards:", holeCards)
	//holeCards.add(card)
	//fmt.Println("Holecards:", holeCards)
	//
	//os.Exit(3)

	//for numericRank := range [13]int{} {
	//	newCard := NewCard("D", numericRank+2)
	//	//fmt.Println("New card:", newCard)
	//	cardSet.add(newCard)
	//}
	//fmt.Println(cardSet.getStatus())
	//cardSet.shuffle()
	//fmt.Println(cardSet.getStatus())
	/*
		H12 S4 D2 D3 D4 D5 D6 D7 D8 D9 D10 D11 D12 D13 D14
		D11 D10 D14 D3 D2 D7 D12 S4 H12 D6 D4 D5 D8 D13 D9
	*/

	for i := 1; i <= 2; i++ {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		fmt.Printf("Starting tournament #%d\n", i)
		runTournament()
	}

}

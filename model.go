package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"time"
	"gopkg.in/inconshreveable/log15.v2"
)

type Card struct {
	suit         string
	numericaRank int
	rank		 string
}

func NewCard(s string, nr int) Card {
	var c Card
	c.suit = s
	c.numericaRank = nr

	if nr == 14 {
		c.rank = "A"  // Aces are aces.
	}

	switch nr {
	case 14: c.rank = "A"  // Aces are aces.
	case 13: c.rank = "K"
	case 12: c.rank = "Q"
	case 11: c.rank = "J"
	case 10: c.rank = "T"
	default: c.rank = strconv.Itoa(nr)
	}

	return c
}

// Should I rename all these getStatus() methods to toString()?
func (c *Card) getStatus() string {
	return fmt.Sprintf("%s%s", c.suit, c.rank)
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
}

func NewCardSet(cards ...Card) CardSet {
	var cs CardSet

	for _, card := range cards {
		cs.cards = append(cs.cards, &card)
	}

	return cs
}

func (cs CardSet) String() string {
	var toString string

	for _, c := range cs.cards {
		if toString == "" {
			toString = c.getStatus()
			continue
		}

		toString += fmt.Sprintf(" %s", c.getStatus())
	}

	return toString
}

func (cs *CardSet) add(c Card) {
	cs.cards = append(cs.cards, &c)
}

/*
This will accept a cardset and combine it with this.
 */
func (cs *CardSet) combine(cs2 CardSet) {
	for _, card := range cs2.cards {
		cs.add(*card)
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

// This should be a CardSet method.
func getEmptyCardSet() CardSet {
	empty := CardSet{}

	return empty
}

/*
Not sure this should be a struct.

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
 */
type Evaluation struct {
	cardSet *CardSet
	humanEval string
	allRanks []string
}

func NewEvaluation() *Evaluation {
	var eval Evaluation
	ecs := getEmptyCardSet()
	eval.cardSet = &ecs

	return &eval
}

func (e Evaluation) String() string {
	toString := fmt.Sprint("%s: %s with ranks: [TBD]", e.cardSet, e.humanEval)



	return toString
}

type HoleCards struct {
	cardset *CardSet
}

func (hc *HoleCards) toString() string {
	return hc.cardset.String()
}

func (hc *HoleCards) add(c Card) {
	hc.cardset.add(c)
}

// Maybe call this hc.reset() and it can be used in more than one place.
func (hc *HoleCards) toss() {
	ecs := getEmptyCardSet()
	hc.cardset = &ecs
}

func (hc HoleCards) String() string {
	return hc.cardset.String()
}

type Deck struct {
	cardset *CardSet  // Deck is not a CardSet.
}

func NewDeck() *Deck {
	var d Deck
	ecs := getEmptyCardSet()
	d.cardset = &ecs

	for _, suit := range []string{"S", "H", "D", "C"} {
		// https://stackoverflow.com/questions/21950244/is-there-a-way-to-iterate-over-a-range-of-integers-in-golang
		for numericRank := range [13]int{} {
			newCard := NewCard(suit, numericRank+2)
			//fmt.Println("New card:", newCard)
			d.cardset.add(newCard)
		}
	}

	fmt.Println("Made a new deck.")

	return &d
}

func (d *Deck) getStatus() string {
	status := fmt.Sprintf("Deck has %d cards:\n", len(d.cardset.cards))
	for i, card := range d.cardset.cards {
		if i != 0 && i % 13 == 0 {
			status += "\n"
		}
		status += card.getStatus()
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
	ecs := getEmptyCardSet()

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
	value int  // This could be gotten from summing equity.
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

	for player, value := range p.equity {
		toString += fmt.Sprintf("%s has equity: $%d\n", player.name, value)
	}

	return toString
}

func (p *Pot) addEquity(playerBet int, player *Player) {
	p.value += playerBet
	p.equity[player] += playerBet
}

func (p *Pot) getSegments() map[int][]*Player {
	//fmt.Println("segments wooo hoo")

	// First invert the map.
	invertedMap := make(map[int][]*Player)
	for player, equity := range p.equity {
		invertedMap[equity] = append(invertedMap[equity], player)
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
		for _, player := range invertedMap[equity] {
			fmt.Printf("    %s\n", player.name)
			segments[equity] = append(segments[equity], player)
		}
		fmt.Println()
	}

	//
	//// Get all the values of the equity map.
	//values := make([]int, len(p.equity))
	//for _, player := range p.equity {
	//
	//	// We only care about unique values.
	//	// https://stackoverflow.com/questions/9251234/go-append-if-unique
	//	values = append(values, player)
	//}
	//fmt.Println(values)

	//previous := 0

	//fmt.Println("Returning from getSegments().")
	return segments
}

type PlayerInterface interface {
	fold()
	allIn()
}

type Player struct {
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

// Player constructor
// http://www.golangpatterns.info/object-oriented/constructors
func NewPlayer(name string) Player {
	ecs := getEmptyCardSet()
	hc := HoleCards{cardset: &ecs}
	initialStack := 1000 // dollars
	newPlayer := Player{name, nil, nil, hc, initialStack, 0,
		false, false}
	return newPlayer
}

func (p *Player) getStatus() string {
	status := ""
	status = fmt.Sprintf("%s: [%s] stack=%d bet=%d", p.name, p.holeCards.toString(), p.stack, p.bet)

	return status
}

func (p Player) String() string {
	return fmt.Sprintf("%s: [%s] $%d/$%d", p.name, p.holeCards.toString(), p.bet, p.stack)
}

/*
Preset before each game.

Maybe use "prepare()" instead of "preset()" because the latter implies
something you do afterwards.
 */
func (p *Player) preset() {
	// Maybe NewPlayer can call this?
	ecs := getEmptyCardSet()
	p.holeCards = HoleCards{cardset: &ecs}
	p.bet = 0
	p.hasFolded = false
	p.isAllIn = false
}

func (p *Player) payBlind(blindAmount int) {
	if p.stack > blindAmount {
		p.stack -= blindAmount
		p.bet = blindAmount
	} else {
		p.allIn()
	}
}

func (p *Player) check(t *Table) {
	// This may one day do something.

	return
}

func (p *Player) call(t *Table) {
	maxBet := t.getMaxBet()
	increase := maxBet - p.bet
	stackBefore := p.stack
	betBefore := p.bet

	if increase >= p.stack {
		p.allIn()
	} else {
		p.bet = maxBet
		p.stack -= increase
	}

	fmt.Println("calling - max bet:", maxBet, "increase:", increase, "bet before/after: ", betBefore, "/",
		p.bet, ", stack before/after:", stackBefore, "/", p.stack)
}

func (p *Player) fold() {
	p.hasFolded = true
	p.holeCards.toss()
}

func (p *Player) allIn() {
	p.bet += p.stack
	p.stack = 0
	p.isAllIn = true
}

func (p *Player) checkOrCall(t *Table) {
	if t.getMaxBet() == 0 {
		p.check(t)
		return
	}

	p.call(t)

	return
}

func (p *Player) checkOrFold(t *Table) {
	return
}

func (p *Player) chooseActionPreflop(t *Table) {
	fmt.Println("Using the default preflop action of check-calling")
	p.checkOrCall(t)

	return
}

func (p *Player) chooseActionFlop(t *Table) {
	fmt.Println("Using the default flop action of check-folding")
	p.checkOrFold(t)

	return
}

func (p *Player) chooseActionTurn(t *Table) {
	fmt.Println("Using the default turn action of check-folding")
	p.checkOrFold(t)

	return
}

func (p *Player) chooseActionRiver(t *Table) {
	fmt.Println("Using the default river action of check-folding")
	p.checkOrFold(t)

	return
}

func (p *Player) chooseAction(t *Table) {

	// This is handled by Table() but be redundant for clearness.
	if p.hasFolded {
		fmt.Println("I,", p.name, "have already folded so no action.")
		return
	}
	if p.isAllIn {
		fmt.Println("I,", p.name, "am all in so no action.")
		return
	}

	fmt.Println(p.getStatus())

	switch t.bettingRound {
	case "PREFLOP": p.chooseActionPreflop(t)
	case "FLOP": p.chooseActionFlop(t)
	case "TURN": p.chooseActionTurn(t)
	case "RIVER": p.chooseActionRiver(t)
	}
	return
}

type CallingStationPlayer struct {
	Player
}

// Repeating the constructor is kinda lame.
func NewCallingStationPlayer(name string) CallingStationPlayer {
	ecs := getEmptyCardSet()
	hc := HoleCards{cardset: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(CallingStationPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack

	return *newPlayer
}


/**
This breaks my brain.
 */
type Table struct {
	players          []*Player
	gameCtr          int

	// The below get preset before each game.
	community        Community
	button           Player
	smallBlindValue  int
	smallBlindPlayer *Player
	bigBlindValue    int
	bigBlindPlayer   *Player
	bettingRound     string
	deck			 *Deck
	pot				 Pot
}

func (t *Table) getStatus() string {
	status := "------\n"
	status += fmt.Sprintf("%d players\n", len(t.players))

	for _, player := range t.players {
		status += fmt.Sprintf("%s\n", player)
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

	for _, p := range t.players{
		p.preset()
	}

}

//   receiver   name      inputs         return type
//func (t *Table) addPlayer(player PlayerInterface) {  # TODO this needs a lot of work.
func (t *Table) addPlayer(player Player) {
	if len(t.players) == 0 {
		t.players = append(t.players, &player)
		return
	}

	lastPlayer := t.players[len(t.players)-1]

	lastPlayer.nextPlayer = &player
	player.previousPlayer = lastPlayer
	player.nextPlayer = t.players[0]
	t.players[0].previousPlayer = &player

	t.players = append(t.players, &player)

	fmt.Printf("CURRENT PLAYER: %s > %s > %s\n",
		player.previousPlayer.name, player.name, player.nextPlayer.name)
	return
}

func (t Table) printPlayerList() {
	//fmt.Println("Here's the table: ", t)

	fmt.Println("Players: ")
	for _, p := range t.players {
		fmt.Println(p.name, p.nextPlayer.name, p.nextPlayer.nextPlayer.name)
	}
}

func (t Table) printLinkList(reverse bool, p *Player) {
	// the zero value of a bool is false

	if p == nil {
		p = t.players[0]
	}

	if reverse == false {
		if p.nextPlayer == t.players[0] {
			fmt.Println(p.name)
			return
		}
	} else {
		if p.previousPlayer == t.players[0] {
			fmt.Println(p.name)
			return
		}
	}
	//if p.nextPlayer == t.players[0] {
	//	fmt.Println(p.name)
	//	return
	//}

	fmt.Printf("%s -> ", p.name)

	if reverse == false {
		t.printLinkList(reverse, p.nextPlayer)
	} else {
		t.printLinkList(reverse, p.previousPlayer)
	}
	return
}

func (t *Table) assignInitialButtonAndBlinds() {
	n := rand.Int() % len(t.players)
	t.button = *t.players[n]
	fmt.Println("Assigning the button to:", t.button.name)

	t.smallBlindPlayer = t.button.nextPlayer
	fmt.Println("Assigning SB to:", t.smallBlindPlayer.name)
	t.bigBlindPlayer = t.smallBlindPlayer.nextPlayer
	fmt.Println("Assigning BB to:", t.bigBlindPlayer.name)
}

func (t *Table) defineBlinds(sb int) {
	t.smallBlindValue = sb
	t.bigBlindValue = sb * 2
}

func (t *Table) postBlinds() (table Table) {
	t.bigBlindPlayer.payBlind(t.bigBlindValue)
	fmt.Println(t.bigBlindPlayer.name, "just paid the blind of $", t.bigBlindPlayer.bet, "and has $",
		t.bigBlindPlayer.stack, "left.")

	t.smallBlindPlayer.payBlind(t.smallBlindValue)
	fmt.Println(t.smallBlindPlayer.name, "just paid the blind of $", t.smallBlindPlayer.bet, "and has $",
		t.smallBlindPlayer.stack, "left.")

	return
}

func (t *Table) dealHoleCards() {
	for _, player := range t.players {
		player.holeCards.add(*t.deck.getCard())  // You need two hole cards.
		player.holeCards.add(*t.deck.getCard())
	}
}

func (t *Table) getMaxBet() int {
	maxBet := 0

	for _, player := range t.players {
		if player.bet > maxBet {
			maxBet = player.bet
		}
	}

	return maxBet
}

func (t *Table) getPlayerAction(player *Player) {
	if player.hasFolded {
		fmt.Println(player.name, "has folded so no action.")
		return
	}

	if player.isAllIn {
		fmt.Println(player.name, "is all-in so no action.")
		return
	}

	fmt.Println(player.name, "has action - finding it.")
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
		if p.hasFolded || p.isAllIn {
			continue
		}

		if p.bet != maxBet {
			fmt.Println(p.name, "needs to take action.  Current bet is $", p.bet, "which is less than the max " +
				"bet of $", maxBet)
			return false
		}
	}

	return true
}

func (t *Table) genericBet(firstBetter *Player) {
	//firstBetter := t.bigBlindPlayer.nextPlayer
	fmt.Println("The first better is", firstBetter.name)
	log15.Info("The first better is", firstBetter.name)

	better := firstBetter
	t.getPlayerAction(better)
	better = better.nextPlayer

	// First go around the table.
	for better != firstBetter {
		fmt.Println(better.name, "is the better.")

		if t.checkForOnePlayer() {
			fmt.Println("We are down to one player.")
			return
		}

		t.getPlayerAction(better)
		better = better.nextPlayer
	}

	fmt.Println("After going around the table once, we have:")
	fmt.Println(t.getStatus())

	// There may be raises and re-raises so handle that.
	for {
		if t.checkForOnePlayer() {
			return
		}

		// These players have no action.
		if better.hasFolded || better.isAllIn {
			continue
		}

		if t.checkBetParity() {
			fmt.Println("Everyone had a chance to bet and everyone is all-in, has folded or has called.")
			break
		}

		t.getPlayerAction(better)
		better = better.nextPlayer
	}
}

func (t *Table) preFlopBet() {
	firstBetter := t.bigBlindPlayer.nextPlayer
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
		fmt.Println(p.name, "had bet $", p.bet)
		t.pot.addEquity(p.bet, p)
		p.bet = 0
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
		for _, player := range segmentPlayers {
			fmt.Printf("%s, ", player.name)
		}
		fmt.Println()
		//fmt.Println(segmentAmount, "->", segmentPlayers)
	}

}

func runTournament() {
	var table Table
	table.initialize()

	table.addPlayer(NewPlayer("Adam"))
	table.addPlayer(NewPlayer("Bert"))
	table.addPlayer(NewCallingStationPlayer("Cail"))
	table.addPlayer(NewPlayer("Dale"))
	table.addPlayer(NewPlayer("Eyor"))
	table.printPlayerList()
	table.printLinkList(false, nil)
	table.printLinkList(true, nil)
	fmt.Print("\n\n")

	// Set an initial small blind value.
	//table.assignInitialButtonAndBlinds()
	table.defineBlinds(25)

	for i := 1; i <= 2; i++ {
		fmt.Println("============================")
		table.assignInitialButtonAndBlinds()
		table.preset()
		fmt.Printf("This is game #%d.\n", table.gameCtr)
		table.bettingRound = "PREFLOP"
		table.postBlinds()
		fmt.Print(table.getStatus())
		table.dealHoleCards()
		table.preFlopBet()
		table.moveBetsToPot()
		fmt.Println(table.getStatus())

		fmt.Println("Dealing the flop.")
		table.dealFlop()
		table.postPreFlopBet()
		fmt.Println(table.getStatus())

		fmt.Println("Dealing the turn.")
		table.dealTurn()
		table.postPreFlopBet()
		fmt.Println(table.getStatus())

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
	//fmt.Println(cardSet.getStatus())
	//card2 := NewCard("S", 4)
	//cardSet.add(card2)
	//fmt.Println(cardSet.getStatus())
	//fmt.Println(cardSet)



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

	for i := 1; i <= 1; i++ {
		runTournament()
	}



}

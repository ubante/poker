package main

import (
	"fmt"
	"math/rand"
	"time"
	"strconv"
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

func (cs *CardSet) getStatus() string {
	var status string

	for _, c := range cs.cards {
		if status == "" {
			status = c.getStatus()
			continue
		}

		status += fmt.Sprintf(" %s", c.getStatus())
	}

	return status
}

func (cs *CardSet) add(c Card) {
	cs.cards = append(cs.cards, &c)
}

func (cs *CardSet) shuffle() {
	rand.Shuffle(cs.length(), func(i, j int) {
		cs.cards[i], cs.cards[j] = cs.cards[j], cs.cards[i]
	})
}

func (cs CardSet) length() int {
	return len(cs.cards)
}

func getEmptyCardSet() CardSet {
	empty := CardSet{}

	return empty
}

type HoleCards struct {
	cardset CardSet
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
	status := ""
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

type Community struct {
	cards []Card
}

type Player struct {
	name           string
	nextPlayer     *Player
	previousPlayer *Player

	// The below get reset after each game.
	holeCards	   HoleCards
	stack		   int
	bet            int
	hasFolded      bool
	isAllIn        bool
}

func (p *Player) getStatus() string {
	status := ""
	status = fmt.Sprintf("%s: (%s/%s) stack=%d bet=%d", p.name, p.previousPlayer.name, p.nextPlayer.name,
		p.stack, p.bet)

	return status
}

func (p *Player) payBlind(blindAmount int) {
	if p.stack > blindAmount {
		p.stack -= blindAmount
		p.bet = blindAmount
	} else {
		p.allIn()
	}
}

func (p *Player) allIn() {
	p.bet += p.stack
	p.stack = 0
	p.isAllIn = true
}

// Player constructor
// http://www.golangpatterns.info/object-oriented/constructors
func NewPlayer(name string) Player {
	ecs := getEmptyCardSet()
	hc := HoleCards{cardset: ecs}
	initialStack := 1000  // dollars
	newPlayer := Player{name, nil, nil, hc, initialStack, 0,
	false, false}
	return newPlayer
}

type Pot struct {
	value int
}

/**
This breaks my brain.
 */
type Table struct {
	players          []*Player
	gameCtr          int

	// The below get reset after each game.
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
	status := "------"
	status = fmt.Sprintf("%d players\n", len(t.players))

	for _, player := range t.players {
		status += fmt.Sprintf("%s: $%d ($%d)\n", player.name, player.stack, player.bet)
	}
	return status
}

/*
This happens at the start of tournaments.
 */
func (t *Table) initialize() {
	t.gameCtr = 0
}

/*
This happens at the start of games.
 */
func (t *Table) reset() {
	t.gameCtr++

	t.deck = NewDeck()
	t.deck.shuffle()
}

//   receiver   name      inputs         return type
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
	// https://stackoverflow.com/questions/33994677/pick-a-random-value-from-a-go-slice
	rand.Seed(time.Now().Unix())
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
	fmt.Println(t.bigBlindPlayer.name, "just bet the blind of $", t.bigBlindPlayer.bet, "and has $",
		t.bigBlindPlayer.stack, "left.")

	t.smallBlindPlayer.payBlind(t.smallBlindValue)
	fmt.Println(t.smallBlindPlayer.name, "just bet the blind of $", t.smallBlindPlayer.bet, "and has $",
		t.smallBlindPlayer.stack, "left.")

	return
}

func (t *Table) dealHoleCards() {

}

func runTournament() {
	var table Table
	table.initialize()

	table.addPlayer(NewPlayer("Adam"))
	table.addPlayer(NewPlayer("Bert"))
	table.addPlayer(NewPlayer("Cail"))
	table.addPlayer(NewPlayer("Dale"))
	table.addPlayer(NewPlayer("Eyor")) // way #3
	table.printPlayerList()
	table.printLinkList(false, nil)
	table.printLinkList(true, nil)
	fmt.Print("\n\n")

	// Set an initial small blind value.
	table.assignInitialButtonAndBlinds()
	table.defineBlinds(25)

	for i := 1; i <= 2; i++ {
		table.reset()

		fmt.Printf("This is game #%d.\n", table.gameCtr)
		table.bettingRound = "PREFLOP"
		table.postBlinds()
		fmt.Print(table.getStatus())
		table.dealHoleCards()

		table.bettingRound = "FLOP"

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

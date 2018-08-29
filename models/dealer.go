package models

import (
	"fmt"
	"math/rand"
	"strconv"

	"gopkg.in/inconshreveable/log15.v2"
	"os"
	"sort"
	"time"
	//"goven/poker/models"
)

/*
Trying to stick to the definitions in https://www.pokernews.com/pokerterms/reraise.htm
and the rules in http://www.pokercoach.us/RobsPkrRules11.doc

In addition to those definitions:
  - hand: the five cards a player constructs
  - round: a betting stage; there are four per game: preflop/flop/turn/river
  - game: what happens between the shuffling of the deck
  - tournament: from the first game until there's only one player left
*/

func NewCard(s string, nr int) Card {
	var c Card
	c.Suit = s
	c.NumericalRank = nr

	if nr == 14 {
		c.Rank = "A" // Aces are aces.
	}

	switch nr {
	case 14:
		c.Rank = "A" // Aces are aces.
	case 13:
		c.Rank = "K"
	case 12:
		c.Rank = "Q"
	case 11:
		c.Rank = "J"
	case 10:
		c.Rank = "T"
	default:
		c.Rank = strconv.Itoa(nr)
	}

	return c
}

func NewCardSet(cards ...Card) CardSet {
	var cs CardSet

	for _, card := range cards {
		cs.Add(card)
	}

	return cs
}

func NewEvaluation(cardSet CardSet) *Evaluation {
	var eval Evaluation
	eval.cardSet = &cardSet
	eval.evaluate()

	// How to do the below in one line?
	eval.allRanks[0] = &eval.primaryRank
	eval.allRanks[1] = &eval.secondaryRank
	eval.allRanks[2] = &eval.tertiaryRank
	eval.allRanks[3] = &eval.quaternaryRank
	eval.allRanks[4] = &eval.quinaryRank
	eval.allRanks[5] = &eval.senaryRank
	eval.flattenScore()

	return &eval
}

func NewCommunity() Community {
	ecs := NewCardSet()

	var c Community
	c.cards = &ecs

	return c
}

func NewDeck() *Deck {
	var d Deck
	ecs := NewCardSet()
	d.cardSet = &ecs

	for _, suit := range []string{"S", "H", "D", "C"} {
		// https://stackoverflow.com/questions/21950244/is-there-a-way-to-iterate-over-a-range-of-integers-in-golang
		for numericRank := range [13]int{} {
			newCard := NewCard(suit, numericRank+2)
			d.cardSet.Add(newCard)
		}
	}

	return &d
}

type SubPot struct {
	amount              int
	contributingPlayers []*Player
}

func (sp SubPot) String() string {
	toString := fmt.Sprintf("$%d: ", sp.amount)
	for _, cp := range sp.contributingPlayers {
		cpp := *cp
		toString += fmt.Sprintf("%s, ", cpp.getName())
	}
	return toString
}

func (sp SubPot) contains(player *Player) bool {
	// This could have been a map but then I'd have to convert the slice
	// in NewPot.
	for _, cp := range sp.contributingPlayers {
		if cp == player {
			return true
		}
	}
	return false
}

func (sp *SubPot) deposit(player *Player, amount int) {
	pp := *player
	if !sp.contains(player) {
		sp.contributingPlayers = append(sp.contributingPlayers, player)
	}
	sp.amount += amount
	pp.setBet(pp.getBet() - amount)
}

type Pot struct {
	subPots     []*SubPot // The first SubPot is the main pot.
	subPotIndex int
}

func (pot Pot) String() string {
	toString := fmt.Sprintf("POT total is $%d", pot.getValue())
	for i, sp := range pot.subPots {
		toString += fmt.Sprintf("\npot #%d, %s", i, *sp)
	}
	return toString
}

func NewPot(players []*Player) Pot {
	var pot Pot

	aSubPot := SubPot{contributingPlayers: players, amount: 0}
	pot.subPots = append(pot.subPots, &aSubPot)
	pot.subPotIndex = 0

	return pot
}

func (pot *Pot) getValue() int {
	value := 0
	for _, sp := range pot.subPots {
		value += sp.amount
	}

	return value
}

func (pot *Pot) recordRoundBets(players []*Player) {
	// First find all-in players with bets.
	allInBetAmountMap := make(map[int]int) // To record unique amounts.
	for _, p := range players {
		pp := *p
		if pp.checkIsAllIn() && pp.getBet() > 0 {
			allInBetAmountMap[pp.getBet()]++
		}
	}

	// Loop through them in ascending order by Player.bet.
	var allInBetAmounts []int
	for amount := range allInBetAmountMap {
		allInBetAmounts = append(allInBetAmounts, amount)
	}
	sort.Sort(sort.IntSlice(allInBetAmounts))

	previousAllInBetAmount := 0
	for _, bet := range allInBetAmounts {
		// Loop through all the players with bets.
		for _, p := range players {
			pp := *p
			if pp.getBet() == 0 {
				continue
			}

			// Deposit their bets (up to the current all-in player's
			// bet) into the current subPot.
			margin := bet - previousAllInBetAmount
			if pp.getBet() > margin {
				pot.subPots[pot.subPotIndex].deposit(p, margin)
			} else {
				pot.subPots[pot.subPotIndex].deposit(p, pp.getBet())
			}
		}

		// Create a new subPot and move the index.
		fmt.Printf("Pot #%d is closed at $%d\n", pot.subPotIndex, bet)
		pot.subPots = append(pot.subPots, &SubPot{contributingPlayers: nil, amount: 0})
		pot.subPotIndex++
		previousAllInBetAmount = bet
	}

	// Catch the rest of the bets by looping over all players with bets.
	for _, p := range players {
		pp := *p
		if pp.getBet() == 0 {
			continue
		}

		// Deposit their bets in its entirety into the current subPot.
		pot.subPots[pot.subPotIndex].deposit(p, pp.getBet())
	}
}

func (pot *Pot) getWinnings(playerScores map[*Player]int) map[*Player]int {
	winnings := make(map[*Player]int)

	// Loop through all the subPots.
	for i, sp := range pot.subPots {
		// It is possible to have an empty subPot.  eg, some goes all
		// in, someone calls and no one raises until the end of the
		// game.
		if sp.amount == 0 {
			continue
		}

		// Find the contributingPlayer(s) with the highest score.
		highestScore := 0
		var subPotWinners []*Player

		for _, cp := range sp.contributingPlayers {
			if playerScores[cp] > highestScore {
				highestScore = playerScores[cp]
				subPotWinners = nil
				subPotWinners = append(subPotWinners, cp)
			} else if playerScores[cp] == highestScore {
				subPotWinners = append(subPotWinners, cp)
			}
		}
		fmt.Printf("We have %d winners of a subPot worth $%d (each $%d):\n", len(subPotWinners), sp.amount,
			sp.amount/len(subPotWinners))
		fmt.Printf("SubPot #%d was won by ", i)
		for _, spw := range subPotWinners {
			spwp := *spw
			fmt.Printf("%s, ", spwp.getName())

			// Tally their winnings.
			if _, ok := winnings[spw]; ok {
				winnings[spw] += sp.amount / len(subPotWinners)
			} else {
				winnings[spw] = sp.amount / len(subPotWinners)
			}
		}
		fmt.Println()
	}

	return winnings
}

/**
This breaks my brain.
*/
type Table struct {
	players []*Player
	gameCtr int
	bustLog map[int][]*Player

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

// getStatus is more verbose than ToString.
func (t *Table) getStatus() string {
	status := "------\n"
	status += fmt.Sprintf("%s -- %d players left (game #%d)\n", t.bettingRound, len(t.players), t.gameCtr)

	betTotal := 0
	stackTotal := 0
	for _, p := range t.players {
		pp := *p
		status += fmt.Sprintf("%s\n", *p)
		betTotal += pp.getBet()
		stackTotal += pp.getStack()
	}

	status += fmt.Sprintf("Pot: %d\n", t.pot.getValue())
	status += fmt.Sprintf("Community: %s\n", t.community)
	status += fmt.Sprintf("Bet totals: %d\n", betTotal)
	status += fmt.Sprintf("Stack totals: %d\n", stackTotal)
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
	t.bustLog = make(map[int][]*Player)
}

/*
This happens at the start of games.
*/
func (t *Table) preset() {
	t.gameCtr++

	t.deck = NewDeck()
	t.deck.shuffle()
	t.pot = NewPot(t.players)
	t.community = NewCommunity()

	for _, p := range t.players {
		player := *p
		player.preset()
	}
}

func (t *Table) addPlayer(player Player) {
	if len(t.players) == 0 {
		t.players = append(t.players, &player)
		return
	}

	initialPlayer := *t.players[0]
	lastPlayerPtr := t.players[len(t.players)-1]
	lastPlayer := *lastPlayerPtr
	lastPlayer.setNextPlayer(&player)
	player.setPreviousPlayer(lastPlayerPtr)
	player.setNextPlayer(t.players[0])
	initialPlayer.setPreviousPlayer(&player)
	t.players = append(t.players, &player)

	return
}

func (t Table) printPlayerList() {
	if len(t.players) == 1 {
		lonePlayer := *t.players[0]
		fmt.Println("There is only one player at the table:", lonePlayer.getName())
		return
	}

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

	if reverse == false {
		fmt.Printf("%s -> ", player.getName())
	} else {
		fmt.Printf("%s <- ", player.getName())
	}

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

	if len(t.players) == 2 {
		fmt.Println("Since we're head to head, the button is the small blind.")
		t.smallBlindPlayer = &t.button
		smallBlindDerefd := *t.smallBlindPlayer
		fmt.Println("So assigning SB to:", smallBlindDerefd.getName())
		t.bigBlindPlayer = t.button.getNextPlayer()
		bigBlindPlayerDerefd := *t.bigBlindPlayer
		fmt.Println("And assigning BB to:", bigBlindPlayerDerefd.getName())
	} else {
		t.smallBlindPlayer = t.button.getNextPlayer()
		smallBlindDerefd := *t.smallBlindPlayer
		fmt.Println("Assigning SB to:", smallBlindDerefd.getName())
		t.bigBlindPlayer = smallBlindDerefd.getNextPlayer()
		bigBlindPlayerDerefd := *t.bigBlindPlayer
		fmt.Println("Assigning BB to:", bigBlindPlayerDerefd.getName())
	}
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
	remainingPlayers := len(t.players)
	if remainingPlayers > 1 {
		return false
	}

	return true
}

func (t *Table) moveBetsToPot() {
	fmt.Println("Moving bets to pot.")
	t.pot.recordRoundBets(t.players)
	fmt.Println(t.pot)
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

	// Find all the players still in it and find their hand strength.
	fmt.Println("Finding the still active players.")
	var activePlayers []*Player
	playerScores := make(map[*Player]int)
	for _, p := range t.players {
		player := *p
		if player.checkHasFolded() {
			continue
		}
		activePlayers = append(activePlayers, p)

		// Evaluate the players' hand strengths.
		fmt.Println("Evaluating the hand of", player.getName())
		hc := player.getHoleCards()
		combinedCardSet := hc.Combine(*t.community.cards) // 7 cards.
		combinedCardSet.FindBestHand()
		fmt.Printf("%s's best hand is: %s\n", player.getName(), combinedCardSet.bestEval)
		playerScores[p] = combinedCardSet.bestEval.flattenedScore
	}

	// Send all the hand strengths to the pot to find payouts.
	payments := t.pot.getWinnings(playerScores)

	// Pay the players.
	fmt.Println("Paying the winners finally.")
	for p := range payments {
		pp := *p
		fmt.Println(pp.getName(), "gets", payments[p])
		pp.addToStack(payments[p])
	}
}

func (t *Table) removePlayer(player *Player) {
	for i, p := range t.players {
		pp := *p
		if p == player {
			fmt.Println("Removing", pp.getName(), "from position", i, "on the table.")
			t.players = append(t.players[:i], t.players[i+1:]...)
			return
		}
	}

	playerDerefd := *player
	fmt.Println("Could not find", playerDerefd.getName(), "so exiting program.")
	os.Exit(11)
}

func (t *Table) removeBustedPlayers() {
	var bustedPlayers []*Player // We need this temp slice so not to affect the below range.
	for _, p := range t.players {
		pp := *p
		if pp.getStack() > 0 {
			continue
		}
		bustedPlayers = append(bustedPlayers, p)
	}

	if len(bustedPlayers) == 0 {
		return
	}

	for _, bp := range bustedPlayers {
		bpp := *bp

		// Remove the busted player.
		fmt.Println(bpp.getName(), "has busted out so removing from the table.")
		fmt.Printf("From: ")
		t.printLinkList(false, nil)
		fmt.Printf("      ")
		t.printLinkList(true, nil)
		previousPlayer := bpp.getPreviousPlayer()
		ppp := *previousPlayer
		nextPlayer := bpp.getNextPlayer()
		npp := *nextPlayer

		ppp.setNextPlayer(nextPlayer)
		npp.setPreviousPlayer(previousPlayer)

		t.removePlayer(bp)
		fmt.Printf("To: ")
		t.printLinkList(false, nil)
		fmt.Printf("    ")
		t.printLinkList(true, nil)
	}

	// Record this busting.
	t.bustLog[t.gameCtr] = bustedPlayers
}

// Is this still needed?
func (t *Table) payWinnersForSegment(segmentValue int, players []*Player) {
	// It is possible that a segment has no valid players because they
	// have all folded.  For example, if the small blind folds, he may
	// be the only player in the segment equal to t.smallBlindValue.
	// Another example is a single person folds after the flop.  The
	// preflop segment would only have one person and he folded so would
	// be invalid.

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
	var segmentWinningPlayers []Player // Ties happen.
	var segmentWinningEvaluation Evaluation
	for _, ap := range activePlayers {
		fmt.Println("-", ap.getName())

		// I will pay the cost of reevaluating these hands so I don't
		// have to Add more methods to the Player interface.
		aphc := ap.getHoleCards()
		combinedCardset := aphc.Combine(*t.community.cards)
		combinedCardset.FindBestHand()
		fmt.Printf("%s's best hand is: %s\n", ap.getName(), combinedCardset.bestEval)
		thisEval := *combinedCardset.bestEval

		if len(segmentWinningPlayers) == 0 {
			segmentWinningPlayers = []Player{ap}
			segmentWinningEvaluation = thisEval
			fmt.Println("Initialing the best player with", ap.getName())
			continue
		}
		switch segmentWinningEvaluation.Compare(thisEval) {
		case -1:
			thisPlayer := ap
			segmentWinningPlayers = []Player{thisPlayer}
			segmentWinningEvaluation = thisEval
			fmt.Println("YYYYY: We have a new best player:", ap.getName())
		case 0:
			segmentWinningPlayers = append(segmentWinningPlayers, ap)
			fmt.Println("OOOOO:", ap.getName(), "has tied the best hand.")
		default:
			fmt.Println("NNNNNN: We do not have a new best player - the reign continues.")
			fmt.Printf("(%s) remains the best\n  over %s (%s)\n", segmentWinningEvaluation,
				ap.getName(), thisEval)
		}
	}

	fmt.Println("\nThe community cards were:", t.community)
	fmt.Println("The winning hand:", segmentWinningEvaluation)
	payout := segmentValue / len(segmentWinningPlayers)
	fmt.Printf("The winners of the $%d segment, each winning $%d:\n", segmentValue, payout)
	for i, p := range segmentWinningPlayers {
		fmt.Printf("%d: %v\n", i, p)
		p.addToStack(payout)
		fmt.Printf("   %v\n", p)
	}
	os.Exit(4)
}

func RunTournament() {
	var table Table
	table.initialize()

	temp := NewAllInAlwaysPlayer("Adam")
	//temp.stack = 100 // To test Pot, we need two all-in short stack players.
	table.addPlayer(&temp)
	//temp1 := NewAllInAlwaysPlayer("Alma")
	//temp1.stack = 125 // To test Pot, we need two all-in short stack players.
	//table.addPlayer(&temp1)
	temp2 := NewGenericPlayer("Bert")
	table.addPlayer(&temp2)
	tempCSP := NewCallingStationPlayer("Cali")
	table.addPlayer(&tempCSP)
	temp4 := NewGenericPlayer("Dale")
	table.addPlayer(&temp4)
	temp5 := NewGenericPlayer("Eyor")
	table.addPlayer(&temp5)
	temp6 := NewFoldingPlayer("Fred")
	table.addPlayer(&temp6)
	testCTF := NewCallToFivePlayer("Greg")
	table.addPlayer(&testCTF)
	tempStreetFlopper := NewStreetTestPlayer("Flow", "FLOP")
	table.addPlayer(&tempStreetFlopper)
	testStretTurner := NewStreetTestPlayer("Turk", "TURN")
	table.addPlayer(&testStretTurner)
	testStreetRiverer := NewStreetTestPlayer("Rivv", "RIVER")
	table.addPlayer(&testStreetRiverer)
	tempMinRaiser := NewMinRaisingPlayer("Ming")
	table.addPlayer(&tempMinRaiser)
	tempSMP5 := NewSklanskyMalmuthPlayer("Stan", 5)
	table.addPlayer(&tempSMP5)
	fmt.Print("\n\n")

	// Set an initial small blind value.
	table.defineBlinds(25)

	//for i := 1; i <= 20; i++ {
	for {
		fmt.Println("============================")
		table.preset()
		fmt.Printf("This is game #%d.\n", table.gameCtr)
		table.assignInitialButtonAndBlinds()
		table.bettingRound = "PREFLOP"

		table.postBlinds()
		fmt.Println(table.getStatus())
		table.dealHoleCards()
		table.preFlopBet()
		table.moveBetsToPot()
		fmt.Println(table.getStatus())

		table.bettingRound = "FLOP"
		fmt.Println("---------------------- Dealing the flop.")
		table.dealFlop()
		table.postPreFlopBet()
		table.moveBetsToPot()
		fmt.Println(table.getStatus())

		table.bettingRound = "TURN"
		fmt.Println("---------------------- Dealing the turn.")
		table.dealTurn()
		table.postPreFlopBet()
		table.moveBetsToPot()
		fmt.Println(table.getStatus())

		table.bettingRound = "RIVER"
		fmt.Println("---------------------- Dealing the river.")
		table.dealRiver()

		// Mock the community cards for testing Evaluation()
		//mockedCardSet := NewCommunity()
		//mockedCardSet.Add(NewCard("H", 11))
		//mockedCardSet.Add(NewCard("H", 6))
		//mockedCardSet.Add(NewCard("C", 11))
		//mockedCardSet.Add(NewCard("S", 11))
		//mockedCardSet.Add(NewCard("C", 6))
		//table.community = mockedCardSet

		table.postPreFlopBet()
		table.moveBetsToPot()
		fmt.Println(table.getStatus())

		fmt.Println("Finding and paying the winners.")
		table.payWinners()

		fmt.Println("Removing busted players.")
		table.removeBustedPlayers()

		fmt.Println("At the end of the game, the table looks like this:")
		fmt.Println(table.getStatus())

		// Look for a winner.
		if len(table.players) == 1 {
			fmt.Println("After", table.gameCtr, "games, we have a winner.")
			break
		}

		//time.Sleep(1 * time.Second)
	}

	fmt.Println("=============================================================")
	fmt.Println("=============================================================")
	fmt.Println("=============================================================")
	fmt.Println()
	winner := *table.players[0]
	fmt.Println(winner.getName(), "has won with a stack of $", winner.getStack())
	fmt.Println("Here's the bust log:")
	var sortedGameCtrs []int
	for gameCtr := range table.bustLog {
		sortedGameCtrs = append(sortedGameCtrs, gameCtr)
	}
	sort.Sort(sort.IntSlice(sortedGameCtrs))
	for _, bustedGameCtr := range sortedGameCtrs {
		fmt.Printf("Round %d: ", bustedGameCtr)
		for _, bustedPlayer := range table.bustLog[bustedGameCtr] {
			bpp := *bustedPlayer
			fmt.Printf("%s, ", bpp.getName())
		}
		fmt.Println()
	}

}

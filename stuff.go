package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Card struct {
	suit string
	rank int
}

type CardSet struct {
	cards []Card
}

func emptyCardSet() CardSet {
	empty := CardSet{}

	return empty
}

type HoleCards struct {
	cardset CardSet
}

type Deck struct {
	cardset CardSet
}

type Community struct {
	cards []Card
}

type Player struct {
	name           string
	nextPlayer     *Player
	previousPlayer *Player
	holeCards	   HoleCards
}

// Player constructor
func NewPlayer(name string) Player {
	ecs := emptyCardSet()
	hc := HoleCards{cardset: ecs}

	newPlayer := Player{name, nil, nil, hc}
	//newPlayer := Player{name, nil, nil, nil}
	return newPlayer
}

type Pot struct {
	value int
}

type Table struct {
	players          []Player
	community        Community
	button           Player
	smallBlindValue  int
	smallBlindPlayer Player
	bigBlindValue    int
	bigBlindPlayer   Player
	gameCtr          int
	bettingRound     string
	deck			 Deck
	pot				 Pot
}

func (t *Table) initialize() {
	t.gameCtr = 1
}

//   receiver  name       inputs         return
func (t *Table) addPlayer(player Player) {
	if len(t.players) == 0 {
		t.players = append(t.players, player)
		return
	}

	lastPlayer := t.players[len(t.players)-1]
	lastPlayer.nextPlayer = &player
	player.previousPlayer = &lastPlayer
	player.nextPlayer = &t.players[0]
	t.players[0].previousPlayer = &player
	t.players = append(t.players, player)

	fmt.Printf("%s is after %s and before %s\n",
		player.name, player.previousPlayer.name, player.nextPlayer.name)
	return
}

func (t *Table) reset() {
	t.gameCtr++
}

func (t Table) printPlayerList() {
	//fmt.Println("Here's the table: ", t)

	fmt.Println("Players: ")
	for _, p := range t.players {
		fmt.Println(p.name)
	}
}

func (t *Table) assignButton() {
	// https://stackoverflow.com/questions/33994677/pick-a-random-value-from-a-go-slice
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(t.players)
	t.button = t.players[n]
	fmt.Println("Assigning the button to:", t.button.name)

	t.smallBlindPlayer = *t.button.nextPlayer
	//fmt.Println("Assigning SB to:", t.smallBlindPlayer.name)
	//t.bigBlindPlayer = *t.smallBlindPlayer.nextPlayer
	fmt.Println("Assigning BB to:", &t.bigBlindPlayer.name)
}

func (t *Table) defineBlinds(sb int) {
	t.smallBlindValue = sb
	t.bigBlindValue = sb * 2
}

func (t *Table) postBlinds() (table Table) {

	return
}

func runTournament() {
	var table Table
	table.initialize()

	var p1 Player
	p1.name = "Bert"
	table.addPlayer(p1)                       // way #1
	table.addPlayer(NewPlayer("Cail"))
	table.addPlayer(NewPlayer("Dale"))
	table.addPlayer(NewPlayer("Eyor")) // way #3
	table.printPlayerList()

	// Set an initial small blind value.
	table.assignButton()
	table.defineBlinds(25)

	for i := 1; i <= 2; i++ {
		fmt.Printf("This is game #%d.\n", table.gameCtr)
		table.postBlinds()
		table.bettingRound = "PREFLOP"

		table.bettingRound = "FLOP"

		table.reset()
	}
}

func main() {
	for i := 1; i <= 1; i++ {
		runTournament()
	}
	//
	////var localPlayers []string
	//fmt.Println("== make([]Player, 3) ==")
	//localPlayers := make([]Player, 3)
	//fmt.Println(localPlayers)
	//fmt.Println("== append(localPlayers, Player{'hokeydokey''} ==")
	//localPlayers = append(localPlayers, Player{"hokeydokey"})
	//fmt.Println(localPlayers)
	//
	//fmt.Println("== make([]int, 0) ==")
	//a := make([]int, 0)
	//fmt.Println(a)
	//fmt.Println("== append(a, 1, 2, 3) ==")
	//a = append(a, 1, 2, 3)
	//fmt.Println(a)
	//

}

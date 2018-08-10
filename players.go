package main

import "fmt"

type FoldingPlayer struct {
	GenericPlayer
}

func NewFoldingPlayer(name string) FoldingPlayer {
	ecs := getEmptyCardSet()
	hc := HoleCards{cardset: &ecs}
	initialStack := 1000 // dollars
	newPlayer := FoldingPlayer{GenericPlayer{name, nil, nil, hc, initialStack, 0,
		false, false}}
	return newPlayer
}

func main() {
	fmt.Println("Players playas")

	fp := NewFoldingPlayer("Fernie")

	fmt.Println(fp)
}

//
//type GenericPlayer struct {
//	name           string
//	nextPlayer     *GenericPlayer
//	previousPlayer *GenericPlayer
//
//	// The below get preset after each game.
//	holeCards HoleCards
//	stack     int
//	bet       int
//	hasFolded bool
//	isAllIn   bool
//}
//
//func (p *GenericPlayer) getStatus() string {
//	status := ""
//	status = fmt.Sprintf("%s: [%s] stack=%d bet=%d", p.name, p.holeCards.toString(), p.stack, p.bet)
//
//	return status
//}
//
//func (p *GenericPlayer) payBlind(blindAmount int) {
//	if p.stack > blindAmount {
//		p.stack -= blindAmount
//		p.bet = blindAmount
//	} else {
//		p.allIn()
//	}
//}
//
//func (p *GenericPlayer) allIn() {
//	p.bet += p.stack
//	p.stack = 0
//	p.isAllIn = true
//}
//
//// GenericPlayer constructor
//// http://www.golangpatterns.info/object-oriented/constructors
//func NewGenericPlayer(name string) GenericPlayer {
//	ecs := getEmptyCardSet()
//	hc := HoleCards{cardset: &ecs}
//	initialStack := 1000 // dollars
//	newPlayer := GenericPlayer{name, nil, nil, hc, initialStack, 0,
//		false, false}
//	return newPlayer
//}
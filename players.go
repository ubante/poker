package main

// I fail at imports.

//
//type Player struct {
//	name           string
//	nextPlayer     *Player
//	previousPlayer *Player
//
//	// The below get preset after each game.
//	holeCards HoleCards
//	stack     int
//	bet       int
//	hasFolded bool
//	isAllIn   bool
//}
//
//func (p *Player) getStatus() string {
//	status := ""
//	status = fmt.Sprintf("%s: [%s] stack=%d bet=%d", p.name, p.holeCards.toString(), p.stack, p.bet)
//
//	return status
//}
//
//func (p *Player) payBlind(blindAmount int) {
//	if p.stack > blindAmount {
//		p.stack -= blindAmount
//		p.bet = blindAmount
//	} else {
//		p.allIn()
//	}
//}
//
//func (p *Player) allIn() {
//	p.bet += p.stack
//	p.stack = 0
//	p.isAllIn = true
//}
//
//// Player constructor
//// http://www.golangpatterns.info/object-oriented/constructors
//func NewPlayer(name string) Player {
//	ecs := getEmptyCardSet()
//	hc := HoleCards{cardset: &ecs}
//	initialStack := 1000 // dollars
//	newPlayer := Player{name, nil, nil, hc, initialStack, 0,
//		false, false}
//	return newPlayer
//}
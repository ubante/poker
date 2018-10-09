package models

// This player is based on OddsComputingPlayer.  It will compute the
// current odds at each street and bet accordingly.

type VariableOddsPlayer struct {
	OddsComputingPlayer
}

func NewVariableOddsPlayer(name string, smLevel int, computedOddsPercentageLevel float64) VariableOddsPlayer {
	ecs := NewCardSet()
	hc := HoleCards{cardSet: &ecs}
	initialStack := 1000 // dollars

	newPlayer := new(VariableOddsPlayer)
	newPlayer.name = name
	newPlayer.holeCards = hc
	newPlayer.stack = initialStack
	newPlayer.smThreshold = smLevel
	newPlayer.oddsThreshold = computedOddsPercentageLevel

	return *newPlayer
}
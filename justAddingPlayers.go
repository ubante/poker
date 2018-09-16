package main

import (
	"goven/poker/models"
)

func main() {
	var c models.Card
	c.NumericalRank = 3

	models.TestTournament()

}
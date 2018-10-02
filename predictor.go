package main

import (
	"goven/poker/models"
	"fmt"
)

// Given the community rounds, find the odds of a single opponent having
// the different hands.

func main() {
	fmt.Println("Starting to guess....")

	// I tried to do it all in the main package but got sick of
	// capitalizing fields.
	models.VillainGuesser()

	// Compute the chance of winning given your hole cards and the flop.
	fmt.Println("================================================== Compute post-flop")
	models.ComputeFlop("H4C4", "S4HQCA")
}
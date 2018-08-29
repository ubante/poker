package main

import (
	"fmt"
	"goven/poker/models"
)

/*
We will run multiple tournaments to find the best player type.
Each tournament has multiple poker _games_.
Each game has multiple betting _rounds_.
Each winner of each game has the best _hand_ - multiple winners are possible.
*/

// This is just to cast a GenericPlayer to Player.
//func castToPlayer(player Player) Player {
//	return player
//}

func main() {
	for i := 1; i <= 1; i++ {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		fmt.Printf("Starting tournament #%d\n", i)
		models.RunTournament()
	}
}

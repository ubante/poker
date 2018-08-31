package main

import (
	"fmt"
	"goven/poker/models"
	"time"
	"sort"
)

/*
We will run multiple sit-n-go tournaments to find the best player type.
Each tournament has multiple poker _games_.
Each game has multiple betting _rounds_.
The winner of each game has the best _hand_.
Multiple winners are possible for a game but for a tournament,
there can only be one.
*/

// This is just to cast a GenericPlayer to Player.
//func castToPlayer(player Player) Player {
//	return player
//}

func main() {
	startTime := time.Now().Unix()
	winRecord := make(map[string][]int)
	numberOfTournaments := 49

	for i := 1; i <= numberOfTournaments; i++ {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		fmt.Printf("Starting tournament #%d\n", i)
		tournamentPlacing := models.RunTournament()

		// Record the placings.
		for playerName := range tournamentPlacing {
			if _, ok := winRecord[playerName]; ok {
				winRecord[playerName] = append(winRecord[playerName], tournamentPlacing[playerName])
			} else {
				winRecord[playerName] = []int{tournamentPlacing[playerName]}
			}
		}
	}
	endTime := time.Now().Unix()

	// Find the best player overall.
	overallPlacings := make(map[float64]string)
	fmt.Println("\nAfter", endTime-startTime, "seconds and", numberOfTournaments,
		"tournaments, this is the standing record (lower is better):")
	for playerName := range winRecord {
		sum := 0
		for _, placing := range winRecord[playerName] {
			sum += placing
		}
		overallPlacings[float64(sum)/float64(numberOfTournaments)] = playerName
	}

	// Print players from best to worst.  Note that these are just
	// standings.  If we sorted by payout, it could be very different.
	var sortedStandings []float64
	for standing := range overallPlacings {
		sortedStandings = append(sortedStandings, standing)
	}
	sort.Sort(sort.Float64Slice(sortedStandings))

	for _, sortedStanding := range sortedStandings {
		if numberOfTournaments > 40 {
			fmt.Printf("%s: %5.2f %v....\n", overallPlacings[sortedStanding], sortedStanding,
				winRecord[overallPlacings[sortedStanding]][0:40]) // A couple printed do.
		} else {
			fmt.Printf("%s: %5.2f %v\n", overallPlacings[sortedStanding], sortedStanding,
				winRecord[overallPlacings[sortedStanding]])
		}
	}

}

/*
AllInAlwaysPlayer: Adam
CallingStationPlayer: Cali
GenericPlayer: Dale, Jenn
FoldingPlayer: Fred
CallToFivePlayer: Carl
StreetTestPlayer: Flow, Turk, Rivv
MinRaisingPlayer: Ming
SklanskyMalmuthPlayer: Stan

After 543 seconds and 499 tournaments, this is the standing record (lower is better):
Fred:  2.06 [2 2 2 2 2 2 2 2 2 2 2 3 1 4 2 2 2 2 2 2 2 2 2 2 2 2 2 3 1 2 2 2 2 2 2 2 2 3 2 2]
Adam:  3.38 [3 4 4 4 4 4 3 4 3 4 3 5 3 6 4 3 4 3 4 3 4 3 4 4 3 4 4 4 3 4 4 4 3 3 1 3 4 4 4 3]
Cali:  4.04 [4 5 5 5 5 5 4 1 4 1 4 6 4 2 5 4 5 4 5 4 5 4 5 5 4 5 1 5 4 5 1 5 4 4 3 4 5 5 5 4]
Carl:  4.07 [6 3 3 3 3 3 6 3 5 3 6 4 2 5 3 6 3 5 3 5 3 6 3 3 6 3 3 7 2 3 3 3 6 6 5 6 3 7 3 6]
Dale:  4.78 [5 6 6 6 6 6 5 5 1 5 5 2 5 7 6 5 6 1 6 1 6 5 6 6 5 1 5 6 5 6 5 6 5 5 4 5 6 6 6 5]
Jenn:  5.96 [7 1 7 7 7 7 7 6 6 6 7 7 6 1 1 7 7 6 7 6 7 7 7 7 1 6 6 2 6 7 6 1 7 7 6 7 7 2 7 7]
Flow:  6.75 [1 7 1 8 8 8 8 7 7 7 8 8 7 8 7 8 8 7 8 7 8 8 8 8 7 7 7 8 7 8 7 7 1 1 7 1 8 8 1 8]
Turk:  7.76 [8 8 8 9 9 1 9 8 8 8 9 9 8 9 8 9 9 8 9 8 9 9 1 9 8 8 8 1 8 9 8 8 8 8 8 8 9 9 8 9]
Rivv:  8.06 [9 9 9 1 1 9 1 9 9 9 10 10 9 10 9 10 10 9 1 9 10 10 9 10 9 9 9 9 9 1 9 9 9 9 9 9 1 10 9 1]
Ming:  9.15 [10 10 10 10 10 10 10 10 10 10 11 11 10 11 10 11 11 10 10 10 1 1 10 11 10 10 10 10 10 10 10 10 10 10 10 10 10 11 10 10]
Stan:  9.99 [11 11 11 11 11 11 11 11 11 11 1 1 11 3 11 1 1 11 11 11 11 11 11 1 11 11 11 11 11 11 11 11 11 11 11 11 11 1 11 11]

After 25 seconds and 41 tournaments, this is the standing record (lower is better):
Fred:  2.12 [2 2 3 2 2 2 2 2 2 2 2 2 2 2 2 2 4 3 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 3 2 2 2 2 2]
Adam:  3.32 [3 3 5 4 3 3 3 3 3 3 4 3 4 4 1 3 5 4 3 4 1 4 4 3 1 3 3 4 3 4 4 3 4 3 4 1 4 4 3 4]
Cali:  4.07 [4 4 6 5 4 4 4 4 4 1 5 4 5 5 4 4 6 5 4 5 4 5 1 4 4 4 4 1 4 5 5 4 5 4 5 4 5 5 1 5]
Carl:  4.41 [6 6 4 3 6 6 6 5 6 5 3 6 3 3 3 5 3 7 6 3 3 3 3 6 3 6 6 3 6 3 3 6 3 6 6 3 3 3 5 3]
Dale:  4.80 [5 5 1 6 5 5 5 1 5 4 6 5 6 6 5 1 7 6 5 6 5 6 5 5 5 5 5 5 5 6 6 5 6 5 1 5 1 6 4 6]
Flow:  6.54 [8 8 8 7 8 8 8 7 8 7 1 1 1 8 7 7 9 9 1 1 7 8 7 8 7 8 1 7 8 8 8 8 8 8 8 7 7 1 7 8]
Jenn:  6.61 [7 7 7 1 7 7 7 6 7 6 7 7 7 7 6 6 8 8 7 7 6 7 6 7 6 7 7 6 7 7 7 7 7 7 7 6 6 7 6 7]
Turk:  7.61 [9 9 9 8 9 9 9 8 9 8 8 8 8 9 8 8 10 2 8 8 8 1 8 9 8 9 8 8 9 9 1 9 9 1 2 8 8 8 8 9]
Rivv:  8.37 [1 10 2 9 10 10 10 9 10 9 9 9 9 1 9 9 11 10 9 9 9 9 9 10 9 10 9 9 10 1 9 1 10 9 9 9 9 9 9 10]
Ming:  8.59 [10 1 10 10 11 11 11 10 1 10 10 10 10 10 10 10 1 1 10 10 10 10 10 1 10 11 10 10 1 10 10 10 1 10 10 10 10 10 10 11]
Stan:  9.56 [11 11 11 11 1 1 1 11 11 11 11 11 11 11 11 11 2 11 11 11 11 11 11 11 11 1 11 11 11 11 11 11 11 11 11 11 11 11 11 1]

After 38 seconds and 49 tournaments, this is the standing record (lower is better):
Fred:  2.00 [3 2 2 2 3 2 2 2 2 2 2 2 2 2 2 2 2 1 2 2 2 2 2 2 2 3 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 1 2 2 1 2 2 2 2]
Adam:  2.90 [5 4 1 4 4 3 4 1 4 3 3 3 4 3 3 4 1 3 3 4 3 1 4 3 4 4 1 4 3 3 1 3 1 1 4 4 4 3 1 1 3 3 1 4 3 1 4 3 3]
Cali:  3.73 [2 5 4 5 5 4 1 3 5 4 4 1 5 4 4 5 3 4 4 5 4 3 1 4 5 5 3 5 4 4 4 4 4 3 5 5 5 1 3 4 1 4 3 1 4 4 5 4 4]
Carl:  4.20 [4 3 3 3 2 6 3 5 3 6 6 5 3 6 6 3 5 2 6 3 6 5 3 5 3 7 5 3 6 6 3 6 3 5 3 3 3 5 5 3 5 2 5 3 2 3 3 6 6]
Dale:  4.82 [6 6 5 1 6 5 5 4 6 5 5 4 6 5 5 6 4 5 5 6 5 4 5 1 6 6 4 6 5 5 5 5 5 4 6 6 6 4 4 5 4 5 4 5 5 5 1 5 5]
Jenn:  5.80 [7 1 6 6 7 7 6 6 7 7 7 6 7 7 7 1 6 6 7 7 7 6 6 6 1 8 6 7 7 7 6 7 6 6 1 7 7 6 6 6 6 6 6 6 6 6 6 1 1]
Flow:  6.98 [1 7 7 7 8 8 7 7 1 8 8 7 8 8 8 7 7 7 8 8 1 7 7 7 7 9 7 8 8 8 7 8 7 7 7 8 8 7 7 7 7 7 7 7 7 7 7 7 7]
Turk:  7.51 [8 8 8 8 9 9 8 8 8 9 1 8 1 9 1 8 8 8 1 9 8 8 8 8 8 1 8 9 9 9 8 9 8 8 8 9 9 8 8 8 8 8 8 8 8 8 8 8 8]
Rivv:  8.69 [9 9 9 9 10 10 9 9 9 10 9 9 9 10 9 9 9 9 9 1 9 9 9 9 9 10 9 1 10 10 9 1 9 9 9 10 10 9 9 9 9 9 9 9 9 9 9 9 9]
Ming:  9.16 [10 10 10 10 1 1 10 10 10 1 10 10 10 1 10 10 10 10 10 10 10 10 10 10 10 11 10 10 11 1 10 10 10 10 10 11 11 10 10 10 10 10 10 10 10 10 10 10 10]
Stan: 10.20 [11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 11 2 11 11 1 11 11 11 11 11 11 1 1 11 11 11 11 11 11 11 11 11 11 11 11]

After 5 seconds and 49 tournaments, this is the standing record (lower is better):
Fred:  2.06 [2 2 2 1 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 3 4 2 3 2 3 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 1 2 2 2]
Adam:  3.29 [1 3 4 3 4 4 4 4 1 1 4 3 3 3 3 4 4 1 3 3 3 3 4 4 5 3 6 3 4 4 3 4 3 1 3 4 3 4 4 4 4 3 3 4 3 3 4 4 1]
Carl:  4.22 [5 6 3 2 3 3 3 3 3 5 3 6 6 6 6 3 3 3 5 5 6 6 3 2 3 6 2 7 7 3 6 3 5 5 6 3 6 3 3 3 3 6 6 3 6 2 3 3 5]
Cali:  4.31 [3 4 5 4 5 5 5 5 4 3 5 4 4 4 4 5 5 4 1 1 4 4 5 5 6 4 7 5 5 5 4 1 4 3 4 5 4 5 5 5 5 4 4 5 4 4 5 7 3]
Dale:  5.20 [4 5 6 5 6 6 6 6 5 4 1 5 5 5 5 6 6 5 4 4 5 5 6 6 7 5 8 6 6 6 5 5 1 4 5 6 5 6 6 6 6 5 5 6 5 5 6 5 4]
Jenn:  6.14 [6 7 1 6 7 7 7 7 6 6 6 1 7 7 7 7 7 6 6 6 7 1 7 7 2 7 4 8 8 7 7 6 6 6 7 7 7 1 7 7 7 7 7 7 7 6 7 8 6]
Flow:  7.00 [7 1 7 7 8 1 8 8 7 7 7 7 8 8 8 8 8 7 7 7 8 7 1 8 8 8 9 9 9 8 8 7 7 7 8 8 8 7 8 1 8 8 8 8 8 7 8 1 7]
Turk:  7.10 [8 8 8 8 9 8 1 9 8 8 8 8 1 9 9 9 9 8 8 8 1 8 8 9 9 1 5 4 2 1 9 8 8 8 9 1 9 8 9 8 9 9 9 1 9 8 9 9 8]
Ming:  8.29 [10 10 10 10 1 10 10 1 10 10 10 10 10 11 10 1 11 10 10 10 10 10 10 1 1 10 1 11 1 10 1 10 10 10 10 10 1 10 1 10 11 11 10 10 10 10 10 11 10]
Rivv:  8.51 [9 9 9 9 10 9 9 10 9 9 9 9 9 10 1 10 10 9 9 9 9 9 9 10 10 9 10 10 10 9 10 9 9 9 1 9 10 9 10 9 10 10 1 9 1 9 1 10 9]
Stan:  9.88 [11 11 11 11 11 11 11 11 11 11 11 11 11 1 11 11 1 11 11 11 11 11 11 11 11 11 11 1 11 11 11 11 11 11 11 11 11 11 11 11 1 1 11 11 11 11 11 6 11]


 */
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

	fmt.Println("=============================================================")
	fmt.Println(time.Now())
	fmt.Println("After", endTime-startTime, "seconds and", numberOfTournaments,
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

=============================================================
2018-09-23 23:38:41.3469576 -0700 PDT m=+102.264056901
After 102 seconds and 49 tournaments, this is the standing record (lower is better):
Fred:  2.20 [2 3 3 2 2 3 2 2 2 2 2 3 2 2 2 2 2 2 2 2 2 2 2 3 2 2 2 2 2 2 2 2 3 2 2 4 4 2 2 2]....
Carl:  4.43 [3 4 8 9 7 2 3 4 10 9 3 10 3 3 9 3 3 3 6 3 3 3 3 2 4 10 11 3 3 7 3 3 4 3 3 5 2 3 3 3]....
Turk:  6.31 [5 8 7 1 9 11 8 11 6 1 5 1 5 8 3 5 4 7 7 9 1 7 6 8 8 9 9 11 8 9 6 5 2 8 1 3 9 10 7 1]....
Dale:  6.33 [4 2 2 10 10 10 5 5 5 6 6 8 11 10 6 11 10 1 1 1 8 4 4 9 1 3 4 10 6 11 1 11 5 5 8 2 10 5 6 5]....
Flow:  6.35 [9 5 9 8 8 8 9 1 1 8 8 2 1 9 1 10 8 4 3 10 5 1 11 6 9 5 8 8 10 3 8 6 11 9 5 7 6 7 8 9]....
Ming:  6.63 [7 9 4 4 6 1 4 7 4 11 11 11 7 6 11 6 9 9 4 5 10 10 10 1 10 1 1 1 5 5 7 9 1 6 6 8 1 4 11 11]....
Rivv:  6.76 [1 7 10 3 11 4 7 10 8 3 1 5 9 11 5 7 1 6 5 8 11 5 9 7 11 11 5 6 1 1 4 10 7 10 9 6 11 8 5 4]....
Stan:  6.84 [6 6 11 7 3 5 1 9 7 5 9 6 4 1 7 9 6 11 9 11 7 8 7 5 3 4 6 4 9 10 5 8 6 11 11 9 3 6 4 7]....
Adam:  7.10 [8 10 5 11 1 6 10 3 3 7 7 9 6 5 8 8 5 10 10 6 6 11 1 11 6 8 10 7 7 4 9 4 10 7 10 10 5 11 10 8]....
=============================================================
2018-09-23 23:31:25.0175639 -0700 PDT m=+51.612138901
After 52 seconds and 49 tournaments, this is the standing record (lower is better):
Fred:  2.08 [2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 3 2 2 2 2 4 2 2 2 2 2 2 2 2 2 2 2 2]....
Carl:  3.98 [3 3 3 3 3 3 3 4 3 7 3 3 9 10 6 5 3 3 3 3 3 3 4 3 3 9 6 2 3 3 3 3 10 3 3 3 3 3 3 3]....
Jenn:  5.90 [5 4 9 4 4 7 10 1 4 3 9 8 6 6 5 7 6 1 5 10 11 7 1 9 5 1 1 7 4 11 4 5 7 6 7 4 1 8 8 5]....
Turk:  6.35 [1 10 4 9 8 10 6 8 11 1 11 7 4 8 11 8 4 4 1 1 1 4 10 5 6 7 9 5 11 5 11 9 8 4 10 9 8 6 1 9]....
Dale:  6.57 [9 7 6 1 9 11 1 6 10 11 8 5 7 1 4 9 5 6 7 9 8 1 9 10 7 10 10 3 9 8 1 6 1 1 11 1 5 10 9 6]....
Adam:  6.61 [10 6 11 11 10 5 5 3 5 6 5 1 8 7 1 3 9 8 11 11 6 10 11 1 10 6 4 9 6 6 8 1 9 5 1 10 9 1 5 7]....
Rivv:  7.04 [8 8 1 10 11 6 9 5 9 8 4 4 1 11 8 11 7 10 8 4 9 8 5 6 1 3 3 11 8 10 5 7 5 9 4 8 11 4 6 10]....
Flow:  7.08 [7 11 8 5 1 8 8 11 7 9 6 10 10 5 3 4 10 5 9 8 4 9 8 11 9 8 11 6 10 7 6 8 3 7 9 5 6 11 7 1]....
Cali:  7.24 [4 9 10 7 7 4 4 7 6 10 10 11 3 9 10 6 1 9 4 6 10 6 7 8 4 5 5 10 1 1 9 10 4 11 8 11 10 7 10 8]....
=============================================================
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
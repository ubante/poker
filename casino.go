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
	numberOfTournaments := 1
	//numberOfTournaments := 2
	//numberOfTournaments := 10
	//numberOfTournaments := 100
	//numberOfTournaments := 1000
	//numberOfTournaments := 100000  // 100k

	for tournamentNumber := 1; tournamentNumber <= numberOfTournaments; tournamentNumber++ {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Tournament #", tournamentNumber)
		fmt.Printf("Starting tournament #%d\n", tournamentNumber)
		tournamentPlacing := models.RunTournament(tournamentNumber)

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
SklanskyMalmuthPlayer: Stan (5), Saul (1)
SklanskyMalmultModifiedPlayer: Mits (5), Muts (1)
OddsComputingPlayer: Odom (5,80), Otis (6,70), Omar (2,70)
=============================================================
2018-10-07 02:19:28.1664666 -0700 PDT m=+3268.024353701
After 3268 seconds and 1000 tournaments, this is the standing record (lower is better):
Omar:  3.10 [2 3 2 3 2 4 1 1 5 2 4 4 5 2 2 2 2 3 3 2 4 2 1 8 4 2 3 3 10 3 5 11 5 1 2 3 2 4 3 2]....
Fred:  3.27 [5 2 3 5 3 2 3 3 2 3 5 3 3 3 4 3 4 4 2 5 3 3 3 2 2 3 2 5 4 5 4 3 2 3 3 2 4 3 4 3]....
Muts:  3.46 [4 4 4 2 5 3 2 2 3 4 2 2 2 4 6 4 1 2 5 4 2 4 2 3 3 5 6 4 2 4 3 2 3 2 1 1 5 6 2 6]....
Odom:  4.05 [1 5 6 1 1 1 7 7 4 7 3 6 4 6 1 5 5 9 1 1 6 10 6 6 5 8 9 2 1 2 2 4 1 5 6 11 1 2 9 1]....
Saul:  5.01 [3 1 8 7 6 5 5 5 6 5 1 1 1 5 3 6 3 1 7 9 1 6 5 4 6 7 5 7 3 6 7 6 4 7 5 6 6 1 1 9]....
Carl:  6.61 [11 6 7 6 10 6 11 6 8 6 7 7 6 11 5 7 6 7 6 6 7 5 7 5 10 6 4 6 5 7 6 5 6 4 4 5 7 9 6 5]....
Adam:  8.02 [8 10 1 11 9 8 10 9 7 8 8 5 8 8 11 11 10 10 9 10 5 7 4 10 9 11 1 9 9 11 11 9 8 11 7 4 10 7 10 7]....
Cali:  8.04 [6 7 10 10 4 7 8 4 9 9 10 11 11 10 8 9 11 6 4 8 11 8 9 1 1 4 8 10 6 9 8 7 10 6 9 7 9 8 8 4]....
Turk:  8.06 [7 8 11 4 7 9 6 10 10 1 9 8 9 7 9 8 8 11 10 3 9 9 11 11 7 9 11 8 8 10 9 8 11 10 11 10 11 5 7 10]....
Ming:  8.11 [10 9 9 9 8 11 4 11 1 11 11 9 10 1 7 10 7 5 8 7 8 11 10 9 8 10 7 11 11 8 10 1 7 8 8 9 8 10 5 8]....
Jenn:  8.28 [9 11 5 8 11 10 9 8 11 10 6 10 7 9 10 1 9 8 11 11 10 1 8 7 11 1 10 1 7 1 1 10 9 9 10 8 3 11 11 11]....
=============================================================
2018-10-07 00:03:17.979018 -0700 PDT m=+3172.499521801
After 3172 seconds and 1000 tournaments, this is the standing record (lower is better):
Fred:  3.06 [2 2 2 4 4 3 2 4 2 3 4 3 3 3 2 2 2 4 3 2 4 2 3 4 3 3 2 2 3 2 3 3 4 3 3 4 2 3 3 4]....
Muts:  3.40 [3 3 5 3 3 4 9 7 3 4 2 2 7 4 3 4 1 3 2 4 3 3 2 5 2 2 3 4 4 5 5 4 5 5 2 2 3 2 1 2]....
Odom:  3.92 [4 8 6 1 2 1 1 2 1 7 3 1 2 7 6 3 5 2 4 1 2 5 5 3 1 4 5 6 2 8 4 7 1 2 6 3 4 7 2 3]....
Otis:  4.31 [6 1 3 5 6 10 6 1 9 1 5 6 5 2 5 7 8 10 5 5 1 9 4 1 6 6 6 1 1 6 2 2 2 7 1 1 7 1 6 5]....
Saul:  4.64 [1 4 1 6 1 2 5 3 4 5 8 7 1 6 7 6 6 6 6 3 5 6 11 2 10 5 1 5 6 3 10 6 11 4 4 5 5 4 7 1]....
Carl:  6.51 [10 5 7 7 8 6 4 6 6 10 6 5 6 5 4 5 4 5 9 6 6 4 6 8 4 7 4 7 7 7 6 5 6 6 5 8 6 6 4 6]....
Ming:  7.73 [8 6 4 10 10 8 3 10 7 6 1 10 4 11 8 11 3 8 1 11 10 11 1 11 5 1 11 3 9 1 11 1 7 8 8 7 10 8 10 7]....
Adam:  7.91 [9 11 9 8 7 9 11 8 5 8 9 8 11 9 1 9 7 7 11 7 9 10 7 7 7 11 9 10 10 11 7 11 10 9 7 9 1 11 9 8]....
Jenn:  8.02 [7 7 10 2 5 11 8 9 11 2 11 11 9 1 10 10 10 1 8 8 7 1 8 6 9 9 7 8 11 9 1 9 3 10 11 6 9 9 5 10]....
Turk:  8.24 [11 10 8 9 11 7 7 5 10 11 7 9 8 8 11 8 9 9 10 9 11 8 9 9 8 8 10 9 5 10 9 10 8 11 9 11 8 5 8 9]....
Cali:  8.24 [5 9 11 11 9 5 10 11 8 9 10 4 10 10 9 1 11 11 7 10 8 7 10 10 11 10 8 11 8 4 8 8 9 1 10 10 11 10 11 11]....
=============================================================
2018-10-06 21:20:23.949638 -0700 PDT m=+1888.292182301
After 1888 seconds and 1000 tournaments, this is the standing record (lower is better):
Fred:  2.72 [3 3 3 2 3 2 2 2 3 2 2 2 3 3 4 3 3 2 2 3 2 3 3 2 3 2 3 3 3 3 2 1 3 3 3 2 3 3 3 3]....
Muts:  3.11 [2 2 2 3 2 6 4 4 2 3 4 5 4 5 2 2 6 6 3 2 3 5 2 5 2 3 2 2 2 5 5 2 2 2 2 3 2 5 2 4]....
Otis:  3.85 [1 1 5 4 5 3 1 3 9 4 6 3 10 1 1 12 10 5 4 5 6 1 1 8 1 1 11 7 4 2 6 4 5 11 1 8 5 2 10 2]....
Saul:  4.46 [5 6 1 1 4 1 5 5 1 1 3 4 5 6 5 4 4 1 5 11 1 4 5 1 6 4 6 1 5 1 1 9 6 1 5 5 6 4 5 5]....
Carl:  6.10 [11 5 4 5 8 5 11 6 5 6 5 6 6 4 12 11 5 4 12 4 5 6 6 4 5 7 9 4 12 4 3 3 4 5 6 6 4 7 4 7]....
Jenn:  8.14 [7 7 10 9 7 9 7 11 10 5 10 8 12 10 6 8 2 3 11 7 9 2 9 3 10 11 1 10 7 9 8 11 8 9 9 9 1 1 11 6]....
Cali:  8.16 [8 9 6 8 9 8 6 12 4 7 7 11 11 2 10 7 1 7 8 12 4 12 8 10 7 5 4 6 11 11 11 8 12 12 11 4 8 9 9 9]....
Ming:  8.16 [9 11 11 11 6 12 8 7 11 10 11 10 8 7 9 10 7 9 10 10 10 10 7 12 4 10 8 11 1 7 4 10 10 7 8 1 7 8 1 11]....
Adam:  8.21 [6 8 9 10 11 7 9 8 6 11 12 1 9 8 11 9 8 11 6 1 8 8 10 11 11 9 5 5 9 6 12 12 1 10 7 12 9 10 7 12]....
Turk:  8.23 [12 10 8 12 10 4 3 1 8 12 1 7 7 11 8 5 9 8 9 8 12 9 11 7 8 8 7 8 10 8 7 6 9 6 10 7 11 6 6 1]....
Flow:  8.40 [10 12 7 6 12 11 10 9 12 9 9 9 2 12 3 1 12 12 1 6 11 11 4 9 9 12 12 9 8 10 9 5 11 8 4 11 10 12 12 10]....
Rivv:  8.45 [4 4 12 7 1 10 12 10 7 8 8 12 1 9 7 6 11 10 7 9 7 7 12 6 12 6 10 12 6 12 10 7 7 4 12 10 12 11 8 8]....
=============================================================
2018-10-06 20:47:00.1581101 -0700 PDT m=+79.555106801
After 80 seconds and 49 tournaments, this is the standing record (lower is better):
Fred:  2.69 [3 4 2 3 2 2 3 2 3 3 2 2 4 2 3 3 3 3 3 3 2 3 3 3 2 2 2 2 3 3 3 2 2 3 4 2 3 3 2 3]....
Muts:  2.80 [4 3 4 2 5 3 2 3 2 2 5 6 2 4 2 2 2 2 2 2 3 2 4 2 1 4 3 3 1 4 4 3 3 2 2 6 1 2 3 2]....
Otis:  4.08 [1 1 3 1 1 5 4 8 5 4 4 4 1 11 1 4 9 12 9 5 7 1 1 4 4 6 1 6 4 1 2 1 6 7 3 3 4 4 6 1]....
Saul:  4.82 [11 6 6 5 4 4 1 5 10 1 1 3 5 3 5 5 4 7 5 1 5 6 5 5 12 3 4 4 5 5 6 5 1 1 1 4 2 8 5 6]....
Carl:  6.49 [10 5 11 6 6 6 6 4 4 9 6 5 8 10 6 10 5 11 4 10 4 5 8 12 3 5 6 5 6 12 5 6 5 9 11 5 5 5 4 7]....
Ming:  7.33 [5 10 10 7 9 8 9 6 1 8 10 1 3 9 10 7 1 5 12 12 10 12 7 7 9 1 9 7 2 6 9 7 11 12 9 1 6 1 11 8]....
Flow:  7.78 [2 7 8 12 3 1 12 10 9 5 11 8 12 6 7 9 7 9 8 7 8 4 6 8 10 7 7 8 10 9 7 9 7 10 5 7 12 7 10 4]....
Adam:  8.08 [8 9 1 11 11 11 11 9 11 6 9 9 6 1 9 8 6 6 6 9 9 8 12 1 5 8 12 1 9 10 11 10 9 6 7 11 7 11 7 9]....
Turk:  8.14 [7 12 12 9 7 12 10 1 6 12 3 12 7 8 4 1 12 10 11 4 6 11 2 6 6 9 8 11 7 2 1 12 8 4 12 10 10 10 12 10]....
Cali:  8.43 [6 8 7 10 10 7 5 11 8 7 7 7 10 7 8 6 10 1 1 11 1 9 11 11 11 12 11 9 11 11 12 8 10 11 6 8 11 9 8 5]....
Rivv:  8.65 [12 11 9 4 12 9 8 7 7 11 12 11 9 5 11 11 11 8 7 8 12 7 9 10 7 10 5 12 8 8 10 4 12 5 10 9 8 12 9 11]....
Jenn:  8.71 [9 2 5 8 8 10 7 12 12 10 8 10 11 12 12 12 8 4 10 6 11 10 10 9 8 11 10 10 12 7 8 11 4 8 8 12 9 6 1 12]....
=============================================================
2018-10-04 23:30:25.544134 -0700 PDT m=+51.417203101
After 51 seconds and 49 tournaments, this is the standing record (lower is better):
Fred:  2.69 [2 2 2 3 3 2 2 3 4 3 3 3 3 2 3 3 2 2 4 2 4 2 4 2 3 3 3 3 3 3 1 3 3 3 3 2 2 3 3 2]....
Muts:  2.92 [3 4 1 2 2 3 3 2 3 2 2 2 2 5 6 2 4 4 2 3 3 3 3 4 2 2 2 2 4 2 4 2 2 4 2 3 3 2 2 8]....
Saul:  4.43 [8 6 5 1 1 5 5 1 11 1 1 5 4 1 4 1 1 1 9 5 5 5 5 5 5 4 1 6 2 1 11 1 5 5 6 5 12 5 1 4]....
Otis:  5.02 [4 1 10 7 6 4 11 4 1 9 4 1 7 6 5 12 5 6 1 6 2 6 6 6 4 5 8 4 1 5 3 4 1 7 5 4 4 1 6 1]....
Carl:  6.06 [11 5 4 5 5 10 4 6 8 8 11 6 11 3 12 5 3 5 7 4 6 4 2 3 7 8 5 5 8 12 2 5 6 2 4 12 5 6 5 12]....
Rivv:  7.24 [6 3 11 11 8 1 9 10 6 11 6 7 1 7 7 11 10 7 8 1 11 8 12 8 11 1 11 10 9 8 8 9 4 6 10 6 10 4 4 5]....
Adam:  7.86 [1 12 7 6 4 8 6 8 12 4 5 8 10 11 11 9 9 8 6 9 9 1 9 10 12 7 6 1 7 11 12 7 10 12 11 11 1 7 8 10]....
Turk:  7.96 [10 9 12 8 11 9 10 5 5 5 7 4 5 12 8 4 6 11 12 12 8 10 8 7 6 10 4 12 5 9 9 6 7 9 12 8 7 12 11 3]....
Jenn:  8.16 [9 8 3 12 10 11 1 7 2 10 8 12 6 10 9 7 7 10 3 7 1 11 10 12 9 11 9 8 12 4 5 10 12 11 8 10 9 8 12 9]....
Ming:  8.31 [12 10 9 9 12 7 12 9 10 12 12 9 9 4 1 6 12 3 5 10 7 12 1 1 8 6 12 7 10 6 7 11 8 1 7 9 11 9 7 11]....
Flow:  8.37 [5 11 8 10 9 12 8 12 7 6 10 10 8 8 10 8 11 9 11 8 12 7 7 9 1 12 10 9 6 10 6 8 9 8 1 1 6 11 10 6]....
Cali:  8.98 [7 7 6 4 7 6 7 11 9 7 9 11 12 9 2 10 8 12 10 11 10 9 11 11 10 9 7 11 11 7 10 12 11 10 9 7 8 10 9 7]....
=============================================================
2018-09-26 04:07:19.4363792 -0700 PDT m=+11655.606158301
After 11656 seconds and 100000 tournaments, this is the standing record (lower is better):
Fred:  3.00 [2 2 4 2 3 2 3 4 3 2 4 3 3 4 4 4 2 3 3 4 4 3 3 3 3 2 3 3 3 3 3 2 2 3 3 3 2 3 3 3]....
Muts:  3.41 [3 5 3 3 2 1 6 3 1 3 3 5 5 2 2 2 4 2 4 2 1 7 7 2 2 3 5 1 4 1 4 4 3 2 2 4 1 2 6 5]....
Stan:  3.42 [4 3 2 4 4 4 2 2 4 1 2 2 1 1 3 3 3 5 2 3 3 4 2 4 4 4 4 7 2 4 6 5 4 4 4 2 3 4 4 2]....
Mits:  5.05 [6 10 13 5 5 7 7 7 7 6 1 1 10 5 1 5 7 1 9 1 6 6 5 7 5 6 6 5 6 5 5 3 7 12 5 1 6 1 5 6]....
Saul:  5.08 [5 4 8 9 6 6 4 1 9 4 7 9 4 6 5 1 6 4 6 7 7 8 6 5 1 5 2 6 5 2 13 6 6 1 6 7 4 7 7 4]....
Carl:  7.08 [7 6 10 12 11 5 5 6 5 5 6 12 12 8 6 8 5 11 5 5 5 5 4 6 11 13 7 4 7 8 9 9 5 5 7 6 9 6 8 11]....
Adam:  9.07 [12 11 7 6 10 8 9 10 11 10 8 11 11 7 10 11 10 6 1 11 12 10 9 1 9 10 9 11 1 6 12 13 10 8 9 5 13 10 13 8]....
Ming:  9.13 [8 9 9 10 1 12 1 11 2 9 12 4 8 11 8 7 9 10 12 9 11 13 1 9 6 7 1 10 11 11 8 10 8 10 10 13 10 9 10 10]....
Cali:  9.14 [1 7 11 8 8 13 11 5 13 11 13 10 7 3 12 13 1 7 7 12 2 12 12 8 13 12 11 2 10 12 2 1 12 7 13 10 5 12 2 1]....
Flow:  9.15 [11 8 1 13 9 10 13 8 6 8 11 6 2 9 9 6 8 9 8 13 9 11 8 13 8 1 10 13 13 7 7 11 1 6 11 12 12 13 12 13]....
Turk:  9.15 [9 13 5 1 7 11 12 9 8 12 9 7 9 12 11 10 12 13 13 6 13 2 13 11 10 9 13 12 9 9 1 7 11 11 12 9 7 5 11 12]....
Jenn:  9.15 [13 1 12 7 13 9 10 13 10 7 5 8 6 10 13 12 11 8 11 8 8 1 10 10 7 8 12 8 8 10 11 12 13 9 8 8 8 11 1 9]....
Rivv:  9.16 [10 12 6 11 12 3 8 12 12 13 10 13 13 13 7 9 13 12 10 10 10 9 11 12 12 11 8 9 12 13 10 8 9 13 1 11 11 8 9 7]....
============================================================= Faster with increasing blinds
2018-09-26 00:42:33.9576433 -0700 PDT m=+119.236215201
After 119 seconds and 1000 tournaments, this is the standing record (lower is better):
Fred:  3.04 [2 4 2 2 4 3 2 2 3 2 4 3 3 3 3 4 3 2 3 3 4 3 3 3 3 3 3 3 3 2 2 3 4 4 4 3 3 5 3 3]....
Stan:  3.37 [5 3 6 3 1 2 3 3 2 3 3 1 2 2 2 6 7 4 1 1 2 1 4 2 2 1 2 4 2 6 1 7 3 3 3 2 5 2 4 2]....
Muts:  3.42 [3 2 3 7 3 4 1 4 6 4 2 2 4 5 5 3 4 5 4 2 5 6 2 4 5 2 4 2 5 4 4 2 2 2 1 6 2 4 2 4]....
Saul:  5.00 [4 1 4 5 6 1 12 1 4 1 13 6 6 6 6 7 5 6 2 5 1 5 1 5 1 5 6 1 4 7 6 4 6 1 2 4 4 7 5 5]....
Mits:  5.03 [6 7 5 4 12 7 10 6 5 6 1 4 7 1 4 2 1 3 9 7 3 7 6 6 7 6 5 5 6 5 5 1 1 6 7 1 1 6 6 9]....
Carl:  7.18 [10 12 11 6 5 6 4 7 12 13 6 7 5 10 10 5 6 11 5 13 6 4 10 7 8 7 7 10 10 3 13 11 7 8 6 7 7 3 12 7]....
Adam:  9.01 [8 13 7 10 10 11 11 5 7 9 8 12 10 12 8 1 10 10 11 12 9 12 9 12 6 12 1 7 1 13 11 12 8 5 12 11 11 13 9 1]....
Ming:  9.03 [11 11 1 11 11 8 8 9 10 12 10 8 11 11 7 8 13 9 12 4 11 10 13 8 12 10 10 8 12 1 8 13 10 11 10 12 8 9 11 13]....
Jenn:  9.04 [7 5 13 8 7 5 13 10 8 5 12 13 12 13 11 12 2 8 13 8 10 11 5 11 13 13 13 13 13 12 9 6 12 7 9 13 6 8 10 11]....
Flow:  9.10 [12 9 8 12 2 13 9 13 1 10 7 10 8 4 1 11 12 1 10 6 13 8 7 10 11 9 11 9 7 10 10 5 13 13 11 8 10 12 8 6]....
Cali:  9.12 [9 6 10 1 8 9 5 11 13 8 11 11 13 7 13 13 9 13 7 11 7 13 11 1 10 8 9 12 8 11 12 9 5 10 8 9 9 11 1 12]....
Turk:  9.27 [1 10 9 13 13 10 7 8 11 11 9 5 9 8 9 9 8 12 6 9 12 2 12 9 9 11 8 6 11 8 3 8 11 12 5 10 12 10 7 10]....
Rivv:  9.38 [13 8 12 9 9 12 6 12 9 7 5 9 1 9 12 10 11 7 8 10 8 9 8 13 4 4 12 11 9 9 7 10 9 9 13 5 13 1 13 8]....
============================================================= SMM (1) is slightly better than SM (1)
2018-09-26 00:30:27.4145662 -0700 PDT m=+861.318089601
After 861 seconds and 1000 tournaments, this is the standing record (lower is better):
Fred:  2.67 [2 3 2 2 4 3 2 3 3 3 2 2 2 2 4 2 3 4 2 4 3 2 2 4 3 2 3 3 3 2 3 4 3 3 3 2 2 2 4 2]....
Muts:  3.29 [4 2 3 5 2 4 7 2 2 2 5 1 3 1 1 3 1 2 6 1 2 1 3 1 7 4 2 2 4 3 2 3 2 2 2 3 1 3 1 4]....
Stan:  3.31 [3 1 5 1 3 1 3 5 4 4 3 12 4 4 5 4 5 3 7 3 7 3 6 3 4 5 5 6 2 6 4 2 1 4 5 4 5 1 5 3]....
Saul:  5.10 [5 6 6 13 1 5 4 4 1 1 1 5 5 5 2 1 8 1 4 6 4 5 12 2 1 6 4 4 6 7 5 7 4 12 7 13 13 7 6 10]....
Mits:  5.28 [10 7 4 4 12 6 6 13 12 5 6 13 1 6 3 7 6 6 5 2 1 6 4 8 2 7 1 1 1 4 1 1 7 6 6 6 7 5 7 5]....
Carl:  7.15 [7 5 7 11 7 7 5 6 13 6 7 4 7 7 7 6 4 5 3 11 5 10 5 13 6 3 13 5 8 5 7 6 6 9 8 5 6 6 2 6]....
Adam:  9.05 [11 11 12 7 6 2 11 7 7 12 13 8 12 8 6 12 2 11 8 9 10 8 11 12 13 12 12 7 7 8 12 12 10 8 12 7 11 8 10 1]....
Cali:  9.08 [6 12 8 9 10 12 9 9 5 11 11 11 6 10 13 10 7 7 12 13 9 13 10 10 11 1 8 8 9 11 6 11 13 11 4 10 12 13 11 12]....
Rivv:  9.14 [1 4 13 10 8 10 12 1 9 8 10 9 11 13 9 5 12 10 9 12 12 9 7 9 5 8 10 12 10 1 13 9 12 5 11 1 4 4 8 13]....
Jenn:  9.14 [9 8 1 6 13 13 1 11 10 10 12 10 8 12 10 13 10 8 10 7 13 4 8 7 9 10 11 9 13 12 8 13 9 10 10 9 9 9 9 8]....
Ming:  9.23 [13 10 9 8 9 9 10 8 6 7 9 7 10 11 11 11 13 12 13 5 11 7 13 5 10 13 7 10 5 10 11 5 8 1 1 8 10 12 3 9]....
Turk:  9.24 [12 13 11 3 5 8 13 12 11 9 8 3 9 9 8 8 9 13 1 10 8 11 1 6 12 9 6 13 11 13 10 10 5 13 9 12 8 10 12 7]....
Flow:  9.32 [8 9 10 12 11 11 8 10 8 13 4 6 13 3 12 9 11 9 11 8 6 12 9 11 8 11 9 11 12 9 9 8 11 7 13 11 3 11 13 11]....
============================================================= Created SMM player
2018-09-26 00:07:25.902844 -0700 PDT m=+28.688851901
After 28 seconds and 49 tournaments, this is the standing record (lower is better):
Fred:  2.16 [3 3 2 2 2 2 2 2 2 2 2 2 3 3 2 2 2 2 2 2 2 2 2 3 2 2 2 3 2 2 3 2 2 2 2 2 2 2 2 2]....
Stan:  3.98 [4 1 5 1 5 4 3 3 5 3 4 3 9 4 4 4 4 1 1 3 4 4 4 4 5 4 8 5 6 6 4 4 1 1 3 1 3 5 9 1]....
Saul:  4.00 [2 2 1 3 3 3 1 5 6 4 3 1 12 5 5 1 1 3 4 4 5 3 5 5 3 3 1 4 4 8 7 1 3 6 6 4 4 4 4 4]....
Mits:  4.49 [1 5 11 4 4 7 6 4 4 1 5 4 4 1 1 3 3 6 3 5 7 5 3 6 1 5 4 6 5 3 2 3 5 5 4 6 6 3 5 9]....
Carl:  6.16 [7 6 10 11 6 5 5 12 3 5 9 5 7 7 3 5 6 4 5 12 3 6 8 2 4 6 3 2 3 9 10 5 4 3 7 8 5 12 3 12]....
Adam:  7.20 [11 4 7 10 10 1 9 1 7 9 7 9 6 11 7 10 10 10 9 10 1 10 7 1 10 1 12 1 12 12 1 7 6 4 11 7 9 1 11 8]....
Cali:  7.90 [8 8 8 7 9 6 11 8 11 8 6 7 10 10 6 6 12 5 8 7 8 8 6 9 9 10 10 12 8 4 9 8 8 8 10 5 10 10 1 7]....
Flow:  8.16 [6 7 3 9 12 12 7 10 9 10 10 6 1 8 11 7 5 9 11 9 10 9 9 7 7 7 9 11 10 10 11 12 12 9 1 11 8 7 12 10]....
Ming:  8.37 [10 9 9 5 7 10 10 6 10 7 8 8 11 2 9 9 9 12 7 6 11 7 10 8 6 11 11 7 11 1 12 9 11 10 12 9 1 6 6 3]....
Rivv:  8.41 [12 12 4 6 11 8 4 11 8 6 12 11 5 9 8 8 11 7 6 8 9 11 11 10 12 12 6 8 7 11 6 6 9 7 5 3 11 8 7 11]....
Turk:  8.57 [9 10 6 12 1 9 12 9 12 11 11 12 2 6 10 12 7 11 12 1 6 12 12 12 8 9 5 9 9 7 5 10 10 12 9 12 7 11 10 5]....
Jenn:  8.59 [5 11 12 8 8 11 8 7 1 12 1 10 8 12 12 11 8 8 10 11 12 1 1 11 11 8 7 10 1 5 8 11 7 11 8 10 12 9 8 6]....
============================================================= SM player 5 is good
2018-09-25 23:21:52.0977769 -0700 PDT m=+596.611054001
After 597 seconds and 1000 tournaments, this is the standing record (lower is better):
Fred:  2.07 [2 2 2 3 2 2 2 2 2 2 2 2 2 2 3 2 2 2 2 2 2 2 2 2 2 2 3 2 3 2 2 2 2 2 3 2 2 2 2 2]....
Stan:  3.83 [5 1 1 4 3 4 1 4 6 5 1 5 8 4 1 4 4 3 3 4 5 4 1 4 4 4 5 1 5 5 3 7 4 4 1 4 3 1 4 3]....
Saul:  3.92 [4 5 4 5 4 11 4 6 4 1 10 4 1 5 2 12 1 8 1 10 4 5 4 5 1 1 6 3 2 4 4 1 5 5 2 9 5 3 3 4]....
Carl:  5.50 [3 10 5 2 11 3 5 3 3 4 3 3 4 3 8 3 3 10 10 3 3 3 11 3 5 3 4 4 4 3 8 8 3 3 8 3 4 4 10 5]....
Jenn:  7.67 [9 6 12 8 7 12 8 5 1 6 7 7 7 8 9 9 8 4 4 9 11 7 5 10 6 9 12 12 9 12 9 3 9 8 7 8 12 9 9 8]....
Adam:  7.71 [6 9 10 10 8 9 3 11 7 10 9 6 10 6 4 11 5 12 9 1 6 6 3 8 12 11 11 6 6 9 6 5 1 10 11 1 1 10 7 1]....
Dale:  7.77 [10 8 9 6 9 10 12 12 11 8 8 11 11 9 12 1 6 7 6 6 1 8 9 7 7 12 10 7 7 1 12 4 12 9 4 12 11 12 6 9]....
Rivv:  7.77 [8 4 11 11 10 6 11 7 12 11 4 9 3 10 10 8 9 1 7 11 10 9 7 12 9 8 2 11 12 7 11 12 7 6 5 7 10 7 1 10]....
Ming:  7.83 [1 3 3 1 6 8 9 10 9 12 12 8 9 1 7 6 11 6 11 12 8 12 8 9 8 5 8 8 1 10 10 9 11 12 12 6 6 8 8 11]....
Turk:  7.87 [11 11 6 9 5 7 7 8 8 9 11 10 6 7 5 10 12 11 8 7 9 1 10 1 11 6 7 5 8 11 7 10 8 1 10 11 7 5 12 7]....
Cali:  7.90 [12 7 8 12 12 5 6 1 10 7 5 1 5 11 6 5 10 5 12 5 12 10 6 6 3 7 1 9 11 8 5 11 6 11 6 10 8 11 11 12]....
Flow:  8.16 [7 12 7 7 1 1 10 9 5 3 6 12 12 12 11 7 7 9 5 8 7 11 12 11 10 10 9 10 10 6 1 6 10 7 9 5 9 6 5 6]....
=============================================================
2018-09-24 21:39:36.6115527 -0700 PDT m=+6.084523501
After 6 seconds and 49 tournaments, this is the standing record (lower is better):
Fred:  2.00 [2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 1 2 2 2 2]....
Carl:  3.69 [4 3 3 3 3 3 3 8 3 3 3 11 3 3 7 3 6 3 3 3 4 3 3 4 3 4 7 3 3 3 3 3 3 3 3 2 3 3 3 3]....
Adam:  6.00 [1 5 1 6 8 4 11 7 9 11 8 7 7 7 1 4 1 8 6 10 1 11 4 11 4 1 1 8 1 11 9 7 6 6 7 5 6 10 10 5]....
Jenn:  6.29 [5 11 9 10 1 9 5 4 11 7 10 4 4 4 11 10 9 11 5 1 5 6 10 3 11 6 3 4 4 1 10 5 8 1 11 6 5 6 5 10]....
Cali:  6.53 [6 10 5 1 10 7 1 10 4 10 7 9 9 1 3 11 11 9 1 11 8 8 6 5 6 10 5 1 5 8 6 4 4 4 1 4 9 1 9 6]....
Flow:  6.69 [3 4 4 11 11 10 8 6 7 4 1 6 11 6 9 9 3 4 4 5 9 1 5 10 10 7 10 5 9 4 1 10 7 10 6 7 4 4 4 11]....
Turk:  6.86 [7 6 11 4 5 5 7 1 1 8 4 3 1 11 10 6 8 5 10 8 11 5 8 1 9 11 8 10 10 6 11 8 9 9 10 8 7 5 1 8]....
Stan:  6.88 [10 1 8 8 6 6 4 11 5 6 6 1 8 8 6 7 7 6 7 6 7 9 9 7 5 5 6 7 11 5 4 1 1 11 8 11 1 8 11 9]....
Rivv:  6.94 [9 7 10 7 9 1 10 3 8 1 9 10 10 9 4 1 5 7 9 9 6 7 1 8 7 8 11 6 8 7 7 9 10 5 4 9 8 7 8 4]....
Dale:  7.04 [11 9 6 9 7 8 9 5 10 9 5 5 5 5 8 8 4 1 8 7 10 10 11 6 1 3 4 11 6 10 5 11 11 7 5 3 11 11 7 1]....
Ming:  7.08 [8 8 7 5 4 11 6 9 6 5 11 8 6 10 5 5 10 10 11 4 3 4 7 9 8 9 9 9 7 9 8 6 5 8 9 10 10 9 6 7]....

 */
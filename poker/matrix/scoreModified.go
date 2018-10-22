package matrix

import (
	"sync"
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
)

// This is a singletonModified that contains the default Sklansky Malmuth hole
// card rating matrix.  It will read the text file once and create a
// map of maps that represent the SM matrix.  The singletonModified class will
// be called from different places but the file will be read in just
// once.

// We shouldn't need a separate class for this.

// https://medium.com/@MrToBe/the-singleton-object-oriented-design-pattern-in-golang-9f6ce75c21f7
// http://marcio.io/2015/07/singleton-pattern-in-go/

type ScoreModified struct {
	filename string
	matrixMap map[string]map[string]int
}

var singletonModified *ScoreModified
var onceModified sync.Once

func (s *ScoreModified) readInFile() {
	fmt.Println("Oh snap; reading in:", s.filename)
	file, err := os.Open(s.filename)
	if err != nil {
		fmt.Println("Error opening file:", s.filename)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	ranks := []string{"A", "K", "Q", "J", "T", "9", "8", "7", "6", "5", "4", "3", "2"}

	s.matrixMap = map[string]map[string]int{
		// How to do this better?
		"A": {},
		"K": {},
		"Q": {},
		"J": {},
		"T": {},
		"9": {},
		"8": {},
		"7": {},
		"6": {},
		"5": {},
		"4": {},
		"3": {},
		"2": {},
	}

	for lineNumber, line := range lines {
		if lineNumber == 0 {
			continue  // Skip header line.
		}
		i := lineNumber - 1

		// Remove the hint before the colon and the colon itself.
		// https://programming.guide/go/split-string-into-slice.html
		colonSeparated := strings.Split(line, ":")

		// These are the Sklansky Malmuth values.
		commaSeparated := strings.Split(colonSeparated[1], ",")
		//fmt.Println("                    ", commaSeparated[2])

		// Assign these values to the map.
		for j := 0; j < len(commaSeparated); j++ {
			num, err := strconv.Atoi(strings.TrimSpace(commaSeparated[j]))
			if err != nil {
				fmt.Printf("Error converting value from string to int: [%s%s] -> %s\n", ranks[i], ranks[j],
					commaSeparated[j])
				os.Exit(3)
			}
			//smMap[ranks[i]][ranks[j]] = num
			s.matrixMap[ranks[i]][ranks[j]] = num

		}
	}
}

// Given a two character string that represents the holecards, return
// the SM value in an int.
func (s ScoreModified) GetScoreOfString(hcString string) int {
	// https://stackoverflow.com/questions/15018545/how-to-index-characters-in-a-golang-string
	value1 := string([]rune(hcString)[0])
	value2 := string([]rune(hcString)[1])

	return s.matrixMap[value1][value2]
}

// This class cannot know of HoleCards because that would create an
// import cycle.  Instead, accept the HoleCards.ToString().  Inputs
// should look like "S4" or "HQ".
func (s ScoreModified) GetScoreOfHoleCardStrings(hc1, hc2 string) int {
	hc1suit := string([]rune(hc1)[0])
	hc1value := string([]rune(hc1)[1])
	hc2suit := string([]rune(hc2)[0])
	hc2value := string([]rune(hc2)[1])

	// If the hole cards are suited, then put the higher card first.
	if hc1suit == hc2suit {
		if isFirstGreater(hc1value, hc2value) {
			return s.GetScoreOfString(hc1value + hc2value)
		} else {
			return s.GetScoreOfString(hc2value + hc1value)
		}
	} else {
		if isFirstGreater(hc1value, hc2value) {
			//fmt.Println("Testing3:", hc2value + hc1value)
			return s.GetScoreOfString(hc2value + hc1value)
		} else {
			//fmt.Println("Testing4:", hc1value + hc2value)
			return s.GetScoreOfString(hc1value + hc2value)
		}
	}

	return 1999
}

func (s *ScoreModified) GetFilename() string {
	return s.filename
}

func (s *ScoreModified) SetFilename(fn string) {
	s.filename = fn
}

func GetSMModifiedScore() *ScoreModified {
	onceModified.Do(func() {
		singletonModified = &ScoreModified{}
		singletonModified.SetFilename("poker/matrix/holeCardValues_SklanskyMalmuthModified.txt")
		singletonModified.readInFile()
	})
	return singletonModified
}





package history

import (
	"sync"
	"os"
	"fmt"
	"bufio"
)

// This is a singleton that will lazily read the history.txt file.  That
// file will contain the record of every action played in every game in
// every tournament.  Ever.
// It will also contain:
// - all community cards dealt
// - all tournament results

// We'll see how this goes.

type History struct {
	filename string
}

var singleton *History
var once sync.Once

func (s *History) readInFile() {
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

	// Do something here.
	// Start by just writing to the file.
}

func GetHistory() *History {
	once.Do(func() {
		singleton = &History{}
		singleton.filename = "poker/history/history.txt"
		singleton.readInFile()
	})
	return singleton
}
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
	alreadyReadFile bool
}

var singleton *History
var once sync.Once

// We'll read the history file in lazily because it will eventually
// get big.
func (s *History) ReadInFile() {
	if s.alreadyReadFile {
		return
	}

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
	s.alreadyReadFile = true

	// Do something here.
	// Start by just writing to the file.
}

func (s History) String() string {
	returnedString := fmt.Sprintf("The history file is at %s", s.filename)

	if s.alreadyReadFile {
		returnedString += fmt.Sprintf(" and it has already been read.")
	} else {
		returnedString += fmt.Sprintf(" and it has not yet been read.")
	}

	return returnedString
}

func (s History) ToString() string {
	return s.String()
}

func (s *History) SetFilename(newFilename string) {
	s.filename = newFilename
}

func GetHistory() *History {
	once.Do(func() {
		singleton = &History{}
		singleton.filename = "poker/history/history.txt"
		singleton.alreadyReadFile = false
	})
	return singleton
}
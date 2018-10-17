package history

import (
	"sync"
	"os"
	"fmt"
	"bufio"
	"time"
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
func (h *History) ReadInFile() {
	if h.alreadyReadFile {
		return
	}

	fmt.Println("Oh snap; reading in:", h.filename)
	f, err := os.Open(h.filename)
	if err != nil {
		fmt.Println("Error opening hFile:", h.filename)
		os.Exit(1)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	h.alreadyReadFile = true

	// Do something here.
	// Start by just writing to the hFile.
}

func (h *History) Write(entryType string, kvps string) {
	// First check that we have a valid file handle.
	//if ! h.alreadyReadFile {
	//	h.ReadInFile()
	//}

	now := int32(time.Now().Unix())
	entryLine := fmt.Sprintf("%d,%s,%s\n", now, entryType, kvps)
	fmt.Printf("Writing to file: %s", entryLine)

	// For now, just append to the file.
	// TODO batch writes

	f, err := os.OpenFile(h.filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("FATAL: Could not open file.")
		panic(err)
	}

	// https://stackoverflow.com/questions/7151261/append-to-a-file-in-go
	_, err = f.WriteString(entryLine)
	if err != nil {
		fmt.Println("FATAL: Could not write to history file.")
		panic(err)
	}

}

func (h History) String() string {
	returnedString := fmt.Sprintf("The history file is at %h", h.filename)

	if h.alreadyReadFile {
		returnedString += fmt.Sprintf(" and it has already been read.")
	} else {
		returnedString += fmt.Sprintf(" and it has not yet been read.")
	}

	return returnedString
}

func (h History) ToString() string {
	return h.String()
}

func (h *History) SetFilename(newFilename string) {
	h.filename = newFilename
}

func GetHistory() *History {
	once.Do(func() {
		singleton = &History{}
		singleton.filename = "poker/history/history.txt"
		singleton.alreadyReadFile = false
	})
	return singleton
}
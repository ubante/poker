package main

import (
	"fmt"
	"goven/poker/models"
	"goven/poker/history"
)

// Just a stub to test the history.
//
// The history file will be CSV format with these fields:
// datetime,entry_type,KVP
//
// datetime: epoch time to seconds
// entry_type: supported are (tournament)
// KVP: key value pairs separated by semicolons
//
// Examples:

// TODO learn how to do unit tests in golang

func main() {
	fmt.Println("Testing the history file.")

	allHistory := history.GetHistory()
	fmt.Println("1. This is the file:", allHistory)

	// Prove singletoniness.
	theOtherHistory := history.GetHistory()
	fmt.Println("2. This should be the same file:", theOtherHistory)
	//allHistory.SetFilename("a/b/c/d")
	//fmt.Println("1. After changing the filepath of allHistory:", allHistory)
	//fmt.Println("2. Did this change too?", theOtherHistory)  // Yes, it did.

	// Read in the history file exactly once.
	allHistory.ReadInFile()
	fmt.Println("1.", allHistory)
	fmt.Println("2.", theOtherHistory)

	// Right now, a player joins a tournament when the player is added
	// to the Table().  This could be cleaned up in the future with
	// making a queue of players interest in playing a tournament.
	var table models.Table
	table.Initialize()

	temp := models.NewAllInAlwaysPlayer("Alex")
	table.AddPlayerWithHistory(&temp)

	fmt.Println(table.GetStatus())

	allHistory.Write("tournament", "start=x;stop=y")
	tid := history.GetUniqueId()
	allHistory.Write("tournament", fmt.Sprintf("name=%s;state=starting", tid))




}

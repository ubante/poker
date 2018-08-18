package main

import (
	"sync"
	"fmt"
)

// This is similar to python's itertools.combinations()
// https://play.golang.org/p/KS6TcIafc3


func generate(alphabet string) <-chan string {
	c := make(chan string, len(alphabet))

	go func() {
		defer close(c)

		if len(alphabet) == 0 {
			return
		}

		// Use a sync.WaitGroup to spawn permutation
		// goroutines and allow us to wait for them all
		// to finish.
		var wg sync.WaitGroup
		wg.Add(len(alphabet))

		for i := 1; i <= len(alphabet); i++ {
			go func(i int) {
				// Perform permutation.
				Word(alphabet[:i]).Permute(c)

				// Signal Waitgroup that we are done.
				wg.Done()
			}(i)
		}

		// Wait for all routines to finish.
		wg.Wait()
	}()

	return c
}

type Word []rune

// Permute generates all possible combinations of
// the current word. This assumes the runes are sorted.
func (w Word) Permute(out chan<- string) {
	if len(w) <= 1 {
		out <- string(w)
		return
	}

	// Write first result manually.
	out <- string(w)

	// Find and print all remaining permutations.
	for w.next() {
		out <- string(w)
	}
}

// next performs a single permutation by shuffling characters around.
// Returns false if there are no more new permutations.
func (w Word) next() bool {
	var left, right int

	left = len(w) - 2
	for w[left] >= w[left+1] && left >= 1 {
		left--
	}

	if left == 0 && w[left] >= w[left+1] {
		return false
	}

	right = len(w) - 1
	for w[left] >= w[right] {
		right--
	}

	w[left], w[right] = w[right], w[left]

	left++
	right = len(w) - 1

	for left < right {
		w[left], w[right] = w[right], w[left]
		left++
		right--
	}

	return true
}

// https://www.geeksforgeeks.org/print-all-possible-combinations-of-r-elements-in-a-given-array-of-size-n/
func makePossibleHands(cards []string) [][]string {
	var possible [][]string
	ctr := 0

	for a := 0; a < len(cards)-4; a++ {
		for b := a+1; b < len(cards)-3; b++ {
			for c := b+1; c < len(cards)-2; c++ {
				for d := c+1; d < len(cards)-1; d++ {
					for e := d+1; e < len(cards); e++ {
						fmt.Println(cards[a], cards[b], cards[c], cards[d], cards[e])
						possible = append(possible, []string{cards[a], cards[b], cards[c], cards[d], cards[e]})
						ctr++
					}
				}
			}
		}
	}
	fmt.Println(ctr, "different combinations.")
	return possible
}

func main() {
	// Make sure the alphabet is sorted.
	const alphabet = "abcde"

	//for str := range generate(alphabet) {
	//	fmt.Println(str)
	//}

	cards := []string{"a1", "b2", "c3", "e4", "d5", "f6", "g7"}
	fmt.Println(cards)

	fmt.Println("==========================")
	possibleHands := makePossibleHands(cards)
	for i, hand := range possibleHands {
		fmt.Println(i, hand)
	}
}

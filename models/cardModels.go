package models

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

type Card struct {
	Suit          string
	NumericalRank int
	Rank          string
}

func (c *Card) ToString() string {
	return fmt.Sprintf("%s%s", c.Suit, c.Rank)
}

func (c Card) String() string {
	return c.ToString()
}

func (c Card) IsSuited(c2 Card) bool {
	if c.Suit == c2.Suit {
		return true
	}

	return false
}

func (c Card) IsPaired(c2 Card) bool {
	if c.NumericalRank == c2.NumericalRank {
		return true
	}

	return false
}

type CardSet struct {
	//cards         []*models.Card
	cards         []*Card
	bestHand      *CardSet
	bestEval      *Evaluation
	possibleHands []*CardSet
}

// I should rename ToString to String and explicitly call
// SomeType.String() when I want a string outside of Println()
func (cs CardSet) ToString() string {
	var toString string

	for _, c := range cs.cards {
		if toString == "" {

			toString = c.ToString()
			continue
		}

		toString += fmt.Sprintf(" %s", c)
	}

	return toString
}

func (cs CardSet) String() string {
	return cs.ToString()
}

//func (cs *CardSet) Add(c models.Card) {
func (cs *CardSet) Add(c Card) {
	cs.cards = append(cs.cards, &c)
}

/*
This will accept a cardset and return the union of it with this.
*/
func (cs *CardSet) Combine(cs2 CardSet) CardSet {
	var combined CardSet

	for _, card := range cs.cards {
		combined.Add(*card)
	}
	for _, card := range cs2.cards {
		combined.Add(*card)
	}

	return combined
}

// Find all the possible combinations.
// I like it gross.
func (cs *CardSet) SetPossibleHands() {
	cards := cs.cards
	for a := 0; a < len(cards)-4; a++ {
		for b := a + 1; b < len(cards)-3; b++ {
			for c := b + 1; c < len(cards)-2; c++ {
				for d := c + 1; d < len(cards)-1; d++ {
					for e := d + 1; e < len(cards); e++ {
						possibleHand := NewCardSet(*cards[a], *cards[b], *cards[c], *cards[d], *cards[e])
						cs.possibleHands = append(cs.possibleHands, &possibleHand)
					}
				}
			}
		}
	}
}

func (cs *CardSet) FindBestHand() {
	cs.SetPossibleHands()
	for _, ph := range cs.possibleHands {
		eval := NewEvaluation(*ph)
		//fmt.Println(eval)

		if cs.bestHand == nil {
			cs.bestHand = ph
			cs.bestEval = eval
			//fmt.Println("Here's the initial bestHand:", cs.bestHand)
			continue
		}

		// Maybe stash this eval so no need to recompute it?
		currentEval := NewEvaluation(*cs.bestHand)
		if eval.IsBetterThan(*currentEval) {
			cs.bestHand = ph
			cs.bestEval = eval
		}
	}
}

/*
I couldn't find a Pop function - weird.
https://groups.google.com/forum/#!topic/Golang-nuts/obZI4uyZTe0
*/
//func (cs *CardSet) Pop() *models.Card {
func (cs *CardSet) Pop() *Card {
	card := cs.cards[0]
	copy(cs.cards, cs.cards[1:])
	cs.cards = cs.cards[:len(cs.cards)-1]

	return card
}

func (cs *CardSet) Shuffle() {
	rand.Shuffle(cs.Length(), func(i, j int) {
		cs.cards[i], cs.cards[j] = cs.cards[j], cs.cards[i]
	})
}

func (cs CardSet) Length() int {
	return len(cs.cards)
}

func (cs CardSet) GetReverseOrderedNumericRanks() []int {
	var orderedRanks []int
	for _, card := range cs.cards {
		orderedRanks = append(orderedRanks, card.NumericalRank)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(orderedRanks)))

	return orderedRanks
}

/*
Primary ranks are:
  9: Straight flush
  8: Four of a kind
  7: Full house
  6: Flush
  5: Straight
  4: Three of a kind
  3: Two pair
  2: One pair
  1: High card

Secondary ranks vary with the primary rank.  For example,  for
a straight, the secondary rank will be the highest card.  For
a full house, the secondary rank will be the numeric value of
the three of a kind.

Tertiary rank also vary's with the primary rank.  For a flush,
this would be the second highest card.  For a full house, this
would be the numeric value of the pair.

This just goes on and on.  For a flush and high card, it is
possible to have a rank for each of the five cards.  For the
other types of hands, these will remain at zero.

It is possible for different hands to have the same evaluation.
*/
type Evaluation struct {
	cardSet        *CardSet
	humanEval      string
	allRanks       [6]*int
	primaryRank    int
	secondaryRank  int
	tertiaryRank   int
	quaternaryRank int
	quinaryRank    int
	senaryRank     int
	flattenedScore int // A simple way to score the hand.  Higher is better.
}

func (e Evaluation) String() string {
	toString := fmt.Sprintf("%s: %s with ranks: [%d %d %d %d %d %d]", e.cardSet, e.humanEval,
		e.primaryRank, e.secondaryRank, e.tertiaryRank, e.quaternaryRank, e.quinaryRank, e.senaryRank)

	return toString
}

func (e Evaluation) isFlush() bool {
	var suit string
	for _, card := range e.cardSet.cards {
		if suit == "" {
			suit = card.Suit
		} else {
			if card.Suit != suit {
				return false
			}
		}
	}

	return true
}

func (e Evaluation) isStraight() bool {
	// This needs to handle a wheel straight, ie A2345.

	orderedRanks := e.cardSet.GetReverseOrderedNumericRanks()
	previous := 0
	for _, rank := range orderedRanks {
		if previous == 0 {
			previous = rank
		} else {
			if previous-rank != 1 {
				return false
			}
			previous = rank
		}
	}

	return true
}

// This will return a map of ints -> ints.
// The key will be the size of the match, ie 4 means quads and 1 means
// unmatched cards.  The value will be the numeric rank(s) of that size.
// For the numeric rank 2, there could be two values, ie two pairs.
//
// If cards have ranks 13 13 11 11 7, then the returned map will be:
//   2 -> [13, 11]  // a pair of kings and jacks
//   1 -> [7]       // a lone 7
func (e Evaluation) hasMatches() map[int][]int {
	frequency := make(map[int]int)
	for _, card := range e.cardSet.cards {
		frequency[card.NumericalRank]++
	}

	// We want the keys to the frequency map to be added to the matches
	// map in descending order.
	var orderedKeys []int
	for k := range frequency {
		orderedKeys = append(orderedKeys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(orderedKeys)))

	matches := make(map[int][]int)
	for _, rank := range orderedKeys {
		matches[frequency[rank]] = append(matches[frequency[rank]], rank)
	}

	return matches
}

// This will flatten the allRanks array into an int preserving
// ordinality. This could be more compact by using powers of 13 but then
// I'd have to cast to float64 to use Pow.  smh
func (e *Evaluation) flattenScore() {
	e.flattenedScore = 0
	for _, subScore := range e.allRanks {
		e.flattenedScore *= 100
		e.flattenedScore += *subScore // this line panics
	}
}

// This is because defer can call functions but cannot call methods.
// Yeah, I did it.  Wot?
//func notAMethodFlattenScore(e *Evaluation) {
//	fmt.Println("flattening:", e)
//	e.flattenScore()
//}

func (e *Evaluation) evaluate() {
	//defer notAMethodFlattenScore(e)
	e.humanEval = "TBDeval"

	if e.isStraight() && e.isFlush() {
		e.humanEval = "staight flush"
		e.primaryRank = 9
		highestRank := e.cardSet.GetReverseOrderedNumericRanks()[0]
		if highestRank == 14 {
			e.humanEval = "royal flush"
		}
		e.secondaryRank = highestRank
		return
	}

	// https://stackoverflow.com/questions/2050391/how-to-check-if-a-map-contains-a-key-in-go
	allMatches := e.hasMatches()
	if rank, ok := allMatches[4]; ok {
		e.humanEval = "quads"
		e.primaryRank = 8
		e.secondaryRank = rank[0]
		e.tertiaryRank = allMatches[1][0]
		return
	}

	// Guess golang if's only allow one initialization statement?
	//if rank3, ok3 := allMatches[3]; rank2, ok2 := allMatches[2]; ok2 {
	var trips, firstPair, secondPair int
	var ok3, ok2 bool
	if _, ok := allMatches[3]; ok {
		ok3 = true
		trips = allMatches[3][0]
	}
	if _, ok := allMatches[2]; ok {
		ok2 = true
		firstPair = allMatches[2][0]
		if len(allMatches[2]) == 2 {
			secondPair = allMatches[2][1]
		}
	}

	if ok3 && ok2 {
		e.humanEval = "full house"
		e.primaryRank = 7
		e.secondaryRank = trips
		e.tertiaryRank = firstPair
		return
	}

	if e.isFlush() {
		e.humanEval = "flush"
		e.primaryRank = 6
		orderedRanks := e.cardSet.GetReverseOrderedNumericRanks()
		e.secondaryRank = orderedRanks[0]
		e.tertiaryRank = orderedRanks[1]
		e.quaternaryRank = orderedRanks[2]
		e.quinaryRank = orderedRanks[3]
		e.senaryRank = orderedRanks[4]
		return
	}

	if e.isStraight() {
		e.humanEval = "straight"
		e.primaryRank = 5
		e.secondaryRank = e.cardSet.GetReverseOrderedNumericRanks()[0]
		return
	}

	// At this point, we're guaranteed to have unmatched cards.
	singles := allMatches[1]
	sort.Sort(sort.Reverse(sort.IntSlice(singles)))

	if trips != 0 {
		e.humanEval = "three of a kind"
		e.primaryRank = 4
		e.secondaryRank = trips
		e.tertiaryRank = singles[0]
		e.quinaryRank = singles[1]
		return
	}

	if firstPair*secondPair != 0 {
		e.humanEval = "two pairs"
		e.primaryRank = 3
		e.secondaryRank = firstPair
		e.tertiaryRank = secondPair
		e.quaternaryRank = singles[0]
		return
	}

	if firstPair != 0 {
		e.humanEval = "pair"
		e.primaryRank = 2
		e.secondaryRank = firstPair
		e.tertiaryRank = singles[0]
		e.quaternaryRank = singles[1]
		e.quinaryRank = singles[2]
		return
	}

	// Don't laugh.
	e.humanEval = "high card"
	e.primaryRank = 1
	e.secondaryRank = singles[0]
	e.tertiaryRank = singles[1]
	e.quaternaryRank = singles[2]
	e.quinaryRank = singles[3]
	e.senaryRank = singles[4]
}

// This will return 1 if this Evaluation is greater than the given
// Evaluation and will return -1 if this Evaluation is lesser than
// the other Evaluation.  If they are even, then this will return 0.
// Maybe use an enum?
//
// Greater values in rank means betterness.
func (e Evaluation) Compare(otherEval Evaluation) int {

	for i := 0; i < 6; i++ {
		//fmt.Println(i, "this:", e.allRanks[i], ", other:", otherEval.allRanks[i])
		if *e.allRanks[i] < *otherEval.allRanks[i] {
			return -1
		}
		if *e.allRanks[i] > *otherEval.allRanks[i] {
			return 1
		}
	}

	return 0
}

// Courtesy method.  Note that comparing two Evaluations has three possible
// results.
func (e Evaluation) IsBetterThan(otherEval Evaluation) bool {
	results := e.Compare(otherEval)

	if results == 1 {
		return true
	}

	return false
}

type HoleCards struct {
	cardSet *CardSet
}

func (hc *HoleCards) ToString() string {
	if hc.cardSet == nil {
		return ""
	}

	return hc.cardSet.ToString()
}

func (hc HoleCards) String() string {
	return hc.ToString()
}

func (hc *HoleCards) Add(c Card) {
	hc.cardSet.Add(c)
}

func (hc *HoleCards) Empty() {
	ecs := NewCardSet()
	hc.cardSet = &ecs
}

func (hc *HoleCards) toss() {
	hc.Empty()
}

func (hc HoleCards) IsSuited() bool {
	if hc.cardSet.cards[0].Suit == hc.cardSet.cards[1].Suit {
		return true
	}

	return false
}

// https://en.wikipedia.org/wiki/Texas_hold_%27em_starting_hands#Sklansky_hand_groups
func (hc HoleCards) getSklanskyMalmuthScore() int {
	var smScore int
	smScore = 10  // Just for stubbing.

	// Read in the txt file and make a map of maps.
	fileMap := make(map[string]map[string]int)
	fileMap["K"]["K"] = 1

	return smScore
}

type Deck struct {
	cardSet *CardSet // Deck is not a CardSet.
}

func (d *Deck) getStatus() string {
	status := fmt.Sprintf("Deck has %d cards:\n", len(d.cardSet.cards))
	for i, card := range d.cardSet.cards {
		if i != 0 && i%13 == 0 {
			status += "\n"
		}
		status += card.ToString()
		status += " "
	}

	return status
}

func (d Deck) length() int {
	return d.cardSet.Length()
}

func (d *Deck) shuffle() {
	d.cardSet.Shuffle()
}

func (d *Deck) getCard() *Card {
	return d.cardSet.Pop()
}

type Community struct {
	cards *CardSet
}

// Will eventually replace all the getStatus() methods.
func (c Community) String() string {
	return c.cards.String()
}

func (c *Community) add(card Card) {
	c.cards.Add(card)
}

func NewCard(s string, nr int) Card {
	var c Card
	c.Suit = s
	c.NumericalRank = nr

	//if nr == 14 {
	//	c.Rank = "A" // Aces are aces.
	//}
	//
	switch nr {
	case 14:
		c.Rank = "A"  // Aces are aces.
	case 13:
		c.Rank = "K"
	case 12:
		c.Rank = "Q"
	case 11:
		c.Rank = "J"
	case 10:
		c.Rank = "T"
	default:
		c.Rank = strconv.Itoa(nr)
	}

	return c
}

func NewCardSet(cards ...Card) CardSet {
	var cs CardSet

	for _, card := range cards {
		cs.Add(card)
	}

	return cs
}

func NewEvaluation(cardSet CardSet) *Evaluation {
	var eval Evaluation
	eval.cardSet = &cardSet
	eval.evaluate()

	// How to do the below in one line?
	eval.allRanks[0] = &eval.primaryRank
	eval.allRanks[1] = &eval.secondaryRank
	eval.allRanks[2] = &eval.tertiaryRank
	eval.allRanks[3] = &eval.quaternaryRank
	eval.allRanks[4] = &eval.quinaryRank
	eval.allRanks[5] = &eval.senaryRank
	eval.flattenScore()

	return &eval
}

func NewCommunity() Community {
	ecs := NewCardSet()

	var c Community
	c.cards = &ecs

	return c
}

func NewDeck() *Deck {
	var d Deck
	ecs := NewCardSet()
	d.cardSet = &ecs

	for _, suit := range []string{"S", "H", "D", "C"} {
		// https://stackoverflow.com/questions/21950244/is-there-a-way-to-iterate-over-a-range-of-integers-in-golang
		for numericRank := range [13]int{} {
			newCard := NewCard(suit, numericRank+2)
			d.cardSet.Add(newCard)
		}
	}

	return &d
}

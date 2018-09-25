package matrix

import "sync"

// https://medium.com/@MrToBe/the-singleton-object-oriented-design-pattern-in-golang-9f6ce75c21f7
// http://marcio.io/2015/07/singleton-pattern-in-go/
type Score struct {
	filename string
	//matrixMap map[string]map[string]int
}

var singleton *Score
var once sync.Once

func (s *Score) GetFilename() string {
	return s.filename
}

func (s *Score) SetFilename(fn string) {
	s.filename = fn
}

func GetScore() *Score {
	once.Do(func() {
		singleton = &Score{"default_filename"}
	})
	return singleton
}

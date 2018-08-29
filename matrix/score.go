package matrix

// https://medium.com/@MrToBe/the-singleton-object-oriented-design-pattern-in-golang-9f6ce75c21f7
// http://marcio.io/2015/07/singleton-pattern-in-go/
type Score struct {
	filename string
	matrixMap map[string]map[string]int
}

package core

type Operator struct {
	Symbol string
	Line   int
	CharAt int
}

var precedenceLevels = map[string]int{
	// ".":   1,
	"*":   2,
	"/":   2,
	"%":   2,
	"+":   3,
	"-":   3,
	"<":   4,
	"<=":  4,
	">":   4,
	">=":  4,
	"==":  5,
	"===": 5,
	"!=":  5,
	"!==": 5,
	"&&":  6,
	"||":  7,
}

func IsOperatorSymbol(s string) bool {
	_, ok := precedenceLevels[s]
	return ok
}

func (o Operator) Compare(oo Operator) int {
	if precedenceLevels[o.Symbol] > precedenceLevels[oo.Symbol] {
		return 1
	} else if precedenceLevels[o.Symbol] < precedenceLevels[oo.Symbol] {
		return -1
	}
	return 0
}

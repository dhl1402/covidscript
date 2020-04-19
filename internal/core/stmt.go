package core

type Statement interface {
	Execute(*ExecutionContext) (Expression, error)
}

type Identifier struct {
	Name   string
	Line   int
	CharAt int
}

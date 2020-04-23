package core

type Statement interface {
	Execute(*ExecutionContext) (Expression, error)
	Clone() Statement
}

type Identifier struct {
	Name   string
	Line   int
	CharAt int
}

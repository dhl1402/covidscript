package core

type Expression interface {
	Evaluate(*ExecutionContext) (Expression, error)
	GetLine() int
	GetCharAt() int
	SetLine(int)
	SetCharAt(int)
	GetType() string
	ToString() string
}

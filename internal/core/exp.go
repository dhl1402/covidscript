package core

type Expression interface {
	Evaluate(*ExecutionContext) (Expression, error)
	IsTruthy() bool
	GetLine() int
	GetCharAt() int
	SetLine(int)
	SetCharAt(int)
	GetType() string
	ToString() string
	Clone() Expression
}

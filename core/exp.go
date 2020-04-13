package core

type Expression interface {
	Evaluate(ExecutionContext) (Expression, error)
	GetCharAt() int
}

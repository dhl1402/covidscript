package core

type ArrayExpression struct {
	Elements []Expression
	Line     int
	CharAt   int
}

func (e ArrayExpression) Evaluate(ec ExecutionContext) (Expression, error) {
	return e, nil
}

func (e ArrayExpression) GetCharAt() int {
	return e.CharAt
}

package core

type FunctionExpression struct {
	Params []Identifier
	Body   []Statement
	Line   int
	CharAt int
}

func (e FunctionExpression) Evaluate(ec ExecutionContext) (Expression, error) {
	return e, nil
}

func (e FunctionExpression) GetCharAt() int {
	return e.CharAt
}

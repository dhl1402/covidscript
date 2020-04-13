package core

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
	Line      int
	CharAt    int
}

func (e CallExpression) Evaluate(ec ExecutionContext) (Expression, error) {
	return e, nil
}

func (e CallExpression) GetCharAt() int {
	return e.CharAt
}

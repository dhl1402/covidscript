package core

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
	Line      int
	CharAt    int
}

func (e *CallExpression) Evaluate(ec ExecutionContext) (Expression, error) {
	return e, nil
}

func (e *CallExpression) GetCharAt() int {
	return e.CharAt
}

func (e *CallExpression) GetLine() int {
	return e.Line
}

func (e *CallExpression) SetLine(i int) {
	e.Line = i
}

func (e *CallExpression) SetCharAt(i int) {
	e.CharAt = i
}

func (e *CallExpression) GetType() string {
	return "call expression"
}

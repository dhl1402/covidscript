package core

type LiteralExpression struct {
	Type   string
	Value  string
	Line   int
	CharAt int
}

func (e LiteralExpression) Evaluate(ec ExecutionContext) (Expression, error) {
	return e, nil
}

func (e LiteralExpression) GetCharAt() int {
	return e.CharAt
}

func (e LiteralExpression) GetLine() int {
	return e.Line
}

func (e LiteralExpression) GetType() string {
	return e.Type
}

package core

type FunctionExpression struct {
	Params []Identifier
	Body   []Statement
	EC     *ExecutionContext
	Line   int
	CharAt int
}

func (e *FunctionExpression) Evaluate(ec *ExecutionContext) (Expression, error) {
	e.EC.Variables = map[string]Expression{}
	for _, p := range e.Params {
		e.EC.Set(p.Name, &LiteralExpression{
			Type:   "undefined",
			Line:   p.Line,
			CharAt: p.CharAt,
		})
	}
	return e, nil
}

func (e *FunctionExpression) GetCharAt() int {
	return e.CharAt
}

func (e *FunctionExpression) GetLine() int {
	return e.Line
}

func (e *FunctionExpression) SetLine(i int) {
	e.Line = i
}

func (e *FunctionExpression) SetCharAt(i int) {
	e.CharAt = i
}

func (e *FunctionExpression) GetType() string {
	return "function"
}

package core

type ArrayExpression struct {
	Elements []Expression
	Line     int
	CharAt   int
}

func (e ArrayExpression) Evaluate(ec ExecutionContext) (Expression, error) {
	elems := []Expression{}
	for _, ee := range e.Elements {
		exp, err := ee.Evaluate(ec)
		if err != nil {
			return nil, err
		}
		elems = append(elems, exp)
	}
	e.Elements = elems
	return e, nil
}

func (e ArrayExpression) GetCharAt() int {
	return e.CharAt
}

func (e ArrayExpression) GetLine() int {
	return e.Line
}

func (e ArrayExpression) GetType() string {
	return "array"
}

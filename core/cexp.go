package core

import "fmt"

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
	Line      int
	CharAt    int
}

func (e *CallExpression) Evaluate(ec ExecutionContext) (Expression, error) {
	callee, err := e.Callee.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	f, ok := callee.(*FunctionExpression)
	if !ok {
		return nil, fmt.Errorf("a is not a function. [%d,%d]", e.Line, e.CharAt) // TODO: e.Callee.ToString()
	}
	for i, p := range f.Params {
		if i < len(e.Arguments) {
			arg, err := e.Arguments[i].Evaluate(ec)
			if err != nil {
				return nil, err
			}
			f.EC.Set(p.Name, arg)
		}
	}
	// TODO: execute function
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

package core

import "fmt"

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
	Line      int
	CharAt    int
}

func (e *CallExpression) Evaluate(ec *ExecutionContext) (Expression, error) {
	callee, err := e.Callee.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	f, ok := callee.(*FunctionExpression)
	if !ok {
		return nil, fmt.Errorf("a is not a function. [%d,%d]", e.Line, e.CharAt) // TODO: e.Callee.ToString()
	}

	for i, argexp := range e.Arguments {
		arg, err := argexp.Evaluate(ec)
		if err != nil {
			return nil, err
		}
		if i < len(f.Params) {
			f.EC.Set(f.Params[i].Name, arg)
		}
		f.EC.Set(fmt.Sprintf("_args%d_", i), arg)
	}
	if f.NativeFunction != nil {
		return f.NativeFunction(f.EC)
	}
	for _, stmt := range f.Body {
		rexp, err := stmt.Execute(f.EC)
		if rexp != nil || err != nil {
			return rexp, err
		}
	}
	return &LiteralExpression{
		Type:   LiteralTypeUndefined,
		Line:   e.Line,
		CharAt: e.CharAt,
	}, nil
}

func (e *CallExpression) IsTruthy() bool {
	return true
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

func (e *CallExpression) ToString() string {
	return ""
}

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
	for i, p := range f.Params {
		if i < len(e.Arguments) {
			arg, err := e.Arguments[i].Evaluate(ec)
			if err != nil {
				return nil, err
			}
			f.EC.Set(p.Name, arg)
		}
	}
	return Execute(f.EC, f.Body)
}

func Execute(ec *ExecutionContext, statement []Statement) (Expression, error) {
	for _, ss := range statement {
		switch s := ss.(type) {
		case VariableDeclaration:
			for _, d := range s.Declarations {
				if d.Init != nil {
					if f, ok := d.Init.(*FunctionExpression); ok {
						f.EC = &ExecutionContext{
							Outer:     ec,
							Variables: map[string]Expression{},
						}
					}
					value, err := d.Init.Evaluate(ec)
					if err != nil {
						return nil, err
					}
					ec.Set(d.ID.Name, value)
				} else {
					ec.Set(d.ID.Name, &LiteralExpression{
						Type:   "undefined",
						Line:   d.ID.Line,
						CharAt: d.ID.CharAt,
					})
				}
			}
		case ReturnStatement:
			return s.Argument.Evaluate(ec)
		}
	}
	return nil, nil
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

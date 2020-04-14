package core

import (
	"fmt"
)

type VariableExpression struct {
	Name   string
	Line   int
	CharAt int
}

func (e *VariableExpression) Evaluate(ec *ExecutionContext) (Expression, error) {
	if exp, ok := ec.Get(e.Name); ok {
		result, err := exp.Evaluate(ec)
		if err != nil {
			return nil, err
		}
		result.SetLine(e.Line)
		result.SetCharAt(e.CharAt)
		return result, nil
	}
	return nil, fmt.Errorf("%s is not defined. [%d,%d]", e.Name, e.Line, e.CharAt)
}

func (e *VariableExpression) GetCharAt() int {
	return e.CharAt
}

func (e *VariableExpression) GetLine() int {
	return e.Line
}

func (e *VariableExpression) SetLine(i int) {
	e.Line = i
}

func (e *VariableExpression) SetCharAt(i int) {
	e.CharAt = i
}

func (e *VariableExpression) GetType() string {
	return "variable"
}

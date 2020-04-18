package core

import (
	"fmt"

	"github.com/dhl1402/covidscript/utils"
)

type UnaryExpression struct {
	Expression
	Line   int
	CharAt int
}

func (e *UnaryExpression) Evaluate(ec *ExecutionContext) (Expression, error) {
	exp, err := e.Expression.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	return &LiteralExpression{
		Type:   LiteralTypeBoolean,
		Value:  utils.ToBoolStr(exp.IsTruthy()),
		Line:   e.Line,
		CharAt: e.CharAt,
	}, nil
}

func (e *UnaryExpression) IsTruthy() bool {
	return true
}

func (e *UnaryExpression) GetCharAt() int {
	return e.CharAt
}

func (e *UnaryExpression) GetLine() int {
	return e.Line
}

func (e *UnaryExpression) SetLine(i int) {
	e.Line = i
}

func (e *UnaryExpression) SetCharAt(i int) {
	e.CharAt = i
}

func (e *UnaryExpression) GetType() string {
	return "unary expression"
}

func (e *UnaryExpression) ToString() string {
	return fmt.Sprintf("!%s", e.Expression.ToString())
}

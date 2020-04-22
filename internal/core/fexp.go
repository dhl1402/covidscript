package core

import "fmt"

type FunctionExpression struct {
	Params         []Identifier
	Body           []Statement
	NativeFunction func(*ExecutionContext) (Expression, error)
	EC             *ExecutionContext
	Line           int
	CharAt         int
}

func (e *FunctionExpression) Evaluate(ec *ExecutionContext) (Expression, error) {
	if e.EC == nil {
		e.EC = &ExecutionContext{
			Outer:     ec,
			Variables: map[string]Expression{},
		}
	}
	return e, nil
}

func (e *FunctionExpression) IsTruthy() bool {
	return true
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

func (e *FunctionExpression) ToString() string {
	params := ""
	for _, p := range e.Params {
		params = params + p.Name + ", "
	}
	if len(params) > 1 {
		params = params[:len(params)-2]
		return fmt.Sprintf("func(%s)", params)
	}
	return "func()"
}

package core

import (
	"fmt"
	"strconv"
)

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator Operator
	Group    bool
	Line     int
	CharAt   int
}

func (e *BinaryExpression) Evaluate(ec *ExecutionContext) (Expression, error) {
	left, err := e.Left.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	lle, ok := left.(*LiteralExpression)
	if !ok {
		return nil, fmt.Errorf("Cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, left.GetType(), e.Operator.Line, e.Operator.CharAt)
	}
	right, err := e.Right.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	rle, ok := right.(*LiteralExpression)
	if !ok {
		return nil, fmt.Errorf("Cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, right.GetType(), e.Operator.Line, e.Operator.CharAt)
	}
	if e.Operator.Symbol == "-" || e.Operator.Symbol == "*" || e.Operator.Symbol == "/" || e.Operator.Symbol == "%" {
		if lle.Type != LiteralTypeNumber {
			return nil, fmt.Errorf("Cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, lle.GetType(), e.Operator.Line, e.Operator.CharAt)
		}
		if rle.Type != LiteralTypeNumber {
			return nil, fmt.Errorf("Cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, rle.GetType(), e.Operator.Line, e.Operator.CharAt)
		}
	}
	// TODO: handle float
	li, _ := strconv.Atoi(lle.Value)
	ri, _ := strconv.Atoi(rle.Value)
	switch e.Operator.Symbol {
	case "+":
		if lle.Type == LiteralTypeUndefined || lle.Type == LiteralTypeNull {
			return nil, fmt.Errorf("Cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, lle.GetType(), e.Operator.Line, e.Operator.CharAt)
		}
		if rle.Type == LiteralTypeUndefined || rle.Type == LiteralTypeNull {
			return nil, fmt.Errorf("Cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, rle.GetType(), e.Operator.Line, e.Operator.CharAt)
		}
		if lle.Type == LiteralTypeNumber && rle.Type == LiteralTypeNumber {
			// TODO: handle float
			return &LiteralExpression{
				Type:   LiteralTypeNumber,
				Value:  fmt.Sprintf("%d", li+ri),
				Line:   e.Line,
				CharAt: e.CharAt,
			}, nil
		}
		return &LiteralExpression{
			Type:   LiteralTypeString,
			Value:  lle.Value + rle.Value,
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case "-":
		return &LiteralExpression{
			Type:   LiteralTypeNumber,
			Value:  fmt.Sprintf("%d", li-ri),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case "*":
		return &LiteralExpression{
			Type:   LiteralTypeNumber,
			Value:  fmt.Sprintf("%d", li*ri),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case "/":
		if ri == 0 {
			return nil, fmt.Errorf("Cannot divide by zero. [%d,%d]", rle.Line, rle.CharAt)
		}
		return &LiteralExpression{
			Type:   LiteralTypeNumber,
			Value:  fmt.Sprintf("%d", li/ri),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case "%":
		return &LiteralExpression{
			Type:   LiteralTypeNumber,
			Value:  fmt.Sprintf("%d", li%ri),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
		// TODO: handle logic operator
	}
	return nil, fmt.Errorf("Operator %s is not supported. [%d,%d]", e.Operator.Symbol, e.Operator.Line, e.Operator.CharAt)
}

func (e *BinaryExpression) GetCharAt() int {
	return e.CharAt
}

func (e *BinaryExpression) GetLine() int {
	return e.Line
}

func (e *BinaryExpression) SetLine(i int) {
	e.Line = i
}

func (e *BinaryExpression) SetCharAt(i int) {
	e.CharAt = i
}

func (e *BinaryExpression) GetType() string {
	return "expression"
}

func (e *BinaryExpression) ToString() string {
	return fmt.Sprintf("%s %s %s", e.Left.ToString(), e.Operator.Symbol, e.Right.ToString())
}

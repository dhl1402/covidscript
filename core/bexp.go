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

func (e BinaryExpression) Evaluate(ec ExecutionContext) (Expression, error) {
	left, err := e.Left.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	lle, ok := left.(*LiteralExpression)
	if !ok {
		return nil, fmt.Errorf("Can't do operator on a function") // TODO: wording + left.GetTypeName()
	}
	right, err := e.Left.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	rle, ok := right.(*LiteralExpression)
	if !ok {
		return nil, fmt.Errorf("Can't do operator on a function") // TODO: wording + left.GetTypeName()
	}
	li, _ := strconv.Atoi(lle.Value)
	ri, _ := strconv.Atoi(lle.Value)
	switch e.Operator.Symbol {
	case "+":
		if lle.Type == "number" && rle.Type == "number" {
			// TODO: handle float
			return &LiteralExpression{
				Type:   "number",
				Value:  fmt.Sprintf("%d", li+ri),
				Line:   lle.Line,
				CharAt: lle.CharAt,
			}, nil
		}
		return &LiteralExpression{
			Type:   "stringr",
			Value:  fmt.Sprintf("%s%s", lle.Value, rle.Value),
			Line:   lle.Line,
			CharAt: lle.CharAt,
		}, nil
	case "-":
		if lle.Type != "number" || rle.Type != "number" {
			return nil, fmt.Errorf("Can't do operator on a function") // TODO: wording + left.GetTypeName()
		}
		return &LiteralExpression{
			Type:   "number",
			Value:  fmt.Sprintf("%d", li-ri),
			Line:   lle.Line,
			CharAt: lle.CharAt,
		}, nil
	case "*":
		if lle.Type != "number" || rle.Type != "number" {
			return nil, fmt.Errorf("Can't do operator on a function") // TODO: wording + left.GetTypeName()
		}
		return &LiteralExpression{
			Type:   "number",
			Value:  fmt.Sprintf("%d", li*ri),
			Line:   lle.Line,
			CharAt: lle.CharAt,
		}, nil
	case "/":
		if lle.Type != "number" || rle.Type != "number" {
			return nil, fmt.Errorf("Can't do operator on a function") // TODO: wording + left.GetTypeName()
		}
		return &LiteralExpression{
			Type:   "number",
			Value:  fmt.Sprintf("%d", li/ri),
			Line:   lle.Line,
			CharAt: lle.CharAt,
		}, nil
	case "%":
		if lle.Type != "number" || rle.Type != "number" {
			return nil, fmt.Errorf("Can't do operator on a function") // TODO: wording + left.GetTypeName()
		}
		return &LiteralExpression{
			Type:   "number",
			Value:  fmt.Sprintf("%d", li%ri),
			Line:   lle.Line,
			CharAt: lle.CharAt,
		}, nil
	}
	return nil, fmt.Errorf("Operator %s is not supported. [%d,%d]", e.Operator.Symbol, e.Operator.Line, e.Operator.CharAt)
}

func (e BinaryExpression) GetCharAt() int {
	return e.CharAt
}

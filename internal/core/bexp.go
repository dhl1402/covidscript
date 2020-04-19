package core

import (
	"fmt"
	"strconv"

	"github.com/dhl1402/covidscript/internal/utils"
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
	if e.Operator.Symbol == "&&" {
		left, err := e.Left.Evaluate(ec)
		if err != nil {
			return nil, err
		}
		if !left.IsTruthy() {
			return &LiteralExpression{
				Type:   LiteralTypeBoolean,
				Value:  "#f",
				Line:   e.Line,
				CharAt: e.CharAt,
			}, nil
		}
		right, err := e.Right.Evaluate(ec)
		if err != nil {
			return nil, err
		}
		return &LiteralExpression{
			Type:   LiteralTypeBoolean,
			Value:  utils.ToBoolStr(right.IsTruthy()),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	}
	if e.Operator.Symbol == "||" {
		left, err := e.Left.Evaluate(ec)
		if err != nil {
			return nil, err
		}
		if left.IsTruthy() {
			return &LiteralExpression{
				Type:   LiteralTypeBoolean,
				Value:  "#t",
				Line:   e.Line,
				CharAt: e.CharAt,
			}, nil
		}
		right, err := e.Right.Evaluate(ec)
		if err != nil {
			return nil, err
		}
		return &LiteralExpression{
			Type:   LiteralTypeBoolean,
			Value:  utils.ToBoolStr(right.IsTruthy()),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	}
	left, err := e.Left.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	right, err := e.Right.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	if e.Operator.Symbol == "==" {
		return &LiteralExpression{
			Type:   LiteralTypeBoolean,
			Value:  utils.ToBoolStr(isEqual(left, right)),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	}
	if e.Operator.Symbol == "!=" {
		return &LiteralExpression{
			Type:   LiteralTypeBoolean,
			Value:  utils.ToBoolStr(!isEqual(left, right)),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	}
	lle, ok := left.(*LiteralExpression)
	if !ok {
		return nil, fmt.Errorf("Runtime error: cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, left.GetType(), e.Operator.Line, e.Operator.CharAt)
	}
	rle, ok := right.(*LiteralExpression)
	if !ok {
		return nil, fmt.Errorf("Runtime error: cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, right.GetType(), e.Operator.Line, e.Operator.CharAt)
	}
	if e.Operator.Symbol == "-" || e.Operator.Symbol == "*" || e.Operator.Symbol == "/" || e.Operator.Symbol == "%" {
		if lle.Type != LiteralTypeNumber {
			return nil, fmt.Errorf("Runtime error: cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, lle.GetType(), e.Operator.Line, e.Operator.CharAt)
		}
		if rle.Type != LiteralTypeNumber {
			return nil, fmt.Errorf("Runtime error: cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, rle.GetType(), e.Operator.Line, e.Operator.CharAt)
		}
	}
	if e.Operator.Symbol == ">" || e.Operator.Symbol == "<" || e.Operator.Symbol == ">=" || e.Operator.Symbol == "<=" {
		if lle.Type != rle.Type {
			return nil, fmt.Errorf("Runtime error: cannot use '%s' operator with 2 different types. [%d,%d]", e.Operator.Symbol, e.Operator.Line, e.Operator.CharAt)
		}
	}
	ln, _ := strconv.ParseFloat(lle.Value, 64)
	rn, _ := strconv.ParseFloat(rle.Value, 64)
	switch e.Operator.Symbol {
	case "+":
		if lle.Type == LiteralTypeUndefined || lle.Type == LiteralTypeNull {
			return nil, fmt.Errorf("Runtime error: cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, lle.GetType(), e.Operator.Line, e.Operator.CharAt)
		}
		if rle.Type == LiteralTypeUndefined || rle.Type == LiteralTypeNull {
			return nil, fmt.Errorf("Runtime error: cannot use '%s' operator with %s. [%d,%d]", e.Operator.Symbol, rle.GetType(), e.Operator.Line, e.Operator.CharAt)
		}
		if lle.Type == LiteralTypeNumber && rle.Type == LiteralTypeNumber {
			return &LiteralExpression{
				Type:   LiteralTypeNumber,
				Value:  fmt.Sprintf("%v", ln+rn),
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
		// handle number
		return &LiteralExpression{
			Type:   LiteralTypeNumber,
			Value:  fmt.Sprintf("%v", ln-rn),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case "*":
		// handle number
		return &LiteralExpression{
			Type:   LiteralTypeNumber,
			Value:  fmt.Sprintf("%v", ln*rn),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case "/":
		// handle number
		if rn == 0 {
			return nil, fmt.Errorf("Runtime error: cannot divide by zero. [%d,%d]", rle.Line, rle.CharAt)
		}
		return &LiteralExpression{
			Type:   LiteralTypeNumber,
			Value:  fmt.Sprintf("%v", ln/rn),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case "%":
		// handle number
		li, err := strconv.Atoi(lle.Value)
		if err != nil {
			return nil, fmt.Errorf("Runtime error: cannot use '%s' operator with float. [%d,%d]", e.Operator.Symbol, e.Operator.Line, e.Operator.CharAt)
		}
		ri, err := strconv.Atoi(rle.Value)
		if err != nil {
			return nil, fmt.Errorf("Runtime error: cannot use '%s' operator with float. [%d,%d]", e.Operator.Symbol, e.Operator.Line, e.Operator.CharAt)
		}
		return &LiteralExpression{
			Type:   LiteralTypeNumber,
			Value:  fmt.Sprintf("%v", li%ri),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case ">":
		// handle literal, same type
		if lle.Type == LiteralTypeNumber {
			return &LiteralExpression{
				Type:   LiteralTypeBoolean,
				Value:  utils.ToBoolStr(ln > rn),
				Line:   e.Line,
				CharAt: e.CharAt,
			}, nil
		}
		return &LiteralExpression{
			Type:   LiteralTypeBoolean,
			Value:  utils.ToBoolStr(lle.Value > rle.Value),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case "<":
		// handle literal, same type
		if lle.Type == LiteralTypeNumber {
			return &LiteralExpression{
				Type:   LiteralTypeBoolean,
				Value:  utils.ToBoolStr(ln < rn),
				Line:   e.Line,
				CharAt: e.CharAt,
			}, nil
		}
		return &LiteralExpression{
			Type:   LiteralTypeBoolean,
			Value:  utils.ToBoolStr(lle.Value < rle.Value),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case ">=":
		// handle literal, same type
		if lle.Type == LiteralTypeNumber {
			return &LiteralExpression{
				Type:   LiteralTypeBoolean,
				Value:  utils.ToBoolStr(ln >= rn),
				Line:   e.Line,
				CharAt: e.CharAt,
			}, nil
		}
		return &LiteralExpression{
			Type:   LiteralTypeBoolean,
			Value:  utils.ToBoolStr(lle.Value >= rle.Value),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case "<=":
		// handle literal, same type
		if lle.Type == LiteralTypeNumber {
			return &LiteralExpression{
				Type:   LiteralTypeBoolean,
				Value:  utils.ToBoolStr(ln <= rn),
				Line:   e.Line,
				CharAt: e.CharAt,
			}, nil
		}
		return &LiteralExpression{
			Type:   LiteralTypeBoolean,
			Value:  utils.ToBoolStr(lle.Value <= rle.Value),
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	}
	return nil, fmt.Errorf("Runtime error: operator %s is not supported. [%d,%d]", e.Operator.Symbol, e.Operator.Line, e.Operator.CharAt)
}

func isEqual(e1 Expression, e2 Expression) bool {
	le1, ok := e1.(*LiteralExpression)
	if ok {
		le2, ok := e2.(*LiteralExpression)
		if ok {
			// if both is primitive type
			return le1.Type == le2.Type && le1.Value == le2.Value
		}
	}
	// otherwise compare pointer reference
	return e1 == e2
}

func (e *BinaryExpression) IsTruthy() bool {
	return true
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

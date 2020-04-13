package core

import (
	"fmt"
	"strconv"
)

type MemberAccessExpression struct {
	Object             Expression
	PropertyExpression Expression
	PropertyIdentifier Identifier
	Compute            bool
	Line               int
	CharAt             int
}

func (e MemberAccessExpression) Evaluate(ec ExecutionContext) (Expression, error) {
	obj, err := e.Object.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	prop, err := e.PropertyExpression.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	lexp, ok := prop.(LiteralExpression)
	if !ok || (lexp.Type != "string" && lexp.Type != "number") {
		return nil, fmt.Errorf("Property key of type %s is not supported. [%d,%d]", "function", lexp.Line, lexp.GetCharAt()) // TODO: get type
	}
	switch o := obj.(type) {
	case ObjectExpression:
		for _, p := range o.Properties {
			if p.Computed {
				kexp, ok := p.KeyExpression.(LiteralExpression)
				if ok && lexp.Type == kexp.Type && kexp.Value == lexp.Value {
					p.Line = e.Line
					p.CharAt = e.CharAt
					return p.Value, nil // TODO: line and char should be set to e.Line, e.Char
				}
			} else if p.KeyIdentifier.Name == lexp.Value {
				return p.Value, nil // TODO: line and char should be set to e.Line, e.Char
			}
		}
		return LiteralExpression{
			Type:   "undefined",
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case ArrayExpression:
		if lexp.Type != "number" {
			return nil, fmt.Errorf("Index must be number. [%d,%d]", lexp.Line, lexp.GetCharAt())
		}
		if i, err := strconv.Atoi(lexp.Value); err == nil {
			if i >= len(o.Elements) {
				return nil, fmt.Errorf("Index is out of range. [%d.%d]", lexp.Line, lexp.CharAt)
			}
			return o.Elements[i], nil // TODO: line and char should be set to e.Line, e.Char
		}
		return nil, fmt.Errorf("Invalid array index. [%d,%d]", lexp.Line, lexp.CharAt)
	}
	return nil, fmt.Errorf("Can't access property of type %s. [%d,%d]", "number", e.Line, e.CharAt)
}

func (e MemberAccessExpression) GetCharAt() int {
	return e.CharAt
}

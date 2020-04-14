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

func (e *MemberAccessExpression) Evaluate(ec *ExecutionContext) (Expression, error) {
	var pexp *LiteralExpression
	obj, err := e.Object.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	if e.Compute {
		tmpExp, err := e.PropertyExpression.Evaluate(ec)
		if err != nil {
			return nil, err
		}
		lexp, ok := tmpExp.(*LiteralExpression)
		if !ok || (lexp.Type != "string" && lexp.Type != "number") {
			return nil, fmt.Errorf("Property key of type %s is not supported. [%d,%d]", tmpExp.GetType(), lexp.Line, lexp.GetCharAt())
		}
		pexp = lexp
	}
	switch o := obj.(type) {
	case (*ObjectExpression):
		for _, p := range o.Properties {
			if p.Computed {
				kexp, ok := p.KeyExpression.(*LiteralExpression)
				if ok && (e.PropertyIdentifier.Name == kexp.Value || e.Compute && kexp.Type == pexp.Type && kexp.Value == pexp.Value) {
					p.Line = e.Line
					p.CharAt = e.CharAt
					return p.Value, nil // TODO: line and char should be set to e.Line, e.Char
				}
			} else if p.KeyIdentifier.Name == e.PropertyIdentifier.Name || p.KeyIdentifier.Name == pexp.Value {
				return p.Value, nil // TODO: line and char should be set to e.Line, e.Char
			}
		}
		return &LiteralExpression{
			Type:   "undefined",
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case (*ArrayExpression):
		if !e.Compute {
			// TODO: support built-in property
			return &LiteralExpression{
				Type:   "undefined",
				Line:   e.Line,
				CharAt: e.CharAt,
			}, nil
		}
		if pexp.Type != "number" {
			return nil, fmt.Errorf("Index must be number. [%d,%d]", pexp.Line, pexp.GetCharAt())
		}
		if i, err := strconv.Atoi(pexp.Value); err == nil {
			if i >= len(o.Elements) {
				return nil, fmt.Errorf("Index is out of range. [%d.%d]", pexp.Line, pexp.CharAt)
			}
			o.Elements[i].SetLine(e.Line)
			o.Elements[i].SetCharAt(e.CharAt)
			return o.Elements[i], nil
		}
		return nil, fmt.Errorf("Invalid array index. [%d,%d]", pexp.Line, pexp.CharAt)
	}
	return nil, fmt.Errorf("Can't access property of type %s. [%d,%d]", obj.GetType(), e.Line, e.CharAt)
}

func (e *MemberAccessExpression) GetCharAt() int {
	return e.CharAt
}

func (e *MemberAccessExpression) GetLine() int {
	return e.Line
}

func (e *MemberAccessExpression) SetLine(i int) {
	e.Line = i
}

func (e *MemberAccessExpression) SetCharAt(i int) {
	e.CharAt = i
}

func (e *MemberAccessExpression) GetType() string {
	return "member access expression"
}

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
		if !ok || (lexp.Type != LiteralTypeString && lexp.Type != LiteralTypeNumber) {
			return nil, fmt.Errorf("Runtime error: property key of type %s is not supported. [%d,%d]", tmpExp.GetType(), tmpExp.GetLine(), tmpExp.GetCharAt())
		}
		pexp = lexp
	}
	switch o := obj.(type) {
	case (*ObjectExpression):
		for _, p := range o.Properties {
			if p.Computed {
				kexp, ok := p.KeyExpression.(*LiteralExpression)
				if ok && (e.PropertyIdentifier.Name == kexp.Value || e.Compute && kexp.Type == pexp.Type && kexp.Value == pexp.Value) {
					p.Value.SetLine(e.Line)
					p.Value.SetCharAt(e.CharAt)
					return p.Value, nil
				}
			} else if p.KeyIdentifier.Name == e.PropertyIdentifier.Name || (pexp != nil && p.KeyIdentifier.Name == pexp.Value) {
				p.Value.SetLine(e.Line)
				p.Value.SetCharAt(e.CharAt)
				return p.Value, nil
			}
		}
		return &LiteralExpression{
			Type:   LiteralTypeUndefined,
			Line:   e.Line,
			CharAt: e.CharAt,
		}, nil
	case (*ArrayExpression):
		if !e.Compute {
			return &LiteralExpression{
				Type:   LiteralTypeUndefined,
				Line:   e.Line,
				CharAt: e.CharAt,
			}, nil
		}
		if pexp.Type != LiteralTypeNumber {
			return nil, fmt.Errorf("Runtime error: index must be number. [%d,%d]", pexp.Line, pexp.GetCharAt())
		}
		if i, err := strconv.Atoi(pexp.Value); err == nil {
			if i >= len(o.Elements) {
				return nil, fmt.Errorf("Runtime error: index is out of range. [%d.%d]", pexp.Line, pexp.CharAt)
			}
			o.Elements[i].SetLine(e.Line)
			o.Elements[i].SetCharAt(e.CharAt)
			return o.Elements[i], nil
		}
		return nil, fmt.Errorf("Runtime error: invalid array index. [%d,%d]", pexp.Line, pexp.CharAt)
	case (*LiteralExpression):
		if o.Type != LiteralTypeString {
			break
		}
		if !e.Compute {
			return &LiteralExpression{
				Type:   LiteralTypeUndefined,
				Line:   e.Line,
				CharAt: e.CharAt,
			}, nil
		}
		if pexp.Type != LiteralTypeNumber {
			return nil, fmt.Errorf("Runtime error: index must be number. [%d,%d]", pexp.Line, pexp.GetCharAt())
		}
		if i, err := strconv.Atoi(pexp.Value); err == nil {
			if i >= len(o.Value) {
				return nil, fmt.Errorf("Runtime error: index is out of range. [%d.%d]", pexp.Line, pexp.CharAt)
			}
			return &LiteralExpression{
				Type:   LiteralTypeString,
				Value:  string(o.Value[i]),
				Line:   o.Line,
				CharAt: o.CharAt + i + 1,
			}, nil
		}
		return nil, fmt.Errorf("Runtime error: invalid index. [%d,%d]", pexp.Line, pexp.CharAt)
	}
	return nil, fmt.Errorf("Runtime error: can't access property of type %s. [%d,%d]", obj.GetType(), e.Line, e.CharAt)
}

func (e *MemberAccessExpression) IsTruthy() bool {
	return true
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

func (e *MemberAccessExpression) ToString() string {
	return ""
}

func (e *MemberAccessExpression) Clone() Expression {
	var pexp Expression
	if e.PropertyExpression != nil {
		pexp = e.PropertyExpression.Clone()
	}
	return &MemberAccessExpression{
		Object:             e.Object.Clone(),
		PropertyIdentifier: e.PropertyIdentifier,
		PropertyExpression: pexp,
		Compute:            e.Compute,
		Line:               e.Line,
		CharAt:             e.CharAt,
	}
}

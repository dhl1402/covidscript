package core

import (
	"fmt"
	"strconv"
)

type AssignmentStatement struct {
	Left                 Expression
	Right                Expression
	DeclarationShorthand bool
	Line                 int
	CharAt               int
}

func (stmt AssignmentStatement) Execute(ec *ExecutionContext) (Expression, error) {
	right, err := stmt.Right.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	switch left := stmt.Left.(type) {
	case (*VariableExpression):
		if ec.Variables[left.Name] == nil && !stmt.DeclarationShorthand {
			return nil, fmt.Errorf("%s is not defined. [%d,%d]", left.Name, stmt.Line, stmt.CharAt)
		}
		ec.Set(left.Name, right)
	case (*MemberAccessExpression):
		// do not need to handle error in this case because error have been already handled in the following `left.Evaluate(ec)``
		var pexp *LiteralExpression
		_, err := left.Evaluate(ec)
		if err != nil {
			return nil, err
		}
		obj, _ := left.Object.Evaluate(ec)
		if left.Compute {
			tmpExp, _ := left.PropertyExpression.Evaluate(ec)
			lexp, _ := tmpExp.(*LiteralExpression)
			pexp = lexp
		}
		switch o := obj.(type) {
		case (*ObjectExpression):
			for _, p := range o.Properties {
				if p.Computed {
					kexp, ok := p.KeyExpression.(*LiteralExpression)
					if ok && (left.PropertyIdentifier.Name == kexp.Value || left.Compute && kexp.Type == pexp.Type && kexp.Value == pexp.Value) {
						p.Value = right
						return nil, nil
					}
				} else if p.KeyIdentifier.Name == left.PropertyIdentifier.Name || (pexp != nil && p.KeyIdentifier.Name == pexp.Value) {
					p.Value = right
					return nil, nil
				}
			}
			newProp := &ObjectProperty{
				KeyIdentifier: left.PropertyIdentifier,
				KeyExpression: pexp,
				Computed:      pexp != nil,
				Value:         right,
				Line:          right.GetLine(),
				CharAt:        right.GetCharAt(),
			}
			if pexp == nil {
				newProp.KeyExpression = nil
			}
			o.Properties = append(o.Properties, newProp)
			return nil, nil
		case (*ArrayExpression):
			if i, err := strconv.Atoi(pexp.Value); err == nil {
				o.Elements[i] = right
				return nil, nil
			}
		}
	default:
		return nil, fmt.Errorf("Cannot identify variable. [%d,%d]", stmt.Left.GetLine(), stmt.Left.GetCharAt())
	}
	return nil, nil
}

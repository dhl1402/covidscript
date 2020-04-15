package core

import "fmt"

type (
	ObjectProperty struct {
		KeyExpression Expression
		KeyIdentifier Identifier
		Value         Expression
		Computed      bool
		Shorthand     bool
		Method        bool
		Line          int
		CharAt        int
	}
	ObjectExpression struct {
		Properties []*ObjectProperty
		Line       int
		CharAt     int
	}
)

func (e *ObjectExpression) Evaluate(ec *ExecutionContext) (Expression, error) {
	props := []*ObjectProperty{}
	for _, p := range e.Properties {
		if p.Computed {
			kexp, err := p.KeyExpression.Evaluate(ec)
			if err != nil {
				return nil, err
			}
			kexp.SetLine(p.KeyExpression.GetLine())
			kexp.SetCharAt(p.KeyExpression.GetCharAt())
			p.KeyExpression = kexp
		}
		v, err := p.Value.Evaluate(ec)
		v.SetLine(p.Value.GetLine())
		v.SetCharAt(p.Value.GetCharAt())
		if err != nil {
			return nil, err
		}
		p.Value = v
		props = append(props, p)
	}
	e.Properties = props
	return e, nil
}

func (e *ObjectExpression) GetCharAt() int {
	return e.CharAt
}

func (e *ObjectExpression) GetLine() int {
	return e.Line
}

func (e *ObjectExpression) SetLine(i int) {
	e.Line = i
}

func (e *ObjectExpression) SetCharAt(i int) {
	e.CharAt = i
}

func (e *ObjectExpression) GetType() string {
	return "object"
}

func (e *ObjectExpression) ToString() string {
	s := "{"
	for _, p := range e.Properties {
		if p.Computed {
			s = s + fmt.Sprintf("%s: %s, ", p.KeyExpression.ToString(), p.Value.ToString())
		} else {
			s = s + fmt.Sprintf("%s: %s, ", p.KeyIdentifier.Name, p.Value.ToString())
		}
	}
	if len(s) > 1 {
		s = s[:len(s)-2]
	}
	return s + "]"
}

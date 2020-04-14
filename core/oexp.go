package core

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
		Properties []ObjectProperty
		Line       int
		CharAt     int
	}
)

func (e ObjectExpression) Evaluate(ec ExecutionContext) (Expression, error) {
	props := []ObjectProperty{}
	for _, p := range e.Properties {
		if p.Computed {
			kexp, err := p.KeyExpression.Evaluate(ec)
			if err != nil {
				return nil, err
			}
			p.KeyExpression = kexp
		}
		v, err := p.Value.Evaluate(ec)
		if err != nil {
			return nil, err
		}
		p.Value = v
		props = append(props, p)
	}
	e.Properties = props
	return e, nil
}

func (e ObjectExpression) GetCharAt() int {
	return e.CharAt
}

func (e ObjectExpression) GetLine() int {
	return e.Line
}

func (e ObjectExpression) GetType() string {
	return "object"
}

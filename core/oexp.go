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
	return e, nil
}

func (e ObjectExpression) GetCharAt() int {
	return e.CharAt
}

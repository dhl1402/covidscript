package core

type (
	VariableDeclarator struct {
		ID     Identifier
		Init   Expression
		Line   int
		CharAt int
	}
	VariableDeclaration struct {
		Declarations []VariableDeclarator
		Line         int
		CharAt       int
	}
)

func (stmt VariableDeclaration) Execute(ec *ExecutionContext) (Expression, error) {
	for _, d := range stmt.Declarations {
		if d.Init != nil {
			if f, ok := d.Init.(*FunctionExpression); ok {
				f.EC = &ExecutionContext{
					Outer:     ec,
					Variables: map[string]Expression{},
				}
			}
			value, err := d.Init.Evaluate(ec)
			if err != nil {
				return nil, err
			}
			ec.Set(d.ID.Name, value)
		} else {
			ec.Set(d.ID.Name, &LiteralExpression{
				Type:   "undefined",
				Line:   d.ID.Line,
				CharAt: d.ID.CharAt,
			})
		}
	}
	return nil, nil
}

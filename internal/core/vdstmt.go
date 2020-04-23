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
			value, err := d.Init.Evaluate(ec)
			if err != nil {
				return nil, err
			}
			ec.Set(d.ID.Name, value)
		} else {
			ec.Set(d.ID.Name, &LiteralExpression{
				Type:   LiteralTypeUndefined,
				Line:   d.ID.Line,
				CharAt: d.ID.CharAt,
			})
		}
	}
	return nil, nil
}

func (stmt VariableDeclaration) Clone() Statement {
	ds := []VariableDeclarator{}
	for _, vd := range stmt.Declarations {
		var init Expression
		if vd.Init != nil {
			init = vd.Init.Clone()
		}
		ds = append(ds, VariableDeclarator{
			ID:     vd.ID,
			Init:   init,
			Line:   vd.Line,
			CharAt: vd.CharAt,
		})
	}
	return VariableDeclaration{
		Declarations: ds,
		Line:         stmt.Line,
		CharAt:       stmt.CharAt,
	}
}

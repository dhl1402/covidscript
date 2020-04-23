package core

type FunctionDeclaration struct {
	ID     Identifier
	Params []Identifier
	Body   BlockStatement
	Line   int
	CharAt int
}

func (stmt FunctionDeclaration) Execute(ec *ExecutionContext) (Expression, error) {
	fexp := &FunctionExpression{
		Params: stmt.Params,
		Body:   stmt.Body,
		Line:   stmt.Line,
		CharAt: stmt.CharAt,
		EC: &ExecutionContext{
			Outer:     ec,
			Variables: map[string]Expression{},
		},
	}
	ec.Set(stmt.ID.Name, fexp)
	return nil, nil
}

func (stmt FunctionDeclaration) Clone() Statement {
	return FunctionDeclaration{
		ID:     stmt.ID,
		Params: stmt.Params,
		Body:   stmt.Body,
		Line:   stmt.Line,
		CharAt: stmt.CharAt,
	}
}

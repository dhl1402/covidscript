package core

type ReturnStatement struct {
	Argument Expression
	Line     int
	CharAt   int
}

func (stmt ReturnStatement) Execute(ec *ExecutionContext) (Expression, error) {
	if stmt.Argument == nil {
		return &LiteralExpression{
			Type:   LiteralTypeUndefined,
			Line:   stmt.Line,
			CharAt: stmt.CharAt,
		}, nil
	}
	return stmt.Argument.Evaluate(ec)
}

func (stmt ReturnStatement) Clone() Statement {
	var arg Expression
	if stmt.Argument != nil {
		arg = stmt.Argument.Clone()
	}
	return ReturnStatement{
		Argument: arg,
		Line:     stmt.Line,
		CharAt:   stmt.CharAt,
	}
}

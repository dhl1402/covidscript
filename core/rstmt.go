package core

type ReturnStatement struct {
	Argument Expression
	Line     int
	CharAt   int
}

func (stmt ReturnStatement) Execute(ec *ExecutionContext) (Expression, error) {
	return stmt.Argument.Evaluate(ec)
}

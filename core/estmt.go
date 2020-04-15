package core

type ExpressionStatement struct {
	Expression
	Line   int
	CharAt int
}

func (stmt ExpressionStatement) Execute(ec *ExecutionContext) (Expression, error) {
	return nil, nil
}

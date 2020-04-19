package core

type ExpressionStatement struct {
	Expression
	Line   int
	CharAt int
}

func (stmt ExpressionStatement) Execute(ec *ExecutionContext) (Expression, error) {
	_, err := stmt.Expression.Evaluate(ec)
	return nil, err
}

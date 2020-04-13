package core

type MemberAccessExpression struct {
	Object   Expression
	Property Expression
	Line     int
	CharAt   int
}

func (e MemberAccessExpression) Evaluate(ec ExecutionContext) (Expression, error) {
	return e, nil
}

func (e MemberAccessExpression) GetCharAt() int {
	return e.CharAt
}

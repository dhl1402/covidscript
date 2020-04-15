package core

type AssignmentStatement struct {
	Left   Expression
	Right  Expression
	Line   int
	CharAt int
}

func (stmt AssignmentStatement) Execute(ec *ExecutionContext) (Expression, error) {
	return nil, nil
}

package core

type IfStatement struct {
	Test       Expression
	Assignment AssignmentStatement
	EC         *ExecutionContext
	Consequent BlockStatement
	Alternate  *IfStatement
	Line       int
	CharAt     int
}

func (stmt IfStatement) Execute(ec *ExecutionContext) (Expression, error) {
	return nil, nil
}

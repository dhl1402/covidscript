package core

type IfStatement struct {
	Test                Expression
	VariableDeclaration VariableDeclaration
	EC                  *ExecutionContext
	Consequent          []Statement
	Alternate           *IfStatement
}

func (stmt IfStatement) Execute(ec *ExecutionContext) (Expression, error) {
	return nil, nil
}

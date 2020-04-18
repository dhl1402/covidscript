package core 

type ForStatement struct {
	Assignment *AssignmentStatement
	Test Expression
	Update *AssignmentStatement
	Body BlockStatement
	Line int
	CharAt int
}

func (stmt ForStatement) Execute(ec *ExecutionContext) (Expression,error){
	return nil, nil
}
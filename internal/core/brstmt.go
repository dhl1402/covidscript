package core

import "fmt"

type BreakStatement struct {
	Line   int
	CharAt int
}

func (stmt BreakStatement) Execute(ec *ExecutionContext) (Expression, error) {
	return nil, BreakError{
		Message: fmt.Sprintf("break is not in a loop. [%d,%d]", stmt.Line, stmt.CharAt),
	}
}

type BreakError struct {
	Message string
}

func (err BreakError) Error() string {
	return err.Message
}

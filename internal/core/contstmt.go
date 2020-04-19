package core

import "fmt"

type ContinueStatement struct {
	Line   int
	CharAt int
}

func (stmt ContinueStatement) Execute(ec *ExecutionContext) (Expression, error) {
	return nil, ContinueError{
		Message: fmt.Sprintf("continue is not in a loop. [%d,%d]", stmt.Line, stmt.CharAt),
	}
}

type ContinueError struct {
	Message string
}

func (err ContinueError) Error() string {
	return err.Message
}

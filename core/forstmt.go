package core

type ForStatement struct {
	Assignment *AssignmentStatement
	Test       Expression
	Update     *AssignmentStatement
	Body       BlockStatement
	Line       int
	CharAt     int
}

func (stmt ForStatement) Execute(ec *ExecutionContext) (Expression, error) {
	bec := &ExecutionContext{
		Type:      TypeBlockEC,
		Outer:     ec,
		Variables: map[string]Expression{},
	}
	if stmt.Assignment != nil {
		_, err := stmt.Assignment.Execute(bec)
		if err != nil {
			return nil, err
		}
	}
	if stmt.Test == nil {
	l1:
		for {
			for _, s := range stmt.Body.Statements {
				rexp, err := s.Execute(ec)
				if _, ok := err.(BreakError); ok {
					break l1
				}
				if _, ok := err.(ContinueError); ok {
					break
				}
				if rexp != nil || err != nil {
					return rexp, err
				}
			}
			if stmt.Update != nil {
				_, err := stmt.Update.Execute(bec)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		t, err := stmt.Test.Evaluate(bec)
		if err != nil {
			return nil, err
		}
	l2:
		for t.IsTruthy() {
			for _, s := range stmt.Body.Statements {
				rexp, err := s.Execute(bec)
				if _, ok := err.(BreakError); ok {
					break l2
				}
				if _, ok := err.(ContinueError); ok {
					break
				}
				if rexp != nil || err != nil {
					return rexp, err
				}
			}
			if stmt.Update != nil {
				_, err := stmt.Update.Execute(bec)
				if err != nil {
					return nil, err
				}
			}
			t, err = stmt.Test.Evaluate(bec)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

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
		for {
			rexp, err := stmt.Body.Execute(bec)
			if rexp != nil || err != nil {
				return rexp, err
			}
		}
	} else {
		t, err := stmt.Test.Evaluate(bec)
		if err != nil {
			return nil, err
		}
		for t.IsTruthy() {
			rexp, err := stmt.Body.Execute(bec)
			if rexp != nil || err != nil {
				return rexp, err
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

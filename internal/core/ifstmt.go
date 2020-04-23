package core

type IfStatement struct {
	Test       Expression
	Init       Statement
	Consequent BlockStatement
	Alternate  *IfStatement
	Line       int
	CharAt     int
}

func (stmt IfStatement) Execute(ec *ExecutionContext) (Expression, error) {
	bec := &ExecutionContext{
		Type:      TypeBlockEC,
		Outer:     ec,
		Variables: map[string]Expression{},
	}
	s := stmt
	for true {
		if s.Init != nil {
			_, err := s.Init.Execute(bec)
			if err != nil {
				return nil, err
			}
		}
		if s.Test == nil {
			return s.Consequent.Execute(bec)
		}
		if t, err := s.Test.Evaluate(bec); err != nil {
			return nil, err
		} else if t.IsTruthy() {
			return s.Consequent.Execute(bec)
		}
		if s.Alternate == nil {
			break
		}
		s = *s.Alternate
	}
	return nil, nil
}

func (stmt IfStatement) Clone() Statement {
	consequent := stmt.Consequent.Clone().(BlockStatement)
	var alt *IfStatement
	if stmt.Alternate != nil {
		a := stmt.Alternate.Clone().(IfStatement)
		alt = &a
	}
	var init Statement
	if stmt.Init != nil {
		init = stmt.Init.Clone()
	}
	var test Expression
	if stmt.Test != nil {
		test = stmt.Test.Clone()
	}
	return IfStatement{
		Init:       init,
		Test:       test,
		Alternate:  alt,
		Consequent: consequent,
		Line:       stmt.Line,
		CharAt:     stmt.CharAt,
	}
}

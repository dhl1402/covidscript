package core

type BlockStatement struct {
	Statements []Statement
	Line       int
	CharAt     int
}

func (stmt BlockStatement) Execute(ec *ExecutionContext) (Expression, error) {
	for _, s := range stmt.Statements {
		rexp, err := s.Execute(ec)
		if rexp != nil || err != nil {
			return rexp, err
		}
	}
	return nil, nil
}

func (stmt BlockStatement) Clone() Statement {
	stmts := []Statement{}
	for _, s := range stmt.Statements {
		stmts = append(stmts, s.Clone())
	}
	return BlockStatement{
		Statements: stmts,
		Line:       stmt.Line,
		CharAt:     stmt.CharAt,
	}
}

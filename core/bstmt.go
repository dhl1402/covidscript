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
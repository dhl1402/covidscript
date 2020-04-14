package executor

import (
	"gs/core"
)

func Execute(statement []core.Statement) error {
	gec := core.ExecutionContext{
		Variables: map[string]core.Expression{},
	}
	for _, ss := range statement {
		switch s := ss.(type) {
		case core.VariableDeclaration:
			for _, d := range s.Declarations {
				if d.Init != nil {
					value, err := d.Init.Evaluate(gec)
					if err != nil {
						return err
					}
					gec.Set(d.ID.Name, value)
				} else {
					gec.Set(d.ID.Name, &core.LiteralExpression{
						Type:   "undefined",
						Line:   d.ID.Line,
						CharAt: d.ID.CharAt,
					})
				}
			}
			// case core.AssignmentStatement:
			// 	_, isExist := gec.Get(s.Left.Name)
			// 	if !isExist {
			// 		return fmt.Errorf("%s is not defined", s.Left.Name)
			// 	}
			// 	if s.Right != nil {
			// 		value, err := s.Right.Evaluate(gec)
			// 		if err != nil {
			// 			return err
			// 		}
			// 		gec.Set(s.Left.Name, value)
			// 	} else {
			// 		gec.Set(s.Left.Name, nil)
			// 	}
			// }
		}
	}
	return nil
}

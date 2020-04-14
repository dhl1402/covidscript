package executor

import (
	"gs/core"
)

func Execute(statement []core.Statement) error {
	gec := &core.ExecutionContext{
		Variables: map[string]core.Expression{},
	}
	cexp := core.CallExpression{
		Callee: &core.FunctionExpression{
			Body:   statement,
			Params: []core.Identifier{},
			EC:     gec,
		},
	}
	_, err := cexp.Evaluate(gec)
	return err
}

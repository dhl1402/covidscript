package executor

import (
	"github.com/dhl1402/covidscript/core"
)

func Execute(gec *core.ExecutionContext, stmts []core.Statement) error {
	cexp := core.CallExpression{
		Callee: &core.FunctionExpression{
			Body:   stmts,
			Params: []core.Identifier{},
			EC:     gec,
		},
	}
	_, err := cexp.Evaluate(gec)
	return err
}

package builtin

import (
	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func Type(conf config.Config) *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{{Name: "input"}},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			inp, _ := ec.Get("input")
			return &core.LiteralExpression{
				Type:  core.LiteralTypeString,
				Value: inp.GetType(),
			}, nil
		},
	}
}

package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func Len(conf config.Config) *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{{Name: "inp"}},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			arg, _ := ec.Get("inp")
			switch exp := arg.(type) {
			case (*core.ArrayExpression):
				return &core.LiteralExpression{
					Type:  core.LiteralTypeNumber,
					Value: fmt.Sprintf("%d", len(exp.Elements)),
				}, nil
			case (*core.LiteralExpression):
				if exp.Type == core.LiteralTypeString {
					return &core.LiteralExpression{
						Type:  core.LiteralTypeNumber,
						Value: fmt.Sprintf("%d", len(exp.Value)),
					}, nil
				}
			}
			return nil, fmt.Errorf("Runtime error: unexpected %s as argument type of len, expected array or string.", arg.GetType())
		},
	}
}

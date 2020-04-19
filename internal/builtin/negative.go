package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func Negative(conf config.Config) *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{{Name: "num"}},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			arg, _ := ec.Get("num")
			lexp, ok := arg.(*core.LiteralExpression)
			if !ok || lexp.Type != core.LiteralTypeNumber {
				return nil, fmt.Errorf("Runtime error: unexpected %s as argument type of negative, expected number.", arg.GetType())
			}
			return &core.LiteralExpression{
				Type:   core.LiteralTypeNumber,
				Value:  fmt.Sprintf("-%s", lexp.Value),
				Line:   lexp.Line,
				CharAt: lexp.CharAt,
			}, nil
		},
	}
}

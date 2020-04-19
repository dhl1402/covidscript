package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func IndexOf(conf config.Config) *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{{Name: "array"}, {Name: "elem"}},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			arg1, _ := ec.Get("array")
			arexp, ok := arg1.(*core.ArrayExpression)
			if !ok {
				return nil, fmt.Errorf("Runtime error: first argument must be array.")
			}
			arg2, _ := ec.Get("elem")
			for _, elem := range arexp.Elements {
				bexp := &core.BinaryExpression{
					Left:     elem,
					Right:    arg2,
					Operator: core.Operator{Symbol: "=="},
				}
				rexp, _ := bexp.Evaluate(ec)
				if rexp.IsTruthy() {
					return &core.LiteralExpression{
						Type:  core.LiteralTypeNumber,
						Value: fmt.Sprintf("%d", i),
					}, nil
				}
			}
			return &core.LiteralExpression{
				Type:  core.LiteralTypeNumber,
				Value: "-1",
			}, nil
		},
	}
}

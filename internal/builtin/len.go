package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func Len(conf config.Config) *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{{Name: "array"}},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			arg, _ := ec.Get("array")
			arexp, ok := arg.(*core.ArrayExpression)
			if !ok {
				return nil, fmt.Errorf("Runtime error: unexpected type %s, expected array. [%d,%d]", arg.GetType(), arg.GetLine(), arg.GetCharAt())
			}
			return &core.LiteralExpression{
				Type:  core.LiteralTypeNumber,
				Value: fmt.Sprintf("%d", len(arexp.Elements)),
			}, nil
		},
	}
}
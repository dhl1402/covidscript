package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/core"
)

func Values() *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{{Name: "obj"}},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			arg, _ := ec.Get("obj")
			oexp, ok := arg.(*core.ObjectExpression)
			if !ok {
				return nil, fmt.Errorf("Runtime error: unexpected %s as argument type of keys, expected object.", arg.GetType())
			}
			result := []core.Expression{}
			for _, prop := range oexp.Properties {
				result = append(result, prop.Value)
			}
			return &core.ArrayExpression{
				Elements: result,
			}, nil
		},
	}
}

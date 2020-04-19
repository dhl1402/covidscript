package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func Append(conf config.Config) *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{{Name: "array"}},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			arg, _ := ec.Get("array")
			_, ok := arg.(*core.ArrayExpression)
			if !ok {
				return nil, fmt.Errorf("Runtime error: unexpected %s as argument type of append, expected array.", arg.GetType())
			}
			result := &core.ArrayExpression{
				Elements: []core.Expression{},
			}
			for i := 1; ; i++ {
				if arg, ok := ec.Variables[fmt.Sprintf("_args%d_", i)]; ok {
					result.Elements = append(result.Elements, arg)
				} else {
					break
				}
			}
			return result, nil
		},
	}
}

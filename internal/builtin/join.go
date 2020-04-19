package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func Join(conf config.Config) *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{{Name: "array"}, {Name: "separator"}},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			arg1, _ := ec.Get("array")
			arexp, ok := arg1.(*core.ArrayExpression)
			if !ok {
				return nil, fmt.Errorf("Runtime error: first argument must be array.")
			}
			arg2, _ := ec.Get("separator")
			lexp, ok := arg2.(*core.LiteralExpression)
			if !ok || (lexp.Type != core.LiteralTypeString && lexp.Type != core.LiteralTypeUndefined) {
				return nil, fmt.Errorf("Runtime error: second argument must be string.")
			}
			var sep string
			if lexp.Type == core.LiteralTypeUndefined {
				sep = ","
			} else {
				sep = lexp.Value
			}
			result := ""
			for _, elem := range arexp.Elements {
				result = result + fmt.Sprintf("%s%s", elem, sep)
			}
			// if len(arexp.Elements) > 0 {
			// 	result = result[:len(sep)]
			// }
			return &core.LiteralExpression{
				Type:  core.LiteralTypeString,
				Value: result,
			}, nil
		},
	}
}

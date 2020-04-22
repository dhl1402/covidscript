package builtin

import (
	"fmt"
	"math"
	"strconv"

	"github.com/dhl1402/covidscript/internal/core"
)

func Floor() *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{{Name: "num"}},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			arg, _ := ec.Get("num")
			lexp, ok := arg.(*core.LiteralExpression)
			if !ok || lexp.Type != core.LiteralTypeNumber {
				return nil, fmt.Errorf("Runtime error: unexpected %s as argument type of floor, expected number.", arg.GetType())
			}
			f, _ := strconv.ParseFloat(lexp.Value, 64)
			return &core.LiteralExpression{
				Type:   core.LiteralTypeNumber,
				Value:  fmt.Sprintf("%v", math.Floor(f)),
				Line:   lexp.Line,
				CharAt: lexp.CharAt,
			}, nil
		},
	}
}

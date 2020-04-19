package builtin

import (
	"fmt"
	"sort"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func Sort(conf config.Config) *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{
			{Name: "array"},
			{Name: "comparator"},
		},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			arg, _ := ec.Get("array")
			aexp, ok := arg.(*core.ArrayExpression)
			if !ok {
				return nil, fmt.Errorf("Runtime error: first argument of sort must be array.")
			}
			comp, _ := ec.Get("comparator")
			fexp, ok := comp.(*core.FunctionExpression)
			if !ok {
				return nil, fmt.Errorf("Runtime error: second argument of sort must be function.")
			}
			sort.SliceStable(aexp.Elements, func(i, j int) bool {
				cexp := core.CallExpression{
					Callee: fexp,
					Arguments: []core.Expression{
						aexp.Elements[i],
						aexp.Elements[j],
					},
				}
				test, _ := cexp.Evaluate(ec)
				return test != nil && test.IsTruthy()
			})
			return &core.ArrayExpression{
				Elements: aexp.Elements,
			}, nil
		},
	}
}

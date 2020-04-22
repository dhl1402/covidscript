package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/core"
)

func Filter() *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{
			{Name: "input"},
			{Name: "callback"},
		},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			cb, _ := ec.Get("callback")
			fexp, ok := cb.(*core.FunctionExpression)
			if !ok {
				return nil, fmt.Errorf("Runtime error: second argument of filter must be function.")
			}
			inp, _ := ec.Get("input")
			switch exp := inp.(type) {
			case (*core.ArrayExpression):
				result := &core.ArrayExpression{
					Elements: []core.Expression{},
				}
				for i, elem := range exp.Elements {
					cexp := core.CallExpression{
						Callee: fexp,
						Arguments: []core.Expression{
							elem,
							&core.LiteralExpression{
								Type:  core.LiteralTypeNumber,
								Value: fmt.Sprintf("%d", i),
							},
						},
					}
					test, err := cexp.Evaluate(ec)
					if err != nil {
						return nil, err
					}
					if test.IsTruthy() {
						result.Elements = append(result.Elements, elem)
					}
				}
				return result, nil
			case (*core.ObjectExpression):
				result := &core.ObjectExpression{
					Properties: []*core.ObjectProperty{},
				}
				for _, prop := range exp.Properties {
					var key core.Expression
					if !prop.Computed {
						key = &core.LiteralExpression{
							Type:  core.LiteralTypeString,
							Value: prop.KeyIdentifier.Name,
						}
					} else {
						key = prop.KeyExpression
					}
					cexp := core.CallExpression{
						Callee:    fexp,
						Arguments: []core.Expression{prop.Value, key},
					}
					test, err := cexp.Evaluate(ec)
					if err != nil {
						return nil, err
					}
					if test.IsTruthy() {
						result.Properties = append(result.Properties, prop)
					}
				}
				return result, nil
			}
			return nil, fmt.Errorf("Runtime error: first argument of filter must be array or object.")
		},
	}
}

package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func Map(conf config.Config) *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{
			{Name: "input"},
			{Name: "callback"},
		},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			cb, _ := ec.Get("callback")
			fexp, ok := cb.(*core.FunctionExpression)
			if !ok {
				return nil, fmt.Errorf("Runtime error: second argument of map must be function.")
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
					rexp, err := cexp.Evaluate(ec)
					if err != nil {
						return nil, err
					}
					result.Elements = append(result.Elements, rexp)
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
					rexp, err := cexp.Evaluate(ec)
					if err != nil {
						return nil, err
					}
					result.Properties = append(result.Properties, &core.ObjectProperty{
						KeyExpression: prop.KeyExpression,
						KeyIdentifier: prop.KeyIdentifier,
						Computed:      prop.Computed,
						Shorthand:     prop.Shorthand,
						Method:        prop.Method,
						Value:         rexp,
					})
				}
				return result, nil
			}
			return nil, fmt.Errorf("Runtime error: first argument of map must be array or object.")
		},
	}
}

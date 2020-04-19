package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func Reduce(conf config.Config) *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{
			{Name: "input"},
			{Name: "callback"},
			{Name: "init"},
		},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			cb, _ := ec.Get("callback")
			fexp, ok := cb.(*core.FunctionExpression)
			if !ok {
				return nil, fmt.Errorf("Runtime error: second argument of filter must be function.")
			}
			inp, _ := ec.Get("input")
			init, _ := ec.Get("init")
			switch exp := inp.(type) {
			case (*core.ArrayExpression):
				var result = init
				var err error
				for i, elem := range exp.Elements {
					cexp := core.CallExpression{
						Callee: fexp,
						Arguments: []core.Expression{
							result,
							elem,
							&core.LiteralExpression{
								Type:  core.LiteralTypeNumber,
								Value: fmt.Sprintf("%d", i),
							},
						},
					}
					result, err = cexp.Evaluate(ec)
					if err != nil {
						return nil, err
					}
				}
				return result, nil
			case (*core.ObjectExpression):
				var result = init
				var err error
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
						Arguments: []core.Expression{result, prop.Value, key},
					}
					result, err = cexp.Evaluate(ec)
					if err != nil {
						return nil, err
					}
				}
				return result, nil
			}
			return nil, fmt.Errorf("Runtime error: first argument of filter must be array or object.")
		},
	}
}

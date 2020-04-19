package builtin

import (
	"fmt"
	"strconv"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func Delete(conf config.Config) *core.FunctionExpression {
	return &core.FunctionExpression{
		Params: []core.Identifier{
			{Name: "input"},
			{Name: "key"},
		},
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			inp, _ := ec.Get("input")
			arg2, _ := ec.Get("key")
			switch exp := inp.(type) {
			case (*core.ArrayExpression):
				result := &core.ArrayExpression{
					Elements: []core.Expression{},
				}
				lexp, ok := arg2.(*core.LiteralExpression)
				if !ok || lexp.Type != core.LiteralTypeNumber {
					return nil, fmt.Errorf("Runtime error: second argument must be integer when deleting array element.")
				}
				idx, err := strconv.Atoi(lexp.Value)
				if err != nil {
					return nil, fmt.Errorf("Runtime error: second argument must be integer when deleting array element.")
				}
				for i, e := range exp.Elements {
					if i != idx {
						result.Elements = append(result.Elements, e)
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
					bexp := &core.BinaryExpression{
						Left:     arg2,
						Right:    key,
						Operator: core.Operator{Symbol: "=="},
					}
					rexp, _ := bexp.Evaluate(ec)
					if !rexp.IsTruthy() {
						result.Properties = append(result.Properties, prop)
					}
				}
				return result, nil
			}
			return nil, fmt.Errorf("Runtime error: first argument of filter must be array or object.")
		},
	}
}

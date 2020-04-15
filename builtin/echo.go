package builtin

import (
	"covs/core"
	"fmt"
)

func Echo() *core.FunctionExpression {
	return &core.FunctionExpression{
		NativeFunction: func(ec *core.ExecutionContext) (core.Expression, error) {
			s := ""
			for i := 0; ; i++ {
				if arg, ok := ec.Variables[fmt.Sprintf("_args%d_", i)]; ok {
					s = s + fmt.Sprintf("%s ", arg.ToString())
				} else {
					break
				}
			}
			fmt.Println(s)
			return nil, nil
		},
	}
}

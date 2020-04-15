package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/core"
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

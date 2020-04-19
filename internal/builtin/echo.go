package builtin

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
)

func Echo(conf config.Config) *core.FunctionExpression {
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
			fmt.Fprintln(conf.Writer, s)
			return nil, nil
		},
	}
}

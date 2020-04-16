package builtin

import (
	"fmt"
	"io"
	"os"

	"github.com/dhl1402/covidscript/core"
)

var w io.Writer = os.Stdout

func SetWriter(ww io.Writer) {
	w = ww
}

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
			fmt.Fprintln(w, s)
			return nil, nil
		},
	}
}

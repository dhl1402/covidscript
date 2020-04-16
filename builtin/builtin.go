package builtin

import (
	"io"

	"github.com/dhl1402/covidscript/core"
)

type Config struct {
	Writer io.Writer
}

func CreateGlobalEC(conf Config) *core.ExecutionContext {
	return &core.ExecutionContext{
		Type: core.TypeGlobalEC,
		Variables: map[string]core.Expression{
			"echo": Echo(conf),
		},
	}
}

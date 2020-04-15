package builtin

import "github.com/dhl1402/covidscript/core"

func CreateGlobalEC() *core.ExecutionContext {
	return &core.ExecutionContext{
		Type: core.TypeGlobalEC,
		Variables: map[string]core.Expression{
			"echo": Echo(),
		},
	}
}

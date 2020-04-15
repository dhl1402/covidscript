package builtin

import "covs/core"

func CreateGlobalEC() *core.ExecutionContext {
	return &core.ExecutionContext{
		Type: core.TypeGlobalEC,
		Variables: map[string]core.Expression{
			"echo": Echo(),
		},
	}
}

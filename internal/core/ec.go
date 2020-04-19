package core

type ecType string

const (
	TypeGlobalEC   ecType = "GlobalEC"
	TypeFunctionEC ecType = "FunctionEC"
	TypeBlockEC    ecType = "BlockEC"
)

type ExecutionContext struct {
	Type      ecType
	Outer     *ExecutionContext
	Variables map[string]Expression
}

func (ec *ExecutionContext) Get(s string) (Expression, bool) {
	for ec != nil {
		if exp, ok := ec.Variables[s]; ok {
			return exp, ok
		}
		ec = ec.Outer
	}
	return nil, false
}

func (ec *ExecutionContext) Set(s string, exp Expression) {
	ec.Variables[s] = exp
}

func (ec *ExecutionContext) Assign(s string, exp Expression) bool {
	for ec != nil {
		if _, ok := ec.Variables[s]; ok {
			ec.Variables[s] = exp
			return true
		}
		ec = ec.Outer
	}
	return false
}

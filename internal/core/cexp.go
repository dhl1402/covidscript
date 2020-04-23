package core

import "fmt"

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
	Line      int
	CharAt    int
}

func (e *CallExpression) Evaluate(ec *ExecutionContext) (Expression, error) {
	callee, err := e.Callee.Evaluate(ec)
	if err != nil {
		return nil, err
	}
	f, ok := callee.(*FunctionExpression)
	if !ok {
		return nil, fmt.Errorf("Runtime error: %s is not a function. [%d,%d]", e.Callee.ToString(), e.Line, e.CharAt)
	}
	fEC := f.EC
	if f.EC.Type != TypeGlobalEC {
		fEC = f.EC.Clone()
	}
	for i, argexp := range e.Arguments {
		arg, err := argexp.Evaluate(ec)
		if err != nil {
			return nil, err
		}
		if i < len(f.Params) {
			fEC.Set(f.Params[i].Name, arg)
		}
		fEC.Set(fmt.Sprintf("_args%d_", i), arg)
	}
	for _, p := range f.Params {
		if _, exist := fEC.Get(p.Name); !exist {
			fEC.Set(p.Name, &LiteralExpression{
				Type:   LiteralTypeUndefined,
				Line:   p.Line,
				CharAt: p.CharAt,
			})
		}
	}
	if f.NativeFunction != nil {
		rexp, err := f.NativeFunction(fEC)
		if err != nil && string(err.Error()[len(err.Error())-1]) == "." {
			err = fmt.Errorf("%s [%d,%d]", err.Error(), e.Line, e.CharAt)
		}
		return rexp, err
	}
	var fBody Statement = f.Body
	if f.EC.Type != TypeGlobalEC {
		fBody = f.Body.Clone()
	}
	rexp, err := fBody.Execute(fEC)
	if rexp != nil || err != nil {
		return rexp, err
	}
	return &LiteralExpression{
		Type:   LiteralTypeUndefined,
		Line:   e.Line,
		CharAt: e.CharAt,
	}, nil
}

func (e *CallExpression) IsTruthy() bool {
	return true
}

func (e *CallExpression) GetCharAt() int {
	return e.CharAt
}

func (e *CallExpression) GetLine() int {
	return e.Line
}

func (e *CallExpression) SetLine(i int) {
	e.Line = i
}

func (e *CallExpression) SetCharAt(i int) {
	e.CharAt = i
}

func (e *CallExpression) GetType() string {
	return "call expression"
}

func (e *CallExpression) ToString() string {
	return ""
}

func (e *CallExpression) Clone() Expression {
	args := []Expression{}
	for _, arg := range e.Arguments {
		args = append(args, arg.Clone())
	}
	return &CallExpression{
		Callee:    e.Callee.Clone(),
		Arguments: args,
		Line:      e.Line,
		CharAt:    e.Line,
	}
}

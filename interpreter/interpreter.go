package interpreter

import (
	"github.com/dhl1402/covidscript/builtin"
	"github.com/dhl1402/covidscript/core"
	"github.com/dhl1402/covidscript/lexer"
	"github.com/dhl1402/covidscript/parser"
)

func Interpret(script string) error {
	tokens := lexer.Lex(script)
	ast, err := parser.ToAST(tokens)
	if err != nil {
		return err
	}
	gec := builtin.CreateGlobalEC()
	if err = Execute(gec, ast); err != nil {
		return err
	}
	return nil
}

func Execute(gec *core.ExecutionContext, stmts []core.Statement) error {
	cexp := core.CallExpression{
		Callee: &core.FunctionExpression{
			Body:   stmts,
			Params: []core.Identifier{},
			EC:     gec,
		},
	}
	_, err := cexp.Evaluate(gec)
	return err
}

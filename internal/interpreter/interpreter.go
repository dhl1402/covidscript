package interpreter

import (
	"github.com/dhl1402/covidscript/internal/builtin"
	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/core"
	"github.com/dhl1402/covidscript/internal/lexer"
	"github.com/dhl1402/covidscript/internal/parser"
)

func Interpret(script string, conf config.Config) error {
	tokens, err := lexer.Lex(script)
	if err != nil {
		return err
	}
	ast, err := parser.ToAST(tokens)
	if err != nil {
		return err
	}
	gec := createGlobalEC(conf)
	if err = execute(gec, ast); err != nil {
		return err
	}
	return nil
}

func execute(gec *core.ExecutionContext, stmts []core.Statement) error {
	cexp := core.CallExpression{
		Callee: &core.FunctionExpression{
			Body: core.BlockStatement{
				Statements: stmts,
			},
			Params: []core.Identifier{},
			EC:     gec,
		},
	}
	_, err := cexp.Evaluate(gec)
	return err
}

func createGlobalEC(conf config.Config) *core.ExecutionContext {
	return &core.ExecutionContext{
		Type: core.TypeGlobalEC,
		Variables: map[string]core.Expression{
			"echo":    builtin.Echo(conf),
			"len":     builtin.Len(),
			"filter":  builtin.Filter(),
			"map":     builtin.Map(),
			"reduce":  builtin.Reduce(),
			"join":    builtin.Join(),
			"indexOf": builtin.IndexOf(),
			"append":  builtin.Append(),
			"sort":    builtin.Sort(),
			"keys":    builtin.Keys(),
			"values":  builtin.Values(),
			"type":    builtin.Type(),
			"delete":  builtin.Delete(),
			"neg":     builtin.Negative(),
			"floor":   builtin.Floor(),
			"ceil":    builtin.Ceil(),
		},
	}
}

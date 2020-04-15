package main

import (
	"fmt"

	"github.com/dhl1402/covidscript/builtin"
	"github.com/dhl1402/covidscript/executor"
	"github.com/dhl1402/covidscript/lexer"
	"github.com/dhl1402/covidscript/parser"
)

func main() {
	code := `
		var a = func(){
			echo("Hello","world!","This","is","covidscript")
		}
		a()
	`
	tokens := lexer.Lex(code)
	ast, err := parser.ToAST(tokens)
	if err != nil {
		fmt.Println(err)
		return
	}
	gec := builtin.CreateGlobalEC()
	if err = executor.Execute(gec, ast); err != nil {
		fmt.Println(err)
	}
}

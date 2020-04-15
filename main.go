package main

import (
	"covs/builtin"
	"covs/executor"
	"covs/lexer"
	"covs/parser"
	"fmt"
)

func main() {
	code := `
		var a = func(){
			echo("Hello","world!","This","is","covid-script")
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

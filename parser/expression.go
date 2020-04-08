package parser

import "gs/operator"

type Expression interface {
	Evaluate() Expression
}

type VariableExpression struct {
	Name   string
	Line   int
	CharAt int
}

func (e VariableExpression) Evaluate() Expression {
	// TODO
	return e
}

type LiteralExpression struct {
	Type   string
	Value  string
	Line   int
	CharAt int
}

func (e LiteralExpression) Evaluate() Expression {
	return e
}

type (
	ObjectProperty struct {
		KeyExpression Expression
		KeyIdentifier Identifier
		Value         Expression
		Computed      bool
		Shorthand     bool
		Method        bool
		Line          int
		CharAt        int
	}
	ObjectExpression struct {
		Properties []ObjectProperty
		Line       int
		CharAt     int
	}
)

func (e ObjectExpression) Evaluate() Expression {
	return e
}

type ArrayExpression struct {
	Elements []Expression
	Line     int
	CharAt   int
}

func (e ArrayExpression) Evaluate() Expression {
	return e
}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator operator.Operator
	Group    bool
	Nesting  int
	Line     int
	CharAt   int
}

func (e BinaryExpression) Evaluate() Expression {
	// switch case operator
	return LiteralExpression{Value: "TODO", Type: "string"}
}

type FunctionExpression struct {
	Params []Identifier
	Body   []Statement
	Line   int
	CharAt int
}

func (e FunctionExpression) Evaluate() Expression {
	return e
}

type CallExpression struct {
	Callee    Expression
	Arguments []Expression
	Line      int
	CharAt    int
}

func (e CallExpression) Evaluate() Expression {
	return e
}

type MemberExpression struct {
	Object             Expression
	PropertyIdentifier Identifier
	PropertyExpression Expression
	Computed           bool
	Line               int
	CharAt             int
}

func (e MemberExpression) Evaluate() Expression {
	return e
}

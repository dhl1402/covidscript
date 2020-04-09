package parser

import "gs/operator"

type Expression interface {
	Evaluate() Expression
	GetCharAt() int
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

func (e VariableExpression) GetCharAt() int {
	return e.CharAt
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

func (e LiteralExpression) GetCharAt() int {
	return e.CharAt
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

func (e ObjectExpression) GetCharAt() int {
	return e.CharAt
}

type ArrayExpression struct {
	Elements []Expression
	Line     int
	CharAt   int
}

func (e ArrayExpression) Evaluate() Expression {
	return e
}

func (e ArrayExpression) GetCharAt() int {
	return e.CharAt
}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator operator.Operator
	Group    bool
	Line     int
	CharAt   int
}

func (e BinaryExpression) Evaluate() Expression {
	// switch case operator
	return LiteralExpression{Value: "TODO", Type: "string"}
}

func (e BinaryExpression) GetCharAt() int {
	return e.CharAt
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

func (e FunctionExpression) GetCharAt() int {
	return e.CharAt
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

func (e CallExpression) GetCharAt() int {
	return e.CharAt
}

type MemberAccessExpression struct {
	Object             Expression
	PropertyIdentifier Identifier
	PropertyExpression Expression
	Computed           bool
	Line               int
	CharAt             int
}

func (e MemberAccessExpression) Evaluate() Expression {
	return e
}

func (e MemberAccessExpression) GetCharAt() int {
	return e.CharAt
}

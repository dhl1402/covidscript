package core

type Statement interface{}

type Identifier struct {
	Name   string
	Line   int
	CharAt int
}

type (
	VariableDeclarator struct {
		ID     Identifier
		Init   Expression
		Line   int
		CharAt int
	}
	VariableDeclaration struct {
		Declarations []VariableDeclarator
		Line         int
		CharAt       int
	}
)

type ExpressionStatement struct {
	Expression
	Line   int
	CharAt int
}

type AssignmentStatement struct {
	Left   Expression
	Right  Expression
	Line   int
	CharAt int
}

// type IfStatement struct {
// 	Test Expression
// 	Consequent []Statement
// 	Alternate Statement
// }

type FunctionDeclaration struct {
	ID     Identifier
	Params []Identifier
	Body   []Statement
	Line   int
	CharAt int
}

type ReturnStatement struct {
	Argument Expression
	Line     int
	CharAt   int
}

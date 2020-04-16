package core

type IfStatement struct {
	Test       Expression
	Consequent []Statement
	Alternate  Statement
}

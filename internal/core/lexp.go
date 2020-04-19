package core

type PrimitiveType string

const (
	LiteralTypeNumber    PrimitiveType = "number"
	LiteralTypeString    PrimitiveType = "string"
	LiteralTypeBoolean   PrimitiveType = "boolean"
	LiteralTypeNull      PrimitiveType = "null"
	LiteralTypeUndefined PrimitiveType = "undefined"
)

type LiteralExpression struct {
	Type   PrimitiveType
	Value  string
	Line   int
	CharAt int
}

func (e *LiteralExpression) Evaluate(ec *ExecutionContext) (Expression, error) {
	return e, nil
}

func (e *LiteralExpression) IsTruthy() bool {
	switch e.Type {
	case LiteralTypeNull:
	case LiteralTypeUndefined:
		return false
	case LiteralTypeNumber:
		return e.Value != "0"
	case LiteralTypeString:
		return e.Value != ""
	case LiteralTypeBoolean:
		return e.Value == "#t"
	}
	// never get here
	return false
}

func (e *LiteralExpression) GetCharAt() int {
	return e.CharAt
}

func (e *LiteralExpression) GetLine() int {
	return e.Line
}

func (e *LiteralExpression) SetLine(i int) {
	e.Line = i
}

func (e *LiteralExpression) SetCharAt(i int) {
	e.CharAt = i
}

func (e *LiteralExpression) GetType() string {
	return string(e.Type)
}

func (e *LiteralExpression) ToString() string {
	return e.Value
}

package lexer

import (
	"gs/operator"
	"gs/symbol"
	"strconv"
)

type Token struct {
	Value  string
	Line   int
	CharAt int
}

func (t Token) IsIdentifier() bool {
	// TODO: regex to test identifier
	s := t.Value
	return s != "" && !symbol.IsReservedKeyword(s) && !symbol.IsSpecialChars(s) && !symbol.IsStringBoundary(s) && !symbol.IsWhiteSpace(s) && !symbol.IsNewLine(s)
}

func (t Token) IsNumber() bool {
	_, err := strconv.Atoi(t.Value)
	if err == nil {
		return true
	}
	_, err = strconv.ParseFloat(t.Value, 64)
	return err == nil
}

func (t Token) IsString() bool {
	return t.Value != "" && symbol.IsStringBoundary(string(t.Value[0]))
}

func (t Token) IsBoolean() bool {
	return t.Value == "false" || t.Value == "true"
}

func (t Token) IsOperatorSymbol() bool {
	return operator.IsOperatorSymbol(t.Value)
}

func (t Token) IsPrimitiveValue() bool {
	return t.IsBoolean() || t.IsNumber() || t.IsString()
}

func (t Token) ParsePrimitiveType() (string, bool) {
	if t.IsBoolean() {
		return "boolean", true
	}
	if t.IsNumber() {
		return "number", true
	}
	if t.IsString() {
		return "string", true
	}
	return "", false
}

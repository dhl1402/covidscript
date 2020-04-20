package lexer

import (
	"regexp"
	"strconv"

	"github.com/dhl1402/covidscript/internal/core"
	"github.com/dhl1402/covidscript/internal/utils"
)

type Token struct {
	Value  string
	Line   int
	CharAt int
}

func (t Token) IsIdentifier() bool {
	s := t.Value
	if match, err := regexp.MatchString("^[a-zA-Z_$][a-zA-Z_$0-9]*$", s); !match || err != nil {
		return false
	}
	return s != "" && !t.IsOperatorSymbol() && !utils.IsReservedKeyword(s) && !utils.IsSpecialChars(s) && !utils.IsStringBoundary(s) && !utils.IsWhiteSpace(s) && !utils.IsNewLine(s)
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
	return t.Value != "" && utils.IsStringBoundary(string(t.Value[0]))
}

func (t Token) IsBoolean() bool {
	return t.Value == "#f" || t.Value == "#t"
}

func (t Token) IsOperatorSymbol() bool {
	return core.IsOperatorSymbol(t.Value)
}

func (t Token) IsPrimitiveValue() bool {
	return t.IsBoolean() || t.IsNumber() || t.IsString()
}

func (t Token) ParsePrimitiveType() (core.PrimitiveType, bool) {
	if t.Value == "undefined" {
		return "undefined", true
	}
	if t.Value == "null" {
		return "null", true
	}
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

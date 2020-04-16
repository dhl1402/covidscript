package lexer

import (
	"fmt"

	"github.com/dhl1402/covidscript/utils"
)

// todo: lex float, :=
// Lex source code into tokens
func Lex(sc string) ([]Token, error) {
	tokens := []Token{}
	tmp := ""
	line := 1
	charAt := 1
	for len(sc) > 0 {
		c := string(sc[0])
		if utils.IsStringBoundary(c) {
			str, err := lexString(sc)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, Token{Value: str, Line: line, CharAt: charAt})
			sc = sc[len(str):]
			charAt = charAt + len(str)
			continue
		}
		if utils.IsReservedKeyword(c) || utils.IsSpecialChars(c) {
			if tmp != "" {
				tokens = append(tokens, Token{Value: tmp, Line: line, CharAt: charAt - len(tmp)})
				tmp = ""
			}
			tokens = append(tokens, Token{Value: c, Line: line, CharAt: charAt})
		} else if utils.IsWhiteSpace(c) || utils.IsNewLine(c) {
			if tmp != "" {
				tokens = append(tokens, Token{Value: tmp, Line: line, CharAt: charAt - len(tmp)})
				tmp = ""
			}
		} else {
			tmp = tmp + c
		}

		sc = sc[1:]
		if utils.IsNewLine(c) {
			line++
			charAt = 1
		} else if !utils.IsWhiteSpace(c) || charAt != 1 {
			charAt++
		}
	}
	if tmp != "" {
		tokens = append(tokens, Token{Value: tmp, Line: line, CharAt: charAt - len(tmp)})
	}
	return tokens, nil
}

func lexString(sc string) (string, error) {
	result := ""
	openStrChar := ""
	for _, s := range sc {
		c := string(s)
		if openStrChar == "" && utils.IsStringBoundary(c) {
			openStrChar = c
			result = result + c
		} else if openStrChar != "" {
			if c == openStrChar {
				return result + c, nil
			}
			result = result + c
		}
	}
	return result, fmt.Errorf("Missing closing quote.")
}

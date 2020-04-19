package lexer

import (
	"fmt"

	"github.com/dhl1402/covidscript/internal/utils"
)

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
		if op := lexMultipleCharOperator(sc); op != "" {
			if tmp != "" {
				tokens = append(tokens, Token{Value: tmp, Line: line, CharAt: charAt - len(tmp)})
				tmp = ""
			}
			tokens = append(tokens, Token{Value: op, Line: line, CharAt: charAt})
			charAt = charAt + len(op) - 1
			sc = sc[len(op)-1:]
		} else if utils.IsReservedKeyword(c) || utils.IsSpecialChars(c) {
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
		if sc == "" && tmp != "" {
			tokens = append(tokens, Token{Value: tmp, Line: line, CharAt: charAt - len(tmp)})
		}
	}
	result := []Token{}
	var i int
	for i = 0; i < len(tokens); i++ {
		t := tokens[i]
		var nt *Token
		if i+1 < len(tokens) {
			nt = &tokens[i+1]
		}
		var lt *Token
		if len(result) > 0 {
			lt = &result[len(result)-1]
		}
		if t.Value == "." && lt != nil && utils.IsInteger(lt.Value) && nt != nil && utils.IsInteger(nt.Value) {
			result = result[:len(result)-1]
			result = append(result, Token{
				Value:  fmt.Sprintf("%s.%s", lt.Value, nt.Value),
				Line:   lt.Line,
				CharAt: lt.CharAt,
			})
			i++ // nt is processed
		} else {
			result = append(result, t)
		}
	}
	return result, nil
}

func lexMultipleCharOperator(sc string) string {
	operators := []string{":=", "<=", ">=", "===", "==", "!==", "!=", "&&", "||"} // order matter
	for _, op := range operators {
		for i, r := range sc {
			s := string(r)
			if i >= len(op) {
				return op
			}
			if s != string(op[i]) {
				break
			}
		}
	}
	return ""
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

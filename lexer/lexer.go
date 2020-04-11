package lexer

import (
	"gs/symbol"
)

// todo: lex float, :=
// Lex source code into tokens
func Lex(sc string) (tokens []Token) {
	tmp := ""
	line := 1
	charAt := 1
	for len(sc) > 0 {
		c := string(sc[0])
		if symbol.IsStringBoundary(c) {
			str := lexString(sc)
			tokens = append(tokens, Token{Value: str, Line: line, CharAt: charAt})
			sc = sc[len(str):]
			charAt = charAt + len(str)
			continue
		}
		if symbol.IsReservedKeyword(c) || symbol.IsSpecialChars(c) {
			if tmp != "" {
				tokens = append(tokens, Token{Value: tmp, Line: line, CharAt: charAt - len(tmp)})
				tmp = ""
			}
			tokens = append(tokens, Token{Value: c, Line: line, CharAt: charAt})
		} else if symbol.IsWhiteSpace(c) || symbol.IsNewLine(c) {
			if tmp != "" {
				tokens = append(tokens, Token{Value: tmp, Line: line, CharAt: charAt - len(tmp)})
				tmp = ""
			}
		} else {
			tmp = tmp + c
		}

		sc = sc[1:]
		if symbol.IsNewLine(c) {
			line++
			charAt = 1
		} else if !symbol.IsWhiteSpace(c) || charAt != 1 {
			charAt++
		}
	}
	if tmp != "" {
		tokens = append(tokens, Token{Value: tmp, Line: line, CharAt: charAt - len(tmp)})
	}
	return
}

func lexString(sc string) (result string) {
	openStrChar := ""
	for _, s := range sc {
		c := string(s)
		if openStrChar == "" && symbol.IsStringBoundary(c) {
			openStrChar = c
			result = result + c
		} else if openStrChar != "" {
			if c == openStrChar {
				return result + c
			}
			result = result + c
		}
	}
	// TODO: handle missing close string character
	return ""
}

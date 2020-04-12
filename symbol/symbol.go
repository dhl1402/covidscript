package symbol

import (
	"gs/utils"
)

// hardcoding: = , true false { } [] func

func IsReservedKeyword(s string) bool {
	ss := []string{"var", "func", "return", "true", "false"}
	return utils.IncludeStr(ss, s)
}

func IsSpecialChars(s string) bool {
	ss := []string{"=", ":", ",", ".", "(", ")", "{", "}", "[", "]", "\"", "'", "`", "+", "-", "*", "/", "%", ";"}
	return utils.IncludeStr(ss, s)
}

func IsStringBoundary(s string) bool {
	ss := []string{"\"", "'", "`"}
	return utils.IncludeStr(ss, s)
}

func IsWhiteSpace(s string) bool {
	ss := []string{" ", "\t", "\b"}
	return utils.IncludeStr(ss, s)
}

func IsNewLine(s string) bool {
	return s == "\n"
}

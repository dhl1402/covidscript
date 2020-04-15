package utils

func IsReservedKeyword(s string) bool {
	ss := []string{"var", "func", "return", "true", "false", "null", "undefined"}
	return IncludeStr(ss, s)
}

func IsSpecialChars(s string) bool {
	ss := []string{"=", ":", ",", ".", "(", ")", "{", "}", "[", "]", "\"", "'", "`", "+", "-", "*", "/", "%", ";"}
	return IncludeStr(ss, s)
}

func IsStringBoundary(s string) bool {
	ss := []string{"\"", "'", "`"}
	return IncludeStr(ss, s)
}

func IsWhiteSpace(s string) bool {
	ss := []string{" ", "\t", "\b"}
	return IncludeStr(ss, s)
}

func IsNewLine(s string) bool {
	return s == "\n"
}

func IncludeStr(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

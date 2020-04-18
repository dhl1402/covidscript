package utils

import "strconv"

func IsReservedKeyword(s string) bool {
	ss := []string{"var", "func", "return", "if", "else", "elif", "#t", "#f", "null", "undefined"}
	return IncludeStr(ss, s)
}

func IsSpecialChars(s string) bool {
	ss := []string{"=", ":", ",", ".", "(", ")", "{", "}", "[", "]", "\"", "'", "`", "+", "-", "*", "/", "%", "<", ">", ";", "!"}
	return IncludeStr(ss, s)
}

func IsInteger(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
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

func ToBoolStr(b bool) string {
	if b {
		return "#t"
	}
	return "#f"
}

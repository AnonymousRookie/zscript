package utils

import (
	"fmt"
	"os"
	// "runtime"
)

func IsNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func IsAlphanumeric(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func IsDelimiter(r rune) bool {
	return r == ',' || r == '(' || r == ')' || r == '{' || r == '}' || r == ';'
}

func IsOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' || r == '='
}

func Exit(lineNum int, errmsg string) {
	fmt.Printf("[error]line: %d, %s!\n", lineNum, errmsg)

	// var buf [2048]byte
	// runtime.Stack(buf[:], true)
	// fmt.Println(string(buf[:]))

	os.Exit(-1)
}

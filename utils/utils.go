package utils

import (
	"fmt"
	"os"
	// "runtime"
	"bufio"
	"io"
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

func ExitWithErrMsg(err string) {
	fmt.Printf("[error] %s\n", err)
	os.Exit(-1)
}

func Check(cond bool, err string) {
	if !cond {
		fmt.Printf("[error] %s\n", err)
		os.Exit(-1)
	}
}

// 将sourceFilename文件中的所有行读入到lines中
func LoadSourceFile(sourceFilename string) (lines []string) {
	sourceFile, err := os.Open(sourceFilename)
	if err != nil {
		fmt.Println("loadSourceFile err!")
		return lines
	}
	defer sourceFile.Close()

	sourceReader := bufio.NewReader(sourceFile)
	for {
		line, err := sourceReader.ReadString('\n')
		lines = append(lines, line)
		if err == io.EOF {
			return lines
		}
	}
}

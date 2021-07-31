package zcomplier

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func loadSourceFile(sourceFilename string) (lines []string) {
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
	return lines
}

func Complie(sourceFilename string) {
	// 添加2个临时变量
	TempVar0SymbolIndex = addSymbol("_T0", SymbolTypeVar, GlobalScope)
	TempVar1SymbolIndex = addSymbol("_T1", SymbolTypeVar, GlobalScope)

	// 读取源文件
	lines := loadSourceFile(sourceFilename)

	// 词法分析
	lexicalAnalyze(lines)

	// 语法分析
	parse()

	// 代码生成
	generateCode()
}

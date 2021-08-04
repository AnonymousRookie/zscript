package zcomplier

import (
	"../../utils"
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
	TempVar0SymbolIndex = addSymbol(registerT0, SymbolTypeVar, GlobalScope)
	TempVar1SymbolIndex = addSymbol(registerT1, SymbolTypeVar, GlobalScope)

	suffix := sourceFilename[len(sourceFilename)-3:]
	if suffix != srcFileSuffix {
		utils.ExitWithErrMsg("source file suffix should be: " + srcFileSuffix)
	}

	// 读取源文件
	lines := loadSourceFile(sourceFilename)

	// 词法分析
	lexicalAnalyze(lines)

	// 语法分析
	parse()

	var outputFilename string
	outputFilename = sourceFilename[:len(sourceFilename)-3] + outFileSuffix

	// 代码生成
	generateCode(outputFilename)
}

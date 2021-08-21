package zcomplier

import (
	"zscript/utils"
)

func Complie(sourceFilename string) {

	suffix := sourceFilename[len(sourceFilename)-3:]
	if suffix != srcFileSuffix {
		utils.ExitWithErrMsg("source file suffix should be: " + srcFileSuffix)
	}

	// 读取源文件
	lines := utils.LoadSourceFile(sourceFilename)

	// 词法分析
	lexer := NewLexer()
	lexer.lexicalAnalyze(lines)

	// 语法分析
	parser := NewParser(lexer)
	parser.parse()

	// 代码生成
	gen := NewGenerator()
	outputFilename := sourceFilename[:len(sourceFilename)-3] + outFileSuffix
	gen.generateCode(outputFilename)
}

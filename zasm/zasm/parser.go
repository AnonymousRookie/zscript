package zasm

import (
	"fmt"
	"strconv"
)

var curFuncIndex int32 = 0      // 函数索引
var curFuncParamCount int32 = 0 // 当前函数参数个数

func parse() {

	header.isExistMainFunc = false

	for {
		tokenType := parseStatement()
		if tokenType == TokenTypeEndOfStream {
			break
		}
	}
}

func parseStatement() TokenType {
	token := nextToken()

	switch token.Type {
	case TokenTypeEndOfStream:
		fmt.Println("parseStatement() TokenTypeEndOfStream")

	case TokenTypeVar:
		parseVar()

	case TokenTypeParam:
		parseParam()

	case TokenTypeInstr:
		parseInstr()

	case TokenTypeFunc:
		parseFunc()

	case TokenTypeOpenCurlyBrace:
		parseBlock()
	}

	return token.Type
}

func parseInstr() {
	backToken()
	token := nextToken()
	instrType, instrOpCount := getInstrType(&token)

	// fmt.Printf("parseInstr: %+v\n", token)

	var instr Instr
	instr.instrType = instrType
	instr.opCount = instrOpCount

	for i := 0; i < int(instrOpCount); i++ {
		token := nextToken()

		var op Operand

		switch token.Type {
		case TokenTypeInt:
			op.opType = OperandTypeInt
			val, _ := strconv.Atoi(token.Lexem)
			op.opVal = int32(val)
		case TokenTypeFloat:
			op.opType = OperandTypeFloat
			f, _ := strconv.ParseFloat(token.Lexem, 32)
			op.opVal = float32(f)
		case TokenTypeString:
			op.opType = OperandTypeStrIndex
			strindex := addStr(token.Lexem)
			op.opVal = strindex
		case TokenTypeFunc:
			op.opType = OperandTypeFuncIndex
			op.opVal = curFuncIndex
		case TokenTypeIdentifier:
			node := getSymbol(token.Lexem, curFuncIndex)
			if node != nil {
				op.opType = OperandTypeIdentifierIndex
				op.opVal = node.index
			} else {
				funcNode := getFuncNodeByName(token.Lexem)
				op.opType = OperandTypeFuncIndex
				op.opVal = funcNode.index
			}
		case TokenTypeRetVal:
			op.opType = OperandTypeReg
			op.opVal = int32(0)
		default:
			fmt.Printf("unexpected token type, token:%+v\n", token)
		}

		// fmt.Printf("op: %+v\n", op)

		instr.ops = append(instr.ops, op)

		// 操作数之间以逗号分隔
		if i != int(instrOpCount-1) {
			token := nextToken()
			checkToken(&token, TokenTypeComma)
		}
	}

	istrIndex := addInstr(instr)

	// 更新函数入口点（函数第一个指令的索引）
	funcNode := getFuncNodeByIndex(curFuncIndex)
	// fmt.Printf("parseInstr: %+v\n", funcNode)
	if funcNode.entryPoint < 0 {
		funcNode.entryPoint = istrIndex
	}
}

func parseFunc() {
	var token Token
	token = nextToken()
	checkToken(&token, TokenTypeIdentifier)
	funcName := token.Lexem
	funcIndex := addFuncNode(funcName)
	curFuncIndex = funcIndex

	if funcName == "main" {
		header.isExistMainFunc = true
		header.mainFuncIndex = funcIndex
	}

	token = nextToken()
	checkToken(&token, TokenTypeOpenCurlyBrace)

	parseBlock()

	curFuncIndex = 0
	curFuncParamCount = 0
}

func parseVar() {
	token := nextToken()
	checkToken(&token, TokenTypeIdentifier)
	addSymbol(token.Lexem, SymbolTypeVar, curFuncIndex)
	// fmt.Printf("parseVar token:%+v, %v\n", token, curFuncIndex)
}

func parseParam() {
	token := nextToken()
	checkToken(&token, TokenTypeIdentifier)
	addSymbol(token.Lexem, SymbolTypeParam, curFuncIndex)
	funcNode := getFuncNodeByIndex(curFuncIndex)
	funcNode.paramcount++
}

func parseBlock() {
	for {
		token := nextToken()
		if token.Lexem == "}" {
			return
		}

		if token.Type == TokenTypeEndOfStream {
			token := backToken()
			checkToken(&token, TokenTypeCloseCurlyBrace)
		}

		backToken()
		parseStatement()
	}
}

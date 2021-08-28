package zasm

import (
	"fmt"
	"strconv"
)

type Parser struct {
	curFuncIndex      int32 // 函数索引
	curFuncParamCount int32 // 当前函数参数个数
	lexer             *Lexer
}

func NewParser(l *Lexer) *Parser {
	return &Parser{curFuncIndex: 0, curFuncParamCount: 0, lexer: l}
}

func (parser *Parser) parse() {

	header.isExistMainFunc = false

	for {
		tokenType := parser.parseStatement()
		if tokenType == TokenTypeEndOfStream {
			break
		}
	}
}

func (parser *Parser) parseStatement() TokenType {
	token := parser.lexer.nextToken()

	switch token.Type {
	case TokenTypeEndOfStream:
		fmt.Println("parseStatement() TokenTypeEndOfStream")

	case TokenTypeVar:
		parser.parseVar()

	case TokenTypeParam:
		parser.parseParam()

	case TokenTypeInstr:
		parser.parseInstr()

	case TokenTypeFunc:
		parser.parseFunc()

	case TokenTypeOpenCurlyBrace:
		parser.parseBlock()
	}

	return token.Type
}

func (parser *Parser) parseInstr() {
	parser.lexer.backToken()
	token := parser.lexer.nextToken()
	instrType, instrOpCount := parser.lexer.getInstrType(&token)

	// fmt.Printf("parseInstr: %+v\n", token)

	var instr Instr
	instr.instrType = instrType
	instr.opCount = instrOpCount

	for i := 0; i < int(instrOpCount); i++ {
		token := parser.lexer.nextToken()

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
			op.opVal = parser.curFuncIndex
		case TokenTypeHostApiPrint:
			op.opType = OperandTypeHostApiIndex
			hostApiNode := getHostApiNode(token.Lexem)
			op.opVal = hostApiNode.index
		case TokenTypeIdentifier:
			node := getSymbol(token.Lexem, parser.curFuncIndex)
			if node != nil {
				op.opType = OperandTypeIdentifierIndex
				op.opVal = node.index
			} else {
				reg, ok := regTypeMap[token.Lexem]
				if ok {
					op.opType = OperandTypeReg
					op.opVal = reg
				} else {
					funcNode := getFuncNodeByName(token.Lexem)
					op.opType = OperandTypeFuncIndex
					op.opVal = funcNode.index
				}
			}
		default:
			fmt.Printf("unexpected token type, token:%+v\n", token)
		}

		// fmt.Printf("op: %+v\n", op)

		instr.ops = append(instr.ops, op)

		// 操作数之间以逗号分隔
		if i != int(instrOpCount-1) {
			token := parser.lexer.nextToken()
			parser.lexer.checkToken(&token, TokenTypeComma)
		}
	}

	istrIndex := addInstr(instr)

	// 更新函数入口点（函数第一个指令的索引）
	funcNode := getFuncNodeByIndex(parser.curFuncIndex)
	// fmt.Printf("parseInstr: %+v\n", funcNode)
	if funcNode.entryPoint < 0 {
		funcNode.entryPoint = istrIndex
	}
}

func (parser *Parser) parseFunc() {
	var token Token
	token = parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeIdentifier)
	funcName := token.Lexem
	funcIndex := addFuncNode(funcName)
	parser.curFuncIndex = funcIndex

	if funcName == "main" {
		header.isExistMainFunc = true
		header.mainFuncIndex = funcIndex
	}

	token = parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeOpenCurlyBrace)

	parser.parseBlock()

	parser.curFuncIndex = 0
	parser.curFuncParamCount = 0
}

func (parser *Parser) parseVar() {
	token := parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeIdentifier)
	addSymbol(token.Lexem, SymbolTypeVar, parser.curFuncIndex)
	// fmt.Printf("parseVar token:%+v, %v\n", token, parser.curFuncIndex)
}

func (parser *Parser) parseParam() {
	token := parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeIdentifier)
	addSymbol(token.Lexem, SymbolTypeParam, parser.curFuncIndex)
	funcNode := getFuncNodeByIndex(parser.curFuncIndex)
	funcNode.paramcount++
}

func (parser *Parser) parseBlock() {
	for {
		token := parser.lexer.nextToken()
		if token.Lexem == "}" {
			return
		}

		if token.Type == TokenTypeEndOfStream {
			token := parser.lexer.backToken()
			parser.lexer.checkToken(&token, TokenTypeCloseCurlyBrace)
		}

		parser.lexer.backToken()
		parser.parseStatement()
	}
}

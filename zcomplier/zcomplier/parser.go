package zcomplier

import (
	"fmt"
	"strconv"

	"zscript/utils"
)

const (
	GlobalScope = 0 // 函数外
)

type Parser struct {
	curScope int
	lexer    *Lexer
}

func NewParser(l *Lexer) *Parser {
	return &Parser{curScope: GlobalScope, lexer: l}
}

func (parser *Parser) parse() {

	for {
		tokenType := parser.parseStatement()
		if tokenType == TokenTypeEndOfStream {
			break
		}
	}

	// printFuncTable()
	// printSymbolList()
}

func (parser *Parser) parseStatement() TokenType {
	token := parser.lexer.nextToken()

	switch token.Type {
	case TokenTypeEndOfStream:
		fmt.Println("parseStatement() TokenTypeEndOfStream")

	case TokenTypeVar:
		parser.parseVar()

	case TokenTypeIdentifier:
		if getSymbolNode(token.Lexem, parser.curScope) != nil {
			// 赋值语句
			parser.parseAssignment(token)
		} else if getFuncByName(token.Lexem) != nil {
			// 函数调用
			addICodeNodeSourceLine(parser.curScope, token.strLine)
			parser.parseFuncCall(token.Lexem)
		} else {
			utils.Exit(token.LineNumber, "invalid identifier "+token.Lexem)
		}

	case TokenTypeReturn:
		parser.parseReturn()

	case TokenTypeFunc:
		parser.parseFunc()

	case TokenTypeOpenCurlyBrace:
		parser.parseBlock()
	}

	return token.Type
}

func (parser *Parser) parseVar() {
	var token Token

	token = parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeIdentifier)

	// 将变量添加到SymbolTable中
	addSymbol(token.Lexem, SymbolTypeVar, parser.curScope)

	token = parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeSemicolon)
}

func (parser *Parser) parseReturn() {
	var token Token
	token = parser.lexer.nextToken()

	addICodeNodeSourceLine(parser.curScope, token.strLine)

	if token.Lexem != ";" {
		parser.lexer.backToken()

		parser.parseExpr()

		instrIndex := addICodeNodeInstruction(parser.curScope, InstrTypePop)
		addOperandReg(parser.curScope, instrIndex, registerTypeRetVal)

		token = parser.lexer.nextToken()
		parser.lexer.checkToken(&token, TokenTypeSemicolon)
	}
	addICodeNodeInstruction(parser.curScope, InstrTypeRet)
}

// 表达式
func (parser *Parser) parseExpr() {
	parser.parseSubExpr()
	for {

		token := parser.lexer.nextToken()
		if token.Type != TokenTypeOperator {
			parser.lexer.backToken()
			break
		}
		parser.parseSubExpr()
	}
}

// 子表达式(A+B-C)
func (parser *Parser) parseSubExpr() {
	var instrIndex int

	parser.parseTerm()

	for {

		token := parser.lexer.nextToken()
		operatorType := parser.lexer.getOperatorType(token)
		if token.Type != TokenTypeOperator || (operatorType != operatorTypeAdd && operatorType != operatorTypeSub) {
			parser.lexer.backToken()
			break
		}

		parser.parseTerm()

		// 取出第一个操作数存入_T1
		instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePop)
		addOperandReg(parser.curScope, instrIndex, registerTypeT1)

		// 取出第一个操作数存入_T0
		instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePop)
		addOperandReg(parser.curScope, instrIndex, registerTypeT0)

		var instrType InstrType = InstrTypeInvalid
		switch operatorType {
		case operatorTypeAdd:
			instrType = InstrTypeAdd
		case operatorTypeSub:
			instrType = InstrTypeSub
		default:
			instrType = InstrTypeInvalid
		}

		instrIndex = addICodeNodeInstruction(parser.curScope, instrType)
		addOperandReg(parser.curScope, instrIndex, registerTypeT0)
		addOperandReg(parser.curScope, instrIndex, registerTypeT1)

		// 把结果存入_T0
		instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePush)
		addOperandReg(parser.curScope, instrIndex, registerTypeT0)
	}
}

// 子表达式中的项(a*b/c)
func (parser *Parser) parseTerm() {
	var instrIndex int

	parser.parseFactor()

	for {

		token := parser.lexer.nextToken()
		operatorType := parser.lexer.getOperatorType(token)
		if token.Type != TokenTypeOperator || (operatorType != operatorTypeMul && operatorType != operatorTypeDiv) {
			parser.lexer.backToken()
			break
		}

		parser.parseFactor()

		// 取出第一个操作数存入_T1
		instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePop)
		addOperandReg(parser.curScope, instrIndex, registerTypeT1)

		// 取出第一个操作数存入_T0
		instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePop)
		addOperandReg(parser.curScope, instrIndex, registerTypeT0)

		var instrType InstrType = InstrTypeInvalid
		switch operatorType {
		case operatorTypeMul:
			instrType = InstrTypeMul
		case operatorTypeDiv:
			instrType = InstrTypeDiv
		default:
			instrType = InstrTypeInvalid
		}

		utils.Check(instrType != InstrTypeInvalid, "instrType is InstrTypeInvalid!")

		instrIndex = addICodeNodeInstruction(parser.curScope, instrType)
		addOperandReg(parser.curScope, instrIndex, registerTypeT0)
		addOperandReg(parser.curScope, instrIndex, registerTypeT1)

		// 把结果存入_T0
		instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePush)
		addOperandReg(parser.curScope, instrIndex, registerTypeT0)
	}
}

// 项中的因子
func (parser *Parser) parseFactor() {
	var instrIndex int
	var token Token
	token = parser.lexer.nextToken()

	// fmt.Printf("parseFactor: %+v\n", token)

	switch token.Type {
	case TokenTypeInt:
		instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePush)
		i, _ := strconv.Atoi(token.Lexem)
		addOperandInt(parser.curScope, instrIndex, i)
	case TokenTypeFloat:
		instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePush)
		f, _ := strconv.ParseFloat(token.Lexem, 32)
		addOperandFloat(parser.curScope, instrIndex, float32(f))
	case TokenTypeString:
		instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePush)
		strIndex := addstr(token.Lexem)
		addOperandStr(parser.curScope, instrIndex, strIndex)
	case TokenTypeIdentifier:
		symbolNode := getSymbolNode(token.Lexem, parser.curScope)

		if symbolNode != nil {
			// 变量
			instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePush)
			addOperandVar(parser.curScope, instrIndex, symbolNode.index)
		} else {
			// 函数
			funcNode := getFuncByName(token.Lexem)
			if funcNode != nil {
				parser.parseFuncCall(token.Lexem)
				instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePush)
				addOperandReg(parser.curScope, instrIndex, registerTypeRetVal)
			}
		}

	case TokenTypeOpenParen:
		parser.parseExpr()
		token = parser.lexer.nextToken()
		parser.lexer.checkToken(&token, TokenTypeCloseParen)

	default:
		utils.Exit(token.LineNumber, "invalid input: "+token.Lexem)
	}
}

func (parser *Parser) parseFunc() {

	var token Token
	token = parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeIdentifier)
	funcName := token.Lexem
	addFunc(funcName)

	token = parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeOpenParen)

	// 函数参数列表
	var paramList []string
	token = parser.lexer.nextToken()
	for token.Lexem != ")" {
		paramList = append(paramList, token.Lexem)
		token = parser.lexer.nextToken()
		if token.Lexem == ")" {
			break
		}
		parser.lexer.checkToken(&token, TokenTypeComma)
		token = parser.lexer.nextToken()
	}

	funcNode := getFuncByName(funcName)
	if funcNode != nil {
		funcNode.ParamCount = len(paramList)
		parser.curScope = funcNode.FuncIndex
		// fmt.Println("parser.curScope:", funcNode.FuncIndex)
	}

	// 将函数参数添加到SymbolTable中
	for i := 0; i < len(paramList); i++ {
		addSymbol(paramList[len(paramList)-i-1], SymbolTypeParam, parser.curScope)
	}

	token = parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeOpenCurlyBrace)

	parser.parseBlock()

	parser.curScope = GlobalScope
}

func (parser *Parser) parseBlock() {
	for {
		token := parser.lexer.nextToken()
		if token.Lexem == "}" {
			return
		}

		if token.Type == TokenTypeEndOfStream {
			parser.lexer.backToken()
			parser.lexer.checkToken(parser.lexer.curToken(), TokenTypeCloseCurlyBrace)
		}

		parser.lexer.backToken()
		parser.parseStatement()
	}
}

// 函数调用
func (parser *Parser) parseFuncCall(funcName string) {
	var token Token

	funcNode := getFuncByName(funcName)
	if funcNode == nil {
		return
	}

	token = parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeOpenParen)

	var paramCount int = 0

	// 函数参数
	for {
		token = parser.lexer.nextToken()
		if token.Lexem == ")" {
			break
		}

		parser.lexer.backToken()

		parser.parseExpr()
		paramCount++

		token = parser.lexer.nextToken()
		if token.Lexem == ")" {
			break
		}
		parser.lexer.checkToken(&token, TokenTypeComma)
	}

	if paramCount != funcNode.ParamCount {
		var err string = "too few parameters!"
		if paramCount > funcNode.ParamCount {
			err = "too mutch parameters!"
			fmt.Println(paramCount, funcNode.ParamCount)
		}
		utils.Exit(token.LineNumber, err)
	}

	instrIndex := addICodeNodeInstruction(parser.curScope, InstrTypeCall)
	addOperandFuncIndex(parser.curScope, instrIndex, funcNode.FuncIndex)
}

// 赋值语句
func (parser *Parser) parseAssignment(token Token) {

	symbol := getSymbolNode(token.Lexem, parser.curScope)
	// fmt.Printf("%+v\n", symbol)

	token = parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeOperator)

	addICodeNodeSourceLine(parser.curScope, token.strLine)

	if token.Lexem != "=" {
		utils.Exit(token.LineNumber, "invalid operator "+token.Lexem)
	}

	parser.parseExpr()

	token = parser.lexer.nextToken()
	parser.lexer.checkToken(&token, TokenTypeSemicolon)

	var instrIndex = InstrTypeInvalid
	instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypePop)
	addOperandReg(parser.curScope, instrIndex, registerTypeT0)

	instrIndex = addICodeNodeInstruction(parser.curScope, InstrTypeMov)
	addOperandVar(parser.curScope, instrIndex, symbol.index)
	addOperandReg(parser.curScope, instrIndex, registerTypeT0)
}

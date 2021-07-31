package zcomplier

import (
	"../../utils"
	"fmt"
	"strconv"
)

const (
	GlobalScope = 0 // 函数外
)

var curScope int = GlobalScope

func parse() {

	for {
		tokenType := parseStatement()
		if tokenType == TokenTypeEndOfStream {
			break
		}
	}

	printFuncTable()
	printSymbolList()
}

func parseStatement() TokenType {
	token := nextToken()

	switch token.Type {
	case TokenTypeEndOfStream:
		fmt.Println("parseStatement() TokenTypeEndOfStream")

	case TokenTypeVar:
		parseVar()

	case TokenTypeIdentifier:
		// fmt.Println(SymbolTypeVar, SymbolTypeVar, curScope)

		if getSymbolNode(token.Lexem, SymbolTypeVar, curScope) != nil {
			// 赋值语句
			parseAssignment()
		} else if getFuncByName(token.Lexem) != nil {
			// 函数调用
			addICodeNodeSourceLine(curScope, token.strLine)
			parseFuncCall(token.Lexem)
		} else {
			utils.Exit(token.LineNumber, "invalid identifier "+token.Lexem)
		}

	case TokenTypeReturn:
		parseReturn()

	case TokenTypeFunc:
		parseFunc()

	case TokenTypeOpenCurlyBrace:
		parseBlock()
	}

	return token.Type
}

func parseVar() {
	var token Token

	token = nextToken()
	checkToken(&token, TokenTypeIdentifier)

	// 将变量添加到SymbolTable中
	addSymbol(token.Lexem, SymbolTypeVar, curScope)

	token = nextToken()
	checkToken(&token, TokenTypeSemicolon)
}

func parseReturn() {
	var token Token
	token = nextToken()

	if token.Lexem != ";" {
		backToken()

		parseExpr()

		instrIndex := addICodeNodeInstruction(curScope, InstrTypePop)
		addOperandReg(curScope, instrIndex, OperandTypeReg)

		token = nextToken()
		checkToken(&token, TokenTypeSemicolon)
	}
}

// 表达式
func parseExpr() {
	parseSubExpr()
	for {

		token := nextToken()
		if token.Type != TokenTypeOperator {
			backToken()
			break
		}
		parseSubExpr()
	}
}

// 子表达式(A+B-C)
func parseSubExpr() {
	var instrIndex int

	parseTerm()

	for {

		token := nextToken()
		operatorType := getOperatorType(token)
		if token.Type != TokenTypeOperator || (operatorType != operatorTypeAdd && operatorType != operatorTypeSub) {
			backToken()
			break
		}

		parseTerm()

		// 取出第一个操作数存入_T1
		instrIndex = addICodeNodeInstruction(curScope, InstrTypePop)
		addOperandVar(curScope, instrIndex, TempVar1SymbolIndex)

		// 取出第一个操作数存入_T0
		instrIndex = addICodeNodeInstruction(curScope, InstrTypePop)
		addOperandVar(curScope, instrIndex, TempVar0SymbolIndex)

		var instrType InstrType = InstrTypeInvalid
		switch operatorType {
		case operatorTypeAdd:
			instrType = InstrTypeAdd
		case operatorTypeSub:
			instrType = InstrTypeSub
		default:
			instrType = InstrTypeInvalid
		}

		instrIndex = addICodeNodeInstruction(curScope, instrType)
		addOperandVar(curScope, instrIndex, TempVar0SymbolIndex)
		addOperandVar(curScope, instrIndex, TempVar1SymbolIndex)

		// 把结果存入_T0
		instrIndex = addICodeNodeInstruction(curScope, InstrTypePush)
		addOperandVar(curScope, instrIndex, TempVar0SymbolIndex)
	}
}

// 子表达式中的项(a*b/c)
func parseTerm() {
	var instrIndex int

	parseFactor()

	for {

		token := nextToken()
		operatorType := getOperatorType(token)
		if token.Type != TokenTypeOperator || (operatorType != operatorTypeMul && operatorType != operatorTypeDiv) {
			backToken()
			break
		}

		parseFactor()

		// 取出第一个操作数存入_T1
		instrIndex = addICodeNodeInstruction(curScope, InstrTypePop)
		addOperandVar(curScope, instrIndex, TempVar1SymbolIndex)

		// 取出第一个操作数存入_T0
		instrIndex = addICodeNodeInstruction(curScope, InstrTypePop)
		addOperandVar(curScope, instrIndex, TempVar0SymbolIndex)

		var instrType InstrType = InstrTypeInvalid
		switch operatorType {
		case operatorTypeAdd:
			instrType = InstrTypeMul
		case operatorTypeSub:
			instrType = InstrTypeDiv
		default:
			instrType = InstrTypeInvalid
		}

		instrIndex = addICodeNodeInstruction(curScope, instrType)
		addOperandVar(curScope, instrIndex, TempVar0SymbolIndex)
		addOperandVar(curScope, instrIndex, TempVar1SymbolIndex)

		// 把结果存入_T0
		instrIndex = addICodeNodeInstruction(curScope, InstrTypePush)
		addOperandVar(curScope, instrIndex, TempVar0SymbolIndex)
	}
}

// 项中的因子
func parseFactor() {
	var instrIndex int
	var token Token
	token = nextToken()

	switch token.Type {
	case TokenTypeInt:
		instrIndex = addICodeNodeInstruction(curScope, InstrTypePush)
		i, _ := strconv.Atoi(token.Lexem)
		addOperandInt(curScope, instrIndex, i)
	case TokenTypeFloat:
		instrIndex = addICodeNodeInstruction(curScope, InstrTypePush)
		f, _ := strconv.ParseFloat(token.Lexem, 32)
		addOperandFloat(curScope, instrIndex, float32(f))
	case TokenTypeString:
		instrIndex = addICodeNodeInstruction(curScope, InstrTypePush)
		strIndex := addstr(token.Lexem)
		addOperandStr(curScope, instrIndex, strIndex)
	case TokenTypeIdentifier:

		symbolNode := getSymbolNode(token.Lexem, SymbolTypeVar, curScope)
		if symbolNode != nil {
			// 变量
			instrIndex = addICodeNodeInstruction(curScope, InstrTypePush)
			addOperandVar(curScope, instrIndex, symbolNode.index)

		} else {
			// 函数
			funcNode := getFuncByName(token.Lexem)
			if funcNode != nil {
				parseFuncCall(token.Lexem)
				instrIndex = addICodeNodeInstruction(curScope, InstrTypePush)
				addOperandReg(curScope, instrIndex, OperandTypeReg)
			}
		}

	case TokenTypeOpenParen:
		parseExpr()
		token = nextToken()
		checkToken(&token, TokenTypeCloseParen)

	default:
		utils.Exit(token.LineNumber, "invalid input: "+token.Lexem)
	}
}

func parseFunc() {

	var token Token
	token = nextToken()
	checkToken(&token, TokenTypeIdentifier)
	funcName := token.Lexem
	addFunc(funcName)

	token = nextToken()
	checkToken(&token, TokenTypeOpenParen)

	// 函数参数列表
	var paramList []string
	token = nextToken()
	for token.Lexem != ")" {
		paramList = append(paramList, token.Lexem)
		token = nextToken()
		if token.Lexem == ")" {
			break
		}
		checkToken(&token, TokenTypeComma)
		token = nextToken()
	}

	// fmt.Println("paramList:", paramList)

	funcNode := getFuncByName(funcName)
	if funcNode != nil {
		funcNode.ParamCount = len(paramList)
		curScope = funcNode.FuncIndex
		// fmt.Println("curScope:", funcNode.FuncIndex)
	}

	// 将函数参数添加到SymbolTable中
	for i := 0; i < len(paramList); i++ {
		addSymbol(paramList[i], SymbolTypeParam, curScope)
	}

	token = nextToken()
	checkToken(&token, TokenTypeOpenCurlyBrace)

	parseBlock()

	curScope = GlobalScope
}

func parseBlock() {
	for {
		token := nextToken()
		if token.Lexem == "}" {
			return
		}

		if token.Type == TokenTypeEndOfStream {
			backToken()
			checkToken(curToken(), TokenTypeCloseCurlyBrace)
		}

		backToken()
		parseStatement()
	}
}

// 函数调用
func parseFuncCall(funcName string) {
	var token Token

	funcNode := getFuncByName(funcName)
	if funcNode == nil {
		return
	}

	token = nextToken()
	checkToken(&token, TokenTypeOpenParen)

	var paramCount int = 0

	// 函数参数
	for {
		token = nextToken()
		if token.Lexem == ")" {
			break
		}

		backToken()

		parseExpr()
		paramCount++

		token = nextToken()
		if token.Lexem == ")" {
			break
		}
		checkToken(&token, TokenTypeComma)
	}

	if paramCount != funcNode.ParamCount {
		var err string = "too few parameters!"
		if paramCount > funcNode.ParamCount {
			err = "too mutch parameters!"
		}
		utils.Exit(token.LineNumber, err)
	}

	instrIndex := addICodeNodeInstruction(curScope, InstrTypeCall)
	addOperandFuncIndex(curScope, instrIndex, funcNode.FuncIndex)
}

// 赋值语句
func parseAssignment() {

	var token Token

	token = nextToken()
	checkToken(&token, TokenTypeOperator)

	if token.Lexem != "=" {
		utils.Exit(token.LineNumber, "invalid operator "+token.Lexem)
	}

	parseExpr()

	token = nextToken()
	checkToken(&token, TokenTypeSemicolon)
}

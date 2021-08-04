package zcomplier

import (
	"../../utils"
	"fmt"
)

// 有限状态机中的状态
const (
	LexerStateInvalid = iota
	LexerStateStart
	LexerStateInt
	LexerStateFloat
	LexerStateString
	LexerStateIdentifier
	LexerStateOperator
	LexerStateDelimiter
)

const (
	TokenTypeInvalid = iota

	TokenTypeEndOfStream

	// 变量类型
	TokenTypeInt
	TokenTypeFloat
	TokenTypeString

	// 标识符
	TokenTypeIdentifier

	//操作符: +, -, *, /, =
	TokenTypeOperator

	// 关键字
	TokenTypeVar    // var
	TokenTypeFunc   // func
	TokenTypeReturn // return

	// 分隔符
	TokenTypeComma           // ,
	TokenTypeSemicolon       // ;
	TokenTypeOpenParen       // (
	TokenTypeCloseParen      // )
	TokenTypeOpenCurlyBrace  // {
	TokenTypeCloseCurlyBrace // }
)

// 运算符类型
const (
	operatorTypeInvalid = iota
	operatorTypeAdd     // +
	operatorTypeSub     // -
	operatorTypeMul     // *
	operatorTypeDiv     // /
	operatorTypeAssign  // =
)

type TokenType int
type OperatorType int

type Token struct {
	Type       TokenType
	Lexem      string
	LineNumber int
	strLine    string
}

var tokenIndex int = 0
var tokens []Token

var lexemTokenTypeMap = map[string]TokenType{
	"var":    TokenTypeVar,
	"func":   TokenTypeFunc,
	"return": TokenTypeReturn,
	",":      TokenTypeComma,
	";":      TokenTypeSemicolon,
	"(":      TokenTypeOpenParen,
	")":      TokenTypeCloseParen,
	"{":      TokenTypeOpenCurlyBrace,
	"}":      TokenTypeCloseCurlyBrace,
}

var tokenTypeLexemTMap = map[TokenType]string{
	TokenTypeVar:             "var",
	TokenTypeFunc:            "func",
	TokenTypeReturn:          "return",
	TokenTypeComma:           ",",
	TokenTypeSemicolon:       ";",
	TokenTypeOpenParen:       "(",
	TokenTypeCloseParen:      ")",
	TokenTypeOpenCurlyBrace:  "{",
	TokenTypeCloseCurlyBrace: "}",
}

var operatorTypeMap = map[string]OperatorType{
	"+": operatorTypeAdd,
	"-": operatorTypeSub,
	"*": operatorTypeMul,
	"/": operatorTypeDiv,
	"=": operatorTypeAssign,
}

func lexicalAnalyze(lines []string) {

	var curLexem string = ""
	var isAddCurChar bool = false
	var isLexemDone bool = false

	curState := LexerStateStart

	for lineIndex := 0; lineIndex < len(lines); lineIndex++ {
		line := lines[lineIndex]
		for charIndex := 0; charIndex < len(line); charIndex++ {
			r := rune(line[charIndex])

			isAddCurChar = true
			isLexemDone = false

			switch curState {
			case LexerStateStart:
				if utils.IsWhitespace(r) {
					isAddCurChar = false
				} else if utils.IsNumeric(r) {
					curState = LexerStateInt
				} else if utils.IsDelimiter(r) {
					curState = LexerStateDelimiter
				} else if utils.IsAlpha(r) {
					curState = LexerStateIdentifier
				} else if utils.IsOperator(r) {
					curState = LexerStateOperator
				} else if r == '"' {
					curState = LexerStateString
					isAddCurChar = false
				}
			case LexerStateInt:
				if utils.IsNumeric(r) {
					curState = LexerStateInt
				} else if r == '.' {
					curState = LexerStateFloat
				} else {
					isAddCurChar = false
					isLexemDone = true
					charIndex--
				}
			case LexerStateFloat:
				if utils.IsNumeric(r) {
					curState = LexerStateFloat
				} else {
					isAddCurChar = false
					isLexemDone = true
					charIndex--
				}
			case LexerStateString:
				if r == '"' {
					isAddCurChar = false
					isLexemDone = true
				}
			case LexerStateIdentifier:
				if utils.IsAlphanumeric(r) {
					curState = LexerStateIdentifier
				} else {
					isAddCurChar = false
					isLexemDone = true
					charIndex--
				}
			case LexerStateOperator:
				isAddCurChar = false
				isLexemDone = true
				charIndex--
			case LexerStateDelimiter:
				isAddCurChar = false
				isLexemDone = true
				charIndex--

			default:
				fmt.Println("invalid state!")
			}

			if isAddCurChar {
				curLexem += string(r)
			}

			if isLexemDone {

				var token Token
				token.LineNumber = lineIndex + 1
				token.strLine = line
				token.Lexem = curLexem
				updateTokenType(&token, curState)

				tokens = append(tokens, token)

				curState = LexerStateStart
				curLexem = ""
				isAddCurChar = true
			}
		}
	}

	if curLexem != "" {
		var token Token
		token.LineNumber = len(lines)
		token.strLine = lines[len(lines)-1]
		token.Lexem = curLexem
		updateTokenType(&token, curState)
		tokens = append(tokens, token)
	}

	// fmt.Printf("tokens: %+v\n", tokens)
}

func updateTokenType(token *Token, lexerState int) {
	switch lexerState {
	case LexerStateInt:
		token.Type = TokenTypeInt
	case LexerStateFloat:
		token.Type = TokenTypeFloat
	case LexerStateString:
		token.Type = TokenTypeString
	case LexerStateIdentifier:
		t, ok := lexemTokenTypeMap[token.Lexem]
		if ok {
			token.Type = t
		} else {
			token.Type = TokenTypeIdentifier
		}
	case LexerStateOperator:
		token.Type = TokenTypeOperator
	case LexerStateDelimiter:
		t, ok := lexemTokenTypeMap[token.Lexem]
		if ok {
			token.Type = t
		} else {
			token.Type = TokenTypeInvalid
		}
	default:
		token.Type = TokenTypeInvalid
	}
}

func nextToken() Token {
	if tokenIndex < len(tokens) {
		token := tokens[tokenIndex]
		tokenIndex++
		return token
	}
	var token Token
	token.Type = TokenTypeEndOfStream
	return token
}

func curToken() *Token {
	if tokenIndex < len(tokens) {
		return &tokens[tokenIndex]
	}
	var token Token
	token.Type = TokenTypeEndOfStream
	return &token
}

func backToken() Token {
	if tokenIndex == 0 {
		return tokens[tokenIndex]
	}
	tokenIndex--
	return tokens[tokenIndex]
}

func checkToken(token *Token, tokenType TokenType) {
	var errmsg string
	var ok bool
	if token.Type != tokenType {
		switch tokenType {
		case TokenTypeInt:
			errmsg = "int expected"
		case TokenTypeFloat:
			errmsg = "float expected"
		case TokenTypeString:
			errmsg = "string expected"
		case TokenTypeIdentifier:
			errmsg = "identifier expected"
		case TokenTypeOperator:
			errmsg = "operator expected"
		default:
			errmsg, ok = tokenTypeLexemTMap[tokenType]
			if ok {
				errmsg = "\"" + errmsg + "\" expected"
			} else {
				errmsg = "unknow error"
			}
		}
		utils.Exit(token.LineNumber, errmsg)
	}
}

func getOperatorType(token Token) OperatorType {
	operatorType, ok := operatorTypeMap[token.Lexem]
	if ok {
		return operatorType
	}
	return operatorTypeInvalid
}

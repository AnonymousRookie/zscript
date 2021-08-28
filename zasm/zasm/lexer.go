package zasm

import (
	"fmt"
	"strings"

	"zscript/utils"
)

// 有限状态机中的状态
const (
	LexerStateInvalid = iota
	LexerStateStart
	LexerStateInt
	LexerStateFloat
	LexerStateString
	LexerStateIdentifier
	LexerStateInstr
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

	//指令类型
	TokenTypeInstr

	// 关键字
	TokenTypeVar   // Var
	TokenTypeFunc  // Func
	TokenTypeParam // Param

	// 主应用程序API
	TokenTypeHostApiPrint // print

	// 分隔符
	TokenTypeComma           // ,
	TokenTypeOpenCurlyBrace  // {
	TokenTypeCloseCurlyBrace // }
)

var instrMap = map[string]InstrType{
	"Mov":         InstrTypeMov,
	"Add":         InstrTypeAdd,
	"Sub":         InstrTypeSub,
	"Mul":         InstrTypeMul,
	"Div":         InstrTypeDiv,
	"Jmp":         InstrTypeJmp,
	"Push":        InstrTypePush,
	"Pop":         InstrTypePop,
	"Call":        InstrTypeCall,
	"CallHostApi": InstrTypeCallHostApi,
	"Ret":         InstrTypeRet,
}

var instrOpCountMap = map[InstrType]int{
	InstrTypeMov:         2,
	InstrTypeAdd:         2,
	InstrTypeSub:         2,
	InstrTypeMul:         2,
	InstrTypeDiv:         2,
	InstrTypeJmp:         1,
	InstrTypePush:        1,
	InstrTypePop:         1,
	InstrTypeCall:        1,
	InstrTypeCallHostApi: 1,
	InstrTypeRet:         0,
}

type TokenType int

type Token struct {
	Type       TokenType
	Lexem      string
	LineNumber int
	strLine    string
}

var lexemTokenTypeMap = map[string]TokenType{
	"Var":   TokenTypeVar,
	"Func":  TokenTypeFunc,
	"Param": TokenTypeParam,
	"print": TokenTypeHostApiPrint,
	",":     TokenTypeComma,
	"{":     TokenTypeOpenCurlyBrace,
	"}":     TokenTypeCloseCurlyBrace,
}

var tokenTypeLexemTMap = map[TokenType]string{
	TokenTypeVar:             "Var",
	TokenTypeFunc:            "Func",
	TokenTypeParam:           "Param",
	TokenTypeHostApiPrint:    "print",
	TokenTypeComma:           ",",
	TokenTypeOpenCurlyBrace:  "{",
	TokenTypeCloseCurlyBrace: "}",
}

type Lexer struct {
	tokenIndex int
	tokens     []Token
}

func NewLexer() *Lexer {
	return &Lexer{tokenIndex: 0}
}

func (lexer *Lexer) lexicalAnalyze(lines []string) {
	var curLexem string = ""
	var isAddCurChar bool = false
	var isLexemDone bool = false

	curState := LexerStateStart

	for lineIndex := 0; lineIndex < len(lines); lineIndex++ {
		line := lines[lineIndex]

		line = strings.TrimLeft(line, "\t")
		line = strings.TrimLeft(line, " ")
		// fmt.Printf("[%s]\n", line)

		// 去除注释语句
		if len(line) > 0 && line[0] == ';' {
			continue
		}

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
				lexer.updateTokenType(&token, curState)

				lexer.tokens = append(lexer.tokens, token)

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
		lexer.updateTokenType(&token, curState)
		lexer.tokens = append(lexer.tokens, token)
	}

	// fmt.Printf("tokens: %+v\n", tokens)
}

func (lexer *Lexer) updateTokenType(token *Token, lexerState int) {
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
			_, ok := instrMap[token.Lexem]
			if ok {
				token.Type = TokenTypeInstr
			} else {
				token.Type = TokenTypeIdentifier
			}
		}
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

func (lexer *Lexer) nextToken() Token {
	if lexer.tokenIndex < len(lexer.tokens) {
		token := lexer.tokens[lexer.tokenIndex]
		lexer.tokenIndex++
		return token
	}
	var token Token
	token.Type = TokenTypeEndOfStream
	return token
}

func (lexer *Lexer) backToken() Token {
	if lexer.tokenIndex == 0 {
		return lexer.tokens[lexer.tokenIndex]
	}
	lexer.tokenIndex--
	return lexer.tokens[lexer.tokenIndex]
}

func (lexer *Lexer) checkToken(token *Token, tokenType TokenType) {
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

func (lexer *Lexer) getInstrType(token *Token) (instrType InstrType, instrOpCount int32) {
	instrType = InstrTypeInvalid
	instrOpCount = 0
	tp, ok := instrMap[token.Lexem]
	if ok {
		instrType = tp
		cnt, ok := instrOpCountMap[tp]
		if ok {
			instrOpCount = int32(cnt)
		}
	}
	return instrType, instrOpCount
}

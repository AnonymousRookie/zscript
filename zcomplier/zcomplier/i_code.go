package zcomplier

import (
	// "fmt"
	"strings"
)

type InstrType int

// 指令类型
const (
	InstrTypeInvalid = iota
	InstrTypeMov
	InstrTypeAdd
	InstrTypeSub
	InstrTypeMul
	InstrTypeDiv

	InstrTypeJmp

	InstrTypePush
	InstrTypePop

	InstrTypeCall
	InstrTypeRet
	InstrTypeExit
)

// 操作数类型
const (
	OperandTypeInvalid   = iota
	OperandTypeInt       // 整数
	OperandTypeFloat     // 浮点数
	OperandTypeString    // 字符串
	OperandTypeVar       // 变量
	OperandTypeFuncIndex // FuncIndex
	OperandTypeReg       // 寄存器
)

const (
	RegisterTypeInvalid = iota
	RegisterTypeRetVal
)

// 操作数
type Operand struct {
	OperandType int         // 操作数类型
	Val         interface{} // 操作数
}

const (
	ICodeNodeTypeInvalid    = iota
	ICodeNodeTypeInstr      // 中间代码指令
	ICodeNodeTypeSourceLine // 标注指令的源代码行
)

type ICodeNode struct {
	ICodeNodeType int
	Val           interface{} // 指令 or 源代码行
}

type ICodeInstr struct {
	instr    InstrType // 指令类型
	operands []Operand // 操作数
}

type ICodeNodeList []ICodeNode

var FuncICodeNodes map[int]ICodeNodeList = make(map[int]ICodeNodeList)

var instrMap = map[InstrType]string{
	InstrTypeMov:  "Mov",
	InstrTypeAdd:  "Add",
	InstrTypeSub:  "Sub",
	InstrTypeMul:  "Mul",
	InstrTypeDiv:  "Div",
	InstrTypeJmp:  "Jmp",
	InstrTypePush: "Push",
	InstrTypePop:  "Pop",
	InstrTypeCall: "Call",
	InstrTypeRet:  "Ret",
}

func addICodeNode(funcIndex int, iCodeNode ICodeNode) int {
	iCodeNodeList, ok := FuncICodeNodes[funcIndex]
	iCodeNodeIndex := len(iCodeNodeList)
	if ok {
		iCodeNodeList = append(iCodeNodeList, iCodeNode)
		FuncICodeNodes[funcIndex] = iCodeNodeList
	} else {
		var iCodeNodeList ICodeNodeList
		iCodeNodeList = append(iCodeNodeList, iCodeNode)
		FuncICodeNodes[funcIndex] = iCodeNodeList
	}
	return iCodeNodeIndex
}

func getICodeNode(funcIndex int, iCodeNodeIndex int) *ICodeNode {
	iCodeNodeList, ok := FuncICodeNodes[funcIndex]
	if !ok {
		return nil
	}
	return &iCodeNodeList[iCodeNodeIndex]
}

func getIcodeNodeList(funcIndex int) *ICodeNodeList {
	iCodeNodeList, ok := FuncICodeNodes[funcIndex]
	if !ok {
		return nil
	}
	return &iCodeNodeList
}

// 添加中间代码指令
func addICodeNodeInstruction(funcIndex int, instr InstrType) int {
	iCodeInstr := ICodeInstr{instr, make([]Operand, 0)}
	iCodeNode := ICodeNode{ICodeNodeTypeInstr, iCodeInstr}
	return addICodeNode(funcIndex, iCodeNode)
}

// 添加标注指令的源代码行
func addICodeNodeSourceLine(funcIndex int, sourceLine string) {
	sourceLine = strings.TrimSpace(sourceLine)
	// fmt.Println("sourceLine:", sourceLine)
	iCodeNode := ICodeNode{ICodeNodeTypeSourceLine, sourceLine}
	addICodeNode(funcIndex, iCodeNode)
}

// 添加操作数
func addOperand(funcIndex int, instrIndex int, op Operand) {
	iCodeNode := getICodeNode(funcIndex, instrIndex)
	if iCodeNode == nil || iCodeNode.ICodeNodeType != ICodeNodeTypeInstr {
		return
	}
	iCodeInstr := iCodeNode.Val.(ICodeInstr)
	iCodeInstr.operands = append(iCodeInstr.operands, op)
	iCodeNode.Val = iCodeInstr
}

func addOperandInt(funcIndex int, instrIndex int, val int) {
	operand := Operand{OperandTypeInt, val}
	addOperand(funcIndex, instrIndex, operand)
}

func addOperandFloat(funcIndex int, instrIndex int, val float32) {
	operand := Operand{OperandTypeFloat, val}
	addOperand(funcIndex, instrIndex, operand)
}

func addOperandStr(funcIndex int, instrIndex int, strIndex int) {
	operand := Operand{OperandTypeString, strIndex}
	addOperand(funcIndex, instrIndex, operand)
}

func addOperandVar(funcIndex int, instrIndex int, symbolIndex int) {
	operand := Operand{OperandTypeVar, symbolIndex}
	addOperand(funcIndex, instrIndex, operand)
}

func addOperandFuncIndex(funcIndex int, instrIndex int, opfuncIndex int) {
	operand := Operand{OperandTypeFuncIndex, opfuncIndex}
	addOperand(funcIndex, instrIndex, operand)
}

func addOperandReg(funcIndex int, instrIndex int, regType int) {
	operand := Operand{OperandTypeReg, regType}
	addOperand(funcIndex, instrIndex, operand)
}

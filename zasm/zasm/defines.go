package zasm

const (
	srcFileSuffix = ".zasm"
	outFileSuffix = ".zse"

	registerT0        = "_T0"
	registerT1        = "_T1"
	registerReturnVal = "_RetVal"
)

const (
	registerTypeInvalid = iota
	registerTypeT0
	registerTypeT1
	registerTypeRetVal
)

var regTypeMap = map[string]int32{
	registerT0:        registerTypeT0,
	registerT1:        registerTypeT1,
	registerReturnVal: registerTypeRetVal,
}

// 指令类型
type InstrType int32

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
type OperandType int32

const (
	OperandTypeInvalid         = iota
	OperandTypeInt             // 整数
	OperandTypeFloat           // 浮点数
	OperandTypeStrIndex        // 字符串索引
	OperandTypeIdentifierIndex // 标识符索引
	OperandTypeFuncIndex       // 函数索引
	OperandTypeReg             // 寄存器
)

// .zse可执行文件格式
// Header
// Instruction Stream
// String Table
// Function Table

type Header struct {
	isExistMainFunc bool  // 是否存在main函数
	mainFuncIndex   int32 // main函数在函数表中的索引
}

type Operand struct {
	opType OperandType // 操作数类型
	opVal  interface{} // 操作数
}
type Instr struct {
	index     int32
	instrType InstrType // 指令类型
	opCount   int32     // 操作数个数
	ops       []Operand // 操作数
}
type InstrStream struct {
	count  int32   // 指令的个数
	instrs []Instr // 指令
}

type StrInfo struct {
	index int32
	len   int32
	str   string
}
type StringTable struct {
	count    int32 // 字符串的个数
	strinfos []StrInfo
}

type FuncNode struct {
	index       int32  // 函数索引
	len         int32  // 函数名长度
	funcName    string // 函数名称
	entryPoint  int32  // 函数入口点, 函数第一个指令的索引
	paramcount  int32  // 函数参数个数
	symbolCount int32
	symbolNodes []SymbolNode
}
type FuncTable struct {
	count     int32 // 函数个数
	funcNodes []FuncNode
}

type SymbolType int32

const (
	SymbolTypeInvalid = iota
	SymbolTypeVar     // 变量
	SymbolTypeParam   // 函数参数
)

type SymbolNode struct {
	index      int32
	len        int32
	identifier string
	symbolType SymbolType
	funcIndex  int32
}

var globalSymbolNodes []SymbolNode

// Header
var header Header

// Instruction Stream
var instrStream InstrStream

// String Table
var strTable StringTable

// Function Table
var funcTable FuncTable

func (funcNode *FuncNode) addFuncSymbol(identifier string, symbolType SymbolType) int32 {
	index := len(funcNode.symbolNodes)
	node := SymbolNode{int32(index), int32(len(identifier)), identifier, symbolType, funcNode.index}
	funcNode.symbolNodes = append(funcNode.symbolNodes, node)
	funcNode.symbolCount = int32(len(funcNode.symbolNodes))
	return int32(index)
}

func addGlobalSymbol(identifier string, symbolType SymbolType, funcIndex int32) int32 {
	index := len(globalSymbolNodes)
	node := SymbolNode{int32(index), int32(len(identifier)), identifier, symbolType, funcIndex}
	globalSymbolNodes = append(globalSymbolNodes, node)
	return int32(index)
}

func addSymbol(identifier string, symbolType SymbolType, funcIndex int32) int32 {
	if funcIndex == 0 {
		return addGlobalSymbol(identifier, symbolType, funcIndex)
	} else {
		funcNode := getFuncNodeByIndex(funcIndex)
		return funcNode.addFuncSymbol(identifier, symbolType)
	}
}

func getSymbol(identifier string, funcIndex int32) *SymbolNode {

	funcNode := getFuncNodeByIndex(funcIndex)
	if funcNode != nil {
		for i := 0; i < len(funcNode.symbolNodes); i++ {
			node := funcNode.symbolNodes[i]
			if node.identifier == identifier && node.funcIndex == funcIndex {
				return &node
			}
		}
	}

	for i := 0; i < len(globalSymbolNodes); i++ {
		node := globalSymbolNodes[i]
		if node.identifier == identifier && node.funcIndex == 0 {
			return &node
		}
	}

	return nil
}

// 返回新增字符串的索引
func addStr(str string) int32 {
	var strInfo StrInfo
	strInfo.index = int32(len(strTable.strinfos))
	strInfo.len = int32(len(str))
	strInfo.str = str

	strTable.strinfos = append(strTable.strinfos, strInfo)
	strTable.count = int32(len(strTable.strinfos))

	return strInfo.index
}

// 返回新增函数的索引
func addFuncNode(funcName string) int32 {
	var funcNode FuncNode
	funcNode.index = int32(len(funcTable.funcNodes) + 1)
	funcNode.funcName = funcName
	funcNode.len = int32(len(funcName))
	funcNode.entryPoint = -1

	funcTable.funcNodes = append(funcTable.funcNodes, funcNode)
	funcTable.count = int32(len(funcTable.funcNodes))

	return funcNode.index
}

func getFuncNodeByIndex(index int32) *FuncNode {
	for i := 0; i < int(funcTable.count); i++ {
		if funcTable.funcNodes[i].index == index {
			return &funcTable.funcNodes[i]
		}
	}
	return nil
}

func getFuncNodeByName(name string) *FuncNode {
	for i := 0; i < int(funcTable.count); i++ {
		if funcTable.funcNodes[i].funcName == name {
			return &funcTable.funcNodes[i]
		}
	}
	return nil
}

// 返回新增指令的索引
func addInstr(instr Instr) int32 {
	instr.index = int32(len(instrStream.instrs))
	instrStream.instrs = append(instrStream.instrs, instr)
	instrStream.count = int32(len(instrStream.instrs))
	return instr.index
}

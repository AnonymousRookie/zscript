package zvm

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
	InstrTypeCallHostApi

	InstrTypeRet
	InstrTypeExit
)

// 操作数类型
type OperandType int32

const (
	OperandTypeInvalid      = iota
	OperandTypeInt          // 整数
	OperandTypeFloat        // 浮点数
	OperandTypeString       // 字符串
	OperandTypeVar          // 变量
	OperandTypeFuncIndex    // 函数索引
	OperandTypeReg          // 寄存器
	OperandTypeHostApiIndex // 主应用程序API索引
)

// .zse可执行文件格式
// Header
// Instruction Stream
// String Table
// Function Table
// HostApi Table

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

type HostApiNode struct {
	index int32  // 主应用程序API索引
	len   int32  // 主应用程序API名称长度
	name  string // 主应用程序API名称
}

type HostApiTable struct {
	count        int32 // 主应用程序API个数
	hostApiNodes []HostApiNode
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

const (
	registerTypeInvalid = iota
	registerTypeT0
	registerTypeT1
	registerTypeRetVal
)

const (
	registerT0        = "_T0"
	registerT1        = "_T1"
	registerReturnVal = "_RetVal"
)

const (
	ZValTypeInvalid = iota
	ZValTypeInt
	ZValTypeFloat
	ZValTypeStr
)

type ZVal struct {
	valType    int32
	identifier string
	val        interface{}
}

type FuncZVal map[int32]ZVal

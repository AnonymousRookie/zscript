package zvm

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"zscript/utils"
)

type Zvm struct {
	curFuncIndex        int32
	curInstrIndex       int32
	funcIndexBeforeJmp  int32 // 跳转前所在函数索引
	instrIndexBeforeJmp int32 // 跳转前的指令位置

	header       Header
	instrStream  InstrStream
	strTable     StringTable
	funcTable    FuncTable
	hostApiTable HostApiTable

	runtimeStack *Stack
	regRetVal    ZVal // register _RetVal
	regT0        ZVal // register _T0
	regT1        ZVal // register _T1

	allFuncZVal map[int32]FuncZVal
}

func NewZvm() *Zvm {
	var zvm Zvm
	zvm.curFuncIndex = 0
	zvm.curInstrIndex = 0
	zvm.funcIndexBeforeJmp = 0
	zvm.instrIndexBeforeJmp = 0
	zvm.regT0.identifier = registerT0
	zvm.regT1.identifier = registerT1
	zvm.regRetVal.identifier = registerReturnVal
	return &zvm
}

func (zvm *Zvm) Load(zseFilename string) {

	suffix := zseFilename[len(zseFilename)-4:]
	if suffix != ".zse" {
		utils.ExitWithErrMsg("file suffix should be: .zse")
	}

	fp, err := os.Open(zseFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fp.Close()

	// Header
	err = binary.Read(fp, binary.LittleEndian, &zvm.header.isExistMainFunc)
	utils.Check(err == nil, "read header.isExistMainFunc failed!")
	err = binary.Read(fp, binary.LittleEndian, &zvm.header.mainFuncIndex)
	utils.Check(err == nil, "read header.mainFuncIndex failed!")
	// fmt.Printf("zvm.header:%+v\n", zvm.header)

	// Instruction Stream
	err = binary.Read(fp, binary.LittleEndian, &zvm.instrStream.count)
	utils.Check(err == nil, "read instrStream.count failed!")
	for i := 0; i < int(zvm.instrStream.count); i++ {
		var instr Instr
		err = binary.Read(fp, binary.LittleEndian, &instr.index)
		utils.Check(err == nil, "read instr.index failed!")
		err = binary.Read(fp, binary.LittleEndian, &instr.instrType)
		utils.Check(err == nil, "read instr.instrType failed!")
		err = binary.Read(fp, binary.LittleEndian, &instr.opCount)
		utils.Check(err == nil, "read instr.opCount failed!")

		for j := 0; j < int(instr.opCount); j++ {
			var op Operand
			err = binary.Read(fp, binary.LittleEndian, &op.opType)
			utils.Check(err == nil, "read op.opType failed!")

			if op.opType == OperandTypeFloat {
				var v float32
				err = binary.Read(fp, binary.LittleEndian, &v)
				utils.Check(err == nil, "read op.opVal.(float32) failed!")
				op.opVal = v
			} else {
				var v int32
				err = binary.Read(fp, binary.LittleEndian, &v)
				utils.Check(err == nil, "read op.opVal.(int32) failed!")
				op.opVal = v
			}
			instr.ops = append(instr.ops, op)
		}
		zvm.instrStream.instrs = append(zvm.instrStream.instrs, instr)
	}
	// fmt.Printf("zvm.instrStream:%+v\n", zvm.instrStream)

	// String Table
	err = binary.Read(fp, binary.LittleEndian, &zvm.strTable.count)
	utils.Check(err == nil, "read strTable.count failed!")
	for i := 0; i < int(zvm.strTable.count); i++ {
		var strInfo StrInfo
		err = binary.Read(fp, binary.LittleEndian, &strInfo.index)
		utils.Check(err == nil, "read strInfo.index failed!")

		err = binary.Read(fp, binary.LittleEndian, &strInfo.len)
		utils.Check(err == nil, "read strInfo.len failed!")

		buf := make([]byte, strInfo.len)
		err = binary.Read(fp, binary.LittleEndian, &buf)
		utils.Check(err == nil, "read strInfo.str failed!")
		strInfo.str = string(buf)

		zvm.strTable.strinfos = append(zvm.strTable.strinfos, strInfo)
	}
	// fmt.Printf("zvm.strTable:%+v\n", zvm.strTable)

	// Function Table
	err = binary.Read(fp, binary.LittleEndian, &zvm.funcTable.count)
	utils.Check(err == nil, "read strTable.count failed!")
	for i := 0; i < int(zvm.funcTable.count); i++ {
		var funcNode FuncNode
		err = binary.Read(fp, binary.LittleEndian, &funcNode.index)
		utils.Check(err == nil, "read funcNode.index failed!")

		err = binary.Read(fp, binary.LittleEndian, &funcNode.len)
		utils.Check(err == nil, "read funcNode.len failed!")

		buf := make([]byte, funcNode.len)
		err = binary.Read(fp, binary.LittleEndian, &buf)
		utils.Check(err == nil, "read funcNode.funcName failed!")
		funcNode.funcName = string(buf)

		err = binary.Read(fp, binary.LittleEndian, &funcNode.entryPoint)
		utils.Check(err == nil, "read funcNode.entryPoint failed!")

		err = binary.Read(fp, binary.LittleEndian, &funcNode.paramcount)
		utils.Check(err == nil, "read funcNode.paramcount failed!")

		err = binary.Read(fp, binary.LittleEndian, &funcNode.symbolCount)
		utils.Check(err == nil, "read funcNode.symbolCount failed!")
		for j := 0; j < int(funcNode.symbolCount); j++ {
			var symbolNode SymbolNode
			err = binary.Read(fp, binary.LittleEndian, &symbolNode.index)
			utils.Check(err == nil, "read symbolNode.index failed!")

			err = binary.Read(fp, binary.LittleEndian, &symbolNode.len)
			utils.Check(err == nil, "read symbolNode.len failed!")

			buf := make([]byte, symbolNode.len)
			err = binary.Read(fp, binary.LittleEndian, &buf)
			utils.Check(err == nil, "read symbolNode.identifier failed!")
			symbolNode.identifier = string(buf)

			err = binary.Read(fp, binary.LittleEndian, &symbolNode.symbolType)
			utils.Check(err == nil, "read ymbolNode.symbolType failed!")

			err = binary.Read(fp, binary.LittleEndian, &symbolNode.funcIndex)
			utils.Check(err == nil, "read symbolNode.funcIndex failed!")

			funcNode.symbolNodes = append(funcNode.symbolNodes, symbolNode)
		}
		zvm.funcTable.funcNodes = append(zvm.funcTable.funcNodes, funcNode)
	}
	// fmt.Printf("zvm.funcTable:%+v\n", zvm.funcTable)

	// HostApi Table
	err = binary.Read(fp, binary.LittleEndian, &zvm.hostApiTable.count)
	utils.Check(err == nil, "read hostApiTable.count failed!")
	for i := 0; i < int(zvm.hostApiTable.count); i++ {
		var hostApiNode HostApiNode
		err = binary.Read(fp, binary.LittleEndian, &hostApiNode.index)
		utils.Check(err == nil, "read hostApiNode.index failed!")

		err = binary.Read(fp, binary.LittleEndian, &hostApiNode.len)
		utils.Check(err == nil, "read hostApiNode.len failed!")

		buf := make([]byte, hostApiNode.len)
		err = binary.Read(fp, binary.LittleEndian, &buf)
		utils.Check(err == nil, "read hostApiNode.name failed!")
		hostApiNode.name = string(buf)

		zvm.hostApiTable.hostApiNodes = append(zvm.hostApiTable.hostApiNodes, hostApiNode)
	}

	// fmt.Printf("zvm.hostApiTable:%+v\n", zvm.hostApiTable)
}

func (zvm *Zvm) Run() {
	if !zvm.header.isExistMainFunc {
		fmt.Println("main func not exist!")
		return
	}

	zvm.runtimeStack = NewStack()

	mainFuncNode := zvm.getFuncNodeByIndex(zvm.header.mainFuncIndex)

	zvm.curInstrIndex = mainFuncNode.entryPoint
	zvm.curFuncIndex = mainFuncNode.index

	for {
		instr := zvm.getInstrByIndex(zvm.curInstrIndex)
		if instr == nil {
			break
		}

		switch instr.instrType {
		case InstrTypeMov:
			zvm.processMov(instr)
			zvm.curInstrIndex++
		case InstrTypeAdd:
			zvm.processAdd(instr)
			zvm.curInstrIndex++
		case InstrTypeSub:
			zvm.processSub(instr)
			zvm.curInstrIndex++
		case InstrTypeMul:
			zvm.processMul(instr)
			zvm.curInstrIndex++
		case InstrTypeDiv:
			zvm.processDiv(instr)
			zvm.curInstrIndex++

		case InstrTypePush:
			zvm.processPush(instr)
			zvm.curInstrIndex++

		case InstrTypePop:
			zvm.processPop(instr)
			zvm.curInstrIndex++

		case InstrTypeCall:
			zvm.processCall(instr)

		case InstrTypeCallHostApi:
			zvm.processCallHostApi(instr)
			zvm.curInstrIndex++

		case InstrTypeRet:
			exit := zvm.processRet()
			if exit {
				return
			}

		default:
			fmt.Println("default")
		}
	}
}

func (zvm *Zvm) processPush(instr *Instr) {
	op := instr.getOperantByIndex(0)
	log.Printf("processPush %+v\n", op)

	var zval ZVal
	if op.opType == OperandTypeReg {
		zval = *zvm.getRegVar(op)
	} else if op.opType == OperandTypeInt {
		zval.valType = ZValTypeInt
		zval.val = op.opVal.(int32)
	} else if op.opType == OperandTypeFloat {
		zval.valType = ZValTypeFloat
		zval.val = op.opVal.(float32)
	} else if op.opType == OperandTypeString {
		strIndex := op.opVal.(int32)
		zval.valType = ZValTypeStr
		zval.val = zvm.getStrByIndex(strIndex)
	} else {
		varIndex := zvm.getZVal(zvm.curFuncIndex, op)
		zval = zvm.allFuncZVal[zvm.curFuncIndex][varIndex]
	}
	zvm.runtimeStack.Push(zval)
}

func (zvm *Zvm) processPop(instr *Instr) {
	op := instr.getOperantByIndex(0)
	// fmt.Printf("%+v\n", op)

	utils.Check(OperandTypeReg == op.opType, "operand type should be OperandTypeReg!")
	regType := op.opVal.(int32)
	// fmt.Printf("regType:%+v\n", regType)
	utils.Check(registerTypeT0 == regType || registerTypeT1 == regType || registerTypeRetVal == regType, "unexpect register type!")

	var zval *ZVal = nil
	if registerTypeT0 == regType {
		zval = &zvm.regT0
	} else if registerTypeT1 == regType {
		zval = &zvm.regT1
	} else if registerTypeRetVal == regType {
		zval = &zvm.regRetVal
	}

	pv := zvm.runtimeStack.Pop()
	utils.Check(pv != nil, "runtimeStack is empty, can not pop!")
	log.Printf("pv: %+v\n", pv)

	*zval = pv.(ZVal)

	if registerTypeT0 == regType {
		log.Printf("processPop regT0:%+v\n", zvm.regT0)
	} else if registerTypeT1 == regType {
		log.Printf("processPop regT1:%+v\n", zvm.regT1)
	} else if registerTypeRetVal == regType {
		log.Printf("processPop retVal:%+v\n", zvm.regRetVal)
	}
}

func (zvm *Zvm) processMov(instr *Instr) {
	dstOp := instr.getOperantByIndex(0)
	srcOp := instr.getOperantByIndex(1)

	var srcV ZVal
	if srcOp.opType == OperandTypeReg {
		srcV = *zvm.getRegVar(srcOp)
	} else {
		varIndex := zvm.getZVal(zvm.curFuncIndex, srcOp)
		srcV = zvm.allFuncZVal[zvm.curFuncIndex][varIndex]
	}

	if dstOp.opType == OperandTypeReg {
		dstV := zvm.getRegVar(dstOp)
		dstV.valType = srcV.valType
		dstV.val = srcV.val
	} else {
		varIndex := zvm.getZVal(zvm.curFuncIndex, dstOp)
		zval := zvm.allFuncZVal[zvm.curFuncIndex][varIndex]
		zval.valType = srcV.valType
		zval.val = srcV.val
		zvm.allFuncZVal[zvm.curFuncIndex][varIndex] = zval
	}

	log.Printf("processMov regT0:%+v\n", zvm.regT0)
	log.Printf("processMov regT1:%+v\n", zvm.regT1)
	log.Printf("processMov allFuncZVal:%+v\n", zvm.allFuncZVal)
}

func (zvm *Zvm) processCall(instr *Instr) {
	log.Printf("processCall instr:%+v\n", instr)

	op := instr.getOperantByIndex(0)
	utils.Check(op.opType == OperandTypeFuncIndex, "expected OperandTypeFuncIndex!")

	zvm.funcIndexBeforeJmp = zvm.curFuncIndex
	zvm.curFuncIndex = op.opVal.(int32)
	funcNode := zvm.getFuncNodeByIndex(zvm.curFuncIndex)

	zvm.instrIndexBeforeJmp = zvm.curInstrIndex
	zvm.curInstrIndex = funcNode.entryPoint

	log.Printf("processCall funcNode:%+v\n", funcNode)

	// 函数参数
	for i := 0; i < int(funcNode.symbolCount); i++ {
		symbolNode := funcNode.symbolNodes[i]
		if symbolNode.symbolType == SymbolTypeParam {
			// 传参
			pv := zvm.runtimeStack.Pop()
			log.Printf("processCall pv:%+v\n", pv)

			var op Operand
			op.opType = OperandTypeVar
			op.opVal = symbolNode.index
			varIndex := zvm.getZVal(zvm.curFuncIndex, &op)
			zvm.allFuncZVal[zvm.curFuncIndex][varIndex] = pv.(ZVal)
		}
	}

	log.Printf("processCall allFuncZVal:%+v\n", zvm.allFuncZVal)
}

func (zvm *Zvm) processCallHostApi(instr *Instr) {
	log.Printf("processCallHostApi instr:%+v\n", instr)

	op := instr.getOperantByIndex(0)
	utils.Check(op.opType == OperandTypeHostApiIndex, "expected OperandTypeHostApiIndex!")

	hostApiIndex := op.opVal.(int32)
	hostApiNode := zvm.getHostApiNode(hostApiIndex)

	if hostApiNode.name == "print" {
		pv := zvm.runtimeStack.Pop()
		log.Printf("processCallHostApi pv:%+v\n", pv)
		fmt.Printf("%v = %v\n", pv.(ZVal).identifier, pv.(ZVal).val)
	}
}

func (zvm *Zvm) processAdd(instr *Instr) {
	zvm.calc(instr, "+")

	log.Printf("processAdd regT0:%+v\n", zvm.regT0)
	log.Printf("processAdd regT1:%+v\n", zvm.regT1)
	log.Printf("processAdd allFuncZVal:%+v\n", zvm.allFuncZVal)
}

func (zvm *Zvm) processMul(instr *Instr) {
	zvm.calc(instr, "*")

	log.Printf("processMul regT0:%+v\n", zvm.regT0)
	log.Printf("processMul regT1:%+v\n", zvm.regT1)
	log.Printf("processMul allFuncZVal:%+v\n", zvm.allFuncZVal)
}

func (zvm *Zvm) processDiv(instr *Instr) {
	zvm.calc(instr, "/")

	log.Printf("processDiv regT0:%+v\n", zvm.regT0)
	log.Printf("processDiv regT1:%+v\n", zvm.regT1)
	log.Printf("processDiv allFuncZVal:%+v\n", zvm.allFuncZVal)
}

func (zvm *Zvm) processSub(instr *Instr) {
	zvm.calc(instr, "-")

	log.Printf("processSub regT0:%+v\n", zvm.regT0)
	log.Printf("processSub regT1:%+v\n", zvm.regT1)
	log.Printf("processSub allFuncZVal:%+v\n", zvm.allFuncZVal)
}

func (zvm *Zvm) calc(instr *Instr, operator string) {
	dstOp := instr.getOperantByIndex(0)
	srcOp := instr.getOperantByIndex(1)

	var srcV ZVal
	if srcOp.opType == OperandTypeReg {
		srcV = *zvm.getRegVar(srcOp)
	} else {
		varIndex := zvm.getZVal(zvm.curFuncIndex, srcOp)
		srcV = zvm.allFuncZVal[zvm.curFuncIndex][varIndex]
	}

	var zval ZVal
	if dstOp.opType == OperandTypeReg {
		zval = *zvm.getRegVar(dstOp)
	} else {
		varIndex := zvm.getZVal(zvm.curFuncIndex, dstOp)
		zval = zvm.allFuncZVal[zvm.curFuncIndex][varIndex]
	}

	if ZValTypeInt == zval.valType {
		if ZValTypeFloat == srcV.valType {
			zval.valType = ZValTypeFloat
			if operator == "+" {
				zval.val = float32(zval.val.(int32)) + srcV.val.(float32)
			} else if operator == "-" {
				zval.val = float32(zval.val.(int32)) - srcV.val.(float32)
			} else if operator == "*" {
				zval.val = float32(zval.val.(int32)) * srcV.val.(float32)
			} else if operator == "/" {
				zval.val = float32(zval.val.(int32)) / srcV.val.(float32)
			}
		} else if ZValTypeInt == srcV.valType {
			zval.valType = ZValTypeInt
			if operator == "+" {
				zval.val = zval.val.(int32) + srcV.val.(int32)
			} else if operator == "-" {
				zval.val = zval.val.(int32) - srcV.val.(int32)
			} else if operator == "*" {
				zval.val = zval.val.(int32) * srcV.val.(int32)
			} else if operator == "/" {
				zval.val = zval.val.(int32) / srcV.val.(int32)
			}
		}
	} else if ZValTypeFloat == zval.valType {
		zval.valType = ZValTypeFloat
		if ZValTypeFloat == srcV.valType {
			if operator == "+" {
				zval.val = zval.val.(float32) + srcV.val.(float32)
			} else if operator == "-" {
				zval.val = zval.val.(float32) - srcV.val.(float32)
			} else if operator == "*" {
				zval.val = zval.val.(float32) * srcV.val.(float32)
			} else if operator == "/" {
				zval.val = zval.val.(float32) / srcV.val.(float32)
			}
		} else if ZValTypeInt == srcV.valType {
			if operator == "+" {
				zval.val = zval.val.(float32) + float32(srcV.val.(int32))
			} else if operator == "-" {
				zval.val = zval.val.(float32) - float32(srcV.val.(int32))
			} else if operator == "*" {
				zval.val = zval.val.(float32) * float32(srcV.val.(int32))
			} else if operator == "/" {
				zval.val = zval.val.(float32) / float32(srcV.val.(int32))
			}
		}
	}

	if dstOp.opType == OperandTypeReg {
		dstV := zvm.getRegVar(dstOp)
		*dstV = zval
	} else {
		varIndex := zvm.getZVal(zvm.curFuncIndex, dstOp)
		zvm.allFuncZVal[zvm.curFuncIndex][varIndex] = zval
	}
}

func (zvm *Zvm) processRet() bool {
	funcNode := zvm.getFuncNodeByIndex(zvm.curFuncIndex)
	log.Printf("processRet: %v\n", funcNode.funcName)
	if funcNode.funcName == "main" {
		// main函数返回时程序执行结束
		return true
	}

	zvm.curFuncIndex = zvm.funcIndexBeforeJmp
	zvm.curInstrIndex = zvm.instrIndexBeforeJmp + 1

	return false
}

func (zvm *Zvm) getRegVar(op *Operand) *ZVal {
	if OperandTypeReg == op.opType {
		if registerTypeT0 == op.opVal.(int32) {
			return &zvm.regT0
		} else if registerTypeT1 == op.opVal.(int32) {
			return &zvm.regT1
		} else if registerTypeRetVal == op.opVal.(int32) {
			return &zvm.regRetVal
		}
	}
	return nil
}

func (zvm *Zvm) getZVal(funcIndex int32, op *Operand) int32 {
	if zvm.allFuncZVal == nil {
		zvm.allFuncZVal = make(map[int32]FuncZVal)
	}

	funcZval := zvm.allFuncZVal[funcIndex]
	if funcZval == nil {
		funcZval = make(FuncZVal)
	}

	varIndex := op.opVal.(int32)

	_, ok := funcZval[varIndex]
	if ok {
		return varIndex
	}

	var varName string
	funcNode := zvm.getFuncNodeByIndex(funcIndex)

	symbolNode := funcNode.symbolNodes[varIndex]
	varName = symbolNode.identifier

	funcZval[varIndex] = ZVal{ZValTypeInvalid, varName, 0}
	zvm.allFuncZVal[funcIndex] = funcZval

	return varIndex
}

func (zvm *Zvm) getFuncNodeByIndex(index int32) *FuncNode {
	for i := 0; i < int(zvm.funcTable.count); i++ {
		if zvm.funcTable.funcNodes[i].index == index {
			return &zvm.funcTable.funcNodes[i]
		}
	}
	return nil
}

func (zvm *Zvm) getInstrByIndex(index int32) *Instr {
	for i := 0; i < int(zvm.instrStream.count); i++ {
		if zvm.instrStream.instrs[i].index == index {
			return &zvm.instrStream.instrs[i]
		}
	}
	return nil
}

func (instr *Instr) getOperantByIndex(index int32) *Operand {
	return &instr.ops[index]
}

func (zvm *Zvm) getStrByIndex(index int32) string {
	return zvm.strTable.strinfos[index].str
}

func (zvm *Zvm) getHostApiNode(index int32) *HostApiNode {
	for i := 0; i < int(zvm.hostApiTable.count); i++ {
		if zvm.hostApiTable.hostApiNodes[i].index == index {
			return &zvm.hostApiTable.hostApiNodes[i]
		}
	}
	return nil
}

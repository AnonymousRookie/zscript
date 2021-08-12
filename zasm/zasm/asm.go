package zasm

import (
	"../../utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func Asm(sourceFilename string) {

	suffix := sourceFilename[len(sourceFilename)-5:]
	if suffix != srcFileSuffix {
		utils.ExitWithErrMsg("source file suffix should be: " + srcFileSuffix)
	}

	// 读取源文件
	lines := utils.LoadSourceFile(sourceFilename)

	// 词法分析
	lexicalAnalyze(lines)

	// 语法分析
	parse()

	// fmt.Printf("header: %+v\n", header)
	// fmt.Printf("instrStream: %+v\n", instrStream)
	// fmt.Printf("strTable: %+v\n", strTable)
	// fmt.Printf("funcTable: %+v\n", funcTable)

	var outputFilename string
	outputFilename = sourceFilename[:len(sourceFilename)-5] + outFileSuffix

	fp, err := os.Create(outputFilename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fp.Close()

	buf := new(bytes.Buffer)

	// write header
	binary.Write(buf, binary.LittleEndian, header.isExistMainFunc)
	binary.Write(buf, binary.LittleEndian, header.mainFuncIndex)

	// write Instruction Stream
	binary.Write(buf, binary.LittleEndian, instrStream.count)
	for i := 0; i < int(instrStream.count); i++ {
		instr := instrStream.instrs[i]

		// fmt.Printf("instr:%+v\n", instr)

		binary.Write(buf, binary.LittleEndian, instr.index)
		binary.Write(buf, binary.LittleEndian, instr.instrType)
		binary.Write(buf, binary.LittleEndian, instr.opCount)
		for j := 0; j < int(instr.opCount); j++ {
			op := instr.ops[j]
			binary.Write(buf, binary.LittleEndian, op.opType)
			switch op.opType {
			case OperandTypeInt:
				binary.Write(buf, binary.LittleEndian, op.opVal.(int32))
			case OperandTypeFloat:
				binary.Write(buf, binary.LittleEndian, op.opVal.(float32))
			case OperandTypeStrIndex:
				binary.Write(buf, binary.LittleEndian, op.opVal.(int32))
			case OperandTypeIdentifierIndex:
				binary.Write(buf, binary.LittleEndian, op.opVal.(int32))
			case OperandTypeFuncIndex:
				binary.Write(buf, binary.LittleEndian, op.opVal.(int32))
			case OperandTypeReg:
				binary.Write(buf, binary.LittleEndian, op.opVal.(int32))
			}
		}
	}

	// write String Table
	binary.Write(buf, binary.LittleEndian, strTable.count)
	for i := 0; i < int(strTable.count); i++ {
		strinfo := strTable.strinfos[i]
		binary.Write(buf, binary.LittleEndian, strinfo.index)
		binary.Write(buf, binary.LittleEndian, strinfo.len)
		binary.Write(buf, binary.LittleEndian, []byte(strinfo.str))
	}

	// write Function Table
	binary.Write(buf, binary.LittleEndian, funcTable.count)
	for i := 0; i < int(funcTable.count); i++ {
		funcNode := funcTable.funcNodes[i]
		binary.Write(buf, binary.LittleEndian, funcNode.index)
		binary.Write(buf, binary.LittleEndian, funcNode.len)
		binary.Write(buf, binary.LittleEndian, []byte(funcNode.funcName))
		binary.Write(buf, binary.LittleEndian, funcNode.entryPoint)
		binary.Write(buf, binary.LittleEndian, funcNode.paramcount)
	}

	fp.Write(buf.Bytes())
}

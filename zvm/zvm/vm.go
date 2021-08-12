package zvm

import (
	"../../utils"
	"encoding/binary"
	"fmt"
	"os"
)

func Load(zseFilename string) {

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
	err = binary.Read(fp, binary.LittleEndian, &header.isExistMainFunc)
	utils.Check(err == nil, "read header.isExistMainFunc failed!")
	err = binary.Read(fp, binary.LittleEndian, &header.mainFuncIndex)
	utils.Check(err == nil, "read header.mainFuncIndex failed!")
	fmt.Printf("header:%+v\n", header)

	// Instruction Stream
	err = binary.Read(fp, binary.LittleEndian, &instrStream.count)
	utils.Check(err == nil, "read instrStream.count failed!")
	for i := 0; i < int(instrStream.count); i++ {
		var instr Instr
		err = binary.Read(fp, binary.LittleEndian, &instr.index)
		utils.Check(err == nil, "read instr.index failed!")
		err = binary.Read(fp, binary.LittleEndian, &instr.instrType)
		utils.Check(err == nil, "read instr.instrType failed!")
		err = binary.Read(fp, binary.LittleEndian, &instr.opCount)
		utils.Check(err == nil, "read instr.opCount failed!")

		// fmt.Printf("instr:%+v\n", instr)

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
				if err != nil {
					fmt.Printf("op:%+v\n", op)
				}
				utils.Check(err == nil, "read op.opVal.(int32) failed!")
				op.opVal = v
			}
			instr.ops = append(instr.ops, op)
		}
		instrStream.instrs = append(instrStream.instrs, instr)
	}
	fmt.Printf("instrStream:%+v\n", instrStream)

	// String Table
	err = binary.Read(fp, binary.LittleEndian, &strTable.count)
	utils.Check(err == nil, "read strTable.count failed!")
	for i := 0; i < int(strTable.count); i++ {
		var strInfo StrInfo
		err = binary.Read(fp, binary.LittleEndian, &strInfo.index)
		utils.Check(err == nil, "read strInfo.index failed!")

		err = binary.Read(fp, binary.LittleEndian, &strInfo.len)
		utils.Check(err == nil, "read strInfo.len failed!")

		buf := make([]byte, strInfo.len)
		err = binary.Read(fp, binary.LittleEndian, &buf)
		utils.Check(err == nil, "read strInfo.str failed!")
		strInfo.str = string(buf)

		strTable.strinfos = append(strTable.strinfos, strInfo)
	}
	fmt.Printf("strTable:%+v\n", strTable)

	// Function Table
	err = binary.Read(fp, binary.LittleEndian, &funcTable.count)
	utils.Check(err == nil, "read strTable.count failed!")
	for i := 0; i < int(funcTable.count); i++ {
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

		funcTable.funcNodes = append(funcTable.funcNodes, funcNode)
	}
	fmt.Printf("funcTable:%+v\n", funcTable)
}

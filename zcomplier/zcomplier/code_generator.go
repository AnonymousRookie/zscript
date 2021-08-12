package zcomplier

import (
	"../../utils"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func generateCode(outputFilename string) {
	fmt.Println("Generate: ", outputFilename)
	outputFile, err := os.OpenFile(outputFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		utils.ExitWithErrMsg("generate " + outputFilename + " failed!")
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)

	var ret []string

	ret = append(ret, "; Global variables\n")
	generateSymbols(0, SymbolTypeVar, &ret)

	var mainFunc *FuncNode
	ret = append(ret, "\n; Functions\n")
	for name, node := range FuncTable {
		if name == "main" {
			mainFunc = node
			generateFuncNode(mainFunc, &ret)
		} else {
			generateFuncNode(node, &ret)
		}
	}

	for i := 0; i < len(ret); i++ {
		outputWriter.WriteString(ret[i])
	}
	outputWriter.Flush()
}

func generateSymbols(scope int, symbolType SymbolType, ret *[]string) {
	for e := symbolList.Front(); e != nil; e = e.Next() {
		if sn, ok := e.Value.(*SymbolNode); ok {
			if sn.symbolType == symbolType && sn.socpeIndex == scope {
				if symbolType == SymbolTypeVar {
					if scope != GlobalScope {
						*ret = append(*ret, "\t")
					}
					s := fmt.Sprintf("Var %s\n", sn.identifier)
					*ret = append(*ret, s)
				} else if symbolType == SymbolTypeParam {
					s := fmt.Sprintf("\tParam %s\n", sn.identifier)
					*ret = append(*ret, s)
				}
			}
		}
	}
}

func generateFuncNode(funcNode *FuncNode, ret *[]string) {
	var str string
	str = fmt.Sprintf("Func %s\n{\n", funcNode.FuncName)
	*ret = append(*ret, str)
	// 函数参数
	generateSymbols(funcNode.FuncIndex, SymbolTypeParam, ret)
	// 函数内局部变量
	generateSymbols(funcNode.FuncIndex, SymbolTypeVar, ret)

	iCodeNodeList := getIcodeNodeList(funcNode.FuncIndex)
	if iCodeNodeList == nil {
		return
	}

	for i := 0; i < len(*iCodeNodeList); i++ {
		iCodeNode := getICodeNode(funcNode.FuncIndex, i)

		switch iCodeNode.ICodeNodeType {
		case ICodeNodeTypeInstr:
			iCodeInstr := iCodeNode.Val.(ICodeInstr)
			instrName, ok := instrMap[iCodeInstr.instr]
			if ok {
				str = fmt.Sprintf("\t%s ", instrName)
				*ret = append(*ret, str)

				var instrOps []string

				operands := iCodeInstr.operands
				for i := 0; i < len(operands); i++ {
					operand := operands[i]

					switch operand.OperandType {
					case OperandTypeInt:
						str = strconv.FormatInt(int64(operand.Val.(int)), 10)
					case OperandTypeFloat:
						str = strconv.FormatFloat(float64(operand.Val.(float32)), 'f', 3, 64)
					case OperandTypeString:
						str = "\"" + getStrByIndex(operand.Val.(int)) + "\""
					case OperandTypeVar:
						symbolNode := getSymbolNodeByIndex(operand.Val.(int))
						str = symbolNode.identifier
					case OperandTypeFuncIndex:
						fnode := getFuncByIndex(operand.Val.(int))
						if fnode != nil {
							str = fnode.FuncName
						}
					case OperandTypeReg:
						str = registerReturnVal

					default:
						utils.ExitWithErrMsg("unexpected OperandType: " + strconv.Itoa(operand.OperandType))
					}

					instrOps = append(instrOps, str)
				}

				*ret = append(*ret, strings.Join(instrOps, ", "))
				*ret = append(*ret, "\n")
			}
		case ICodeNodeTypeSourceLine:
			str = fmt.Sprintf("\n\t; %s\n", iCodeNode.Val.(string))
			*ret = append(*ret, str)
		default:
			utils.ExitWithErrMsg("unexpected ICodeNodeType: " + strconv.Itoa(iCodeNode.ICodeNodeType))
		}
	}

	*ret = append(*ret, "}\n")
}

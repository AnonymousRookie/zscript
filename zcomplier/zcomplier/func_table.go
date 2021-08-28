package zcomplier

import (
	"fmt"
	"os"
)

type FuncNode struct {
	FuncName   string // 函数名称
	ParamCount int    // 函数参数个数
	FuncIndex  int    // 函数索引
}

var funcTable []FuncNode

func addFunc(funcName string) {
	for i := 0; i < len(funcTable); i++ {
		funNode := funcTable[i]
		if funNode.FuncName == funcName {
			fmt.Printf("[error] function redefinition, function name: %s!\n", funcName)
			os.Exit(-1)
		}
	}

	var newNode FuncNode
	newNode.FuncName = funcName
	newNode.ParamCount = 0
	newNode.FuncIndex = len(funcTable) + 1
	funcTable = append(funcTable, newNode)
}

func getFuncByName(funcName string) *FuncNode {
	for i := 0; i < len(funcTable); i++ {
		if funcTable[i].FuncName == funcName {
			return &funcTable[i]
		}
	}
	return nil
}

func getFuncByIndex(funcIndex int) *FuncNode {
	for i := 0; i < len(funcTable); i++ {
		funNode := funcTable[i]
		if funNode.FuncIndex == funcIndex {
			return &funNode
		}
	}
	return nil
}

func getHostApiByIndex(index int) string {
	return "print"
}

func printFuncTable() {
	fmt.Println("funcTable:")
	for _, node := range funcTable {
		fmt.Printf("funcNode: %+v\n", node)
	}
}

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

var FuncTable map[string]*FuncNode = make(map[string]*FuncNode)

func addFunc(funcName string) {
	_, ok := FuncTable[funcName]
	if ok {
		fmt.Printf("[error] function redefinition, function name: %s!\n", funcName)
		os.Exit(-1)
	}
	var newNode FuncNode
	newNode.FuncName = funcName
	newNode.ParamCount = 0
	newNode.FuncIndex = len(FuncTable) + 1
	FuncTable[funcName] = &newNode
}

func getFuncByName(funcName string) *FuncNode {
	node, ok := FuncTable[funcName]
	if ok {
		return node
	}
	return nil
}

func getFuncByIndex(funcIndex int) *FuncNode {
	for _, node := range FuncTable {
		if node.FuncIndex == funcIndex {
			return node
		}
	}
	return nil
}

func printFuncTable() {
	fmt.Println("FuncTable:")
	for _, node := range FuncTable {
		fmt.Printf("FuncNode: %+v\n", node)
	}
}

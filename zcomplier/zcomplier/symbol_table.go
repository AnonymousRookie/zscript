package zcomplier

import (
	"container/list"
	"fmt"
)

type SymbolType int

const (
	SymbolTypeInvalid = iota
	SymbolTypeVar     // 变量
	SymbolTypeParam   // 函数参数
)

type SymbolNode struct {
	index      int
	identifier string
	symbolType SymbolType
	socpeIndex int // 在哪个函数中
}

// 临时变量在符号表中的索引
var TempVar0SymbolIndex int = 0
var TempVar1SymbolIndex int = 0

var symbolList *list.List = list.New()

// 返回新增symbol在symbolList中的索引
func addSymbol(identifier string, symbolType SymbolType, socpeIndex int) int {
	index := symbolList.Len()
	node := &SymbolNode{index, identifier, symbolType, socpeIndex}
	symbolList.PushBack(node)
	return index
}

func getSymbolNode(identifier string, socpeIndex int) *SymbolNode {
	for e := symbolList.Front(); e != nil; e = e.Next() {
		if sn, ok := e.Value.(*SymbolNode); ok {
			if sn.identifier == identifier && sn.socpeIndex == socpeIndex {
				return sn
			}
		}
	}
	return nil
}

func getSymbolNodeByIndex(symbolIndex int) *SymbolNode {
	var i int = 0
	for e := symbolList.Front(); e != nil; e = e.Next() {
		if i == symbolIndex {
			sn, ok := e.Value.(*SymbolNode)
			if ok {
				return sn
			}
		}
		i++
	}
	return nil
}

func printSymbolList() {
	fmt.Println("SymbolTable:")
	for e := symbolList.Front(); e != nil; e = e.Next() {
		if sn, ok := e.Value.(*SymbolNode); ok {
			fmt.Printf("SymbolNode: %+v\n", sn)
		}
	}
}

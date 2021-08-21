package main

import (
	"fmt"
	"os"
	"zscript/zasm/zasm"
)

func printUsage() {
	fmt.Println("Usage: zasm zasmfile.zasm")
}

func main() {
	fmt.Println("\nzasm starting...")

	if len(os.Args) != 2 {
		printUsage()
		return
	}

	sourceFilename := os.Args[1]

	fmt.Println("sourceFilename: ", sourceFilename)

	zasm.Asm(sourceFilename)

	fmt.Println("\nzasm stopped!")
}

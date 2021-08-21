package main

import (
	"fmt"
	"os"

	"zscript/zvm/zvm"
)

func printUsage() {
	fmt.Println("Usage: zvm file.zse")
}

func main() {
	fmt.Println("\nzvm starting...")

	if len(os.Args) != 2 {
		printUsage()
		return
	}

	filename := os.Args[1]

	fmt.Println("filename: ", filename)

	zvm := zvm.NewZvm()
	zvm.Load(filename)
	zvm.Run()

	fmt.Println("\nzvm stopped!")
}

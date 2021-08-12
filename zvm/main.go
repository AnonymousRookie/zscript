package main

import (
	"./zvm"
	"fmt"
	"os"
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

	zvm.Load(filename)

	fmt.Println("\nzvm stopped!")
}

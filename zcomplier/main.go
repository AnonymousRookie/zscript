package main

import (
	"fmt"
	"os"

	"zscript/zcomplier/zcomplier"
)

func printUsage() {
	fmt.Println("Usage: zscript source.zs")
}

func main() {
	fmt.Println("\nzcomplier starting...")

	if len(os.Args) != 2 {
		printUsage()
		return
	}

	sourceFilename := os.Args[1]

	fmt.Println("sourceFilename: ", sourceFilename)

	zcomplier.Complie(sourceFilename)

	fmt.Println("\nzcomplier stopped!")
}

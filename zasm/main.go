package main

import (
	"fmt"
	"log"
	"os"
	"zscript/zasm/zasm"
)

func printUsage() {
	fmt.Println("Usage: zasm zasmfile.zasm")
}

func init() {
	logfilename := "./zasm.log"
	logFile, err := os.OpenFile(logfilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lmicroseconds | log.Ldate)
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

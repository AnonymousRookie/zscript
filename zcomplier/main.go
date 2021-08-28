package main

import (
	"fmt"
	"log"
	"os"

	"zscript/zcomplier/zcomplier"
)

func printUsage() {
	fmt.Println("Usage: zscript source.zs")
}

func init() {
	logfilename := "./zcomplier.log"
	logFile, err := os.OpenFile(logfilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lmicroseconds | log.Ldate)
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

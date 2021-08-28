package main

import (
	"fmt"
	"log"
	"os"

	"zscript/zvm/zvm"
)

func printUsage() {
	fmt.Println("Usage: zvm file.zse")
}

func init() {
	logfilename := "./zvm.log"
	logFile, err := os.OpenFile(logfilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Lmicroseconds | log.Ldate)
}

func main() {
	log.Println("\nzvm starting...")

	if len(os.Args) != 2 {
		printUsage()
		return
	}

	filename := os.Args[1]

	log.Println("filename: ", filename)

	zvm := zvm.NewZvm()
	zvm.Load(filename)
	zvm.Run()

	log.Println("\nzvm stopped!")
}

package main

import (
	"os"
	"runtime"

	mack "github.com/acmerocket/mack"
)

func init() {
	if cpu := runtime.NumCPU(); cpu == 1 {
		runtime.GOMAXPROCS(2)
	} else {
		runtime.GOMAXPROCS(cpu)
	}
}

func main() {
	mack := mack.MarkdownAck{Out: os.Stdout, Err: os.Stderr}
	exitCode := mack.Run(os.Args[1:])
	os.Exit(exitCode)
}

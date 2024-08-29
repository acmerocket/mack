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

func run(args []string) int {
	mack := mack.MarkdownAck{Out: os.Stdout, Err: os.Stderr}
	return mack.Run(args)
}

func main() {
	os.Exit(run(os.Args[1:]))
}

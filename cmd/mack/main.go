package main

import (
	"io"
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

func run_main(args []string, stdout, stderr io.Writer) int {
	mack := mack.MarkdownAck{Out: stdout, Err: stderr}
	return mack.Run(args)
}

func main() {
	os.Exit(run_main(os.Args[1:], os.Stdout, os.Stderr))
}

package mack

import (
	"io"
	"os"
	"testing"
)

func BenchmarkMack(b *testing.B) {
	for name, args := range map[string][]string{
		"findOnly": {"-g", ".", os.Getenv("GOPATH")},
		"normal":   {"test", os.Getenv("GOPATH")},
	} {
		args := args
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ret := MarkdownAck{Out: io.Discard, Err: os.Stderr}.Run(args)
				if ret != 0 {
					b.Fatal("failed")
				}
			}
		})
	}
}

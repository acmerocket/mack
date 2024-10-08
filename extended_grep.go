package mack

import (
	"io"
	"log"
	"os"
)

type extendedGrep struct {
	lineGrep
	pattern regex_pattern
}

func (g extendedGrep) grep(path string, buf []byte) {
	f, err := getFileHandler(path)
	if err != nil {
		log.Fatalf("open: %s\n", err)
	}
	defer f.Close()

	if f == os.Stdin {
		// TODO: File type is fixed in ASCII because it can not determine the character code.
		g.grepEachLines(f, ASCII, func(b []byte) bool {
			return g.pattern.regexp.Match(b)
		}, func(b []byte) int {
			return g.pattern.regexp.FindIndex(b)[0] + 1
		})
		return
	}

	c, err := f.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatalf("read: %s\n", err)
	}

	if err == io.EOF {
		return
	}

	// detect encoding.
	limit := c
	if limit > 512 {
		limit = 512
	}

	encoding := detectEncoding(buf[:limit])
	if encoding == ERROR || encoding == BINARY {
		return
	}

	// grep each lines.
	g.grepEachLines(f, encoding, func(b []byte) bool {
		return g.pattern.regexp.Match(b)
	}, func(b []byte) int {
		return g.pattern.regexp.FindIndex(b)[0] + 1
	})
}

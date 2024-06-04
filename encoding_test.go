package mack

import (
	"os"
	"path/filepath"
	"testing"
)

type Assert struct {
	path     string
	fileType int
}

var Asserts = []Assert{
	{"ascii.txt", ASCII},
	{"binary/binary.bin", BINARY},
	{"ja/euc-jp.txt", EUCJP},
	{"ja/shift_jis.txt", SHIFTJIS},
	{"ja/utf8.txt", UTF8},
	{"ja/broken_euc-jp.txt", EUCJP},
	{"ja/broken_shift_jis.txt", SHIFTJIS},
	{"ja/broken_utf8.txt", UTF8},
}

func TestIdentifyType(t *testing.T) {
	for _, f := range Asserts {
		b, _ := os.ReadFile(filepath.Join("files", f.path))
		fileType := detectEncoding(b)
		if fileType != f.fileType {
			t.Errorf("%s should be %d.", f.path, f.fileType)
		}
	}
}

package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func exec_main(T *testing.T, argv []string) []byte {
	var buf bytes.Buffer
	exitcode := run_main(argv, &buf, os.Stderr)
	if exitcode != 0 {
		T.Errorf("Wrong exit code for args: %v, expected: %v, got: %v", argv, 0, exitcode)
	}
	return buf.Bytes()
}

func assert_like_ack(T *testing.T, args string) {
	// given a string of arguments, execute against both `ack` and the local build
	// designed to test compatibility with ack https://beyondgrep.com/

	argv := strings.Split(args, " ")

	// run local, grab results
	test_out := exec_main(T, argv)

	// run against `ack`, grab results
	//ack, err := exec.LookPath("ack")
	//if err != nil {
	//	log.Fatal(err)
	//}
	out, err := exec.Command("ack", argv...).Output()
	if err != nil {
		T.Error(out, err, "ack", args)
	}

	// compare
	if !bytes.Equal(test_out, out) {
		T.Errorf("Output not equal, expected: %v, got: %v", string(out), string(test_out))
	}
}

func TestListFiles(T *testing.T) {
	assert_like_ack(T, "-f")
	//assert_like_ack(T, "-f ../../files") // doesn't work because hidden files work differently.
	assert_like_ack(T, "-t html -f ../../files")
	assert_like_ack(T, "-t markdown -f ../../files")
}

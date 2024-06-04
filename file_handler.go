package mack

import "os"

func getFileHandler(path string) (*os.File, error) {
	if path == "" {
		return os.Stdin, nil
	} else {
		return os.Open(path)
	}
}

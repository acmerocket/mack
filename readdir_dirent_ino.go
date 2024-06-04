//go:build (linux && !appengine) || darwin
// +build linux,!appengine darwin

package mack

import "syscall"

func direntInode(dirent *syscall.Dirent) uint64 {
	return uint64(dirent.Ino)
}

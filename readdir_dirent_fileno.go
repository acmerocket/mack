//go:build freebsd || openbsd || netbsd
// +build freebsd openbsd netbsd

package mack

import "syscall"

func direntInode(dirent *syscall.Dirent) uint64 {
	return uint64(dirent.Fileno)
}

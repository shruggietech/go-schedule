package main

import (
	"os"
	"syscall"
)

func isLinkedPath(path string, info os.FileInfo) (bool, error) {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, err
	}
	attributes, err := syscall.GetFileAttributes(p)
	return info.Mode()&os.ModeSymlink != 0 || attributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT != 0, err
}

//go:build !windows

package main

import "os"

func isLinkedPath(_ string, info os.FileInfo) (bool, error) {
	return info.Mode()&os.ModeSymlink != 0, nil
}

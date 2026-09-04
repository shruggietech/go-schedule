//go:build darwin

package releasegate

import "golang.org/x/sys/unix"

func renameDispositionNoReplace(from, to string) error {
	return unix.RenamexNp(from, to, unix.RENAME_EXCL)
}

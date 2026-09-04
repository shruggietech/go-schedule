//go:build windows

package releasegate

import "golang.org/x/sys/windows"

func renameDispositionNoReplace(from, to string) error {
	fromUTF16, err := windows.UTF16PtrFromString(from)
	if err != nil {
		return err
	}
	toUTF16, err := windows.UTF16PtrFromString(to)
	if err != nil {
		return err
	}
	return windows.MoveFile(fromUTF16, toUTF16)
}

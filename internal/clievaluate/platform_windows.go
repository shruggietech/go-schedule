//go:build windows

package clievaluate

import (
	"errors"
	"fmt"

	"golang.org/x/sys/windows"
)

func workerPrivilege() error {
	if windows.GetCurrentProcessToken().IsElevated() {
		return errors.New("the unattended QA worker must run without elevation")
	}
	return nil
}

func platformDetails() (string, *string) {
	v := windows.RtlGetVersion()
	name := "Windows"
	if v.MajorVersion == 10 && v.BuildNumber >= 22000 {
		name = "Windows 11"
	}
	return fmt.Sprintf("%s (kernel %d.%d, build %d)", name, v.MajorVersion, v.MinorVersion, v.BuildNumber), nil
}

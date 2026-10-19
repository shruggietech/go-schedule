//go:build windows

package clifixture

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

func TestWindowsOwnedChildDoesNotInheritUnrelatedHandles(t *testing.T) {
	if os.Getenv("GO_SCHEDULE_HANDLE_LEAK_CHECK") == "1" {
		value, err := strconv.ParseUint(os.Getenv("GO_SCHEDULE_HANDLE_VALUE"), 10, 64)
		if err != nil {
			os.Exit(3)
		}
		buffer := make([]uint16, 32768)
		length, err := windows.GetFinalPathNameByHandle(windows.Handle(value), &buffer[0], uint32(len(buffer)), 0)
		if err == nil && length < uint32(len(buffer)) {
			actual := strings.TrimPrefix(windows.UTF16ToString(buffer[:length]), `\\?\`)
			if strings.EqualFold(actual, os.Getenv("GO_SCHEDULE_HANDLE_SENTINEL")) {
				os.Exit(2)
			}
		}
		os.Exit(0)
	}
	directory := t.TempDir()
	sentinel := filepath.Join(directory, "unrelated-file")
	file, err := os.OpenFile(sentinel, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var inheritable windows.Handle
	if err := windows.DuplicateHandle(windows.CurrentProcess(), windows.Handle(file.Fd()), windows.CurrentProcess(), &inheritable, 0, true, windows.DUPLICATE_SAME_ACCESS); err != nil {
		t.Fatal(err)
	}
	defer windows.CloseHandle(inheritable)
	p, err := StartProcess(Command{Executable: os.Args[0], Arguments: []string{"-test.run=^TestWindowsOwnedChildDoesNotInheritUnrelatedHandles$"}, Environment: append(os.Environ(), "GO_SCHEDULE_HANDLE_LEAK_CHECK=1", "GO_SCHEDULE_HANDLE_VALUE="+strconv.FormatUint(uint64(inheritable), 10), "GO_SCHEDULE_HANDLE_SENTINEL="+sentinel), Directory: directory, StdoutPath: filepath.Join(directory, "stdout"), StderrPath: filepath.Join(directory, "stderr")})
	if err != nil {
		t.Fatal(err)
	}
	defer p.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.Wait(ctx); err != nil {
		t.Fatalf("unrelated inheritable handle leaked: %v", err)
	}
}

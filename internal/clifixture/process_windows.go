//go:build windows

package clifixture

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

func startContained(command Command, input, stdout, stderr *os.File) (*Process, func() error, error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, nil, err
	}
	limits := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	limits.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if command.Limits != nil {
		limits.BasicLimitInformation.LimitFlags |= windows.JOB_OBJECT_LIMIT_JOB_MEMORY | windows.JOB_OBJECT_LIMIT_ACTIVE_PROCESS
		limits.BasicLimitInformation.ActiveProcessLimit = command.Limits.MaxProcesses
		limits.JobMemoryLimit = uintptr(command.Limits.MemoryBytes)
	}
	if _, err := windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&limits)), uint32(unsafe.Sizeof(limits))); err != nil {
		_ = windows.CloseHandle(job)
		return nil, nil, err
	}
	if command.Limits != nil {
		cpu := struct {
			Flags uint32
			Rate  uint32
		}{1 | 4, uint32(command.Limits.CPUPercent * 100)}
		if _, err := windows.SetInformationJobObject(job, windows.JobObjectCpuRateControlInformation, uintptr(unsafe.Pointer(&cpu)), uint32(unsafe.Sizeof(cpu))); err != nil {
			return nil, nil, errors.Join(err, windows.CloseHandle(job))
		}
	}
	var inherited []windows.Handle
	defer func() {
		for _, h := range inherited {
			_ = windows.CloseHandle(h)
		}
	}()
	duplicate := func(file *os.File) (windows.Handle, error) {
		var h windows.Handle
		err := windows.DuplicateHandle(windows.CurrentProcess(), windows.Handle(file.Fd()), windows.CurrentProcess(), &h, 0, true, windows.DUPLICATE_SAME_ACCESS)
		if err == nil {
			inherited = append(inherited, h)
		}
		return h, err
	}
	in, err := duplicate(input)
	if err != nil {
		_ = windows.CloseHandle(job)
		return nil, nil, err
	}
	out, err := duplicate(stdout)
	if err != nil {
		_ = windows.CloseHandle(job)
		return nil, nil, err
	}
	errout, err := duplicate(stderr)
	if err != nil {
		_ = windows.CloseHandle(job)
		return nil, nil, err
	}
	app, err := windows.UTF16PtrFromString(command.Executable)
	if err != nil {
		_ = windows.CloseHandle(job)
		return nil, nil, err
	}
	args := []string{syscall.EscapeArg(command.Executable)}
	for _, arg := range command.Arguments {
		args = append(args, syscall.EscapeArg(arg))
	}
	line, err := windows.UTF16PtrFromString(strings.Join(args, " "))
	if err != nil {
		_ = windows.CloseHandle(job)
		return nil, nil, err
	}
	directory, err := windows.UTF16PtrFromString(command.Directory)
	if err != nil {
		_ = windows.CloseHandle(job)
		return nil, nil, err
	}
	var environment *uint16
	var block []uint16
	if command.Environment != nil {
		values, err := normalizedWindowsEnvironment(command.Environment)
		if err != nil {
			return nil, nil, errors.Join(err, windows.CloseHandle(job))
		}
		joined := strings.Join(values, "\x00") + "\x00\x00"
		block = utf16.Encode([]rune(joined))
		environment = &block[0]
	}
	attributes, err := windows.NewProcThreadAttributeList(1)
	if err != nil {
		_ = windows.CloseHandle(job)
		return nil, nil, err
	}
	defer attributes.Delete()
	handles := []windows.Handle{in, out, errout}
	if err := attributes.Update(windows.PROC_THREAD_ATTRIBUTE_HANDLE_LIST, unsafe.Pointer(&handles[0]), uintptr(len(handles))*unsafe.Sizeof(handles[0])); err != nil {
		_ = windows.CloseHandle(job)
		return nil, nil, err
	}
	info := windows.StartupInfoEx{StartupInfo: windows.StartupInfo{Cb: uint32(unsafe.Sizeof(windows.StartupInfoEx{})), Flags: windows.STARTF_USESTDHANDLES | windows.STARTF_USESHOWWINDOW, ShowWindow: windows.SW_HIDE, StdInput: in, StdOutput: out, StdErr: errout}, ProcThreadAttributeList: attributes.List()}
	var child windows.ProcessInformation
	flags := uint32(windows.CREATE_NO_WINDOW | windows.CREATE_SUSPENDED | windows.CREATE_UNICODE_ENVIRONMENT | windows.EXTENDED_STARTUPINFO_PRESENT)
	if command.Limits != nil {
		flags |= windows.BELOW_NORMAL_PRIORITY_CLASS
	}
	if err := windows.CreateProcess(app, line, nil, nil, true, flags, environment, directory, &info.StartupInfo, &child); err != nil {
		_ = windows.CloseHandle(job)
		return nil, nil, err
	}
	// The suspended child cannot execute, create grandchildren or escape before assignment.
	if err := windows.AssignProcessToJobObject(job, child.Process); err != nil {
		_ = windows.TerminateProcess(child.Process, 1)
		_ = windows.CloseHandle(child.Thread)
		_ = windows.CloseHandle(child.Process)
		_ = windows.CloseHandle(job)
		return nil, nil, fmt.Errorf("contain owned child before execution: %w", err)
	}
	if _, err := windows.ResumeThread(child.Thread); err != nil {
		_ = windows.TerminateJobObject(job, 1)
		_ = windows.CloseHandle(child.Thread)
		_ = windows.CloseHandle(child.Process)
		_ = windows.CloseHandle(job)
		return nil, nil, err
	}
	_ = windows.CloseHandle(child.Thread)
	p := &Process{PID: int(child.ProcessId)}
	p.stop = func() error {
		terminateErr := windows.TerminateJobObject(job, 1)
		// Termination is asynchronous. Keep the exact job handle until all owned
		// descendants have exited, before callers may remove their runtime files.
		drainErr := waitForEmptyJob(job, 5*time.Second)
		closeErr := windows.CloseHandle(job)
		return errors.Join(terminateErr, drainErr, closeErr)
	}
	wait := func() error {
		defer func() { _ = windows.CloseHandle(child.Process) }()
		if _, err := windows.WaitForSingleObject(child.Process, windows.INFINITE); err != nil {
			return err
		}
		var code uint32
		if err := windows.GetExitCodeProcess(child.Process, &code); err != nil {
			return err
		}
		if code != 0 {
			return childExitError{code: int(code)}
		}
		return nil
	}
	return p, wait, nil
}

type childExitError struct{ code int }

func (e childExitError) Error() string { return fmt.Sprintf("owned child exited with code %d", e.code) }
func (e childExitError) ExitCode() int { return e.code }

func normalizedWindowsEnvironment(values []string) ([]string, error) {
	byName := map[string]string{}
	for _, value := range values {
		if strings.ContainsRune(value, '\x00') {
			return nil, errors.New("owned child environment contains NUL")
		}
		index := strings.IndexByte(value, '=')
		if index == 0 {
			next := strings.IndexByte(value[1:], '=')
			if next >= 0 {
				index = next + 1
			}
		}
		if index <= 0 {
			return nil, errors.New("owned child environment requires name=value entries")
		}
		name := strings.ToUpper(value[:index])
		byName[name] = value
	}
	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)
	result := make([]string, 0, len(names))
	for _, name := range names {
		result = append(result, byName[name])
	}
	return result, nil
}

// jobAccounting matches JOBOBJECT_BASIC_ACCOUNTING_INFORMATION in WinNT.h.
type jobAccounting struct {
	TotalUserTime, TotalKernelTime, PeriodUserTime, PeriodKernelTime int64
	PageFaults, TotalProcesses, ActiveProcesses, TerminatedProcesses uint32
}

func waitForEmptyJob(job windows.Handle, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		var accounting jobAccounting
		if err := windows.QueryInformationJobObject(job, windows.JobObjectBasicAccountingInformation, uintptr(unsafe.Pointer(&accounting)), uint32(unsafe.Sizeof(accounting)), nil); err != nil {
			return fmt.Errorf("inspect owned process tree during cleanup: %w", err)
		}
		if accounting.ActiveProcesses == 0 {
			return nil
		}
		if !time.Now().Before(deadline) {
			return fmt.Errorf("owned process tree still has %d active processes at cleanup deadline", accounting.ActiveProcesses)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

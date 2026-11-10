//go:build windows

package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

type windowsOwned struct {
	process, thread, job windows.Handle
	assigned             bool
}

// startOwnedProcess establishes job membership while the primary thread is
// suspended. Exact native handles remain owned until wait and cleanup finish.
func startOwnedProcess(ctx context.Context, cmd *exec.Cmd) (owned ownedProcess, err error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if cmd.Err != nil {
		return nil, cmd.Err
	}
	application, err := windowsExecutable(cmd)
	if err != nil {
		return nil, err
	}
	app, err := windows.UTF16PtrFromString(application)
	if err != nil {
		return nil, err
	}
	args := make([]string, len(cmd.Args))
	for i, arg := range cmd.Args {
		args[i] = syscall.EscapeArg(arg)
	}
	commandLine := strings.Join(args, " ")
	if cmd.SysProcAttr != nil && cmd.SysProcAttr.CmdLine != "" {
		commandLine = cmd.SysProcAttr.CmdLine
	}
	line, err := windows.UTF16PtrFromString(commandLine)
	if err != nil {
		return nil, err
	}
	var directory *uint16
	if cmd.Dir != "" {
		directory, err = windows.UTF16PtrFromString(cmd.Dir)
		if err != nil {
			return nil, err
		}
	}
	environment, err := windowsEnvironment(cmd.Environ())
	if err != nil {
		return nil, err
	}
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, err
	}
	// Explicit application-requested breakaway remains available. Ordinary
	// descendants still join this job; no silent breakaway or resource limits.
	limits := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	limits.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_BREAKAWAY_OK
	if _, err = windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&limits)), uint32(unsafe.Sizeof(limits))); err != nil {
		return nil, errors.Join(err, windows.CloseHandle(job))
	}
	p := &windowsOwned{job: job}
	var inherited []windows.Handle
	defer func() {
		for _, h := range inherited {
			err = errors.Join(err, windows.CloseHandle(h))
		}
		if owned == nil {
			err = errors.Join(err, windows.CloseHandle(job))
		}
	}()
	duplicate := func(f *os.File) (windows.Handle, error) {
		var h windows.Handle
		err := windows.DuplicateHandle(windows.CurrentProcess(), windows.Handle(f.Fd()), windows.CurrentProcess(), &h, 0, true, windows.DUPLICATE_SAME_ACCESS)
		if err == nil {
			inherited = append(inherited, h)
		}
		return h, err
	}
	in, err := duplicate(cmd.Stdin.(*os.File))
	if err != nil {
		return nil, err
	}
	out, err := duplicate(cmd.Stdout.(*os.File))
	if err != nil {
		return nil, err
	}
	attributes, err := windows.NewProcThreadAttributeList(1)
	if err != nil {
		return nil, err
	}
	handles := []windows.Handle{in, out}
	defer func() { attributes.Delete(); runtime.KeepAlive(handles) }()
	if err = attributes.Update(windows.PROC_THREAD_ATTRIBUTE_HANDLE_LIST, unsafe.Pointer(&handles[0]), uintptr(len(handles))*unsafe.Sizeof(handles[0])); err != nil {
		return nil, err
	}
	info := windows.StartupInfoEx{StartupInfo: windows.StartupInfo{Cb: uint32(unsafe.Sizeof(windows.StartupInfoEx{})), Flags: windows.STARTF_USESTDHANDLES | windows.STARTF_USESHOWWINDOW, ShowWindow: windows.SW_HIDE, StdInput: in, StdOutput: out, StdErr: out}, ProcThreadAttributeList: attributes.List()}
	var child windows.ProcessInformation
	flags := uint32(windows.CREATE_NO_WINDOW | windows.CREATE_SUSPENDED | windows.CREATE_UNICODE_ENVIRONMENT | windows.EXTENDED_STARTUPINFO_PRESENT)
	if cmd.SysProcAttr != nil {
		flags |= cmd.SysProcAttr.CreationFlags
	}
	if err = windows.CreateProcess(app, line, nil, nil, true, flags, &environment[0], directory, &info.StartupInfo, &child); err != nil {
		return nil, err
	}
	p.process = child.Process
	p.thread = child.Thread
	owned = p
	if err = windows.AssignProcessToJobObject(job, child.Process); err != nil {
		return p, fmt.Errorf("establish owned process tree: %w", err)
	}
	p.assigned = true
	if err = ctx.Err(); err != nil {
		return p, err
	}
	if _, err = windows.ResumeThread(child.Thread); err != nil {
		return p, fmt.Errorf("resume owned process: %w", err)
	}
	return p, nil
}

func (p *windowsOwned) wait() nativeResult {
	if _, err := windows.WaitForSingleObject(p.process, windows.INFINITE); err != nil {
		return nativeResult{err: err}
	}
	var value uint32
	if err := windows.GetExitCodeProcess(p.process, &value); err != nil {
		return nativeResult{err: err}
	}
	code := int(value)
	var err error
	if code != 0 {
		err = fmt.Errorf("process exited with code %d", code)
	}
	return nativeResult{code: &code, err: err}
}
func (p *windowsOwned) terminate(deadline time.Time) error {
	var err error
	if p.assigned {
		err = windows.TerminateJobObject(p.job, 1)
	} else {
		err = windows.TerminateProcess(p.process, 1)
	}
	if err != nil {
		empty, queryErr := p.empty()
		if queryErr == nil && empty {
			return nil
		}
		return errors.Join(fmt.Errorf("terminate owned process tree: %w", err), queryErr)
	}
	for {
		empty, queryErr := p.empty()
		if queryErr != nil {
			return queryErr
		}
		if empty {
			return nil
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return errors.New("owned process tree remains active at cancellation deadline")
		}
		time.Sleep(min(10*time.Millisecond, remaining))
	}
}
func (p *windowsOwned) empty() (bool, error) {
	if !p.assigned {
		status, err := windows.WaitForSingleObject(p.process, 0)
		return status == windows.WAIT_OBJECT_0, err
	}
	var counts struct {
		TotalUserTime, TotalKernelTime, PeriodUserTime, PeriodKernelTime int64
		PageFaults, TotalProcesses, ActiveProcesses, TerminatedProcesses uint32
	}
	if err := windows.QueryInformationJobObject(p.job, windows.JobObjectBasicAccountingInformation, uintptr(unsafe.Pointer(&counts)), uint32(unsafe.Sizeof(counts)), nil); err != nil {
		return false, fmt.Errorf("inspect owned process tree: %w", err)
	}
	return counts.ActiveProcesses == 0, nil
}
func (p *windowsOwned) release() error {
	return errors.Join(windows.CloseHandle(p.thread), windows.CloseHandle(p.process), windows.CloseHandle(p.job))
}
func isBrokenInput(err error) bool {
	return errors.Is(err, windows.ERROR_BROKEN_PIPE) || errors.Is(err, windows.ERROR_NO_DATA)
}

func windowsExecutable(cmd *exec.Cmd) (string, error) {
	path := cmd.Path
	if cmd.Dir != "" && !filepath.IsAbs(path) {
		dir, err := filepath.Abs(cmd.Dir)
		if err != nil {
			return "", err
		}
		volume := filepath.VolumeName(path)
		switch {
		case volume == "" && (strings.HasPrefix(path, "\\") || strings.HasPrefix(path, "/")):
			path = filepath.VolumeName(dir) + path
		case volume == "":
			path = filepath.Join(dir, path)
		case strings.EqualFold(volume, filepath.VolumeName(dir)):
			path = filepath.Join(dir, path[len(volume):])
		default:
			path, err = filepath.Abs(path)
			if err != nil {
				return "", err
			}
		}
	}
	return exec.LookPath(path)
}

func windowsEnvironment(values []string) ([]uint16, error) {
	byName := map[string]string{}
	for _, value := range values {
		if strings.ContainsRune(value, '\x00') {
			return nil, errors.New("process environment contains NUL")
		}
		index := strings.IndexByte(value, '=')
		if index == 0 {
			if next := strings.IndexByte(value[1:], '='); next >= 0 {
				index = next + 1
			}
		}
		if index <= 0 {
			return nil, errors.New("process environment requires name=value entries")
		}
		byName[strings.ToUpper(value[:index])] = value
	}
	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)
	normalized := make([]string, 0, len(names))
	for _, name := range names {
		normalized = append(normalized, byName[name])
	}
	return utf16.Encode([]rune(strings.Join(normalized, "\x00") + "\x00\x00")), nil
}

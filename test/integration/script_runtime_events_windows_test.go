//go:build windows

package integration

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/sys/windows"

	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
	"github.com/shruggietech/go-schedule/v2/internal/incident"
)

type runtimeStream struct {
	data []byte
	err  error
}

func scriptRuntimeEvents(ctx context.Context, xpath string) incident.RuntimeEventQueryResult {
	directory, err := windows.GetSystemDirectory()
	if err != nil {
		return incident.RuntimeEventQueryResult{Error: errors.New("windows system directory is unavailable"), Unavailable: true}
	}
	return runRuntimeEventQuery(ctx, filepath.Join(directory, "wevtutil.exe"), []string{"qe", "Application", "/q:" + xpath, "/f:xml", "/e:Events", "/c:" + strconv.Itoa(incident.RuntimeEventLimit+1), "/rd:true", "/uni:true"}, directory)
}

func runRuntimeEventQuery(ctx context.Context, executable string, args []string, directory string) incident.RuntimeEventQueryResult {
	child, pipes, err := clifixture.StartPiped(clifixture.Command{Executable: executable, Arguments: args, Directory: directory, Environment: os.Environ(), Limits: &clifixture.ProcessLimits{CPUPercent: 25, MemoryBytes: 128 << 20, MaxProcesses: 1}})
	if err != nil {
		return incident.RuntimeEventQueryResult{Error: errors.New("windows event query could not launch"), Unavailable: true}
	}
	_ = pipes.Stdin.Close()
	out, diagnostics := make(chan runtimeStream, 1), make(chan runtimeStream, 1)
	read := func(file *os.File, ch chan<- runtimeStream) {
		data, err := io.ReadAll(io.LimitReader(file, incident.RuntimeEventByteLimit+1))
		ch <- runtimeStream{data, err}
	}
	go read(pipes.Stdout, out)
	go read(pipes.Stderr, diagnostics)
	waitErr := child.Wait(ctx)
	code, known := child.ExitCode()
	var stdout, stderr runtimeStream
	gotOut, gotErr := false, false
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	var cleanup error
	for !gotOut || !gotErr {
		select {
		case stdout = <-out:
			gotOut = true
			out = nil
		case stderr = <-diagnostics:
			gotErr = true
			diagnostics = nil
		case <-timer.C:
			cleanup = errors.Join(child.Stop(), pipes.Close())
			if !gotOut {
				stdout = <-out
				gotOut = true
			}
			if !gotErr {
				stderr = <-diagnostics
				gotErr = true
			}
		}
	}
	cleanup = errors.Join(cleanup, child.Stop(), pipes.Close())
	result := incident.RuntimeEventQueryResult{Data: stdout.data, Truncated: len(stdout.data) > incident.RuntimeEventByteLimit || len(stderr.data) > incident.RuntimeEventByteLimit}
	if result.Truncated {
		return result
	}
	if ctx.Err() != nil {
		result.Error = ctx.Err()
		return result
	}
	if !known {
		result.Error = errors.New("event-query native result is unavailable")
		result.Unavailable = true
		return result
	}
	if waitErr != nil || code != 0 {
		result.Error = fmt.Errorf("native event-query exit %d", code)
		result.Unavailable = true
		return result
	}
	if cleanup != nil || stdout.err != nil || stderr.err != nil {
		result.Error = errors.New("event query stream collection or owned cleanup failed")
	}
	return result
}

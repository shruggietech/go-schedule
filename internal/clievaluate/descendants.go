package clievaluate

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
)

func (s *session) descendantDrain(ctx context.Context) execution {
	args := []string{"fixture-child", "--mode=tree"}
	process, pipes, err := clifixture.StartPiped(clifixture.Command{Executable: s.binaries["runner"], Arguments: args, Directory: s.w.Runtime, Environment: quietEnvironment()})
	if err != nil {
		return execution{Blocked: "process_unavailable", Detail: err.Error()}
	}
	lines := make(chan string, 1)
	readDone := make(chan error, 1)
	go func() {
		reader := bufio.NewReader(io.LimitReader(pipes.Stdout, 4096))
		line, e := reader.ReadString('\n')
		lines <- line
		_, drainErr := io.Copy(io.Discard, reader)
		readDone <- errors.Join(e, drainErr)
	}()
	var line string
	select {
	case line = <-lines:
	case <-ctx.Done():
	}
	var descendant int
	_, parseErr := fmt.Sscanf(line, "owned descendant %d", &descendant)
	cleanup := process.Stop()
	if cleanup != nil {
		_ = pipes.Close()
	}
	readErr := <-readDone
	code, done := process.ExitCode()
	closeErr := pipes.Close()
	return execution{Code: 0, Args: args, Out: []byte(line), State: mustObservation(map[string]any{"parent_pid": process.PID, "descendant_pid": descendant, "parent_exit_code": code, "parent_exited": done, "tree_drained": cleanup == nil, "pipe_eof": readErr == nil}), Checks: []bool{parseErr == nil && descendant > 0, done && cleanup == nil && readErr == nil && closeErr == nil}, Cleanup: cleanup}
}

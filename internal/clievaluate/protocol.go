package clievaluate

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/internal/clifixture"
)

type teeCloser struct {
	io.Reader
	io.Closer
}

func (s *session) mcpProtocol(ctx context.Context, permission string, confirmation bool) execution {
	args := []string{"--ipc-path", s.w.IPCPath, "mcp", "serve", "--permission=" + permission, "--name=QA protocol"}
	if confirmation {
		args = append(args, "--require-confirmation")
	}
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	process, pipes, err := clifixture.StartPiped(clifixture.Command{Executable: s.binaries["cli"], Arguments: args, Directory: s.w.Runtime, Environment: quietEnvironment()})
	if err != nil {
		return execution{Blocked: "process_unavailable", Detail: err.Error()}
	}
	output, diagnostic := newMemoryCapture(s.budget.ChannelBytes, cancel), newMemoryCapture(s.budget.ChannelBytes, cancel)
	stderrDone := make(chan error, 1)
	go func() { _, err := io.Copy(diagnostic, pipes.Stderr); stderrDone <- err }()
	observations := map[string]any{}
	good := true
	sdk := mcp.NewClient(&mcp.Implementation{Name: "gosched-qa", Version: Version}, nil)
	connection, err := sdk.Connect(childCtx, &mcp.IOTransport{Reader: teeCloser{io.TeeReader(pipes.Stdout, output), pipes.Stdout}, Writer: pipes.Stdin, MaxLineLength: int(s.budget.ChannelBytes)}, nil)
	if err == nil {
		listed, e := connection.ListResources(childCtx, nil)
		err = errors.Join(err, e)
		if e == nil {
			observations["resources"] = listed
			good = good && len(listed.Resources) == 5
			for _, resource := range listed.Resources {
				read, e := connection.ReadResource(childCtx, &mcp.ReadResourceParams{URI: resource.URI})
				if e != nil {
					err = errors.Join(err, e)
					break
				}
				observations[resource.URI] = read
			}
		}
		tools, e := connection.ListTools(childCtx, nil)
		if permission == "observe" {
			good = good && (e != nil || len(tools.Tools) == 0)
		} else {
			err = errors.Join(err, e)
			if e == nil {
				names := []string{}
				for _, tool := range tools.Tools {
					names = append(names, tool.Name)
				}
				observations["tools"] = names
				good = good && len(names) > 0
			}
		}
		if permission == "manage" && err == nil {
			passed, details, e := s.manageProtocol(childCtx, connection)
			good = good && passed
			observations["manage_protocol"] = details
			err = errors.Join(err, e)
		}
		if permission == "operate" && err == nil {
			passed, details, e := s.operateProtocol(childCtx, connection)
			good = good && passed
			observations["operate_protocol"] = details
			err = errors.Join(err, e)
		}

		err = errors.Join(err, connection.Close())
	}
	_ = pipes.Stdin.Close()
	waitCtx, stop := context.WithTimeout(context.Background(), 30*time.Second)
	waitErr := process.Wait(waitCtx)
	stop()
	cleanup := process.Stop()
	if cleanup != nil {
		_ = pipes.Close()
	}
	stderrErr := <-stderrDone
	closeErr := pipes.Close()
	stdout, outOverflow := output.Snapshot()
	stderr, errOverflow := diagnostic.Snapshot()
	code, done := process.ExitCode()
	good = good && done && code == 0 && !outOverflow && !errOverflow
	err = errors.Join(err, waitErr, cleanup, stderrErr, closeErr)
	observations["diagnostic"] = errorText(err)
	return execution{Code: code, Args: args, Out: stdout, Err: stderr, State: mustObservation(observations), Checks: []bool{good, err == nil, good && err == nil}, Cleanup: cleanup}
}

func (s *session) mcpAuthority(ctx context.Context) execution {
	out := []byte{}
	diagnostic := []byte{}
	states := []json.RawMessage{}
	passed := true
	for _, permission := range []string{"observe", "operate", "manage"} {
		r := s.mcpProtocol(ctx, permission, permission == "manage")
		if r.Blocked != "" {
			return r
		}
		out = append(out, r.Out...)
		diagnostic = append(diagnostic, r.Err...)
		states = append(states, r.State)
		for _, check := range r.Checks {
			passed = passed && check
		}
	}
	return execution{Code: 0, Args: []string{"mcp", "serve", "observe|operate|manage"}, Out: out, Err: diagnostic, State: mustObservation(states), Checks: []bool{passed, passed}}
}

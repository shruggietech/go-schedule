package clievaluate

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/shruggietech/go-schedule/internal/clifixture"
)

func (s *session) secretBoundary(ctx context.Context, definition step) execution {
	first := captureCommand(ctx, clifixture.Command{Executable: s.binaries["cli"], Arguments: definition.Spec.Argv, Directory: s.w.Runtime, Environment: quietEnvironment()}, nil, s.budget.ChannelBytes)
	listener, endpoint, err := secretListener(s.w.Runtime)
	if err != nil {
		return execution{Blocked: "owned_channel_unavailable", Detail: err.Error()}
	}
	defer func() { _ = listener.Close() }()
	stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
	defer stop()
	peerDone := make(chan error, 1)
	go func() {
		connection, err := listener.Accept()
		if err != nil {
			peerDone <- err
			return
		}
		defer func() { _ = connection.Close() }()
		if err := connection.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
			peerDone <- err
			return
		}
		var b [1]byte
		n, err := connection.Read(b[:])
		if n != 0 || !errors.Is(err, io.EOF) {
			peerDone <- errors.New("helper did not reject the foreign parent before requesting a credential frame")
			return
		}
		peerDone <- nil
	}()
	args := append(append([]string{}, definition.Spec.Argv...), endpoint)
	second := captureCommand(ctx, clifixture.Command{Executable: s.binaries["cli"], Arguments: args, Directory: s.w.Runtime, Environment: quietEnvironment()}, nil, s.budget.ChannelBytes)
	_ = listener.Close()
	peerErr := <-peerDone
	good := first.Code == 1 && second.Code == 1 && first.Failure == nil && second.Failure == nil && first.Cleanup == nil && second.Cleanup == nil
	silent := len(first.Out)+len(first.Err)+len(second.Out)+len(second.Err) == 0
	return execution{Code: first.Code, Args: definition.Spec.Argv, Out: append(first.Out, second.Out...), Err: append(first.Err, second.Err...), Checks: []bool{good, silent && peerErr == nil}, State: mustObservation(map[string]any{"unframed_entry_exit": first.Code, "foreign_parent_exit": second.Code, "foreign_parent_args": args, "credential_frame_sent": false, "channel_error": errorText(peerErr)}), Cleanup: errors.Join(first.Cleanup, second.Cleanup)}
}

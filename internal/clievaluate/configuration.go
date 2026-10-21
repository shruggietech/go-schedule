package clievaluate

import (
	"context"
	"os"
	"path/filepath"

	"github.com/shruggietech/go-schedule/internal/clifixture"
)

func (s *session) daemonConfiguration(ctx context.Context, definition step) execution {
	config := s.vars["bad_json"]
	if definition.Hook == "occupied-daemon-config" {
		config = s.w.ConfigPath
	}
	before, err := s.api.RuntimeInfo(ctx)
	if err != nil {
		return execution{Blocked: "owned_daemon_unavailable", Detail: err.Error()}
	}
	capture := captureCommand(ctx, clifixture.Command{Executable: s.binaries["daemon"], Arguments: []string{"--config", config}, Directory: s.w.Runtime, Environment: quietEnvironment()}, nil, s.budget.ChannelBytes)
	after, err := s.api.RuntimeInfo(ctx)
	return execution{Code: capture.Code, Args: []string{s.binaries["daemon"], "--config", config}, Out: capture.Out, Err: capture.Err, State: mustObservation(map[string]any{"before": before, "after": after, "capture_error": errorText(capture.Failure)}), Checks: []bool{capture.Code != 0 && len(capture.Err)+len(capture.Out) > 0, err == nil && equalObservation(before, after) && capture.Cleanup == nil}, Cleanup: capture.Cleanup}
}

func (s *session) readonlySetup() error {
	path := filepath.Join(s.w.Runtime, "exports", "readonly.json")
	s.vars["readonly_file"] = path
	if err := s.writeOwned(path, []byte("owned prior contents\n")); err != nil {
		return err
	}
	return os.Chmod(path, 0400)
}

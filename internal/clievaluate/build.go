package clievaluate

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/internal/clifixture"
)

// BuildBinaries prepares six immutable evaluation inputs using a bounded hidden Go tree.
func BuildBinaries(ctx context.Context, source, directory, goExecutable string, budget Budget) error {
	if runtime.GOOS != "windows" {
		return errors.New("prepare the native and cross-compiled inputs on the Windows coordinator host")
	}
	if err := workerPrivilege(); err != nil {
		return err
	}
	if err := ValidateBudget(budget); err != nil {
		return err
	}
	source, err := filepath.Abs(source)
	if err != nil {
		return err
	}
	directory, err = filepath.Abs(directory)
	if err != nil {
		return err
	}
	if goExecutable == "" {
		goExecutable = "go"
	}
	goPath, err := exec.LookPath(goExecutable)
	if err != nil {
		return errors.New("go is unavailable; use the repository's existing toolchain before preparing evaluation inputs")
	}
	if err := os.Mkdir(directory, 0700); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(budget.TotalTimeoutMS)*time.Millisecond)
	defer cancel()
	for _, platform := range []string{"windows", "linux"} {
		output := filepath.Join(directory, platform)
		if err := os.Mkdir(output, 0700); err != nil {
			return err
		}
		overrides := map[string]string{"GOOS": platform, "GOARCH": "amd64", "CGO_ENABLED": "0", "GOTOOLCHAIN": "local", "GOPROXY": "off", "GOMAXPROCS": "2", "GOFLAGS": "-p=2", "GOMEMLIMIT": "512MiB"}
		environment := []string{}
		for _, entry := range quietEnvironment() {
			name, _, ok := strings.Cut(entry, "=")
			if !ok {
				continue
			}
			replaced := false
			for key := range overrides {
				if strings.EqualFold(name, key) {
					replaced = true
				}
			}
			if !replaced {
				environment = append(environment, entry)
			}
		}
		for key, value := range overrides {
			environment = append(environment, key+"="+value)
		}
		capture := captureCommand(ctx, clifixture.Command{Executable: goPath, Arguments: []string{"build", "-p=2", "-o", output + string(filepath.Separator), "./cmd/gosched", "./cmd/goschedd", "./cmd/gosched-qa"}, Directory: source, Environment: environment, Limits: &clifixture.ProcessLimits{CPUPercent: budget.CPUPercent, MemoryBytes: uint64(budget.MemoryBytes), MaxProcesses: uint32(budget.WindowsProcesses)}}, nil, budget.ChannelBytes)
		if capture.Code != 0 || capture.Failure != nil || capture.Cleanup != nil {
			cause := errors.Join(capture.Failure, capture.Cleanup)
			if cause == nil {
				cause = errors.New("go returned a nonzero build status")
			}
			return fmt.Errorf("%s evaluation build did not complete: exit=%d diagnostic=%s: %w", platform, capture.Code, capture.Err, cause)
		}
	}
	return nil
}

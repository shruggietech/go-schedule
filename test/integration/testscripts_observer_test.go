package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/shruggietech/go-schedule/internal/clifixture"
)

const scriptStreamLimit = 1 << 20

type scriptStream struct {
	text string
	err  error
}

// captureScriptStream retains bounded text while draining the owned pipe.
func captureScriptStream(reader io.Reader) scriptStream {
	data, err := io.ReadAll(io.LimitReader(reader, scriptStreamLimit+1))
	if len(data) > scriptStreamLimit {
		data = data[:scriptStreamLimit]
		_, drainErr := io.Copy(io.Discard, reader)
		err = errors.Join(errors.New("script output exceeded one-MiB stream allowance"), err, drainErr)
	}
	return scriptStream{string(data), err}
}

func observeScript(args []string, harness string, environment []string) bindingResult {
	started := time.Now()
	result := bindingResult{code: -1}
	executable, err := exec.LookPath(args[0])
	if err != nil {
		result.err = err
		return result
	}
	executable, err = filepath.Abs(executable)
	if err != nil {
		result.err = err
		return result
	}
	directory, err := os.Getwd()
	if err != nil {
		result.err = err
		return result
	}
	diagnostic := os.Getenv("GOSCHEDULE_PROBE_DIAGNOSTICS") == "1"
	provenance := map[string]any{}
	if diagnostic {
		provenance["executable"] = executable
		provenance["executable_sha256"] = scriptFileHash(executable)
		provenance["harness"] = harness
		provenance["harness_sha256"] = scriptFileHash(harness)
	}
	process, pipes, err := clifixture.StartPiped(clifixture.Command{Executable: executable, Arguments: args[1:], Environment: environment, Directory: directory})
	if err != nil {
		result.err = err
		return result
	}
	_ = pipes.Stdin.Close()
	stdout, stderr := make(chan scriptStream, 1), make(chan scriptStream, 1)
	go func() { stdout <- captureScriptStream(pipes.Stdout) }()
	go func() { stderr <- captureScriptStream(pipes.Stderr) }()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	result.err = process.Wait(ctx)
	cancel()
	if code, known := process.ExitCode(); known {
		result.code = code
	}
	timer := time.NewTimer(2 * time.Second)
	var out, diagnostics scriptStream
	var gotOut, gotErr bool
	for !gotOut || !gotErr {
		select {
		case out = <-stdout:
			gotOut = true
			stdout = nil
		case diagnostics = <-stderr:
			gotErr = true
			stderr = nil
		case <-timer.C:
			result.collectionErr = errors.New("script root exited but inherited output did not close within two seconds")
			result.cleanupErr = errors.Join(result.cleanupErr, process.Stop(), pipes.Close())
			if !gotOut {
				out = <-stdout
				gotOut = true
			}
			if !gotErr {
				diagnostics = <-stderr
				gotErr = true
			}
		}
	}
	timer.Stop()
	result.cleanupErr = errors.Join(result.cleanupErr, process.Stop(), pipes.Close())
	result.collectionErr = errors.Join(result.collectionErr, out.err, diagnostics.err)
	result.err = errors.Join(result.err, result.collectionErr, result.cleanupErr)
	result.output = out.text + diagnostics.text
	result.duration = time.Since(started)
	if diagnostic {
		provenance["pid"] = process.PID
		provenance["native_exit"] = result.code
		provenance["elapsed_ms"] = result.duration.Milliseconds()
		provenance["collection_error"] = errorText(result.collectionErr)
		provenance["cleanup_error"] = errorText(result.cleanupErr)
		data, encodeErr := json.Marshal(provenance)
		if encodeErr != nil {
			result.err = errors.Join(result.err, encodeErr)
		} else {
			result.provenance = string(data)
		}
	}
	return result
}

func scriptFileHash(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Sprintf("unavailable: %v", err)
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Sprintf("unavailable: %v", err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func errorText(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

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
	"strconv"
	"strings"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
	"github.com/shruggietech/go-schedule/v2/internal/incident"
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

func observeScript(args []string, harness string, environment []string) (result bindingResult) {
	started := time.Now()
	result = bindingResult{code: -1}
	if len(args) == 0 {
		result.err = errors.New("script interpreter command is missing")
		return unavailableScriptLaunch(result, "")
	}
	executable, err := exec.LookPath(args[0])
	if err != nil {
		result.err = err
		return unavailableScriptLaunch(result, args[0])
	}
	executable, err = filepath.Abs(executable)
	if err != nil {
		result.err = err
		return unavailableScriptLaunch(result, args[0])
	}
	directory, err := os.Getwd()
	if err != nil {
		result.err = err
		return unavailableScriptLaunch(result, executable)
	}
	diagnostic := os.Getenv("GOSCHEDULE_PROBE_DIAGNOSTICS") == "1"
	provenance := map[string]any{}
	if diagnostic {
		provenance["executable"] = executable
		provenance["executable_sha256"] = scriptFileHash(executable)
		provenance["harness"] = harness
		provenance["harness_sha256"] = scriptFileHash(harness)
	}
	var incidentCapture *incident.Capture
	var incidentErr error
	var incidentPID int
	var incidentHash string
	if parent := os.Getenv("GOSCHEDULE_INCIDENT_DIR"); parent != "" {
		incidentHash = scriptFileHash(executable)
		incidentCapture, incidentErr = incident.New(parent)
		if incidentErr == nil && (strings.EqualFold(filepath.Base(executable), "pwsh.exe") || strings.EqualFold(filepath.Base(executable), "pwsh")) {
			environment, incidentErr = incidentCapture.CoreCLREnvironment(environment)
		}
		defer func() {
			record := incident.Record{Kind: "go-schedule.incident", Version: 1, Category: "script", Phase: "command", Target: incident.Target{PID: incidentPID, Executable: executable, SHA256: incidentHash}, Original: incident.Original{NativeExit: result.code, Canceled: errors.Is(result.err, context.DeadlineExceeded)}, Collection: incident.Collection{Status: "unavailable", DurationMS: time.Since(started).Milliseconds()}}
			if incidentPID == 0 {
				record.Phase = "launch"
			}
			if incidentErr != nil {
				record.Collection.Detail = "Incident directory could not be created; the original script outcome is unchanged."
			} else {
				record.ID = incidentCapture.ID
				if incidentPID == 0 {
					record.Collection.Detail = "The interpreter did not start. Inspect the original launch error before another fixture attempt; no runtime crash evidence is available."
				} else {
					dump := incidentCapture.Inspect("crash."+strconv.Itoa(incidentPID)+".dmp", 32<<20)
					record.Artifacts = []incident.Artifact{dump, incidentCapture.Inspect("runtime.json", 64<<10), incidentCapture.Inspect("createdump.log", 64<<10)}
					retentionFailed := false
					for _, artifact := range record.Artifacts {
						if artifact.Status == "oversized" {
							limit := int64(64 << 10)
							if strings.HasSuffix(artifact.Name, ".dmp") {
								limit = 32 << 20
							}
							if err := incidentCapture.DiscardOversized(artifact.Name, limit); err != nil {
								retentionFailed = true
								record.Collection.Detail = "An oversized owned artifact could not be removed; inspect the private incident directory."
							}
						}
					}
					record.Collection.Status = dump.Status
					if dump.Status == "absent" {
						record.Collection.Detail = "No native crash dump was produced. An ordinary exit, early crash or unavailable runtime collector remains distinguishable from a collected dump."
					}
					if retentionFailed {
						record.Collection.Status = "failed"
						record.Collection.Detail = "An oversized owned artifact could not be removed. Inspect the private incident directory; the original script outcome is unchanged."
					}
				}
				if err := incidentCapture.Save(record); err != nil {
					record.Collection.Status = "failed"
					record.Collection.Detail = "Incident metadata could not be retained; the original script outcome is unchanged."
				}
			}
			data, _ := json.Marshal(record)
			result.incident = string(data)
		}()
	}
	process, pipes, err := clifixture.StartPiped(clifixture.Command{Executable: executable, Arguments: args[1:], Environment: environment, Directory: directory})
	if err != nil {
		result.err = err
		return result
	}
	incidentPID = process.PID
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

func unavailableScriptLaunch(result bindingResult, executable string) bindingResult {
	if os.Getenv("GOSCHEDULE_INCIDENT_DIR") == "" {
		return result
	}
	record := incident.Record{Kind: "go-schedule.incident", Version: 1, Category: "script", Phase: "launch", Target: incident.Target{Executable: executable}, Original: incident.Original{NativeExit: result.code}, Collection: incident.Collection{Status: "unavailable", Detail: "The interpreter did not start. Inspect the original launch error before another fixture attempt; no runtime crash evidence is available."}}
	data, _ := json.Marshal(record)
	result.incident = string(data)
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

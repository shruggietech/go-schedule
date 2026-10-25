package integration

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/clifixture"
)

func TestPackagedCandidatePreparationIsOncePerTarget(t *testing.T) {
	var cache packageBuildCache
	var calls atomic.Int32
	expected := []byte("exact immutable candidate bytes")
	prepare := func() ([]byte, error) { calls.Add(1); return expected, nil }
	var workers sync.WaitGroup
	for i := 0; i < 16; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			got, err := cache.load("owned-root/goschedd", prepare)
			if err != nil || !bytes.Equal(got, expected) {
				t.Errorf("candidate bytes=%q error=%v", got, err)
			}
		}()
	}
	workers.Wait()
	if calls.Load() != 1 {
		t.Fatalf("equivalent preparation invoked %d times", calls.Load())
	}
	if _, err := cache.load("owned-root/gosched", prepare); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 {
		t.Fatalf("distinct target preparation count %d, want 2", calls.Load())
	}
}

func TestPackagedCandidatePreparationRetainsFailureWithoutRetry(t *testing.T) {
	var cache packageBuildCache
	failure := errors.New("original preparation failure")
	calls := 0
	prepare := func() ([]byte, error) { calls++; return nil, failure }
	for i := 0; i < 2; i++ {
		if _, err := cache.load("owned-root/goschedd", prepare); !errors.Is(err, failure) {
			t.Fatalf("original failure lost: %v", err)
		}
	}
	if calls != 1 {
		t.Fatalf("failed preparation retried %d times", calls)
	}
}

// Each test process prepares a target once. Bytes remain immutable, while every
// runtime fixture receives its own package-shaped copy and independent state.
var packagedCandidates packageBuildCache

type packageBuildResult struct {
	once sync.Once
	data []byte
	err  error
}
type packageBuildCache struct {
	mu      sync.Mutex
	targets map[string]*packageBuildResult
}

func (c *packageBuildCache) load(key string, prepare func() ([]byte, error)) ([]byte, error) {
	c.mu.Lock()
	if c.targets == nil {
		c.targets = map[string]*packageBuildResult{}
	}
	result := c.targets[key]
	if result == nil {
		result = &packageBuildResult{}
		c.targets[key] = result
	}
	c.mu.Unlock()
	result.once.Do(func() { result.data, result.err = prepare() })
	return result.data, result.err
}

func preparePackagedCandidate(t *testing.T, directory, name, target string) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	data, err := packagedCandidates.load(root+"\x00"+target, func() ([]byte, error) {
		goBinary, err := exec.LookPath("go")
		if err != nil {
			return nil, err
		}
		goBinary, err = filepath.Abs(goBinary)
		if err != nil {
			return nil, err
		}
		staging := t.TempDir()
		binary := filepath.Join(staging, "candidate")
		if runtime.GOOS == "windows" {
			binary += ".exe"
		}
		stdout, stderr := filepath.Join(staging, "build-stdout.txt"), filepath.Join(staging, "build-stderr.txt")
		process, err := clifixture.StartProcess(clifixture.Command{Executable: goBinary, Arguments: []string{"build", "-o", binary, "./cmd/" + target}, Directory: root, Environment: append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GCM_INTERACTIVE=never"), StdoutPath: stdout, StderrPath: stderr})
		if err != nil {
			return nil, fmt.Errorf("start owned %s build: %w", target, err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		waitErr := process.Wait(ctx)
		stopErr := process.Stop()
		if err := errors.Join(waitErr, stopErr); err != nil {
			output, outErr := os.ReadFile(stdout)
			diagnostic, diagnosticErr := os.ReadFile(stderr)
			return nil, fmt.Errorf("prepare %s within ten-minute build bound: %w; stdout=%q stderr=%q", target, errors.Join(err, outErr, diagnosticErr), output, diagnostic)
		}
		return os.ReadFile(binary)
	})
	if err != nil {
		t.Fatalf("prepare package-shaped %s: %v", target, err)
	}
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, data, 0755); err != nil {
		t.Fatal(err)
	}
	return path
}

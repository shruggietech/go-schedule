package clievaluate

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/clifixture"
)

func TestMemoryCaptureBoundsOutputAndCancelsOnce(t *testing.T) {
	calls := 0
	c := newMemoryCapture(16, func() { calls++ })
	for i := 0; i < 5; i++ {
		n, err := c.Write(bytes.Repeat([]byte("x"), 10))
		if n != 10 || err != nil {
			t.Fatal("drain contract changed")
		}
	}
	data, exceeded := c.Snapshot()
	if len(data) != 16 || !exceeded || calls != 1 {
		t.Fatalf("bytes=%d exceeded=%v cancellations=%d", len(data), exceeded, calls)
	}
}

func TestCaptureProcessHelper(t *testing.T) {
	mode := os.Args[len(os.Args)-1]
	switch mode {
	case "qa-output-overflow":
		_, _ = os.Stdout.Write(bytes.Repeat([]byte("x"), 1<<20))
		os.Exit(0)
	case "qa-input-liveness":
		fmt.Fprintln(os.Stdout, "started")
		_, _ = io.Copy(io.Discard, os.Stdin)
		fmt.Fprintln(os.Stdout, "graceful EOF cleanup")
		os.Exit(0)
	case "qa-delayed-cleanup":
		_, _ = io.Copy(io.Discard, os.Stdin)
		time.Sleep(150 * time.Millisecond)
		fmt.Fprintln(os.Stdout, "owned delayed cleanup completed")
		os.Exit(0)
	case "qa-deadline":
		time.Sleep(time.Minute)
		os.Exit(0)
	}
}

func TestContainedCaptureBoundsRealChildOutputAndJoinsCancellation(t *testing.T) {
	for _, mode := range []string{"qa-output-overflow", "qa-input-liveness", "qa-deadline"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			command := clifixture.Command{Executable: os.Args[0], Arguments: []string{"-test.run=^TestCaptureProcessHelper$", "--", mode}, Directory: t.TempDir()}
			var observed commandCapture
			if mode == "qa-input-liveness" {
				observed = captureManaged(ctx, command, nil, 64, true)
			} else {
				observed = captureCommand(ctx, command, nil, 64)
			}
			if observed.Cleanup != nil || len(observed.Out) > 64 {
				t.Fatalf("capture cleanup=%v bytes=%d", observed.Cleanup, len(observed.Out))
			}
			if mode == "qa-output-overflow" {
				if !observed.Exceeded {
					t.Fatal("real child output was not bounded")
				}
			} else {
				if observed.Failure == nil || observed.Code == 0 && mode != "qa-input-liveness" {
					t.Fatal("canceled child reported success")
				}
				if mode == "qa-input-liveness" && !strings.Contains(string(observed.Out), "started") {
					t.Fatal("interactive child never attached")
				}
				if mode == "qa-input-liveness" && (observed.Code != 0 || !strings.Contains(string(observed.Out), "graceful EOF cleanup")) {
					t.Fatal("owned worker was killed before its bounded EOF cleanup")
				}
			}
		})
	}
}

func TestWorkerSpecificGraceAllowsCleanupAndStillForcesItsBound(t *testing.T) {
	for _, grace := range []time.Duration{50 * time.Millisecond, 2500 * time.Millisecond} {
		t.Run(grace.String(), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			command := clifixture.Command{Executable: os.Args[0], Arguments: []string{"-test.run=^TestCaptureProcessHelper$", "--", "qa-delayed-cleanup"}, Directory: t.TempDir()}
			result := captureWithGrace(ctx, command, nil, 1024, true, grace, nil)
			if result.Cleanup != nil || result.Failure == nil {
				t.Fatalf("cleanup=%v failure=%v", result.Cleanup, result.Failure)
			}
			completed := strings.Contains(string(result.Out), "owned delayed cleanup completed")
			if grace > time.Second {
				if result.Code != 0 || !completed {
					t.Fatalf("worker cleanup was cut short: %+v", result)
				}
			} else if result.Code == 0 {
				t.Fatal("short grace did not stop delayed worker")
			}
		})
	}
}

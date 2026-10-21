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

	"github.com/shruggietech/go-schedule/internal/clifixture"
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

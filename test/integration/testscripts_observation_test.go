package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// The exact owned test child exits while its descendant retains both streams.
func TestScriptObservationHelper(t *testing.T) {
	var args []string
	for i, arg := range os.Args {
		if arg == "--" {
			args = os.Args[i+1:]
			break
		}
	}
	if len(args) != 4 {
		return
	}
	role, code, release, ready := args[0], args[1], args[2], args[3]
	if role == "overflow" {
		fmt.Print(strings.Repeat("x", scriptStreamLimit+100))
		os.Exit(37)
	}
	if role == "descendant" {
		if err := os.WriteFile(ready, []byte("ready"), 0600); err != nil {
			os.Exit(91)
		}
		for i := 0; i < 12000; i++ {
			if _, err := os.Stat(release); err == nil {
				os.Exit(0)
			}
			time.Sleep(10 * time.Millisecond)
		}
		os.Exit(92)
	}
	child := scriptCommand(os.Args[0], "-test.run=^TestScriptObservationHelper$", "--", "descendant", code, release, ready)
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr
	if err := child.Start(); err != nil {
		os.Exit(93)
	}
	for i := 0; i < 500; i++ {
		if _, err := os.Stat(ready); err == nil {
			n, _ := strconv.Atoi(code)
			os.Exit(n)
		}
		time.Sleep(time.Millisecond)
	}
	os.Exit(94)
}

func TestScriptObservationBoundsInheritedStreams(t *testing.T) {
	for _, code := range []int{0, 37} {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			dir := t.TempDir()
			release := filepath.Join(dir, "release")
			ready := filepath.Join(dir, "ready")
			tw := twin{name: "owned-helper", command: func(string, ...string) []string {
				return []string{os.Args[0], "-test.run=^TestScriptObservationHelper$", "--", "parent", strconv.Itoa(code), release, ready}
			}}
			done := make(chan bindingResult, 1)
			go func() { done <- executeBinding(tw, os.Args[0], "", "") }()
			var result bindingResult
			select {
			case result = <-done:
			case <-time.After(15 * time.Second):
				t.Error("root exit was held indefinitely by descendant output handles")
				_ = os.WriteFile(release, []byte("release"), 0600)
				result = <-done
			}
			_ = os.WriteFile(release, []byte("release"), 0600)
			if result.collectionErr == nil {
				t.Error("inherited-stream collection failure was lost")
			}
			if result.cleanupErr != nil {
				t.Errorf("owned cleanup: %v", result.cleanupErr)
			}
			if result.code != code {
				t.Errorf("native exit=%d, want %d; %v", result.code, code, result.err)
			}
		})
	}
}

func TestScriptObservationSeparatesOutputLimitFromNativeExit(t *testing.T) {
	args := []string{os.Args[0], "-test.run=^TestScriptObservationHelper$", "--", "overflow", "37", "unused", "unused"}
	result := observeScript(args, os.Args[0], os.Environ())
	if result.code != 37 || result.collectionErr == nil || len(result.output) != scriptStreamLimit || result.cleanupErr != nil {
		t.Fatalf("native=%d bytes=%d collection=%v cleanup=%v", result.code, len(result.output), result.collectionErr, result.cleanupErr)
	}
}

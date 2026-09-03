//go:build windows

package commandline

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/shruggietech/go-schedule/internal/platform"
)

func TestNativeWindowsArgumentVector(t *testing.T) {
	if os.Getenv("GO_SCHEDULE_NATIVE_ARG_HELPER") == "windows" {
		writeNativeArgs(t)
		return
	}
	assertNativeArgs(t, "windows", "^TestNativeWindowsArgumentVector$")
}

func writeNativeArgs(t *testing.T) {
	t.Helper()
	for i, arg := range os.Args {
		if arg == "--" {
			_ = json.NewEncoder(os.Stdout).Encode(os.Args[i+1:])
			return
		}
	}
	t.Fatal("helper separator missing")
}

func assertNativeArgs(t *testing.T, host, testPattern string) {
	t.Helper()
	invocation, err := Parse(`helper '' "hello world" "héllo 世界" "tabs\there" "lines
here" C:\trail\ '$HOME' '%PATH%' '|' '>' '*.txt' ';' '&'`)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(os.Args[0], append([]string{"-test.run=" + testPattern, "--"}, invocation.Args...)...)
	cmd.Env = append(os.Environ(), "GO_SCHEDULE_NATIVE_ARG_HELPER="+host)
	platform.HideConsole(cmd)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("native helper: %v", err)
	}
	var got []string
	if err := json.NewDecoder(bytes.NewReader(out)).Decode(&got); err != nil {
		t.Fatalf("decode %q: %v", out, err)
	}
	if !reflect.DeepEqual(got, invocation.Args) {
		t.Fatalf("received = %#v, want %#v", got, invocation.Args)
	}
}

package integration

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Inject samples into a disposable script copy, never the host clock.
func TestHeartbeatReversedClockPreservesRecord(t *testing.T) {
	sqlite := requireSqlite(t)
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			name := "Test-Heartbeat.sh"
			if tw.name == "powershell" {
				name = "Test-Heartbeat.ps1"
			}
			original := scriptPath(t, name)
			data, err := os.ReadFile(original)
			if err != nil {
				t.Fatal(err)
			}
			source := string(data)
			if tw.name == "powershell" {
				source = strings.ReplaceAll(source, "$startedMs = Get-EpochMs", "$startedMs = 1700000000100")
				source = strings.ReplaceAll(source, "-FinishedMs (Get-EpochMs)", "-FinishedMs 1700000000000")
			} else {
				source = strings.ReplaceAll(source, "STARTED=\"$(now_ms)\"", "STARTED=1700000000100")
				source = strings.ReplaceAll(source, "write_beat 1 \"$STARTED\" \"$(now_ms)\"", "write_beat 1 \"$STARTED\" 1700000000000")
			}
			if source == string(data) {
				t.Fatal("clock samples were not injected")
			}
			copies := t.TempDir()
			if err := os.Mkdir(filepath.Join(copies, "lib"), 0700); err != nil {
				t.Fatal(err)
			}
			for _, helper := range []string{"sqlite.sh", "Sqlite.ps1"} {
				raw, err := os.ReadFile(filepath.Join(filepath.Dir(original), "lib", helper))
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(copies, "lib", helper), raw, 0600); err != nil {
					t.Fatal(err)
				}
			}
			script := filepath.Join(copies, name)
			if err := os.WriteFile(script, []byte(source), 0600); err != nil {
				t.Fatal(err)
			}
			for _, code := range []int{0, 3} {
				t.Run(strconv.Itoa(code), func(t *testing.T) {
					dir := t.TempDir()
					args := []string{tw.flag("sqlite-exe"), sqlite, tw.flag("quiet")}
					if code != 0 {
						args = append(args, tw.flag("fail-with"), strconv.Itoa(code))
					}
					out, got := runScript(t, tw, dir, script, args...)
					if got != code {
						t.Fatalf("exit=%d want=%d: %s", got, code, out)
					}
					for _, text := range []string{"1700000000000", "1700000000100", "zero duration", "clock"} {
						if !strings.Contains(out, text) {
							t.Fatalf("missing timing diagnostic %q: %s", text, out)
						}
					}
					db := filepath.Join(dir, "heartbeat.db")
					if n := queryInt(t, sqlite, db, "SELECT COUNT(*) FROM beat WHERE started_ms=1700000000100 AND finished_ms=started_ms AND duration_ms=0;"); n != 1 {
						t.Fatalf("preserved records=%d", n)
					}
					if got := queryInt(t, sqlite, db, "SELECT exit_code FROM beat;"); got != code {
						t.Fatalf("stored exit=%d want=%d", got, code)
					}
				})
			}
		})
	}
}

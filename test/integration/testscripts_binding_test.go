package integration

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// bindingHarness exercises the real adapters, including the sqlite CLI's own
// dot-command parser. An SQL driver would bypass the defective boundary.
func bindingHarness(t *testing.T, tw twin, sql string) string {
	t.Helper()
	dir := t.TempDir()
	var text, name string
	if tw.name == "powershell" {
		name = "binding.ps1"
		library := strings.ReplaceAll(scriptPath(t, "lib/Sqlite.ps1"), "'", "''")
		text = "param([string]$Database)\n$ErrorActionPreference='Stop'\n. '" + library + "'\nInitialize-TestSqlite\n$value=[string][Environment]::GetEnvironmentVariable('SCRIPT_BIND_VALUE')\nInvoke-Sqlite -Database $Database -Sql '" + strings.ReplaceAll(sql, "'", "''") + "' -Parameters @{ value=$value; empty=''; absent=$null }\n"
	} else {
		name = "binding.sh"
		library := strings.ReplaceAll(bashCompatiblePath(scriptPath(t, "lib/sqlite.sh")), "'", "'\"'\"'")
		text = "#!/usr/bin/env bash\nset -euo pipefail\n. '" + library + "'\ninit_test_sqlite '' 0\nsqlite_exec \"$1\" list '" + strings.ReplaceAll(sql, "'", "'\"'\"'") + "' \"value=$SCRIPT_BIND_VALUE\" 'empty=' 'absent=__NULL__'\n"
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

type bindingResult struct {
	output   string
	code     int
	err      error
	duration time.Duration
}

func executeBinding(tw twin, harness, database, value string) bindingResult {
	args := tw.command(harness, database)
	cmd := scriptCommand(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), "SCRIPT_BIND_VALUE="+value)
	started := time.Now()
	out, err := cmd.CombinedOutput()
	result := bindingResult{output: string(out), err: err, duration: time.Since(started)}
	if err != nil {
		result.code = -1
		if e, ok := err.(interface{ ExitCode() int }); ok {
			result.code = e.ExitCode()
		}
	}
	return result
}

func TestScriptsBindingsPreserveLiteralText(t *testing.T) {
	sqlite := requireSqlite(t)
	values := []string{"", "0000123", "1e999999999999999999999999999999", "2e999999999999999999999999999999", "1+2", "O'Brien", "line one\n.quit\nline three\n", "tab\tand \"quotes\" \\ path", "zażółć 東京 👋", "0); DROP TABLE literal_values; --"}
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			sql := "CREATE TABLE IF NOT EXISTS literal_values (value TEXT, empty_value TEXT, absent_value TEXT); INSERT INTO literal_values VALUES(:value,:empty,:absent);"
			harness := bindingHarness(t, tw, sql)
			for i, value := range values {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					database := filepath.Join(t.TempDir(), "literal.db")
					result := executeBinding(tw, harness, database, value)
					if result.code != 0 {
						t.Fatalf("literal write exited %d (%v):\n%s", result.code, result.err, result.output)
					}
					got := queryScalar(t, sqlite, database, "SELECT hex(value)||'|'||typeof(value)||'|'||typeof(empty_value)||'|'||length(empty_value)||'|'||typeof(absent_value) FROM literal_values;")
					want := strings.ToUpper(hex.EncodeToString([]byte(value))) + "|text|text|0|null"
					if got != want {
						t.Errorf("literal parameter changed: got %q, want %q", got, want)
					}
				})
			}
		})
	}
}

func TestScriptsBindingsKeepConcurrentSessionIdentities(t *testing.T) {
	sqlite := requireSqlite(t)
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			var runtimeArgs []string
			if tw.name == "powershell" {
				runtimeArgs = []string{"pwsh", "-NoProfile", "-NonInteractive", "-Command", "[Console]::WriteLine($PSVersionTable.PSVersion.ToString() + ' | ' + [System.Runtime.InteropServices.RuntimeInformation]::FrameworkDescription + ' | ' + [System.Runtime.InteropServices.RuntimeInformation]::OSDescription + ' | ' + [System.Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture.ToString())"}
			} else {
				runtimeArgs = []string{"bash", "-c", "printf '%s\\n' \"$BASH_VERSION\"; uname -a"}
			}
			identity, err := scriptCommand(runtimeArgs[0], runtimeArgs[1:]...).CombinedOutput()
			if err != nil {
				t.Fatalf("interpreter identity: %v %s", err, identity)
			}
			t.Logf("interpreter identity: %s", strings.TrimSpace(string(identity)))
			for _, initialized := range []bool{false, true} {
				t.Run(fmt.Sprintf("initialized=%t", initialized), func(t *testing.T) {
					database := filepath.Join(t.TempDir(), "sessions.db")
					if initialized {
						if out, err := scriptCommand(sqlite, database, "PRAGMA journal_mode=WAL; CREATE TABLE sessions(session_id TEXT NOT NULL, sequence INTEGER NOT NULL, UNIQUE(session_id,sequence));").CombinedOutput(); err != nil {
							t.Fatalf("seed schema: %v\n%s", err, out)
						}
					}
					harness := bindingHarness(t, tw, "CREATE TABLE IF NOT EXISTS sessions(session_id TEXT NOT NULL, sequence INTEGER NOT NULL, UNIQUE(session_id,sequence)); INSERT INTO sessions VALUES(:value,1);")
					const writers = 6
					results := make([]bindingResult, writers)
					var wg sync.WaitGroup
					for i := 0; i < writers; i++ {
						wg.Add(1)
						go func(i int) {
							defer wg.Done()
							// Distinct valid 32-hex IDs previously collapsed to Inf.
							id := strconv.Itoa(i+1) + "e" + strings.Repeat("9", 30)
							results[i] = executeBinding(tw, harness, database, id)
						}(i)
					}
					wg.Wait()
					for i, result := range results {
						if result.code != 0 {
							t.Errorf("writer %d exited %d (0x%08x, %v), elapsed=%s, phase=script invocation, host=%s/%s %s, interpreter=%s:\n%s", i, result.code, uint32(result.code), result.err, result.duration, runtime.GOOS, runtime.GOARCH, runtime.Version(), tw.command(harness, database)[0], result.output)
						}
					}
					if n := queryInt(t, sqlite, database, "SELECT COUNT(DISTINCT session_id) FROM sessions;"); n != writers {
						t.Errorf("distinct sessions = %d, want %d", n, writers)
					}
					if n := queryInt(t, sqlite, database, "SELECT COUNT(*) FROM sessions;"); n != writers {
						t.Errorf("persisted records = %d, want %d", n, writers)
					}
				})
			}
		})
	}
}

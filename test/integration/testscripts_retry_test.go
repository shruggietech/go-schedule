package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// retryHarness replaces only the sqlite process and sleeping. The real helper
// still chooses retry versus failure and emits its own runtime diagnostics.
func retryHarness(t *testing.T, tw twin) string {
	t.Helper()
	var name, text string
	if tw.name == "powershell" {
		name = "retry.ps1"
		library := strings.ReplaceAll(scriptPath(t, "lib/Sqlite.ps1"), "'", "''")
		text = fmt.Sprintf("param([string]$State,[string]$Waits,[string]$Scenario)\n. '%s'\n", library) + `
function Invoke-RetrySqlite {
 process { }
 end {
  $count = [int][IO.File]::ReadAllText($State) + 1
  [IO.File]::WriteAllText($State, [string]$count)
  if ($Scenario -eq 'transient' -and $count -eq 3) {
   $global:LASTEXITCODE = 0
   Write-Output 'fixture write succeeded'
   return
  }
  $kind = if ($Scenario -eq 'transient') {
   if ($count -eq 1) { 'locked' } else { 'busy' }
  } elseif ($Scenario -eq 'other') { 'readonly' } else { $Scenario }
  $global:LASTEXITCODE = 1
  Write-Output ("fixture database {0}: original writer diagnostic {1}" -f $kind, $count)
 }
}
function Start-Sleep {
 param([int]$Milliseconds)
 [IO.File]::AppendAllText($Waits, ([string]$Milliseconds + [Environment]::NewLine))
}
$script:ResolvedSqliteExe = 'Invoke-RetrySqlite'
Invoke-Sqlite -Database 'disposable.db' -Sql 'SELECT 1;'
`
	} else {
		name = "retry.sh"
		library := strings.ReplaceAll(bashCompatiblePath(scriptPath(t, "lib/sqlite.sh")), "'", "'\"'\"'")
		text = fmt.Sprintf("#!/usr/bin/env bash\nset -euo pipefail\n. '%s'\n", library) + `
state="$1"; waits="$2"; scenario="$3"
fake_sqlite() {
 while IFS= read -r line; do :; done
 local count kind
 read -r count < "$state"
 count=$((count + 1))
 printf '%s\n' "$count" > "$state"
 if [ "$scenario" = transient ] && [ "$count" -eq 3 ]; then
  printf '%s\n' 'fixture write succeeded'
  return 0
 fi
 kind="$scenario"
 if [ "$scenario" = transient ]; then
  if [ "$count" -eq 1 ]; then kind=locked; else kind=busy; fi
 elif [ "$scenario" = other ]; then
  kind=readonly
 fi
 printf 'fixture database %s: original writer diagnostic %s\n' "$kind" "$count" >&2
 return 1
}
sleep() { printf '%s\n' "$1" >> "$waits"; }
SQLITE_EXE=fake_sqlite
sqlite_exec disposable.db list 'SELECT 1;'
`
	}
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestScriptsSqliteRetryContract(t *testing.T) {
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			harness := retryHarness(t, tw)
			for _, scenario := range []string{"transient", "locked", "busy", "other"} {
				t.Run(scenario, func(t *testing.T) {
					dir := t.TempDir()
					state, waits := filepath.Join(dir, "attempts"), filepath.Join(dir, "waits")
					if err := os.WriteFile(state, []byte("0\n"), 0o600); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(waits, nil, 0o600); err != nil {
						t.Fatal(err)
					}
					argv := tw.command(harness, state, waits, scenario)
					output, err := scriptCommand(argv[0], argv[1:]...).CombinedOutput()
					code := 0
					if err != nil {
						if e, ok := err.(interface{ ExitCode() int }); ok {
							code = e.ExitCode()
						} else {
							t.Fatalf("run retry harness: %v\n%s", err, output)
						}
					}
					wantCode, wantAttempts, wantWarnings := 1, 3, 3
					wantWaits := []string{"250", "500", "750"}
					if tw.name == "posix" {
						wantWaits = []string{"0.25", "0.50", "0.75"}
					}
					switch scenario {
					case "transient":
						wantCode, wantWarnings = 0, 2
						wantWaits = wantWaits[:2]
						if !strings.Contains(string(output), "fixture write succeeded") {
							t.Errorf("successful result lost:\n%s", output)
						}
					case "other":
						wantAttempts, wantWarnings, wantWaits = 1, 0, nil
					}
					if code != wantCode {
						t.Errorf("exit = %d, want %d (%v)\n%s", code, wantCode, err, output)
					}
					data, readErr := os.ReadFile(state)
					if readErr != nil {
						t.Fatal(readErr)
					}
					if got := strings.TrimSpace(string(data)); got != strconv.Itoa(wantAttempts) {
						t.Errorf("attempts = %q, want %d\n%s", got, wantAttempts, output)
					}
					data, readErr = os.ReadFile(waits)
					if readErr != nil {
						t.Fatal(readErr)
					}
					if got := strings.Fields(string(data)); strings.Join(got, ",") != strings.Join(wantWaits, ",") {
						t.Errorf("wait schedule = %q, want %q", got, wantWaits)
					}
					if got := strings.Count(string(output), "Database contention on attempt"); got != wantWarnings {
						t.Errorf("contention warnings = %d, want %d\n%s", got, wantWarnings, output)
					}
					if scenario != "transient" {
						kind := scenario
						if scenario == "other" {
							kind = "readonly"
						}
						diagnostic := fmt.Sprintf("fixture database %s: original writer diagnostic %d", kind, wantAttempts)
						if !strings.Contains(string(output), diagnostic) {
							t.Errorf("original failure diagnostic discarded, want %q:\n%s", diagnostic, output)
						}
						if scenario != "other" && !strings.Contains(string(output), "The record was NOT written.") {
							t.Errorf("exhaustion did not identify unwritten record:\n%s", output)
						}
						if strings.Contains(string(output), "fixture write succeeded") {
							t.Errorf("failed operation advertised success:\n%s", output)
						}
					}
				})
			}
		})
	}
}

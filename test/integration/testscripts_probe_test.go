package integration

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestScriptsSQLiteCandidateFailureDiagnostics(t *testing.T) {
	goTool, err := exec.LookPath("go")
	if err != nil {
		t.Skip("Go is required to build the owned native sqlite probe fixture")
	}
	dir := t.TempDir()
	source := filepath.Join(dir, "main.go")
	candidate := filepath.Join(dir, "sqlite-probe")
	if runtime.GOOS == "windows" {
		candidate += ".exe"
	}
	program := "package main\nimport(\"fmt\";\"os\";\"time\";\"strings\")\nfunc main(){switch os.Getenv(\"QA_SQLITE_PROBE\"){case \"exit\":fmt.Fprintln(os.Stderr,\"original candidate failure\");os.Exit(42);case \"empty\":return;case \"malformed\":fmt.Println(\"not-a-version\");case \"old\":fmt.Println(\"3.32.0\");case \"minimum\":fmt.Println(\"3.33.0\");case \"new\":fmt.Println(\"3.47.2\");case \"oversized\":fmt.Print(strings.Repeat(\"x\",8192));case \"timeout\":time.Sleep(time.Minute)}}\n"
	if err := os.WriteFile(source, []byte(program), 0600); err != nil {
		t.Fatal(err)
	}
	if output, err := scriptCommand(goTool, "build", "-o", candidate, source).CombinedOutput(); err != nil {
		t.Fatalf("build owned fixture: %v\n%s", err, output)
	}
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			var text, name string
			if tw.name == "powershell" {
				name = "probe.ps1"
				text = fmt.Sprintf("param([string]$Candidate)\n. '%s'\nTest-SqliteCandidate -Candidate $Candidate\n", strings.ReplaceAll(scriptPath(t, "lib/Sqlite.ps1"), "'", "''"))
			} else {
				name = "probe.sh"
				text = fmt.Sprintf("#!/usr/bin/env bash\n. '%s'\ntest_sqlite_candidate \"$1\"\n", strings.ReplaceAll(bashCompatiblePath(scriptPath(t, "lib/sqlite.sh")), "'", "'\"'\"'"))
			}
			harness := filepath.Join(t.TempDir(), name)
			if err := os.WriteFile(harness, []byte(text), 0600); err != nil {
				t.Fatal(err)
			}
			t.Run("selection", func(t *testing.T) { testSQLiteCandidateSelection(t, tw, candidate) })
			for _, tc := range []struct{ mode, version, diagnostic string }{{"exit", "", "exited 42"}, {"empty", "", "empty or malformed"}, {"malformed", "", "empty or malformed"}, {"old", "", "below the required"}, {"minimum", "3.33.0", ""}, {"new", "3.47.2", ""}, {"oversized", "", "exceeded the 4096"}, {"timeout", "", "timed out after 10000 ms"}, {"missing", "", ""}} {
				t.Run(tc.mode, func(t *testing.T) {
					selected := candidate
					if tc.mode == "missing" {
						selected += "-absent"
					}
					args := tw.command(harness, selected)
					cmd := scriptCommand(args[0], args[1:]...)
					cmd.Env = append(os.Environ(), "QA_SQLITE_PROBE="+tc.mode)
					var stderr bytes.Buffer
					cmd.Stderr = &stderr
					stdout, err := cmd.Output()
					if err != nil {
						t.Fatalf("probe harness: %v\n%s", err, stderr.String())
					}
					if strings.TrimSpace(string(stdout)) != tc.version {
						t.Fatalf("version=%q want=%q stderr=%s", stdout, tc.version, stderr.String())
					}
					if tc.diagnostic != "" && !strings.Contains(stderr.String(), tc.diagnostic) {
						t.Fatalf("candidate failure lost: %s", stderr.String())
					}
					if tc.mode == "exit" && !strings.Contains(stderr.String(), "original candidate failure") {
						t.Fatal("original stderr discarded")
					}
				})
			}
		})
	}
}

func testSQLiteCandidateSelection(t *testing.T, tw twin, candidate string) {
	t.Helper()
	requireSqlite(t)
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".bin"), 0700); err != nil {
		t.Fatal(err)
	}
	name := "sqlite3"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	stale := filepath.Join(root, ".bin", name)
	data, err := os.ReadFile(candidate)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stale, data, 0700); err != nil {
		t.Fatal(err)
	}
	var text, script string
	if tw.name == "powershell" {
		script = "selection.ps1"
		text = fmt.Sprintf("param([string]$Root,[string]$Candidate,[string]$Mode)\n. '%s'\n$script:TestScriptsRoot=$Root\nif($Mode -eq 'explicit'){Resolve-Sqlite -Explicit $Candidate}else{Resolve-Sqlite}\n", strings.ReplaceAll(scriptPath(t, "lib/Sqlite.ps1"), "'", "''"))
	} else {
		script = "selection.sh"
		text = fmt.Sprintf("#!/usr/bin/env bash\n. '%s'\nTEST_SCRIPTS_ROOT=\"$1\"\nif [ \"$3\" = explicit ]; then resolve_sqlite \"$2\"; else resolve_sqlite ''; fi\n", strings.ReplaceAll(bashCompatiblePath(scriptPath(t, "lib/sqlite.sh")), "'", "'\"'\"'"))
	}
	harness := filepath.Join(t.TempDir(), script)
	if err := os.WriteFile(harness, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"implicit", "explicit"} {
		t.Run(mode, func(t *testing.T) {
			args := tw.command(harness, root, stale, mode)
			cmd := scriptCommand(args[0], args[1:]...)
			cmd.Env = append(os.Environ(), "QA_SQLITE_PROBE=old")
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			stdout, err := cmd.Output()
			if mode == "explicit" {
				if err == nil || len(bytes.TrimSpace(stdout)) != 0 {
					t.Fatalf("explicit old candidate fell back: %q %v", stdout, err)
				}
				if e, ok := err.(interface{ ExitCode() int }); !ok || e.ExitCode() != 2 {
					t.Fatalf("explicit refusal exit: %v %s", err, stderr.String())
				}
			} else {
				if err != nil || len(bytes.TrimSpace(stdout)) == 0 || strings.Contains(string(stdout), filepath.Base(root)) {
					t.Fatalf("implicit candidate did not fall through: %q %v %s", stdout, err, stderr.String())
				}
			}
		})
	}
}

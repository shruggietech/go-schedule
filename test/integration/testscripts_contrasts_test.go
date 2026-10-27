package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScriptsRuntimeAndSelectionContrasts(t *testing.T) {
	sqlite := requireSqlite(t)
	t.Setenv("GOSCHEDULE_PROBE_DIAGNOSTICS", "1")
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			for _, operation := range []string{"metadata", "writer"} {
				t.Run(operation, func(t *testing.T) {
					original := bindingHarness(t, tw, "CREATE TABLE identities(value TEXT); INSERT INTO identities VALUES(:value);")
					data, err := os.ReadFile(original)
					if err != nil {
						t.Fatal(err)
					}
					text := string(data)
					if operation == "metadata" {
						if tw.name == "powershell" {
							text = text[:strings.Index(text, "$value=")] + "[Console]::WriteLine($PSVersionTable.PSVersion.ToString() + ' | ' + [System.Runtime.InteropServices.RuntimeInformation]::FrameworkDescription)\n"
						} else {
							at := strings.Index(text, "phase=write-started")
							at = strings.LastIndex(text[:at], "if [")
							text = text[:at] + "printf '%s\\n' \"$BASH_VERSION\"\n"
						}
					}
					var selected string
					for _, selection := range []string{"implicit", "explicit"} {
						t.Run(selection, func(t *testing.T) {
							harness := filepath.Join(t.TempDir(), filepath.Base(original))
							body := text
							if selection == "explicit" {
								if selected == "" {
									t.Fatal("implicit candidate identity missing")
								}
								if tw.name == "powershell" {
									body = strings.Replace(body, "Initialize-TestSqlite\n", "Initialize-TestSqlite -SqliteExe '"+strings.ReplaceAll(selected, "'", "''")+"'\n", 1)
								} else {
									body = strings.Replace(body, "init_test_sqlite '' 0", "init_test_sqlite '"+strings.ReplaceAll(selected, "'", "'\"'\"'")+"' 0", 1)
								}
							}
							if err := os.WriteFile(harness, []byte(body), 0600); err != nil {
								t.Fatal(err)
							}
							database := filepath.Join(t.TempDir(), "contrast.db")
							result := executeBinding(tw, harness, database, "owned literal")
							logScriptReplayObservation(t, operation+"/"+selection, -1, result)
							t.Logf("condition=%s/%s/%s outer=%s output=%q", tw.name, operation, selection, result.provenance, result.output)
							if result.code != 0 || result.err != nil {
								t.Fatalf("contrast failed: %+v", result)
							}
							const prefix = "phase=prerequisite-resolved sqlite="
							at := strings.Index(result.output, prefix)
							if at < 0 {
								t.Fatal("missing candidate phase")
							}
							candidate := strings.Split(strings.TrimSpace(result.output[at+len(prefix):]), "\n")[0]
							candidate = strings.TrimSpace(candidate)
							if selection == "implicit" {
								selected = candidate
							} else if candidate != selected {
								t.Fatalf("selection changed candidate: %q != %q", candidate, selected)
							}
							if operation == "writer" {
								if got := queryScalar(t, sqlite, database, "SELECT value FROM identities;"); got != "owned literal" {
									t.Fatalf("writer identity=%q", got)
								}
							} else if _, err := os.Stat(database); !os.IsNotExist(err) {
								t.Fatalf("metadata-only experiment created database: %v", err)
							}
						})
					}
				})
			}
		})
	}
}

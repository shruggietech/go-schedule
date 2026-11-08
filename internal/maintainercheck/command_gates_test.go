package maintainercheck

import (
	"strings"
	"testing"
)

func TestNativeCommandCorpusIntegration(t *testing.T) {
	const workflow = `jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    steps:
      - name: Maintainer command contracts
        shell: bash
        run: sh test/scripts/maintainer-commands_test.sh
`
	for _, tc := range []struct {
		name, text string
		want       bool
	}{
		{"all native platforms", workflow, false},
		{"missing corpus", strings.Replace(workflow, "sh test/scripts/maintainer-commands_test.sh", "echo skipped", 1), true},
		{"wrong job", strings.Replace(workflow, "  test:", "  other:", 1), true},
		{"missing Windows", strings.Replace(workflow, ", windows-latest", "", 1), true},
		{"wrong runner", strings.Replace(workflow, "runs-on: ${{ matrix.os }}", "runs-on: ubuntu-latest", 1), true},
		{"excluded platform", strings.Replace(workflow, "    steps:", "        exclude: [{os: windows-latest}]\n    steps:", 1), true},
		{"conditional job", strings.Replace(workflow, "  test:", "  test:\n    if: false", 1), true},
		{"conditional step", strings.Replace(workflow, "        shell: bash", "        if: false\n        shell: bash", 1), true},
		{"soft failure", strings.Replace(workflow, "        shell: bash", "        continue-on-error: true\n        shell: bash", 1), true},
		{"job soft failure", strings.Replace(workflow, "  test:", "  test:\n    continue-on-error: true", 1), true},
		{"wrong shell", strings.Replace(workflow, "shell: bash", "shell: pwsh", 1), true},
		{"wrong directory", strings.Replace(workflow, "        shell: bash", "        working-directory: desktop\n        shell: bash", 1), true},
		{"suppressed command", strings.Replace(workflow, "sh test/scripts/maintainer-commands_test.sh", "sh test/scripts/maintainer-commands_test.sh || true", 1), true},
		{"comment decoy", strings.Replace(workflow, "run: sh test/scripts/maintainer-commands_test.sh", "run: echo skipped # sh test/scripts/maintainer-commands_test.sh", 1), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			put(t, root, "ci.yml", tc.text)
			c := newChecker(root)
			c.ci("ci.yml", c.yaml("ci.yml"), "v2.17.1")
			got := strings.Contains(strings.Join(c.findings, "\n"), "native maintainer command corpus")
			if got != tc.want {
				t.Fatalf("want integration finding %v, got %v: %v", tc.want, got, c.findings)
			}
		})
	}
}

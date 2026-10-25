package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const validWorkflow = "jobs:\n  verify:\n    runs-on: ubuntu-latest\n    timeout-minutes: 360\n    steps:\n      - run: go test -timeout=60m -race ./...\n      - uses: actions/deploy-pages@v5\n        with:\n          timeout: '3600000'\n"

func TestWorkflowCompletionPolicyRejectsReducedAndMissingAllowances(t *testing.T) {
	cases := []struct{ name, source, want string }{
		{"valid", validWorkflow, ""},
		{"missing job cap", strings.ReplaceAll(validWorkflow, "    timeout-minutes: 360\n", ""), "timeout-minutes: 360"},
		{"short job cap", strings.ReplaceAll(validWorkflow, "timeout-minutes: 360", "timeout-minutes: 5"), "timeout-minutes: 360"},
		{"default package cap", strings.ReplaceAll(validWorkflow, " -timeout=60m", ""), "-timeout=60m"},
		{"short package cap", strings.ReplaceAll(validWorkflow, "-timeout=60m", "-timeout=10m"), "-timeout=60m"},
		{"short Pages cap", strings.ReplaceAll(validWorkflow, "3600000", "600000"), "3600000"},
		{"invalid YAML", "jobs: [", "invalid workflow"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, problems := workflowProblems("fixture.yml", []byte(tc.source))
			got := strings.Join(problems, "\n")
			if tc.want == "" && got != "" || tc.want != "" && !strings.Contains(got, tc.want) {
				t.Fatalf("problems=%q want=%q", got, tc.want)
			}
		})
	}
}

func TestCanonicalShellCommandsNeedTheSamePackageAllowance(t *testing.T) {
	for _, source := range []string{`CGO_ENABLED=1 "$GO" test -race ./...`, `go test -count=1 ./...`} {
		if len(goTestProblems("fixture", source)) != 1 {
			t.Fatalf("missing rejection for %s", source)
		}
	}
	if got := goTestProblems("fixture", `CGO_ENABLED=1 "$GO" test -timeout=60m -race ./...`); len(got) != 0 {
		t.Fatal(got)
	}
}

func TestRepositoryCompletionPolicyAndNestedMutation(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("missing source location")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	jobs, problems := checkRoot(root)
	if jobs == 0 || len(problems) > 0 {
		t.Fatalf("jobs=%d problems=%s", jobs, strings.Join(problems, "\n"))
	}
	fixture := t.TempDir()
	paths := map[string]bool{".github/workflows/fixture.yml": true, "scripts/verify.sh": true, "scripts/coverage-gate.sh": true, "Makefile": true}
	for _, required := range nestedValues {
		paths[required.path] = true
	}
	for path := range paths {
		target := filepath.Join(fixture, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
			t.Fatal(err)
		}
		var content []byte
		var err error
		if path == ".github/workflows/fixture.yml" {
			content = []byte(validWorkflow)
		} else {
			content, err = os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
			if err != nil {
				t.Fatal(err)
			}
		}
		if err = os.WriteFile(target, content, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if _, problems = checkRoot(fixture); len(problems) > 0 {
		t.Fatal(problems)
	}
	target := filepath.Join(fixture, "desktop", "frontend", "vite.config.ts")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.ReplaceAll(string(data), "60_000", "5_000"))
	if err = os.WriteFile(target, data, 0600); err != nil {
		t.Fatal(err)
	}
	_, problems = checkRoot(fixture)
	if !strings.Contains(strings.Join(problems, "\n"), "testTimeout: 60000") {
		t.Fatalf("shortened nested timeout accepted: %v", problems)
	}
}

func TestMakefileCommandsNeedTheSamePackageAllowance(t *testing.T) {
	if got := goTestProblems("Makefile", "\t$(GO) test ./...\n\t$(GO) test -bench=. ./internal/engine/...\n"); len(got) != 2 {
		t.Fatalf("unqualified Makefile commands accepted: %v", got)
	}
	if got := goTestProblems("Makefile", "\t$(GO) test -timeout=60m ./...\n"); len(got) != 0 {
		t.Fatal(got)
	}
}

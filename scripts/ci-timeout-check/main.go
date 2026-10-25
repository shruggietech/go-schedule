// Command ci-timeout-check keeps hosted and local verification completion allowances aligned.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

type workflow struct {
	Jobs map[string]struct {
		Timeout int `yaml:"timeout-minutes"`
		Steps   []struct {
			Run  string
			Uses string
			With map[string]string
		}
	}
}

var goTests = regexp.MustCompile(`(?:^|[[:space:];|&])(?:go|"\$GO"|\$\(GO\)) test\b[^\n]*`)

func workflowProblems(name string, data []byte) (int, []string) {
	var document workflow
	if err := yaml.Unmarshal(data, &document); err != nil {
		return 0, []string{fmt.Sprintf("%s: invalid workflow: %v", name, err)}
	}
	var problems []string
	if len(document.Jobs) == 0 {
		problems = append(problems, name+": no jobs found")
	}
	for key, job := range document.Jobs {
		if job.Timeout != 360 {
			problems = append(problems, fmt.Sprintf("%s job %s: declare timeout-minutes: 360 (hosted maximum)", name, key))
		}
		for _, step := range job.Steps {
			problems = append(problems, goTestProblems(name+" job "+key, step.Run)...)
			if strings.HasPrefix(step.Uses, "actions/deploy-pages@") && step.With["timeout"] != "3600000" {
				problems = append(problems, name+": Pages deployment allowance must be 3600000 ms")
			}
		}
	}
	sort.Strings(problems)
	return len(document.Jobs), problems
}

func goTestProblems(name, source string) []string {
	var problems []string
	for _, command := range goTests.FindAllString(source, -1) {
		if !strings.Contains(command, "-timeout=60m") {
			problems = append(problems, fmt.Sprintf("%s: Go verification command requires -timeout=60m: %s", name, strings.TrimSpace(command)))
		}
	}
	return problems
}

type requiredValue struct{ path, text string }

var nestedValues = []requiredValue{
	{"desktop/frontend/vite.config.ts", "testTimeout: 60000"},
	{"desktop/frontend/vite.config.ts", "hookTimeout: 60000"},
	{"desktop/frontend/vite.config.ts", "teardownTimeout: 60000"},
	{"desktop/frontend/src/test-setup.ts", "asyncUtilTimeout: 10000"},
	{"desktop/frontend/playwright.config.ts", "timeout: 120000"},
	{"desktop/frontend/playwright.config.ts", "expect: { timeout: 30000 }"},
	{"desktop/frontend/playwright.config.ts", "timeout: 300000"},
	{"site/playwright.config.ts", "timeout: 120000"},
	{"site/playwright.config.ts", "expect: { timeout: 30000 }"},
	{"site/next.config.mjs", "staticPageGenerationTimeout: 300"},
	{"site/scripts/browser-test.mjs", "Static test server did not become ready')), 60000"},
	{".github/workflows/release.yml", "seq 1 1200"},
}

func checkRoot(root string) (int, []string) {
	files, err := filepath.Glob(filepath.Join(root, ".github", "workflows", "*.yml"))
	if err != nil {
		return 0, []string{err.Error()}
	}
	extra, extraErr := filepath.Glob(filepath.Join(root, ".github", "workflows", "*.yaml"))
	if extraErr != nil {
		return 0, []string{extraErr.Error()}
	}
	files = append(files, extra...)
	var problems []string
	total := 0
	if len(files) == 0 {
		problems = append(problems, "no workflow files found")
	}
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			problems = append(problems, err.Error())
			continue
		}
		count, found := workflowProblems(filepath.ToSlash(path), data)
		total += count
		problems = append(problems, found...)
	}
	for _, name := range []string{"scripts/verify.sh", "scripts/coverage-gate.sh", "Makefile"} {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
		if err != nil {
			problems = append(problems, err.Error())
			continue
		}
		problems = append(problems, goTestProblems(name, string(data))...)
	}
	for _, required := range nestedValues {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(required.path)))
		if err != nil {
			problems = append(problems, err.Error())
			continue
		}
		// Numeric separators do not change TypeScript/JavaScript allowance values.
		if !strings.Contains(strings.ReplaceAll(string(data), "_", ""), required.text) {
			problems = append(problems, fmt.Sprintf("%s: completion policy requires %s", required.path, required.text))
		}
	}
	sort.Strings(problems)
	return total, problems
}

func main() {
	root := "."
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: ci-timeout-check [repository-root]")
		os.Exit(2)
	}
	if len(os.Args) == 2 {
		root = os.Args[1]
	}
	jobs, problems := checkRoot(root)
	if len(problems) > 0 {
		for _, problem := range problems {
			fmt.Fprintln(os.Stderr, problem)
		}
		os.Exit(1)
	}
	fmt.Printf("ci-timeout-check: OK - %d jobs and nested verification allowances aligned\n", jobs)
}

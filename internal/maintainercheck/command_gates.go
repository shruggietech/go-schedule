package maintainercheck

import (
	"strings"

	"go.yaml.in/yaml/v3"
)

// nativeCommandCorpus scopes the shared command contract to the actual native
// test job. Comments, other jobs and suppressed or conditional steps do not count.
func (c *checker) nativeCommandCorpus(file string, n *yaml.Node) {
	job := field(n, "jobs", "test")
	matrix := field(job, "strategy", "matrix")
	platforms := items(field(matrix, "os"))
	valid := len(platforms) == 3
	for _, os := range []string{"ubuntu-latest", "macos-latest", "windows-latest"} {
		valid = valid && containsScalar(field(matrix, "os"), os)
	}
	valid = valid && field(matrix, "exclude") == nil && field(matrix, "include") == nil
	valid = valid && value(field(job, "runs-on")) == "${{ matrix.os }}"
	valid = valid && field(job, "if") == nil && !softFailure(job)
	for _, parent := range []*yaml.Node{n, job} {
		dir := value(field(parent, "defaults", "run", "working-directory"))
		valid = valid && (dir == "" || dir == ".")
	}
	found := false
	for _, step := range items(field(job, "steps")) {
		if strings.TrimSpace(value(field(step, "run"))) != "sh test/scripts/maintainer-commands_test.sh" {
			continue
		}
		dir := value(field(step, "working-directory"))
		if value(field(step, "shell")) == "bash" && field(step, "if") == nil && !softFailure(step) && (dir == "" || dir == ".") {
			found = true
		}
	}
	c.require(file, valid && found, "native maintainer command corpus must run unconditionally under Bash in jobs.test on Ubuntu, macOS and Windows")
}

func softFailure(n *yaml.Node) bool {
	setting := field(n, "continue-on-error")
	return setting != nil && value(setting) != "false"
}

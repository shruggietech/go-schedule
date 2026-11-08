package maintainercheck

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"
)

func (c *checker) release(file string, n *yaml.Node, version string) {
	if n == nil {
		c.report(file, "release workflow not found or invalid; release source preflight and comparison review cannot be established")
		return
	}
	c.rules(file, n, requirereleasetext, version, "")
	c.nodePins(file, "every release Node runtime pin", n)
	text := semanticText(n)
	for _, step := range steps(n) {
		run := value(field(step, "run"))
		if gitCommand.MatchString(run) {
			c.report(file, "release workflow must not invoke git")
		}
		if value(field(step, "uses")) == "actions/checkout@v7" {
			ref := value(field(step, "with", "ref"))
			if ref == "main" || ref == "refs/heads/main" {
				c.report(file, "release validation must inspect the tagged checkout")
			}
		}
	}
	c.require(file, joined(field(n, "jobs", "release-state", "needs")) == "readme-version, ci-success" && joined(field(n, "jobs", "binaries", "needs")) == "release-state", "binaries job missing README and exact-commit CI preflight dependencies through release-state")
	source := field(n, "jobs", "readme-version")
	setup, command, comparison := -1, -1, -1
	count := 0
	bypass := false
	expected := `go run ./scripts/release-source-check -root . -repository "$GITHUB_REPOSITORY" -tag "$GITHUB_REF_NAME"`
	for index, step := range items(field(source, "steps")) {
		if strings.HasPrefix(value(field(step, "uses")), "actions/setup-go@") {
			setup = index
		}
		run := strings.TrimSpace(value(field(step, "run")))
		if strings.Contains(run, "scripts/release-source-check") {
			count++
			if run == expected {
				command = index
			}
		}
		if strings.Contains(run, "scripts/comparison-check") {
			comparison = index
		}
		bypass = bypass || has(step, "if") || has(step, "continue-on-error")
	}
	c.require(file, count == 1 && setup >= 0 && command > setup && comparison > command && !bypass, "release source preflight must run once after Go setup and before comparison without bypass, using exact root/repository/tag arguments")
	uploads := 0
	for _, step := range steps(n) {
		if value(field(step, "uses")) != "softprops/action-gh-release@v3" {
			continue
		}
		uploads++
		with := field(step, "with")
		if value(field(with, "draft")) != "true" {
			c.report(file, "every staged release upload must explicitly use draft: true")
		}
		if value(field(with, "generate_release_notes")) == "true" {
			c.report(file, "generated release notes must remain disabled")
		}
	}
	c.require(file, uploads > 0, "every staged release upload must explicitly use draft: true")
	if strings.Contains(text, "--draft=false") || strings.Contains(text, "SHA256SUMS.txt") {
		c.report(file, "staging must not publish or create final checksums")
	}
	if retiredRelease.MatchString(text) {
		c.report(file, "retired Fyne release path remains")
	}
}
func (c *checker) promotion(file string, n *yaml.Node) {
	if n == nil {
		c.report(file, "release promotion workflow not found or invalid")
		return
	}
	c.rules(file, n, requirepromotiontext, "", "")
	text := semanticText(n)
	for _, step := range steps(n) {
		if rebuild.MatchString(value(field(step, "run"))) {
			c.report(file, "promotion must not rebuild release artifacts")
		}
	}
	gate, checksum, promote := strings.Index(text, "windows-release-gate verify-bundle"), strings.Index(text, "LC_ALL=C sha256sum *"), strings.Index(text, "-F draft=false")
	c.require(file, gate >= 0 && checksum > gate && promote > checksum, "gate, final checksums, and promotion are not ordered safely")
}

var comparisonCommand = regexp.MustCompile(`(?m)^\s*go run \./scripts/comparison-check -tag [^\n]+`)
var mutationCommand = regexp.MustCompile(`gh release upload|draft=false|release-promotion[^\n]*\s(ensure-checksum|publish)(\s|$)`)

// comparison follows actual needs edges as well as command order inside each job.
// A review in an unrelated concurrent job cannot establish a publication prerequisite.
func (c *checker) comparison(file string, n *yaml.Node) {
	jobs := unalias(field(n, "jobs"))
	if jobs == nil {
		c.report(file, "comparison review must run before release mutation without bypass")
		return
	}
	type jobPolicy struct {
		gate, mutation bool
		needs          []string
	}
	policies := map[string]jobPolicy{}
	binding := `-commit "$TAG_COMMIT"`
	if strings.HasSuffix(file, "/release.yml") {
		binding = `-commit "$GITHUB_SHA"`
	}
	anyMutation := false
	for i := 0; i < len(jobs.Content); i += 2 {
		name, job := jobs.Content[i].Value, jobs.Content[i+1]
		p := jobPolicy{}
		need := unalias(field(job, "needs"))
		if need != nil {
			if need.Kind == yaml.SequenceNode {
				for _, v := range need.Content {
					p.needs = append(p.needs, v.Value)
				}
			} else {
				p.needs = []string{need.Value}
			}
		}
		for _, step := range items(field(job, "steps")) {
			run := value(field(step, "run"))
			matches := comparisonCommand.FindStringIndex(run)
			mutation := mutationCommand.FindStringIndex(run)
			if value(field(step, "uses")) == "softprops/action-gh-release@v3" {
				mutation = []int{0, 0}
			}
			if matches != nil && !p.mutation && !has(job, "if") && !has(job, "continue-on-error") && !has(step, "if") && !has(step, "continue-on-error") {
				command := run[matches[0]:matches[1]]
				if strings.Contains(command, binding) && !strings.Contains(run, "||") && (mutation == nil || matches[0] < mutation[0]) {
					p.gate = true
				}
			}
			if mutation != nil {
				p.mutation = true
				anyMutation = true
			}
		}
		policies[name] = p
	}
	var protected func(string, map[string]bool) bool
	protected = func(name string, visiting map[string]bool) bool {
		if visiting[name] {
			return false
		}
		visiting[name] = true
		defer delete(visiting, name)
		p, ok := policies[name]
		if !ok {
			return false
		}
		if p.gate {
			return true
		}
		for _, dependency := range p.needs {
			if protected(dependency, visiting) {
				return true
			}
		}
		return false
	}
	for name, p := range policies {
		if p.mutation && !protected(name, map[string]bool{}) {
			c.report(file, "comparison review must run before release mutation without bypass and bind the exact source (job "+name+")")
		}
	}
	if !anyMutation {
		c.report(file, "comparison review requires a release mutation entry point")
	}
}

var noteTag = regexp.MustCompile(`^v[0-9].*\.[0-9].*\.[0-9].*\.md$`)
var noteLink = regexp.MustCompile(`^Read the \[full changelog\]\(.+\) for every change\.$`)

func (c *checker) releaseNotes() {
	dir := ".github/release-notes"
	paths, err := filepath.Glob(filepath.Join(c.root, dir, "*.md"))
	if err != nil {
		c.report(dir, err.Error())
		return
	}
	found := 0
	for _, path := range paths {
		base := filepath.Base(path)
		if !noteTag.MatchString(base) {
			continue
		}
		found++
		file := filepath.ToSlash(filepath.Join(dir, base))
		text := c.read(file)
		lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
		bullets, headings := 0, 0
		for _, line := range lines {
			if strings.HasPrefix(line, "- ") {
				bullets++
			}
			if strings.HasPrefix(line, "#") {
				headings++
			}
		}
		c.require(file, strings.Count(text, "## Highlights\n") == 1, "expected exactly one Highlights heading")
		c.require(file, bullets >= 4 && bullets <= 6, "expected four to six highlight bullets")
		valid := len(lines) == bullets+4 && len(lines) >= 8 && lines[0] == "## Highlights" && lines[1] == "" && headings == 1
		if valid {
			for _, line := range lines[2 : 2+bullets] {
				valid = valid && strings.HasPrefix(line, "- ")
			}
			valid = valid && lines[2+bullets] == "" && noteLink.MatchString(lines[3+bullets])
		}
		c.require(file, valid, "expected highlights-only release copy")
		prefix := "https://github.com/shruggietech/go-schedule/blob/" + strings.TrimSuffix(base, ".md") + "/CHANGELOG.md#"
		c.require(file, strings.Count(text, "[full changelog](") == 1 && strings.Count(text, prefix) == 1, "expected exactly one tagged full changelog link")
	}
	if found == 0 {
		c.report(dir, fmt.Sprintf("no tag-specific release notes found (%d files)", len(paths)))
	}
}

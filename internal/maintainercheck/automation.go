package maintainercheck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"
)

var approvedActions = map[string]string{
	"actions/checkout": "v7", "actions/setup-go": "v7", "actions/setup-node": "v7",
	"actions/configure-pages": "v6", "actions/upload-pages-artifact": "v5", "actions/deploy-pages": "v5",
	"actions/upload-artifact": "v7", "github/codeql-action/init": "v4", "github/codeql-action/analyze": "v4", "softprops/action-gh-release": "v3",
}
var wailsRuntime = regexp.MustCompile(`(?m)^\s*(?:require\s+)?github\.com/wailsapp/wails/v2\s+(v2\.[0-9]+\.[0-9]+)\s*$`)
var wailsCLI = regexp.MustCompile(`github\.com/wailsapp/wails/v2/cmd/wails@[^\s]+`)
var retiredCI = regexp.MustCompile(`wails-proof:|wails-browser-contract:|Install Fyne|GUI build & test \(cgo\)`)
var retiredRelease = regexp.MustCompile(`Install Fyne|goversioninfo|\./cmd/gosched-gui`)
var gitCommand = regexp.MustCompile(`(^|[^[:alnum:]_])git([^[:alnum:]_]|$)`)
var rebuild = regexp.MustCompile(`(^|\s)(go|wix)\s+build(\s|$)`)
var secrets = regexp.MustCompile(`\$\{\{\s*secrets(\s|\.|\[)`)

const gateManifest = "format\nvet\nlint\nrace\ngui\ncoverage\ndocs\nautomation"

// Automation evaluates repository automation and its bounded existing helper commands.
func Automation(root string) []string {
	c := newChecker(root)
	c.automation()
	return c.findings
}
func exists(path string) bool { _, err := os.Stat(path); return err == nil }
func (c *checker) actions(file string, n *yaml.Node) {
	check := func(ref *yaml.Node) {
		if ref == nil {
			return
		}
		reference := value(ref)
		family, version, ok := strings.Cut(reference, "@")
		expected, known := approvedActions[family]
		switch {
		case !known:
			c.report(file, fmt.Sprintf("line %d: unaudited action reference %s", ref.Line, reference))
		case !ok || version != expected:
			c.report(file, fmt.Sprintf("line %d: action reference %s is not approved; expected %s@%s", ref.Line, reference, family, expected))
		}
	}
	jobs := unalias(field(n, "jobs"))
	if jobs == nil {
		return
	}
	for i := 1; i < len(jobs.Content); i += 2 {
		job := jobs.Content[i]
		check(field(job, "uses"))
		for _, step := range items(field(job, "steps")) {
			check(field(step, "uses"))
		}
	}
}
func (c *checker) rules(file string, n *yaml.Node, rules []textRule, version, prefix string) {
	text := semanticText(n)
	for _, rule := range rules {
		required := strings.TrimSpace(strings.ReplaceAll(rule.text, "$WAILS_VERSION", version))
		if !strings.Contains(text, required) {
			c.report(file, "missing "+prefix+rule.description)
		}
	}
}

func (c *checker) automation() {
	workflowDir := filepath.Join(c.root, ".github", "workflows")
	entries, err := os.ReadDir(workflowDir)
	if err != nil {
		c.report(".github/workflows", "workflow directory not found: "+err.Error())
		return
	}
	workflows := map[string]*yaml.Node{}
	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())
		if entry.IsDir() || (ext != ".yml" && ext != ".yaml") {
			continue
		}
		file := ".github/workflows/" + entry.Name()
		workflows[file] = c.yaml(file)
		c.actions(file, workflows[file])
	}
	if len(workflows) == 0 {
		c.report(".github/workflows", "no .yml or .yaml workflow files found")
	}
	ci, release, promote := ".github/workflows/ci.yml", ".github/workflows/release.yml", ".github/workflows/promote-release.yml"
	version := ""
	versions := wailsRuntime.FindAllStringSubmatch(c.read("desktop/go.mod"), -1)
	if len(versions) != 1 {
		c.report("desktop/go.mod", "missing or ambiguous Wails runtime version")
	} else {
		version = versions[0][1]
	}
	for _, file := range []string{ci, release, "scripts/verify.sh"} {
		text := semanticText(workflows[file])
		if file == "scripts/verify.sh" {
			text = c.read(file)
		}
		pins := wailsCLI.FindAllString(text, -1)
		if strings.Contains(text, "./scripts/desktop-build") {
			managed := wailsCLI.FindAllString(c.read("scripts/desktop-build/main.go"), -1)
			for _, pin := range managed {
				pins = append(pins, strings.TrimRight(pin, "\"',"))
			}
		}
		if len(pins) == 0 {
			c.report(file, "missing explicit Wails CLI pin")
		}
		for _, pin := range pins {
			if pin != "github.com/wailsapp/wails/v2/cmd/wails@"+version {
				c.report(file, "Wails CLI pin "+pin+" does not match runtime "+version)
			}
		}
	}
	makeVersion := ""
	makeAssignments := 0
	for _, line := range strings.Split(c.read("Makefile"), "\n") {
		words := strings.Fields(line)
		if len(words) >= 2 && words[0] == "WAILS_VERSION" && words[1] == "=" {
			makeAssignments++
			makeVersion = ""
			if len(words) >= 3 {
				makeVersion = words[2]
			}
		}
	}

	c.require("Makefile", makeAssignments == 1 && makeVersion == version && version != "", "Wails build version does not match runtime "+version)
	c.ci(ci, workflows[ci], version)
	c.codeql(workflows[".github/workflows/codeql.yml"])
	c.dependabot()
	if !exists(filepath.Join(c.root, "scripts/spec-lifecycle-check.sh")) {
		c.report("scripts/spec-lifecycle-check.sh", "specification lifecycle checker not found")
	}
	c.lifecycle()
	c.release(release, workflows[release], version)
	c.promotion(promote, workflows[promote])
	c.releaseNotes()
	c.need("specs/048-v100-release-cut/contracts/publication.md", "reviewed S049 merge commit", "S049 tag-boundary identity")
	c.need("specs/048-v100-release-cut/contracts/publication.md", "windows-release-gate render-dispositions", "formal disposition generation step")
	c.need("test/windows/README.md", "windows-release-gate render-dispositions", "operator disposition command")
	c.need("test/windows/README.md", "does not update or close", "non-authority boundary")
	if !exists(filepath.Join(c.root, "scripts/brand-check")) {
		c.report("scripts/brand-check", "brand integrity command not found")
	} else {
		goExe := os.Getenv("GO")
		if goExe == "" {
			goExe = "go"
		}
		c.child("scripts/brand-check", goExe, "run", "./scripts/brand-check")
	}
	observed := c.child("scripts/verify.sh", "sh", filepath.Join(c.root, "scripts/verify.sh"), "list")
	if observed != gateManifest {
		c.report("scripts/verify.sh", "gate manifest differs; expected:\n"+gateManifest+"\nobserved:\n"+observed)
	}
	for _, file := range []string{release, promote, ".github/workflows/promote-release-waived.yml"} {
		c.comparison(file, workflows[file])
	}
}
func (c *checker) nodePins(file, description string, n *yaml.Node) {
	for _, step := range steps(n) {
		pin := field(step, "with", "node-version")
		if pin != nil && value(pin) != "26" {
			c.report(file, fmt.Sprintf("line %d: %s must use Node 26; found %s", pin.Line, description, value(pin)))
		}
	}
}
func (c *checker) ci(file string, n *yaml.Node, version string) {
	if n == nil {
		c.report(file, "canonical CI workflow not found or invalid")
		return
	}
	c.rules(file, n, requirecitext, version, "")
	c.nodePins(file, "every CI Node runtime pin", n)
	c.nativeCommandCorpus(file, n)
	c.require(file, field(n, "jobs", "wails-desktop") != nil, "missing production Wails desktop job")
	job := field(n, "jobs", "v13-release-qualification")
	if job == nil {
		c.report(file, "missing named v1.3 release qualification job")
	}
	c.rules(file, job, requirev13text, version, "v13-release-qualification job missing ")
	if retiredCI.MatchString(semanticText(n)) {
		c.report(file, "retired or parallel desktop validation remains")
	}
}
func (c *checker) codeql(n *yaml.Node) {
	file := ".github/workflows/codeql.yml"
	if n == nil {
		c.report(file, "canonical CodeQL workflow not found or invalid")
		return
	}
	c.require(file, value(field(n, "name")) == "CodeQL", "missing stable CodeQL workflow name")
	for _, event := range []string{"push", "pull_request"} {
		c.require(file, joined(field(n, "on", event, "branches")) == "main", "missing "+event+" trigger targeting main")
	}
	schedules := items(field(n, "on", "schedule"))
	scheduleOK := false
	for _, s := range schedules {
		scheduleOK = scheduleOK || value(field(s, "cron")) == "17 4 * * 1"
	}
	c.require(file, scheduleOK, "missing weekly schedule trigger (expected '17 4 * * 1')")
	permissions := unalias(field(n, "permissions"))
	c.require(file, permissions != nil && len(permissions.Content) == 4 && value(field(permissions, "contents")) == "read" && value(field(permissions, "security-events")) == "write", "permissions must be exactly contents: read and security-events: write")
	jobs := unalias(field(n, "jobs"))
	if jobs != nil {
		for i := 1; i < len(jobs.Content); i += 2 {
			if has(jobs.Content[i], "permissions") {
				c.report(file, "job-level permissions override is not allowed")
			}
		}
	}
	for _, rule := range []textRule{{"runs-on: ubuntu-latest", "Ubuntu analysis runner"}, {"go-version-file: go.mod", "module-selected Go version"}, {"languages: go", "CodeQL Go language"}, {"build-mode: manual", "CodeQL manual build mode"}, {"uses: github/codeql-action/init@v4", "CodeQL init step"}, {"run: CGO_ENABLED=0 go build ./...", "cgo-free headless build step"}, {"uses: github/codeql-action/analyze@v4", "CodeQL analyze step"}} {
		c.requireText(file, rule.text, rule.description, n)
	}
	if secrets.MatchString(semanticText(n)) {
		c.report(file, "CodeQL workflow must not consume repository or organization secrets")
	}
}
func (c *checker) dependabot() {
	file := ".github/dependabot.yml"
	doc := c.yaml(file)
	updates := items(field(doc, "updates"))
	c.require(file, len(updates) == 5, "expected root Go, desktop Go, desktop npm, website npm, and GitHub Actions entries")
	counts := map[string]int{}
	byRoot := map[string]*yaml.Node{}
	directories := map[string]bool{}
	for _, u := range updates {
		ecosystem := value(field(u, "package-ecosystem"))
		counts[ecosystem]++
		if ecosystem != "gomod" && ecosystem != "npm" && ecosystem != "github-actions" {
			c.report(file, "unapproved package ecosystem: "+ecosystem)
		}
		dir := value(field(u, "directory"))
		directories[dir] = true
		if dir == "/" {
			byRoot[ecosystem] = u
		}
	}
	for _, ecosystem := range []string{"gomod", "github-actions"} {
		want := 1
		if ecosystem == "gomod" {
			want = 2
		}
		c.require(file, counts[ecosystem] == want, fmt.Sprintf("expected %d %s entry or entries", want, ecosystem))
		u := byRoot[ecosystem]
		c.require(file, u != nil, ecosystem+" root directory")
		c.require(file, value(field(u, "open-pull-requests-limit")) == "5", ecosystem+" five-PR limit")
		cadence := "weekly"
		if ecosystem == "github-actions" {
			cadence = "monthly"
		}
		c.require(file, value(field(u, "schedule", "interval")) == cadence, ecosystem+" "+cadence+" cadence")
		groups := unalias(field(u, "groups"))
		routine := field(groups, "routine-minor-and-patch")
		c.require(file, routine != nil, ecosystem+" routine update group")
		types := field(routine, "update-types")
		c.require(file, containsScalar(types, "minor"), ecosystem+" grouped minor updates")
		c.require(file, containsScalar(types, "patch"), ecosystem+" grouped patch updates")
		labels := field(u, "labels")
		c.require(file, containsScalar(labels, "dependencies"), ecosystem+" dependencies label")
		if groups != nil {
			for i := 1; i < len(groups.Content); i += 2 {
				g := groups.Content[i]
				if containsScalar(field(g, "update-types"), "major") {
					c.report(file, ecosystem+" routine group must not group major updates")
				}
				if value(field(g, "applies-to")) == "security-updates" {
					c.report(file, ecosystem+" routine group must not group security updates")
				}
			}
		}
		if ecosystem == "gomod" {
			for _, group := range []struct{ name, pattern, description string }{{"storage-minor-and-patch", "modernc.org/sqlite", "gomod storage dependency pattern"}, {"platform-minor-and-patch", "github.com/kardianos/service", "gomod service dependency pattern"}, {"platform-minor-and-patch", "golang.org/x/sys", "gomod system dependency pattern"}} {
				g := field(groups, group.name)
				c.require(file, g != nil, "gomod isolated "+group.name+" update group")
				c.require(file, containsScalar(field(g, "patterns"), group.pattern), group.description)
			}
		}
	}
	for _, dir := range []string{"/desktop", "/desktop/frontend", "/site"} {
		c.require(file, directories[dir], "missing production desktop dependency boundary: "+dir)
	}
	text := semanticText(doc)
	c.require(file, strings.Contains(text, "github.com/wailsapp/wails/v2"), "missing production desktop dependency boundary: github.com/wailsapp/wails/v2")
	if strings.Contains(text, "fyne.io/*") {
		c.report(file, "retired Fyne dependency group remains")
	}
}

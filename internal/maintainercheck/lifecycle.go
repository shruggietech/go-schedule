package maintainercheck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var specPath = regexp.MustCompile(`^[0-9]{3}-.+$`)
var openTask = regexp.MustCompile(`(?m)^- \[ \] T[0-9]+`)
var doneTask = regexp.MustCompile(`(?mi)^- \[[xX]\] T[0-9]+`)
var deliveryEvidence = regexp.MustCompile(`(?i)(commit\s|releases?\s|pull/[0-9]+|PR #[0-9]+|review branch\s)`)
var inventoryRow = regexp.MustCompile(`(?m)^\| [0-9]{3}-[^|]+ \|`)

// Lifecycle evaluates all feature metadata under root and reports the discovered count.
func Lifecycle(root string) (int, []string) {
	c := newChecker(root)
	count := c.lifecycle()
	return count, c.findings
}
func metadata(text, key string) string {
	prefix := "**" + key + "**:"
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}
func (c *checker) lifecycle() int {
	inventory := c.read("specs/README.md")
	entries, err := os.ReadDir(filepath.Join(c.root, "specs"))
	if err != nil {
		c.report("specs", fmt.Sprintf("read specification directory: %v", err))
		return 0
	}
	found := 0
	for _, entry := range entries {
		name := entry.Name()
		if !specPath.MatchString(name) {
			continue
		}
		dir := filepath.Join("specs", name)
		info, err := os.Stat(filepath.Join(c.root, dir))
		if err != nil {
			if !os.IsNotExist(err) {
				c.report(dir, fmt.Sprintf("inspect specification directory: %v", err))
			}
			continue
		}
		if !info.IsDir() {
			continue
		}
		file := filepath.ToSlash(filepath.Join(dir, "spec.md"))
		if !c.regularFile(file) {
			continue
		}
		found++
		text := c.read(file)
		state := metadata(text, "Status")
		delivery := metadata(text, "Delivery")
		switch state {
		case "Draft", "Ready", "In Progress", "Implemented", "Deferred", "Superseded", "Abandoned":
		case "":
			c.report(file, "missing status")
		default:
			c.report(file, "invalid status: "+state)
		}
		if !strings.Contains(inventory, "| "+name+" | "+state+" |") {
			c.report(file, "missing inventory row for "+name+" with state "+state)
		}
		tasksFile := filepath.Join("specs", name, "tasks.md")
		tasks := ""
		if exists(filepath.Join(c.root, tasksFile)) {
			tasks = c.read(tasksFile)
		}
		open, checked := len(openTask.FindAllString(tasks, -1)), len(doneTask.FindAllString(tasks, -1))
		switch state {
		case "Draft", "Ready":
			if open == 0 {
				c.report(file, state+" has no actionable task")
			}
		case "In Progress":
			if open == 0 || checked == 0 {
				c.report(file, "In Progress must have completed and actionable tasks")
			}
		case "Implemented":
			if delivery == "" {
				c.report(file, "missing delivery evidence")
			} else if !deliveryEvidence.MatchString(delivery) {
				c.report(file, "invalid delivery evidence: "+delivery)
			}
			if open != 0 {
				c.report(file, fmt.Sprintf("Implemented has %d unresolved task(s)", open))
			}
		case "Deferred", "Superseded", "Abandoned":
			if metadata(text, "Disposition") == "" {
				c.report(file, state+" is missing a disposition")
			}
		}
	}
	if found == 0 {
		c.report("specs", "no feature specifications found")
	}
	count := len(inventoryRow.FindAllString(inventory, -1))
	if count != found {
		c.report("specs/README.md", fmt.Sprintf("inventory has %d feature row(s), discovered %d", count, found))
	}
	return found
}

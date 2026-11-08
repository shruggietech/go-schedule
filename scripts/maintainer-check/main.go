// Command maintainer-check evaluates existing automation or specification lifecycle policy.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shruggietech/go-schedule/v2/internal/maintainercheck"
)

func run(args []string, out, errOut io.Writer) int {
	if len(args) < 1 || len(args) > 2 || (args[0] != "automation" && args[0] != "lifecycle") {
		fmt.Fprintln(errOut, "usage: maintainer-check {automation|lifecycle} [root]")
		return 2
	}
	root := "."
	if len(args) == 2 {
		root = args[1]
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(errOut, "maintainer-check: resolve root: %v\n", err)
		return 1
	}
	var findings []string
	prefix := "automation-check"
	count := 0
	if args[0] == "lifecycle" {
		prefix = "spec-lifecycle-check"
		count, findings = maintainercheck.Lifecycle(absolute)
	} else {
		findings = maintainercheck.Automation(absolute)
	}
	if len(findings) != 0 {
		for _, finding := range findings {
			fmt.Fprintln(errOut, finding)
		}
		fmt.Fprintf(errOut, "%s: FAILED with %d issue(s)\n", prefix, len(findings))
		return 1
	}
	if args[0] == "lifecycle" {
		fmt.Fprintf(out, "%s: OK - %d specification(s) are lifecycle-consistent\n", prefix, count)
	} else {
		fmt.Fprintln(out, "automation-check: OK - actions, CodeQL, Dependabot, release operations, release notes, brand, lifecycle, and 8 gates")
	}
	return 0
}
func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

package clievaluate

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/url"
	"sort"
	"strings"
)

func markdown(value string) string {
	value = html.EscapeString(value)
	value = strings.ReplaceAll(value, "|", "&#124;")
	value = strings.ReplaceAll(value, "`", "&#96;")
	return strings.ReplaceAll(strings.ReplaceAll(value, "\r\n", "\n"), "\n", "<br>")
}
func artifactLink(path string) string {
	parts := strings.Split(path, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return "[" + markdown(path) + "](" + strings.Join(parts, "/") + ")"
}

// Render derives the fixed human outline from canonical data without clock or file I/O.
// External inputs must first pass ValidateReport; capturing checkpoints are internal.
func Render(r Report) ([]byte, error) {
	if r.Summary != DeriveSummary(r) {
		return nil, errors.New("cannot render contradictory summary")
	}
	var out strings.Builder
	fmt.Fprintf(&out, "# CLI QA after-action report\n\n## 1. Run identity and outcome\n\nRun: %s. Outcome: %s. Phase: %s. Completion: %s.\n\nSubject: %s at %s, dirty=%t, version %s. Started: %s. Finished: %s.\n\n", markdown(r.RunID), markdown(r.Summary.Outcome), markdown(r.Phase), markdown(r.CompletionState), markdown(r.Subject.Repository), markdown(r.Subject.Commit), r.Subject.Dirty, markdown(r.Subject.Version), markdown(r.StartedAt), markdown(r.FinishedAt))
	fmt.Fprintf(&out, "## 2. Environment and budget\n\nCPU: %d%% (denominator recorded per environment). Memory: %d bytes. Windows processes: %d. Linux tasks including threads: %d. Per-case deadline: %d ms. Total coordinator deadline: %d ms. Per-channel cap: %d bytes. Artifact cap: %d bytes.\n\n", r.Runner.Budget.CPUPercent, r.Runner.Budget.MemoryBytes, r.Runner.Budget.WindowsProcesses, r.Runner.Budget.LinuxTasks, r.Runner.Budget.CaseTimeoutMS, r.Runner.Budget.TotalTimeoutMS, r.Runner.Budget.ChannelBytes, r.Runner.Budget.ArtifactBytes)
	envs := append([]Environment{}, r.Environments...)
	sort.Slice(envs, func(i, j int) bool { return envs[i].ID < envs[j].ID })
	for _, e := range envs {
		fmt.Fprintf(&out, "- %s: %s/%s, %s (%s).", markdown(e.ID), markdown(e.OS), markdown(e.Arch), markdown(e.Status), markdown(e.ReasonCode))
		for _, v := range e.Limits {
			fmt.Fprintf(&out, " %s=%s;", markdown(v.Name), markdown(v.Value))
		}
		out.WriteString("\n")
	}
	fmt.Fprintf(&out, "\n## 3. Coverage table\n\n| Planned | Passed | Failed | Blocked | Not applicable | Findings | Executed pass | Executed fail |\n| --- | --- | --- | --- | --- | --- | --- | --- |\n| %d | %d | %d | %d | %d | %d | %d | %d |\n\n", r.Summary.Planned, r.Summary.Passed, r.Summary.Failed, r.Summary.Blocked, r.Summary.NotApplicable, r.Summary.Findings, r.Summary.ExecutionPassed, r.Summary.ExecutionFailed)
	out.WriteString("## 4. Functional results\n\n| Boundary | Planned | Executed | Execution passed | Execution failed | Effective failed |\n| --- | --- | --- | --- | --- | --- |\n")
	for _, kind := range []string{"packaged_cli", "command_harness", "deterministic_engine"} {
		planned, executed, executionPass, executionFail, effectiveFail := 0, 0, 0, 0, 0
		for _, result := range r.Cases {
			if result.ExecutionKind != kind {
				continue
			}
			planned++
			if result.Observed.ExitCode != nil {
				executed++
			}
			if result.ExecutionStatus == "pass" {
				executionPass++
			}
			if result.ExecutionStatus == "fail" {
				executionFail++
			}
			if result.Status == "fail" {
				effectiveFail++
			}
		}
		fmt.Fprintf(&out, "| %s | %d | %d | %d | %d | %d |\n", kind, planned, executed, executionPass, executionFail, effectiveFail)
	}
	out.WriteString("\n<details>\n<summary>Complete case ledger with execution boundaries and evidence</summary>\n\n| Environment / case | Command / scenario | Boundary | Execution / effective outcome | Reason / evidence |\n| --- | --- | --- | --- | --- |\n")
	cases := append([]CaseResult{}, r.Cases...)
	sort.Slice(cases, func(i, j int) bool {
		return caseKey(cases[i].EnvironmentID, cases[i].ID) < caseKey(cases[j].EnvironmentID, cases[j].ID)
	})
	for _, c := range cases {
		links := []string{}
		for _, p := range []*string{c.StdoutArtifact, c.StderrArtifact} {
			if p != nil {
				links = append(links, artifactLink(*p))
			}
		}
		fmt.Fprintf(&out, "| %s / %s | %s / %s | %s | %s / %s | %s; %s |\n", markdown(c.EnvironmentID), markdown(c.ID), markdown(c.CommandPath), markdown(c.ScenarioClass), markdown(c.ExecutionKind), markdown(c.ExecutionStatus), markdown(c.Status), markdown(c.ReasonCode), strings.Join(links, ", "))
	}
	out.WriteString("\n</details>\n\n## 5. Help, error and log review\n\n")
	if r.Phase != "finalized" {
		out.WriteString("Communication review is pending. Capture-only results cannot establish overall pass.\n\n")
	}
	reviews := append([]Review{}, r.CommunicationReviews...)
	sort.Slice(reviews, func(i, j int) bool { return reviews[i].ID < reviews[j].ID })
	passed, failed, silent := 0, 0, 0
	for _, review := range reviews {
		switch review.Verdict {
		case "pass":
			passed++
		case "fail":
			failed++
		}
	}
	fmt.Fprintf(&out, "Reviewed messages: %d; passing: %d; failing: %d. Complete excerpts and all seven rubric records are in %s.\n\n", len(reviews), passed, failed, artifactLink("report.json"))
	if r.Phase == "finalized" {
		for _, result := range r.Cases {
			if result.Observed.ExitCode != nil && len(result.CommunicationReviewIDs) == 0 {
				silent++
			}
		}
		fmt.Fprintf(&out, "Executed cases with no emitted message: %d. The frozen %s preserves their explicit no-message records.\n\n", silent, artifactLink("corpus.json"))
	}
	for _, review := range reviews {
		if review.Verdict != "fail" {
			continue
		}
		fmt.Fprintf(&out, "- %s: %s, channel %s, message %s, findings %s.\n", markdown(review.ID), markdown(review.Verdict), markdown(review.Channel), markdown(review.MessageID), markdown(strings.Join(review.FindingIDs, ", ")))
		for _, item := range review.Rubric {
			if item.Verdict != "fail" {
				continue
			}
			fmt.Fprintf(&out, "  - %s: %s (%s).\n", markdown(item.Name), markdown(item.Verdict), markdown(item.Reason))
		}
	}
	out.WriteString("\n## 6. Defects\n\n")
	findings := append([]Finding{}, r.Findings...)
	sort.Slice(findings, func(i, j int) bool { return findings[i].ID < findings[j].ID })
	if len(findings) == 0 {
		out.WriteString("No defect records in this report.\n\n")
	}
	for _, f := range findings {
		argv, err := json.Marshal(f.Argv)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(&out, "### %s: %s\n\nDeduplication key: %s. Severity: %s. Category: %s. Target: %s.\n\nAffected subject: %s / %s; platforms: %s.\n\n", markdown(f.ID), markdown(f.Title), markdown(f.DeduplicationKey), markdown(f.Severity), markdown(f.Category), markdown(f.Target), markdown(f.Commit), markdown(f.Version), markdown(strings.Join(f.Platforms, ", ")))
		for i, ref := range f.Cases {
			if i == 0 && len(f.Cases) > 10 {
				fmt.Fprintf(&out, "<details>\n<summary>Affected cases (%d)</summary>\n\n", len(f.Cases))
			}
			fmt.Fprintf(&out, "- Case: %s / %s.\n", markdown(ref.EnvironmentID), markdown(ref.CaseID))
		}
		if len(f.Cases) > 10 {
			out.WriteString("\n</details>\n")
		}
		fmt.Fprintf(&out, "\nPreconditions: %s.\n\nRedacted argv: %s.\n\nReproduction: %s.\n\nExpected: %s.\n\nActual: %s.\n\nLayperson impact: %s.\n\n", markdown(strings.Join(f.Preconditions, "; ")), markdown(string(argv)), markdown(strings.Join(f.Reproduction, "; ")), markdown(f.Expected), markdown(f.Actual), markdown(f.Impact))
		for _, span := range f.Evidence {
			fmt.Fprintf(&out, "- Evidence: %s, lines %d-%d, SHA-256 %s.\n", artifactLink(span.Artifact), span.StartLine, span.EndLine, markdown(span.SHA256))
		}
		fmt.Fprintf(&out, "\nSuggested functional acceptance criteria: %s.\n\n", markdown(strings.Join(f.AcceptanceCriteria, "; ")))
		existing := "null"
		if f.ExistingIssueURL != nil {
			existing = *f.ExistingIssueURL
		}
		fmt.Fprintf(&out, "Existing issue: %s.\n\n", markdown(existing))
	}
	out.WriteString("## 7. Cleanup and artifact index\n\n")
	for _, cleanup := range r.Cleanup {
		fmt.Fprintf(&out, "- %s: %d cleanup actions; remaining owned resources: %s; retained: %s.\n", markdown(cleanup.EnvironmentID), len(cleanup.Actions), markdown(strings.Join(cleanup.Remaining, ", ")), markdown(strings.Join(cleanup.Retained, ", ")))
		for _, action := range cleanup.Actions {
			fmt.Fprintf(&out, "  - %s: passed=%t, expected %s, actual %s.\n", markdown(action.Name), action.Passed, markdown(action.Expected), markdown(action.Actual))
		}
	}
	artifacts := append([]Artifact{}, r.Artifacts...)
	sort.Slice(artifacts, func(i, j int) bool { return artifacts[i].Path < artifacts[j].Path })
	out.WriteString("\n<details>\n<summary>Complete artifact index and hashes</summary>\n\n| Artifact | Media type | Bytes | SHA-256 | Redaction |\n| --- | --- | --- | --- | --- |\n")
	for _, a := range artifacts {
		fmt.Fprintf(&out, "| %s | %s | %d | %s | %s |\n", artifactLink(a.Path), markdown(a.MediaType), a.Bytes, markdown(a.SHA256), markdown(a.RedactionPolicy))
	}
	out.WriteString("\n</details>\n")
	return []byte(out.String()), nil
}

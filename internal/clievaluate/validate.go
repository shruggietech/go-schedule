package clievaluate

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"
)

func caseKey(environment, id string) string { return environment + "/" + id }
func utc(value string) (time.Time, error) {
	if !strings.HasSuffix(value, "Z") {
		return time.Time{}, errors.New("timestamp must use UTC Z")
	}
	return time.Parse(time.RFC3339Nano, value)
}

// DeriveSummary recomputes totals without accepting a caller's pass claim.
func DeriveSummary(r Report) Summary {
	s := Summary{Planned: len(r.Cases), Findings: len(r.Findings), Outcome: "pass"}
	for _, c := range r.Cases {
		switch c.Status {
		case "pass":
			s.Passed++
		case "fail":
			s.Failed++
		case "blocked":
			s.Blocked++
		case "not_applicable":
			s.NotApplicable++
		}
		if c.ExecutionStatus == "pass" {
			s.ExecutionPassed++
		}
		if c.ExecutionStatus == "fail" {
			s.ExecutionFailed++
		}
	}
	if s.Failed > 0 {
		s.Outcome = "fail"
	}
	if s.Failed == 0 && (s.Blocked > 0 || s.Planned == s.NotApplicable) {
		s.Outcome = "blocked"
	}
	for _, c := range r.Cleanup {
		if len(c.Remaining) > 0 {
			s.Outcome = "blocked"
		}
		for _, a := range c.Actions {
			if !a.Passed {
				s.Outcome = "blocked"
			}
		}
	}
	if r.Phase != "finalized" || r.CompletionState != "completed" {
		s.Outcome = "blocked"
	}
	return s
}

// ValidateReport checks structure, complete planned coverage and verified evidence.
func ValidateReport(store *ArtifactStore, r Report, m Manifest, corpus Corpus) error {
	if err := ValidateBudget(r.Runner.Budget); err != nil {
		return err
	}
	for _, v := range []struct {
		name  string
		value any
	}{{"report", r}, {"manifest", m}, {"corpus", corpus}} {
		data, err := Encode(v.value)
		if err != nil {
			return err
		}
		if err := ValidateSchema(v.name, data); err != nil {
			return err
		}
	}
	start, err := utc(r.StartedAt)
	if err != nil {
		return err
	}
	finish, err := utc(r.FinishedAt)
	if err != nil || finish.Before(start) {
		return errors.New("invalid run timestamp ordering")
	}
	if r.RunID == "" || corpus.RunID != r.RunID {
		return errors.New("corpus/run identity mismatch")
	}
	manifestData, err := Encode(m)
	if err != nil {
		return err
	}
	corpusData, err := Encode(corpus)
	if err != nil {
		return err
	}
	if Digest(manifestData) != r.Runner.ManifestSHA256 || Digest(corpusData) != r.CorpusSHA256 {
		return errors.New("manifest/corpus digest mismatch")
	}
	if err := ValidateManifest(m); err != nil {
		return err
	}
	envs := map[string]Environment{}
	for _, e := range r.Environments {
		if e.ID == "" {
			return errors.New("empty environment identity")
		}
		if _, exists := envs[e.ID]; exists {
			return errors.New("duplicate environment identity")
		}
		envs[e.ID] = e
	}
	if len(envs) == 0 {
		return errors.New("run must plan at least one environment")
	}
	artifacts := map[string]Artifact{}
	total := int64(0)
	for _, a := range r.Artifacts {
		if _, exists := artifacts[a.Path]; exists {
			return errors.New("duplicate artifact path")
		}
		if _, err := store.Read(a); err != nil {
			return err
		}
		total += a.Bytes
		if total > r.Runner.Budget.ArtifactBytes {
			return errors.New("artifact total exceeds recorded budget")
		}
		artifacts[a.Path] = a
	}
	findings := map[string]Finding{}
	dedup := map[string]bool{}
	for _, f := range r.Findings {
		if f.ID == "" || findings[f.ID].ID != "" || dedup[f.DeduplicationKey] || f.DeduplicationKey == "" {
			return errors.New("duplicate/empty finding identity")
		}
		if len(f.Cases) == 0 || len(f.Evidence) == 0 || len(f.Reproduction) == 0 || len(f.AcceptanceCriteria) == 0 || f.Expected == "" || f.Actual == "" || f.Impact == "" || f.Commit != r.Subject.Commit || f.Version != r.Subject.Version {
			return errors.New("finding lacks exact subject or actionable reproduction")
		}
		findings[f.ID] = f
		dedup[f.DeduplicationKey] = true
		for _, span := range f.Evidence {
			a, ok := artifacts[span.Artifact]
			if !ok || span.SHA256 != a.SHA256 || span.StartLine < 1 || span.EndLine < span.StartLine {
				return errors.New("invalid finding evidence reference")
			}
			data, err := store.Read(a)
			if err != nil {
				return err
			}
			lines := strings.Count(strings.TrimSuffix(string(data), "\n"), "\n") + 1
			if span.EndLine > lines {
				return errors.New("finding evidence span exceeds artifact")
			}
		}
	}
	specs := map[string]CaseSpec{}
	for _, spec := range m.Cases {
		specs[spec.ID] = spec
	}
	seen := map[string]bool{}
	results := map[string]CaseResult{}
	for _, c := range r.Cases {
		key := caseKey(c.EnvironmentID, c.ID)
		if seen[key] {
			return errors.New("duplicate case/environment result")
		}
		seen[key] = true
		results[key] = c
		env, exists := envs[c.EnvironmentID]
		spec, known := specs[c.ID]
		if !exists || !known {
			return errors.New("unplanned case/environment result")
		}
		if c.CommandPath != spec.CommandPath || c.ScenarioClass != spec.ScenarioClass || c.ExecutionKind != spec.ExecutionKind || c.ExpectedExit != spec.ExpectedExit || !reflect.DeepEqual(c.ExpectedChecks, spec.Checks) {
			return errors.New("case contract differs from frozen manifest")
		}
		beg, err := utc(c.StartedAt)
		if err != nil {
			return err
		}
		end, err := utc(c.FinishedAt)
		if err != nil || end.Before(beg) || beg.Before(start) || end.After(finish) {
			return errors.New("invalid case timestamp ordering")
		}
		applicable := slices.Contains(spec.Platforms, env.OS)
		if !applicable {
			if c.ExecutionStatus != "not_applicable" || c.Status != "not_applicable" || c.ReasonCode != spec.ApplicabilityReason || spec.ApplicabilityReason == "" {
				return errors.New("unsupported not-applicable outcome")
			}
		} else if c.ExecutionStatus == "not_applicable" {
			return errors.New("applicable case cannot be skipped")
		}
		for _, p := range []*string{c.StdinArtifact, c.StdoutArtifact, c.StderrArtifact} {
			if p != nil {
				if _, ok := artifacts[*p]; !ok {
					return errors.New("case references absent artifact")
				}
			}
		}
		for _, p := range c.LogArtifacts {
			if _, ok := artifacts[p]; !ok {
				return errors.New("case references absent log artifact")
			}
		}
		if c.ExecutionStatus == "pass" || c.ExecutionStatus == "fail" || capturedCase(c) {
			if env.Status != "ready" || c.Observed.ExitCode == nil || c.StdoutArtifact == nil || c.StderrArtifact == nil {
				return errors.New("executed case lacks ready environment or channel/exit evidence")
			}
			if len(c.Observed.Checks) != len(spec.Checks) {
				return errors.New("executed case is missing required checks")
			}
			valid := *c.Observed.ExitCode == spec.ExpectedExit
			for i, check := range c.Observed.Checks {
				if check.Name != spec.Checks[i].Name || check.Expected != spec.Checks[i].Expected || len(check.Evidence) == 0 {
					return errors.New("check differs from expected contract or lacks evidence")
				}
				for _, p := range check.Evidence {
					if _, ok := artifacts[p]; !ok {
						return errors.New("check references absent evidence")
					}
				}
				if check.Passed != (check.Actual == check.Expected) {
					return errors.New("check verdict contradicts actual observation")
				}
				valid = valid && check.Passed
			}
			for _, check := range c.SideEffects {
				if check.Passed != (check.Actual == check.Expected) {
					return errors.New("side effect verdict contradicts actual observation")
				}
				valid = valid && check.Passed
				for _, p := range check.Evidence {
					if _, ok := artifacts[p]; !ok {
						return errors.New("side effect references absent evidence")
					}
				}
			}
			if c.ExecutionStatus != "blocked" && (c.ExecutionStatus == "pass") != valid {
				return errors.New("functional execution status contradicts observed checks")
			}
		}
		if c.ExecutionStatus == "blocked" && (c.ReasonCode == "" || c.Observed.ExitCode != nil && c.ReasonCode != "run_interrupted") {
			return errors.New("blocked case has invalid reason or executed exit")
		}
		if c.ReasonCode == "run_interrupted" && (r.CompletionState != "interrupted" || c.ExecutionStatus != "blocked") {
			return errors.New("interrupted case contradicts run completion state")
		}
		for _, id := range c.FindingIDs {
			f, ok := findings[id]
			if !ok || !slices.Contains(f.Cases, CaseRef{c.EnvironmentID, c.ID}) {
				return errors.New("unrelated case finding")
			}
		}
		if r.Phase == "finalized" && c.ExecutionStatus == "fail" && len(c.FindingIDs) == 0 {
			return errors.New("finalized failed case has no actionable finding")
		}
	}
	if len(seen) != len(envs)*len(specs) {
		return errors.New("missing planned case/environment results")
	}
	for _, f := range r.Findings {
		for _, ref := range f.Cases {
			if _, ok := results[caseKey(ref.EnvironmentID, ref.CaseID)]; !ok {
				return errors.New("finding references unknown case")
			}
		}
	}
	if err := validateCorpusResults(corpus, results); err != nil {
		return err
	}
	if err := validateCorpusEvidence(store, corpus, results, artifacts); err != nil {
		return err
	}
	if err := validateReportReviews(r, corpus, results, findings); err != nil {
		return err
	}
	cleaned := map[string]bool{}
	for _, c := range r.Cleanup {
		if _, ok := envs[c.EnvironmentID]; !ok || cleaned[c.EnvironmentID] {
			return errors.New("invalid cleanup environment")
		}
		cleaned[c.EnvironmentID] = true
	}
	if len(cleaned) != len(envs) {
		return errors.New("cleanup ledger incomplete")
	}
	if r.Summary != DeriveSummary(r) {
		return errors.New("summary counts or aggregate outcome contradict evidence")
	}
	return nil
}

// ValidateManifest rejects duplicate IDs and metadata-only behavioral coverage claims.
func ValidateManifest(m Manifest) error {
	if len(m.Commands) == 0 || len(m.Cases) == 0 {
		return errors.New("manifest must contain commands and behavioral cases")
	}
	ids := map[string]bool{}
	commands := map[string]Command{}
	for _, c := range m.Commands {
		if _, ok := commands[c.Path]; ok {
			return errors.New("duplicate command metadata")
		}
		commands[c.Path] = c
	}
	behavior := map[string]bool{}
	for _, c := range m.Cases {
		if c.ID == "" || ids[c.ID] || len(c.Platforms) == 0 || len(c.Checks) == 0 {
			return errors.New("invalid/duplicate manifest case")
		}
		ids[c.ID] = true
		if _, exists := commands[c.CommandPath]; !exists && c.CommandPath != "@engine" && c.CommandPath != "@transport" && c.CommandPath != "@isolation" {
			return fmt.Errorf("scenario names an unknown command: %s", c.CommandPath)
		}
		if c.ScenarioClass != "help" && c.ScenarioClass != "metadata" && c.ScenarioClass != "unknown-flag" && c.ScenarioClass != "arity" {
			behavior[c.CommandPath] = true
		}
	}
	for _, c := range m.Commands {
		if c.Runnable && !c.Hidden && !behavior[c.Path] {
			return fmt.Errorf("command has no behavioral scenario: %s", c.Path)
		}
	}
	if missing := MissingFlagCoverage(m); len(missing) != 0 {
		return fmt.Errorf("manifest has uncovered effective flags: %s", strings.Join(missing, ", "))
	}
	covered := map[string]bool{}
	for _, c := range m.Cases {
		for _, item := range c.Coverage {
			covered[item] = true
		}
	}
	for _, command := range m.Commands {
		for _, alias := range command.Aliases {
			if !covered[command.Path+" alias="+alias] {
				return fmt.Errorf("manifest has uncovered command alias: %s", alias)
			}
		}
	}
	return nil
}

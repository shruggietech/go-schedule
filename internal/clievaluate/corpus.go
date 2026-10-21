package clievaluate

import (
	"errors"
	"slices"
	"sort"
	"strings"
)

// AddMessage deduplicates only identical channel/text/command/scenario contexts.
func AddMessage(c *Corpus, result CaseResult, channel, text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if text == "" {
		return ""
	}
	id := "msg-" + Digest([]byte(result.CommandPath+"\x00"+result.ScenarioClass+"\x00"+channel+"\x00"+text))
	ref := CaseRef{result.EnvironmentID, result.ID}
	c.NoMessages = slices.DeleteFunc(c.NoMessages, func(candidate CaseRef) bool { return candidate == ref })
	for i := range c.Messages {
		if c.Messages[i].ID == id {
			if !slices.Contains(c.Messages[i].Cases, ref) {
				c.Messages[i].Cases = append(c.Messages[i].Cases, ref)
			}
			return id
		}
	}
	c.Messages = append(c.Messages, Message{ID: id, CommandPath: result.CommandPath, ScenarioClass: result.ScenarioClass, Channel: channel, Text: text, Cases: []CaseRef{ref}})
	return id
}

// SortCorpus freezes stable record order before hashing or review.
func SortCorpus(c *Corpus) {
	sort.Slice(c.Messages, func(i, j int) bool { return c.Messages[i].ID < c.Messages[j].ID })
	for i := range c.Messages {
		sort.Slice(c.Messages[i].Cases, func(a, b int) bool {
			return caseKey(c.Messages[i].Cases[a].EnvironmentID, c.Messages[i].Cases[a].CaseID) < caseKey(c.Messages[i].Cases[b].EnvironmentID, c.Messages[i].Cases[b].CaseID)
		})
	}
	sort.Slice(c.NoMessages, func(i, j int) bool {
		return caseKey(c.NoMessages[i].EnvironmentID, c.NoMessages[i].CaseID) < caseKey(c.NoMessages[j].EnvironmentID, c.NoMessages[j].CaseID)
	})
}

func validateCorpusResults(c Corpus, results map[string]CaseResult) error {
	seen := map[string]bool{}
	covered := map[string]bool{}
	for _, m := range c.Messages {
		expected := "msg-" + Digest([]byte(m.CommandPath+"\x00"+m.ScenarioClass+"\x00"+m.Channel+"\x00"+m.Text))
		if m.ID != expected || seen[m.ID] || m.Text == "" || len(m.Cases) == 0 {
			return errors.New("corpus message identity/content mismatch")
		}
		seen[m.ID] = true
		refs := map[string]bool{}
		for _, ref := range m.Cases {
			key := caseKey(ref.EnvironmentID, ref.CaseID)
			r, ok := results[key]
			if !ok || refs[key] || r.CommandPath != m.CommandPath || r.ScenarioClass != m.ScenarioClass || !capturedCase(r) {
				return errors.New("corpus references unrelated/unexecuted case")
			}
			refs[key] = true
			covered[key] = true
		}
	}
	silent := map[string]bool{}
	for _, ref := range c.NoMessages {
		key := caseKey(ref.EnvironmentID, ref.CaseID)
		r, ok := results[key]
		if !ok || covered[key] || silent[key] || !capturedCase(r) {
			return errors.New("invalid no-message occurrence")
		}
		silent[key] = true
	}
	for key, r := range results {
		if capturedCase(r) {
			if !covered[key] && !silent[key] {
				return errors.New("executed case missing communication occurrence")
			}
		}
	}
	return nil
}

func capturedCase(r CaseResult) bool {
	return r.Observed.ExitCode != nil && (r.ExecutionStatus == "pass" || r.ExecutionStatus == "fail" || r.ExecutionStatus == "blocked" && r.ReasonCode == "run_interrupted")
}

func validateCorpusEvidence(store *ArtifactStore, c Corpus, results map[string]CaseResult, artifacts map[string]Artifact) error {
	covered := map[string]bool{}
	for _, m := range c.Messages {
		for _, ref := range m.Cases {
			r := results[caseKey(ref.EnvironmentID, ref.CaseID)]
			paths := []string{}
			switch m.Channel {
			case "help", "description", "stdout":
				if r.StdoutArtifact != nil {
					paths = append(paths, *r.StdoutArtifact)
				}
			case "stderr":
				if r.StderrArtifact != nil {
					paths = append(paths, *r.StderrArtifact)
				}
			case "log":
				paths = append(paths, r.LogArtifacts...)
			}
			matched := false
			for _, path := range paths {
				a, ok := artifacts[path]
				if !ok {
					return errors.New("corpus lacks captured channel evidence")
				}
				data, err := store.Read(a)
				if err != nil {
					return err
				}
				if strings.ReplaceAll(string(data), "\r\n", "\n") == m.Text {
					matched = true
					covered[caseKey(ref.EnvironmentID, ref.CaseID)+"\x00"+path] = true
				}
			}
			if !matched {
				return errors.New("corpus text does not equal its captured channel evidence")
			}
		}
	}
	for _, ref := range c.NoMessages {
		r := results[caseKey(ref.EnvironmentID, ref.CaseID)]
		paths := append([]string{}, r.LogArtifacts...)
		if r.StdoutArtifact != nil {
			paths = append(paths, *r.StdoutArtifact)
		}
		if r.StderrArtifact != nil {
			paths = append(paths, *r.StderrArtifact)
		}
		for _, path := range paths {
			if artifacts[path].Bytes != 0 {
				return errors.New("no-message claim contradicts nonempty channel evidence")
			}
		}
	}
	for key, result := range results {
		if result.Observed.ExitCode == nil {
			continue
		}
		paths := append([]string{}, result.LogArtifacts...)
		for _, p := range []*string{result.StdoutArtifact, result.StderrArtifact} {
			if p != nil {
				paths = append(paths, *p)
			}
		}
		for _, path := range paths {
			if artifacts[path].Bytes > 0 && !covered[key+"\x00"+path] {
				return errors.New("captured nonempty channel is missing from the review corpus")
			}
		}
	}
	return nil
}

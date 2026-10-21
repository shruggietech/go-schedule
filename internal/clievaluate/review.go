package clievaluate

import (
	"errors"
	"reflect"
	"slices"
)

var rubricNames = []string{"plain_language", "action_outcome", "actionable_correction", "severity", "terminology", "identifiers_paths_units", "secret_absence"}

func validateReviewRecords(reviews []Review, corpus Corpus, findings map[string]Finding, complete bool) error {
	messages := map[string]Message{}
	for _, m := range corpus.Messages {
		messages[m.ID] = m
	}
	seen := map[string]bool{}
	ids := map[string]bool{}
	for _, r := range reviews {
		m, ok := messages[r.MessageID]
		if !ok || seen[r.MessageID] || r.ID == "" || ids[r.ID] || r.Channel != m.Channel || r.Excerpt != m.Text || !reflect.DeepEqual(r.Cases, m.Cases) {
			return errors.New("stale, duplicate or unrelated review record")
		}
		seen[r.MessageID] = true
		ids[r.ID] = true
		if len(r.Rubric) != len(rubricNames) {
			return errors.New("incomplete communication rubric")
		}
		valid := true
		names := map[string]bool{}
		for _, item := range r.Rubric {
			if !slices.Contains(rubricNames, item.Name) || names[item.Name] || item.Reason == "" || !slices.Contains([]string{"pass", "fail", "not_applicable"}, item.Verdict) {
				return errors.New("invalid communication rubric")
			}
			names[item.Name] = true
			valid = valid && item.Verdict != "fail"
		}
		if (r.Verdict == "pass") != valid || r.Verdict == "not_applicable" {
			return errors.New("review verdict contradicts rubric")
		}
		if !valid && len(r.FindingIDs) == 0 {
			return errors.New("failed communication review requires actionable finding")
		}
		for _, id := range r.FindingIDs {
			f, ok := findings[id]
			if !ok {
				return errors.New("review references absent finding")
			}
			for _, ref := range r.Cases {
				if !slices.Contains(f.Cases, ref) {
					return errors.New("review finding omits affected case")
				}
			}
		}
	}
	if complete && len(seen) != len(messages) {
		return errors.New("pending_review: missing communication reviews")
	}
	return nil
}

func validateReportReviews(r Report, c Corpus, results map[string]CaseResult, findings map[string]Finding) error {
	if err := validateReviewRecords(r.CommunicationReviews, c, findings, r.Phase == "finalized"); err != nil {
		return err
	}
	required := map[string][]string{}
	failed := map[string]bool{}
	byID := map[string]Review{}
	for _, review := range r.CommunicationReviews {
		byID[review.ID] = review
		for _, ref := range review.Cases {
			key := caseKey(ref.EnvironmentID, ref.CaseID)
			required[key] = append(required[key], review.ID)
			failed[key] = failed[key] || review.Verdict == "fail"
		}
	}
	for key, result := range results {
		want := required[key]
		got := append([]string{}, result.CommunicationReviewIDs...)
		slices.Sort(want)
		slices.Sort(got)
		if !slices.Equal(want, got) {
			return errors.New("case communication review index mismatch")
		}
		for _, id := range got {
			if _, ok := byID[id]; !ok {
				return errors.New("unknown review reference")
			}
		}
		status := result.ExecutionStatus
		reason := result.ReasonCode
		if status == "pass" {
			if r.Phase != "finalized" {
				status = "blocked"
				reason = "pending_review"
			} else if failed[key] {
				status = "fail"
				reason = "communication_finding"
			}
		}
		if result.Status != status {
			return errors.New("effective case status contradicts execution/review")
		}
		if status == "blocked" && result.ExecutionStatus == "pass" && reason != "pending_review" {
			return errors.New("missing pending-review reason")
		}
	}
	return nil
}

// ApplyReviews validates an agent submission against the immutable corpus.
func ApplyReviews(r *Report, c Corpus, set ReviewSet) error {
	if r.Phase != "pending_review" {
		return errors.New("only a completed capture with pending_review phase can accept a review submission")
	}
	if set.RunID != r.RunID || set.CorpusSHA256 != r.CorpusSHA256 || set.SchemaVersion != Version {
		return errors.New("review run/corpus identity mismatch")
	}
	data, err := Encode(set)
	if err != nil {
		return err
	}
	if err := ValidateSchema("reviews", data); err != nil {
		return err
	}
	findings := map[string]Finding{}
	for _, f := range r.Findings {
		findings[f.ID] = f
	}
	for _, f := range set.Findings {
		if _, exists := findings[f.ID]; exists {
			return errors.New("duplicate submitted finding")
		}
		findings[f.ID] = f
	}
	if err := validateReviewRecords(set.Reviews, c, findings, true); err != nil {
		return err
	}
	r.CommunicationReviews = set.Reviews
	r.Findings = append(r.Findings, set.Findings...)
	r.Phase = "finalized"
	for i := range r.Cases {
		result := &r.Cases[i]
		result.CommunicationReviewIDs = []string{}
		for _, finding := range set.Findings {
			if slices.Contains(finding.Cases, CaseRef{result.EnvironmentID, result.ID}) && !slices.Contains(result.FindingIDs, finding.ID) {
				result.FindingIDs = append(result.FindingIDs, finding.ID)
			}
		}
		result.Status = result.ExecutionStatus
		if result.ExecutionStatus == "pass" {
			result.ReasonCode = "executed"
		}
		for _, review := range set.Reviews {
			if slices.Contains(review.Cases, CaseRef{result.EnvironmentID, result.ID}) {
				result.CommunicationReviewIDs = append(result.CommunicationReviewIDs, review.ID)
				if review.Verdict == "fail" && result.ExecutionStatus == "pass" {
					result.Status = "fail"
					result.ReasonCode = "communication_finding"
				}
				for _, id := range review.FindingIDs {
					if !slices.Contains(result.FindingIDs, id) {
						result.FindingIDs = append(result.FindingIDs, id)
					}
				}
			}
		}
	}
	r.Summary = DeriveSummary(*r)
	return nil
}

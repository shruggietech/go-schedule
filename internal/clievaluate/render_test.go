package clievaluate

import (
	"bytes"
	"strings"
	"testing"
)

func TestHumanReportIsDeterministicAndUsesFixedOutline(t *testing.T) {
	_, r, _, _ := reportFixture(t)
	r.Subject.Version = "<script>|bad\nline"
	first, err := Render(r)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Render(r)
	if err != nil || !bytes.Equal(first, second) {
		t.Fatal("render is not deterministic")
	}
	for _, heading := range []string{"1. Run identity and outcome", "2. Environment and budget", "3. Coverage table", "4. Functional results", "5. Help, error and log review", "6. Defects", "7. Cleanup and artifact index"} {
		if !strings.Contains(string(first), heading) {
			t.Fatalf("missing section %q", heading)
		}
	}
	if strings.Contains(string(first), "<script>") {
		t.Fatal("unescaped Markdown HTML")
	}
}

func TestHumanReportSurfacesFailedReviewsWithoutDumpingSuccessfulRubrics(t *testing.T) {
	_, report, _, _ := reportFixture(t)
	report.CommunicationReviews = []Review{{ID: "review-pass", Verdict: "pass", Rubric: []Rubric{{Name: "plain_language", Verdict: "pass", Reason: "successful-review-detail-kept-in-canonical-json"}}}, {ID: "review-fail", Verdict: "fail", FindingIDs: []string{"F007"}, Rubric: []Rubric{{Name: "plain_language", Verdict: "fail", Reason: "operation description is missing"}}}}
	data, err := Render(report)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if strings.Contains(text, "successful-review-detail-kept-in-canonical-json") || !strings.Contains(text, "operation description is missing") || !strings.Contains(text, "Reviewed messages: 2; passing: 1; failing: 1") {
		t.Fatal("human review summary buries failures in successful rubric records")
	}
	if !strings.Contains(text, "<details>") || !strings.Contains(text, "report.json") {
		t.Fatal("complete detail is not discoverable")
	}
}

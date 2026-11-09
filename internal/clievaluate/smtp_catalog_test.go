package clievaluate

import "testing"

func TestSMTPTaskScenarioRetainsSelectedAcknowledgementContract(t *testing.T) {
	for _, s := range behavioralSteps() {
		if s.Spec.ID != "smtp.task-delivery" {
			continue
		}
		if !s.JSON || s.OutputContains != "" {
			t.Fatal("SMTP task case must require structured acknowledgement in JSON mode")
		}
		if s.Hook != "run-smtp" || len(s.Spec.Checks) != 3 {
			t.Fatal("SMTP delivery and independent outcome checks must remain")
		}
		return
	}
	t.Fatal("SMTP task case missing")
}

func TestJSONAcknowledgementCasesRequireStructuredOutput(t *testing.T) {
	paths := map[string]bool{"alerts ack": false, "notification channel rm": false, "task enable": false, "task disable": false, "task rm": false, "task run-now": false, "group enable": false, "group disable": false, "group rm": false, "chain rm": false, "gui": false, "service install": false, "service uninstall": false, "service start": false, "service stop": false, "service restart": false, "service status": false}
	for _, s := range behavioralSteps() {
		if _, ok := paths[s.Spec.CommandPath]; !ok {
			continue
		}
		paths[s.Spec.CommandPath] = true
		if !s.JSON || s.OutputContains != "" {
			t.Errorf("%s skips selected JSON contract", s.Spec.ID)
		}
	}
	for path, seen := range paths {
		if !seen {
			t.Errorf("missing acknowledgement coverage for %s", path)
		}
	}
}

func TestRunAcknowledgementBindsSelectedModeAndTask(t *testing.T) {
	for _, tc := range []struct {
		out              string
		structured, want bool
	}{
		{"triggered run for task owned-id\n", false, true},
		{"triggered run for task owned-id-other\n", false, false},
		{"{\"id\":\"owned-id\",\"status\":\"accepted\"}", true, true},
		{"{\"id\":\"other\",\"status\":\"accepted\"}", true, false},
		{"{\"id\":\"owned-id\",\"status\":\"failed\"}", true, false},
		{"triggered run for task owned-id", true, false},
		{"{\"id\":\"owned-id\",\"status\":\"accepted\"} trailing", true, false},
	} {
		if got := runAcknowledgesTask([]byte(tc.out), "owned-id", tc.structured); got != tc.want {
			t.Errorf("structured=%t output=%q got=%t", tc.structured, tc.out, got)
		}
	}
}

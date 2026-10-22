package clievaluate

import "testing"

func TestSMTPTaskScenarioRetainsRunAcknowledgementContract(t *testing.T) {
	for _, s := range behavioralSteps() {
		if s.Spec.ID != "smtp.task-delivery" {
			continue
		}
		if s.JSON || s.OutputContains != "triggered run for task" {
			t.Fatal("SMTP task case must require its text acknowledgement")
		}
		if s.Hook != "run-smtp" || len(s.Spec.Checks) != 3 {
			t.Fatal("SMTP delivery and independent outcome checks must remain")
		}
		return
	}
	t.Fatal("SMTP task case missing")
}

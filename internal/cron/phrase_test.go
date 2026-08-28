package cron

import "testing"

func TestExplainCompositeCronDescriptions(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"0 9,17 * * *", "at minute 00 during hours 9 and 17 every day"},
		{"*/10 9-17 * * *", "every 10 minutes during hours 9 through 17 every day"},
		{"0 0 1,15 JAN,MAR *", "at 00:00 on dates 1 and 15 in January and March"},
		{"0 12 * * MON,WED,FRI", "at 12:00 on Monday, Wednesday, and Friday"},
		{"10-20/2 * * * *", "at minutes 10, 12, 14, 16, 18, and 20 every hour every day"},
		{"0 9 */2 * *", "at 09:00 on dates 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, and 31"},
		{"0 9 * */2 *", "at 09:00 every day in January, March, May, July, September, and November"},
		{"0 9 * * */2", "at 09:00 on Sunday, Tuesday, Thursday, and Saturday"},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			got, bad, err := Explain(tt.expr)
			if err != nil || bad.Reason != "" || got != tt.want {
				t.Fatalf("Explain=%q refusal=%q err=%v, want %q", got, bad.Reason, err, tt.want)
			}
		})
	}
}

func TestExplainKeepsEstablishedConcisePhrases(t *testing.T) {
	for expr, want := range map[string]string{
		"0 9 * * *":    "every day at 09:00",
		"0 9 * * 1-5":  "weekdays at 09:00",
		"*/15 * * * *": "every 15 minutes starting at 00:00",
	} {
		got, bad, err := Explain(expr)
		if err != nil || bad.Reason != "" || got != want {
			t.Errorf("Explain(%q)=%q refusal=%q err=%v, want %q", expr, got, bad.Reason, err, want)
		}
	}
}

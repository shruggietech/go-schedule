package cron

import (
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/internal/schedule"
)

func TestRebootExplainCompileAndConvert(t *testing.T) {
	parsed, err := Parse("@reboot")
	if err != nil || !parsed.OK || !parsed.Startup || parsed.Bad.Reason != "" {
		t.Fatalf("Parse = %+v, err=%v", parsed, err)
	}
	phrase, bad, err := Explain(" @REBOOT ")
	if err != nil || bad.Reason != "" || phrase != schedule.StartupPhrase {
		t.Fatalf("Explain = %q, bad=%q, err=%v", phrase, bad.Reason, err)
	}
	sch, bad, err := Compile("@reboot", "UTC", time.Now())
	if err != nil || bad.Reason != "" || !schedule.IsStartup(sch) || sch.Expression != "@reboot" {
		t.Fatalf("Compile = %+v, bad=%q, err=%v", sch, bad.Reason, err)
	}
	toHuman, err := Convert("@reboot", "")
	if err != nil || toHuman.Output != schedule.StartupPhrase {
		t.Fatalf("cron conversion = %+v, err=%v", toHuman, err)
	}
	toCron, err := Convert(schedule.StartupPhrase, "")
	if err != nil || toCron.Output != "@reboot" {
		t.Fatalf("human conversion = %+v, err=%v", toCron, err)
	}
}

package incident

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
	"unicode/utf16"
)

func eventFixture(pid int, id int, provider, app, when string) string {
	return fmt.Sprintf(`<Event xmlns="http://schemas.microsoft.com/win/2004/08/events/event"><System><Provider Name="%s"/><EventID>%d</EventID><TimeCreated SystemTime="%s"/><EventRecordID>42</EventRecordID><Execution ProcessID="%d"/></System><EventData><Data>Application: %s
CoreCLR Version: 10.0.1226.41902
.NET Version: 10.0.12
Message: SECRET-UNRELATED-PAYLOAD
Stack: private stack
</Data></EventData></Event>`, provider, id, when, pid, app)
}
func utf16EventXML(s string) []byte {
	b := []byte{0xff, 0xfe}
	for _, w := range utf16.Encode([]rune(s)) {
		b = binary.LittleEndian.AppendUint16(b, w)
	}
	return b
}
func TestRuntimeEventAttribution(t *testing.T) {
	start := time.Date(2026, 10, 22, 17, 43, 35, 0, time.UTC)
	end := start.Add(time.Second)
	valid := eventFixture(123, 1023, ".NET Runtime", "pwsh.exe", start.Format(time.RFC3339Nano))
	for _, tc := range []struct {
		name, xml, want string
		count           int
	}{
		{"runtime", valid, "collected", 1},
		{"failfast", strings.Replace(valid, "1023", "1025", 1), "collected", 1},
		{"case", strings.Replace(valid, "pwsh.exe", "PWSH.EXE", 1), "collected", 1},
		{"wrong pid", strings.Replace(valid, `ProcessID="123"`, `ProcessID="124"`, 1), "absent", 0},
		{"wrong provider", strings.Replace(valid, ".NET Runtime", "Other Runtime", 1), "absent", 0},
		{"wrong event", strings.Replace(valid, "1023", "1000", 1), "absent", 0},
		{"missing application", strings.Replace(valid, "Application: pwsh.exe\n", "", 1), "failed", 0},
		{"wrong application", strings.Replace(valid, "pwsh.exe", "other.exe", 1), "absent", 0},
		{"reused pid before start", strings.Replace(valid, start.Format(time.RFC3339Nano), start.Add(-time.Nanosecond).Format(time.RFC3339Nano), 1), "absent", 0},
		{"after exit", strings.Replace(valid, start.Format(time.RFC3339Nano), end.Add(time.Nanosecond).Format(time.RFC3339Nano), 1), "absent", 0},
		{"at exit", strings.Replace(valid, start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano), 1), "collected", 1},
		{"invalid timestamp", strings.Replace(valid, start.Format(time.RFC3339Nano), "invalid", 1), "failed", 0},
		{"ambiguous application", strings.Replace(valid, "CoreCLR Version:", "Application: other.exe\nCoreCLR Version:", 1), "failed", 0},
		{"malformed", valid + "<broken>", "failed", 0},
		{"empty", "", "absent", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := collectRuntimeEvents(Target{PID: 123, Executable: `C:\Program Files\PowerShell\pwsh.exe`}, start, end, func(context.Context, string) RuntimeEventQueryResult {
				return RuntimeEventQueryResult{Data: utf16EventXML("<Events>" + tc.xml + "</Events>")}
			})
			if got.Status != tc.want || len(got.Events) != tc.count {
				t.Fatalf("want %s/%d, got %+v", tc.want, tc.count, got)
			}
			b, _ := json.Marshal(got)
			if strings.Contains(string(b), "SECRET") || strings.Contains(string(b), "private stack") {
				t.Fatal("arbitrary event content leaked")
			}
		})
	}
}
func TestRuntimeEventBoundsAndFailures(t *testing.T) {
	start := time.Now().UTC()
	end := start.Add(time.Second)
	event := eventFixture(123, 1023, ".NET Runtime", "pwsh.exe", start.Format(time.RFC3339Nano))
	for _, tc := range []struct {
		name   string
		result RuntimeEventQueryResult
		want   string
	}{
		{"query unavailable", RuntimeEventQueryResult{Error: errors.New("native exit 5"), Unavailable: true}, "unavailable"},
		{"timeout", RuntimeEventQueryResult{Error: context.DeadlineExceeded}, "failed"},
		{"output limit", RuntimeEventQueryResult{Truncated: true}, "truncated"},
		{"raw byte overflow", RuntimeEventQueryResult{Data: make([]byte, RuntimeEventByteLimit+1)}, "truncated"},
		{"odd utf16", RuntimeEventQueryResult{Data: []byte{255, 254, 60}}, "failed"},
		{"invalid surrogate", RuntimeEventQueryResult{Data: []byte{255, 254, 0, 216}}, "failed"},
		{"count overflow", RuntimeEventQueryResult{Data: utf16EventXML("<Events>" + strings.Repeat(event, RuntimeEventLimit+1) + "</Events>")}, "truncated"},
		{"wrong root", RuntimeEventQueryResult{Data: utf16EventXML("<Other/>")}, "failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := collectRuntimeEvents(Target{PID: 123, Executable: "pwsh.exe"}, start, end, func(ctx context.Context, q string) RuntimeEventQueryResult {
				if _, ok := ctx.Deadline(); !ok {
					t.Error("query has no deadline")
				}
				if !strings.Contains(q, "123") || !strings.Contains(q, "1025") {
					t.Error("query lost exact target or event scope")
				}
				return tc.result
			})
			if got.Status != tc.want || len(got.Events) > RuntimeEventLimit {
				t.Fatalf("want %s, got %+v", tc.want, got)
			}
		})
	}
	calls := 0
	got := collectRuntimeEvents(Target{}, start, end, func(context.Context, string) RuntimeEventQueryResult { calls++; return RuntimeEventQueryResult{} })
	if calls != 0 || got.Status != "failed" {
		t.Fatal("invalid target queried")
	}
}
func TestRuntimeEventReceiptIsSeparateArtifact(t *testing.T) {
	parent := canonicalTempDir(t)
	c, err := New(parent)
	if err != nil {
		t.Fatal(err)
	}
	a, err := c.saveRuntimeEvents(RuntimeEventReceipt{Kind: "go-schedule.runtime-events", Version: 1, Status: "absent", Events: []RuntimeEvent{}})
	if err != nil || a.Name != "runtime-events.json" || a.Status != "collected" {
		t.Fatalf("artifact: %+v %v", a, err)
	}
	b, err := c.Read(a.Name, MetadataLimit)
	if err != nil {
		t.Fatal(err)
	}
	var r RuntimeEventReceipt
	if err = json.Unmarshal(b, &r); err != nil || r.Status != "absent" {
		t.Fatalf("receipt status lost: %s %v", b, err)
	}
	if _, err = os.Stat(c.Dir + "/incident.json"); !os.IsNotExist(err) {
		t.Fatal("event receipt replaced incident envelope")
	}
}

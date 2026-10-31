package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/clientobserve"
	"github.com/shruggietech/go-schedule/v2/internal/clientprofile"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestCalendarWindowDefaultsAndRejectsInvalidRanges(t *testing.T) {
	now := time.Date(2026, 10, 19, 12, 0, 0, 0, time.UTC)
	from, to, err := calendarWindow(now, "", "")
	if err != nil || !from.Equal(now) || to.Sub(from) != 7*24*time.Hour {
		t.Fatalf("from=%s to=%s err=%v", from, to, err)
	}
	for _, input := range [][2]string{{"2026-10-19", ""}, {"", "wrong"}, {"2026-10-19T12:00:00", ""}, {"2026-10-19T12:00:00Z", "2026-10-18T12:00:00Z"}, {"2026-10-19T12:00:00Z", "2026-11-19T12:00:00Z"}} {
		if _, _, err := calendarWindow(now, input[0], input[1]); !errors.Is(err, errUsage) {
			t.Fatalf("input=%v err=%v", input, err)
		}
	}
}

func observationFixture(t *testing.T, feature func(http.ResponseWriter, *http.Request)) *client.Client {
	t.Helper()
	return client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/health":
			json.NewEncoder(w).Encode(server.HealthResponse{Status: "ok", Version: "1.5.0"})
		case "/v1/manifest":
			json.NewEncoder(w).Encode(server.ManifestResponse{InstallationID: "owned-daemon", ProductVersion: "1.5.0"})
		case "/v1/access/current":
			json.NewEncoder(w).Encode(domain.Actor{Capability: domain.CapabilityObserve})
		default:
			feature(w, r)
		}
	}))
}

func TestOperationalCommandsUseScopedClientClockAndOutputs(t *testing.T) {
	now := time.Date(2026, 10, 19, 12, 0, 0, 0, time.UTC)
	var paths []string
	api := observationFixture(t, func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.RequestURI())
		switch r.URL.Path {
		case "/v1/calendar":
			if r.URL.Query().Get("from") != "2026-10-19T12:00:00Z" || r.URL.Query().Get("to") != "2026-10-26T12:00:00Z" {
				t.Errorf("wrong default clock/window: %s", r.URL)
			}
			json.NewEncoder(w).Encode(server.CalendarResponse{From: now, To: now.Add(7 * 24 * time.Hour), Occurrences: []server.Occurrence{}})
		case "/v1/runtime-info":
			json.NewEncoder(w).Encode(server.RuntimeInfoResponse{DataDir: "owned-data", IPCPath: "owned-ipc", ProcessID: 123})
		case "/v1/runs/active":
			json.NewEncoder(w).Encode(map[string]any{"runs": []domain.Run{{ID: "active-run", TaskID: "task", ScheduledFor: now, StartedAt: &now}}})
		case "/v1/runs/exact-run":
			json.NewEncoder(w).Encode(domain.Run{ID: "exact-run", TaskID: "task", Output: "captured output", OutputTruncated: true, SourceRunID: "source-run"})
		case "/v1/logs":
			json.NewEncoder(w).Encode(server.LogsResponse{LogPath: "owned-log", Logs: []domain.LogRecord{}})
		default:
			t.Errorf("unexpected endpoint %s", r.URL)
			http.NotFound(w, r)
		}
	})
	for _, test := range []struct {
		args []string
		want string
	}{{[]string{"--json", "calendar"}, `"occurrences": []`}, {[]string{"runtime"}, "Process ID: 123"}, {[]string{"runs", "active"}, "active-run"}, {[]string{"runs", "get", "exact-run"}, "captured output"}, {[]string{"--json", "runs", "get", "exact-run"}, `"source_run_id": "source-run"`}, {[]string{"--json", "logs", "--metadata"}, `"log_path": "owned-log"`}} {
		var out, diagnostics bytes.Buffer
		root := newRootWithDependencies(Dependencies{Client: api, Out: &out, Err: &diagnostics, Now: func() time.Time { return now }})
		root.SetArgs(test.args)
		if err := root.Execute(); err != nil || !strings.Contains(out.String(), test.want) || diagnostics.Len() != 0 {
			t.Errorf("args=%v out=%s stderr=%s err=%v", test.args, out.String(), diagnostics.String(), err)
		}
	}
	if len(paths) != 6 {
		t.Fatalf("paths=%v", paths)
	}
}

func TestSystemsAllKeepsLocalSuccessWhenProfilesCannotBeRead(t *testing.T) {
	profilePath := filepath.Join(t.TempDir(), "broken-profiles.json")
	if err := os.WriteFile(profilePath, []byte("invalid JSON"), 0600); err != nil {
		t.Fatal(err)
	}
	api := observationFixture(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/system-summary" {
			t.Errorf("unexpected endpoint %s", r.URL)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(domain.SystemSummary{Schema: domain.SystemSummarySchema, ObservedAt: time.Now(), ActiveTaskCount: 7})
	})
	var out bytes.Buffer
	root := newRootWithDependencies(Dependencies{Client: api, Profiles: clientprofile.NewStore(profilePath), Out: &out, Err: io.Discard})
	root.SetArgs([]string{"--json", "systems", "summary", "--all"})
	err := root.Execute()
	if !errors.Is(err, clientobserve.ErrPartial) || !strings.Contains(out.String(), `"active_task_count": 7`) || !strings.Contains(out.String(), `"profile_store_unavailable"`) || !strings.Contains(out.String(), `"observed_daemon_id": "owned-daemon"`) {
		t.Fatalf("out=%s err=%v", out.String(), err)
	}
}

func TestSearchCommandUsesKindsLimitAndPreservesSource(t *testing.T) {
	api := observationFixture(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/search" || r.URL.Query().Get("q") != "backup" || r.URL.Query().Get("kind") != "task" || r.URL.Query().Get("limit") != "5" {
			t.Errorf("request=%s", r.URL)
		}
		json.NewEncoder(w).Encode(domain.DaemonSearch{Schema: domain.DaemonSearchSchema, ObservedAt: time.Now(), Query: "backup", Results: []domain.SearchMatch{{Kind: domain.SearchKindTask, ObjectID: "object-id", Name: "Backup"}}})
	})
	var out bytes.Buffer
	root := newRootWithDependencies(Dependencies{Client: api, Out: &out, Err: io.Discard})
	root.SetArgs([]string{"--json", "systems", "search", "backup", "--kind", "task", "--limit", "5"})
	if err := root.Execute(); err != nil || !strings.Contains(out.String(), `"observed_daemon_id": "owned-daemon"`) || !strings.Contains(out.String(), `"object_id": "object-id"`) {
		t.Fatalf("out=%s err=%v", out.String(), err)
	}
}

func TestObservationBadInputAndCanceledParentNeverReadFeatures(t *testing.T) {
	reads := 0
	api := observationFixture(t, func(w http.ResponseWriter, r *http.Request) { reads++; http.NotFound(w, r) })
	for _, args := range [][]string{{"calendar", "--from", "yesterday"}, {"systems", "search", " ", "--kind", "task"}, {"systems", "search", "x", "--limit", "51"}, {"systems", "search", "x", "--kind", "passwords"}, {"systems", "summary", "--all", "--profile", "production"}, {"systems", "summary", "--all", "--profile", ""}, {"systems", "summary", "--all", "--ipc-path", ""}} {
		root := newRootWithDependencies(Dependencies{Client: api, Out: io.Discard, Err: io.Discard})
		root.SetArgs(args)
		if err := root.Execute(); !errors.Is(err, errUsage) {
			t.Errorf("args=%v err=%v", args, err)
		}
	}
	root := newRootWithDependencies(Dependencies{Client: api, Out: io.Discard, Err: io.Discard})
	root.SetArgs([]string{"systems", "summary"})
	ctx, cancel := context.WithCancel(root.Context())
	cancel()
	if err := root.ExecuteContext(ctx); !errors.Is(err, context.Canceled) {
		t.Errorf("cancel err=%v", err)
	}
	if reads != 0 {
		t.Fatalf("unexpected feature requests=%d", reads)
	}
}

func TestRuntimeRejectsRemoteWithoutReadingIt(t *testing.T) {
	var diagnostics bytes.Buffer
	root := newRootWithDependencies(Dependencies{Client: client.NewUnavailableRemote("https://example.invalid"), Out: io.Discard, Err: &diagnostics})
	root.SetArgs([]string{"runtime"})
	if err := root.Execute(); !errors.Is(err, errUsage) || !strings.Contains(err.Error(), "only through local IPC") {
		t.Fatalf("err=%v", err)
	}
}

func TestHistoryAndLogsZeroLimitPreserveInheritedRequest(t *testing.T) {
	var requests []string
	api := observationFixture(t, func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.URL.RequestURI())
		if r.URL.Query().Has("limit") {
			t.Errorf("zero limit changed inherited query: %s", r.URL)
		}
		switch r.URL.Path {
		case "/v1/runs":
			json.NewEncoder(w).Encode(map[string]any{"runs": []domain.Run{}})
		case "/v1/logs":
			json.NewEncoder(w).Encode(server.LogsResponse{Logs: []domain.LogRecord{}})
		default:
			t.Errorf("unexpected path %s", r.URL)
			http.NotFound(w, r)
		}
	})
	for _, command := range []string{"runs", "logs"} {
		root := newRootWithDependencies(Dependencies{Client: api, Out: io.Discard, Err: io.Discard})
		root.SetArgs([]string{command, "--limit", "0"})
		if err := root.Execute(); err != nil {
			t.Fatalf("%s zero limit rejected: %v", command, err)
		}
		root = newRootWithDependencies(Dependencies{Client: api, Out: io.Discard, Err: io.Discard})
		root.SetArgs([]string{command, "--limit", "-1"})
		if err := root.Execute(); !errors.Is(err, errUsage) {
			t.Fatalf("%s negative limit err=%v", command, err)
		}
	}
	if len(requests) != 2 || requests[0] != "/v1/runs" || requests[1] != "/v1/logs" {
		t.Fatalf("requests=%v", requests)
	}
}

func TestCalendarHumanOutputDistinguishesPredictionsAndRecordedIdentity(t *testing.T) {
	now := time.Date(2026, 10, 19, 12, 0, 0, 0, time.UTC)
	var out bytes.Buffer
	err := writeCalendar(&out, server.CalendarResponse{Occurrences: []server.Occurrence{{TaskID: "future-task", TaskName: "Future", Time: now.Add(time.Hour), Kind: "scheduled"}, {TaskID: "past-task", TaskName: "Past", RunID: "exact-run", Time: now, Kind: "past"}}})
	text := out.String()
	if err != nil || !strings.Contains(text, "predicted") || !strings.Contains(text, "recorded") || !strings.Contains(text, "exact-run") || strings.Index(text, "past-task") > strings.Index(text, "future-task") {
		t.Fatalf("output=%s err=%v", text, err)
	}
}

package server

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/runcontrol"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

type cancelSchedulerStub struct {
	result domain.RunCancellation
	err    error
	calls  int
	id     string
	after  func()
}

func (*cancelSchedulerStub) Reload()             {}
func (*cancelSchedulerStub) RunNow(string) error { return nil }
func (s *cancelSchedulerStub) CancelRun(ctx context.Context, id string) (domain.RunCancellation, error) {
	if err := ctx.Err(); err != nil {
		return domain.RunCancellation{}, err
	}
	s.calls++
	s.id = id
	if s.after != nil {
		s.after()
	}
	return s.result, s.err
}
func TestRunCancellationAPIClassifiesAcknowledgements(t *testing.T) {
	for _, tc := range []struct {
		name   string
		result domain.RunCancellation
		err    error
		code   int
	}{
		{name: "requested", result: domain.RunCancellation{RunID: "run-owned", TaskID: "task-owned", Status: "requested", CancellationStatus: "requested"}, code: http.StatusAccepted},
		{name: "finished", result: domain.RunCancellation{RunID: "run-owned", TaskID: "task-owned", Status: "already_finished", Outcome: domain.OutcomeSuccess}, code: http.StatusOK},
		{name: "missing", err: store.ErrNotFound, code: http.StatusNotFound},
		{name: "nonactive", err: runcontrol.ErrNotActive, code: http.StatusConflict},
		{name: "failure", err: errors.New("private-backend-error"), code: http.StatusInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestServer(t)
			stub := &cancelSchedulerStub{result: tc.result, err: tc.err}
			s.sched = stub
			out := doJSON(t, s, http.MethodPost, "/v1/runs/run-owned/cancel", nil)
			if out.Code != tc.code {
				t.Fatalf("status=%d want=%d body=%s", out.Code, tc.code, out.Body.String())
			}
			if bytes.Contains(out.Body.Bytes(), []byte("private-backend-error")) {
				t.Fatalf("private backend detail leaked: %s", out.Body.String())
			}
			if stub.calls != 1 || stub.id != "run-owned" {
				t.Fatalf("calls=%d id=%q", stub.calls, stub.id)
			}
			if tc.err == nil {
				var got domain.RunCancellation
				if err := json.Unmarshal(out.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if got != tc.result {
					t.Fatalf("result=%+v want=%+v", got, tc.result)
				}
			}
		})
	}
}
func TestRunCancellationAPIEnforcesAuthorityTargetAndAudit(t *testing.T) {
	s := newTestServer(t)
	stub := &cancelSchedulerStub{result: domain.RunCancellation{RunID: "run-owned", TaskID: "task-owned", Status: "requested", CancellationStatus: "requested"}}
	s.sched = stub
	observer, err := s.store.CreateActor(domain.ActorKindJSON, "observer", domain.CapabilityObserve, nil)
	if err != nil {
		t.Fatal(err)
	}
	operator, err := s.store.CreateActor(domain.ActorKindJSON, "operator", domain.CapabilityOperate, nil)
	if err != nil {
		t.Fatal(err)
	}
	actorID := observer.ID
	s.SetActorResolver(func(*http.Request) (string, error) { return actorID, nil })
	denied := doJSON(t, s, http.MethodPost, "/v1/runs/run-owned/cancel", nil)
	if denied.Code != http.StatusForbidden || stub.calls != 0 {
		t.Fatalf("denied=%d calls=%d", denied.Code, stub.calls)
	}
	actorID = operator.ID
	wrong := httptest.NewRequest(http.MethodPost, "/v1/runs/run-owned/cancel", nil)
	wrong.Header.Set(ExpectedDaemonHeader, "wrong-daemon")
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, wrong)
	if rec.Code != http.StatusConflict || stub.calls != 0 {
		t.Fatalf("wrong target=%d calls=%d", rec.Code, stub.calls)
	}
	accepted := doJSON(t, s, http.MethodPost, "/v1/runs/run-owned/cancel", nil)
	if accepted.Code != http.StatusAccepted || stub.calls != 1 {
		t.Fatalf("accepted=%d calls=%d body=%s", accepted.Code, stub.calls, accepted.Body.String())
	}
	if _, err := s.store.RevokeActor(operator.ID); err != nil {
		t.Fatal(err)
	}
	revoked := doJSON(t, s, http.MethodPost, "/v1/runs/run-owned/cancel", nil)
	if revoked.Code != http.StatusForbidden || stub.calls != 1 {
		t.Fatalf("revoked=%d calls=%d", revoked.Code, stub.calls)
	}
	audit, err := s.store.ListAudit(domain.AuditQuery{Operation: "runs.cancel", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(audit) != 4 {
		t.Fatalf("audit=%+v", audit)
	}
	results := map[domain.AuditResult]int{}
	for _, a := range audit {
		if a.TargetID != "run-owned" || a.TargetKind != "run" {
			t.Fatalf("audit target=%+v", a)
		}
		results[a.Result]++
	}
	if results[domain.AuditResultSucceeded] != 1 || results[domain.AuditResultDenied] != 2 || results[domain.AuditResultFailed] != 1 {
		t.Fatalf("audit results=%v", results)
	}
}
func TestRunCancellationAuditCompletesAfterCallerDisconnect(t *testing.T) {
	s := newTestServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stub := &cancelSchedulerStub{result: domain.RunCancellation{RunID: "run-owned", TaskID: "task-owned", Status: "requested"}, after: cancel}
	s.sched = stub
	out := httptest.NewRecorder()
	s.Handler().ServeHTTP(out, httptest.NewRequest(http.MethodPost, "/v1/runs/run-owned/cancel", nil).WithContext(ctx))
	if out.Code != http.StatusAccepted {
		t.Fatalf("status=%d body=%s", out.Code, out.Body.String())
	}
	audit, err := s.store.ListAudit(domain.AuditQuery{Operation: "runs.cancel", Limit: 10})
	if err != nil || len(audit) != 1 || audit[0].Result != domain.AuditResultSucceeded {
		t.Fatalf("audit=%+v err=%v", audit, err)
	}
}
func TestRunCancellationManifestRequiresRuntimeSupport(t *testing.T) {
	s := newTestServer(t)
	for _, supported := range []bool{false, true} {
		if supported {
			s.sched = &cancelSchedulerStub{}
		} else {
			s.sched = nil
		}
		out := doJSON(t, s, http.MethodGet, "/v1/manifest", nil)
		var m ManifestResponse
		if err := json.Unmarshal(out.Body.Bytes(), &m); err != nil {
			t.Fatal(err)
		}
		has := false
		for _, capability := range m.Capabilities {
			has = has || capability == "run-cancellation"
		}
		if has != supported {
			t.Fatalf("supported=%v capabilities=%v", supported, m.Capabilities)
		}
	}
}

func TestRunCancellationRejectsCommandAndPIDOverrides(t *testing.T) {
	s := newTestServer(t)
	stub := &cancelSchedulerStub{}
	s.sched = stub
	out := httptest.NewRecorder()
	s.Handler().ServeHTTP(out, httptest.NewRequest(http.MethodPost, "/v1/runs/run-owned/cancel", bytes.NewBufferString("{\"pid\":123,\"command\":\"other\"}")))
	if out.Code != http.StatusBadRequest || stub.calls != 0 {
		t.Fatalf("status=%d calls=%d body=%s", out.Code, stub.calls, out.Body.String())
	}
}
func TestRunCancellationUnavailableAuditPreventsMutation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "current.db")
	st, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec("DROP TABLE audit_events"); err != nil {
		t.Fatal(err)
	}
	stub := &cancelSchedulerStub{}
	s := New(st, stub, nil, nil, "", slog.New(slog.NewTextHandler(io.Discard, nil)))
	out := doJSON(t, s, http.MethodPost, "/v1/runs/run-owned/cancel", nil)
	if out.Code != http.StatusServiceUnavailable || stub.calls != 0 {
		t.Fatalf("status=%d calls=%d body=%s", out.Code, stub.calls, out.Body.String())
	}
}

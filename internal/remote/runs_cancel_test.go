package remote

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apis "github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/enrollment"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

type remoteCancelScheduler struct{ calls int }

func (*remoteCancelScheduler) Reload()             {}
func (*remoteCancelScheduler) RunNow(string) error { return nil }
func (s *remoteCancelScheduler) CancelRun(ctx context.Context, id string) (domain.RunCancellation, error) {
	if err := ctx.Err(); err != nil {
		return domain.RunCancellation{}, err
	}
	s.calls++
	if id == "missing" {
		return domain.RunCancellation{}, store.ErrNotFound
	}
	return domain.RunCancellation{RunID: id, TaskID: "task-owned", Status: "requested", CancellationStatus: "requested"}, nil
}
func TestRemoteRunCancellationUsesExistingAuthorityAndExactTarget(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	local, err := st.LocalActor()
	if err != nil {
		t.Fatal(err)
	}
	scheduler := &remoteCancelScheduler{}
	api := apis.New(st, scheduler, nil, nil, "", nil)
	api.SetActorResolver(ActorID)
	service := enrollment.New(st)
	handler := NewHandler(api.Handler(), service, local.ID)
	issue := func(capability domain.Capability) domain.IssuedCredential {
		t.Helper()
		pair, err := service.Create("remote cancellation", domain.ActorKindJSON, capability)
		if err != nil {
			t.Fatal(err)
		}
		issued, err := service.Exchange(pair.ID, pair.Phrase, pair.DaemonID, pair.DisplayName, pair.Kind, pair.Capability)
		if err != nil {
			t.Fatal(err)
		}
		return issued
	}
	operator := issue(domain.CapabilityOperate)
	observer := issue(domain.CapabilityObserve)
	request := func(id, token, daemon string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/runs/"+id+"/cancel", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		if daemon != "" {
			req.Header.Set(apis.ExpectedDaemonHeader, daemon)
		}
		out := httptest.NewRecorder()
		handler.ServeHTTP(out, req)
		return out
	}
	denied := request("run-owned", observer.Token, "")
	if denied.Code != http.StatusForbidden || scheduler.calls != 0 {
		t.Fatalf("observer=%d calls=%d body=%s", denied.Code, scheduler.calls, denied.Body.String())
	}
	wrong := request("run-owned", operator.Token, "wrong-daemon")
	if wrong.Code != http.StatusConflict || scheduler.calls != 0 {
		t.Fatalf("identity=%d calls=%d", wrong.Code, scheduler.calls)
	}
	accepted := request("run-owned", operator.Token, "")
	if accepted.Code != http.StatusAccepted || scheduler.calls != 1 {
		t.Fatalf("operator=%d calls=%d body=%s", accepted.Code, scheduler.calls, accepted.Body.String())
	}
	var result domain.RunCancellation
	if err := json.Unmarshal(accepted.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.RunID != "run-owned" || result.TaskID != "task-owned" || result.Status != "requested" || result.Outcome != "" {
		t.Fatalf("acknowledgement=%+v", result)
	}
	missing := request("missing", operator.Token, "")
	if missing.Code != http.StatusNotFound || scheduler.calls != 2 {
		t.Fatalf("missing=%d calls=%d", missing.Code, scheduler.calls)
	}
	if _, err := service.Revoke(operator.ID); err != nil {
		t.Fatal(err)
	}
	revoked := request("run-owned", operator.Token, "")
	if revoked.Code != http.StatusUnauthorized || scheduler.calls != 2 {
		t.Fatalf("revoked=%d calls=%d body=%s", revoked.Code, scheduler.calls, revoked.Body.String())
	}
	for _, operation := range Operations() {
		if operation.ID == "runs.cancel" {
			if operation.Retry != RetryNever || operation.Capability != domain.CapabilityOperate || operation.TargetKind != "run" {
				t.Fatalf("operation=%+v", operation)
			}
			return
		}
	}
	t.Fatal("cancellation operation missing")
}

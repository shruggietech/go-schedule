package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestAuditCompletionSurvivesCallerCancellationAfterMutation(t *testing.T) {
	s := newTestServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("POST /v1/groups", func(w http.ResponseWriter, r *http.Request) {
		group := domain.Group{Name: "Committed before disconnect", Enabled: true}
		if err := s.store.CreateGroup(&group); err != nil {
			t.Error(err)
			w.WriteHeader(500)
			return
		}
		cancel()
		writeJSON(w, http.StatusCreated, group)
	})
	out := httptest.NewRecorder()
	s.Handler().ServeHTTP(out, httptest.NewRequest(http.MethodPost, "/v1/groups", nil).WithContext(ctx))
	if out.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", out.Code, out.Body.String())
	}
	events, err := s.store.ListAudit(domain.AuditQuery{Operation: "groups.create", Limit: 10})
	if err != nil || len(events) != 1 || events[0].Result != domain.AuditResultSucceeded || events[0].TargetID == "" {
		t.Fatalf("completion lost after disconnect: %+v %v", events, err)
	}
	if _, err := s.store.GetGroup(events[0].TargetID); err != nil {
		t.Fatal(err)
	}
}

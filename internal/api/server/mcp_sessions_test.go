package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/store"
)

func createMCPSessionForTest(t *testing.T, s *Server, capability domain.Capability) MCPSessionCredentialResponse {
	t.Helper()
	body, err := json.Marshal(MCPSessionCreateRequest{ClientName: "Codex test", Capability: capability})
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	s.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/v1/mcp/sessions", bytes.NewReader(body)))
	if response.Code != http.StatusCreated {
		t.Fatalf("create session status=%d body=%s", response.Code, response.Body.String())
	}
	var session MCPSessionCredentialResponse
	if err := json.Unmarshal(response.Body.Bytes(), &session); err != nil || session.Credential == "" || session.ActorID == "" {
		t.Fatalf("session=%+v err=%v", session, err)
	}
	return session
}

func performMCPOperation(t *testing.T, s *Server, method, path, credential, daemonID string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, bytes.NewReader(body))
	request.Header.Set(MCPSessionHeader, credential)
	if daemonID != "" {
		request.Header.Set(ExpectedDaemonHeader, daemonID)
	}
	response := httptest.NewRecorder()
	s.Handler().ServeHTTP(response, request)
	return response
}

func TestOperateSessionRunsAndTogglesExistingTaskWithAttributedAudit(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{Name: "MCP task", Command: "echo", Schedule: "every day at 9:00 AM", Timezone: "UTC"})
	identity, err := s.store.DaemonIdentity()
	if err != nil {
		t.Fatal(err)
	}
	session := createMCPSessionForTest(t, s, domain.CapabilityOperate)
	operations := []struct {
		path      string
		operation string
		status    int
	}{
		{path: "/v1/tasks/" + task.Task.ID + "/run-now", operation: "tasks.run_now", status: http.StatusAccepted},
		{path: "/v1/tasks/" + task.Task.ID + "/disable", operation: "tasks.disable", status: http.StatusNoContent},
		{path: "/v1/tasks/" + task.Task.ID + "/enable", operation: "tasks.enable", status: http.StatusNoContent},
	}
	for _, operation := range operations {
		response := performMCPOperation(t, s, http.MethodPost, operation.path, session.Credential, identity.InstallationID, nil)
		if response.Code != operation.status {
			t.Fatalf("%s status=%d body=%s", operation.operation, response.Code, response.Body.String())
		}
		events, err := s.store.ListAudit(domain.AuditQuery{ActorID: session.ActorID, Operation: operation.operation, Limit: 10})
		if err != nil || len(events) != 1 || events[0].TargetID != task.Task.ID || events[0].DaemonID != identity.InstallationID || events[0].Result != domain.AuditResultSucceeded {
			t.Fatalf("%s events=%+v err=%v", operation.operation, events, err)
		}
	}
}

func TestObserveAndRevokedSessionsAreDeniedAndAttributed(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{Name: "MCP task", Command: "echo"})
	identity, _ := s.store.DaemonIdentity()
	observe := createMCPSessionForTest(t, s, domain.CapabilityObserve)
	denied := performMCPOperation(t, s, http.MethodPost, "/v1/tasks/"+task.Task.ID+"/run-now", observe.Credential, identity.InstallationID, nil)
	if denied.Code != http.StatusForbidden {
		t.Fatalf("observe status=%d body=%s", denied.Code, denied.Body.String())
	}
	operate := createMCPSessionForTest(t, s, domain.CapabilityOperate)
	revoke := httptest.NewRecorder()
	s.Handler().ServeHTTP(revoke, httptest.NewRequest(http.MethodDelete, "/v1/mcp/sessions/"+operate.ID, nil))
	if revoke.Code != http.StatusNoContent {
		t.Fatalf("revoke status=%d body=%s", revoke.Code, revoke.Body.String())
	}
	revoked := performMCPOperation(t, s, http.MethodPost, "/v1/tasks/"+task.Task.ID+"/run-now", operate.Credential, identity.InstallationID, nil)
	if revoked.Code != http.StatusForbidden {
		t.Fatalf("revoked status=%d body=%s", revoked.Code, revoked.Body.String())
	}
	for _, actorID := range []string{observe.ActorID, operate.ActorID} {
		events, err := s.store.ListAudit(domain.AuditQuery{ActorID: actorID, Operation: "tasks.run_now", Result: domain.AuditResultDenied, Limit: 10})
		if err != nil || len(events) != 1 || events[0].TargetID != task.Task.ID {
			t.Fatalf("actor=%s events=%+v err=%v", actorID, events, err)
		}
	}
}

func TestExpiredAndMalformedSessionsFailClosed(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{Name: "MCP task", Command: "echo"})
	identity, _ := s.store.DaemonIdentity()
	session := createMCPSessionForTest(t, s, domain.CapabilityOperate)
	expired := time.Now().UTC().Add(-time.Minute)
	expiresAt := &expired
	if _, err := s.store.UpdateActor(session.ActorID, store.ActorUpdate{ExpiresAt: &expiresAt}); err != nil {
		t.Fatal(err)
	}
	response := performMCPOperation(t, s, http.MethodPost, "/v1/tasks/"+task.Task.ID+"/run-now", session.Credential, identity.InstallationID, nil)
	if response.Code != http.StatusForbidden {
		t.Fatalf("expired status=%d body=%s", response.Code, response.Body.String())
	}
	malformed := performMCPOperation(t, s, http.MethodPost, "/v1/tasks/"+task.Task.ID+"/run-now", "malformed", identity.InstallationID, nil)
	if malformed.Code != http.StatusForbidden {
		t.Fatalf("malformed status=%d body=%s", malformed.Code, malformed.Body.String())
	}
}

func TestOperateSessionRejectsWrongDaemonAndManageOperations(t *testing.T) {
	s := newTestServer(t)
	task := newTaskFor(t, s, TaskCreateRequest{Name: "MCP task", Command: "echo"})
	session := createMCPSessionForTest(t, s, domain.CapabilityOperate)
	wrongDaemon := performMCPOperation(t, s, http.MethodPost, "/v1/tasks/"+task.Task.ID+"/run-now", session.Credential, uuid.NewString(), nil)
	if wrongDaemon.Code != http.StatusConflict {
		t.Fatalf("wrong daemon status=%d body=%s", wrongDaemon.Code, wrongDaemon.Body.String())
	}
	events, err := s.store.ListAudit(domain.AuditQuery{ActorID: session.ActorID, Operation: "tasks.run_now", Result: domain.AuditResultFailed, Limit: 10})
	if err != nil || len(events) != 1 || events[0].TargetID != task.Task.ID {
		t.Fatalf("wrong daemon events=%+v err=%v", events, err)
	}
	for _, request := range []struct {
		method string
		path   string
		body   []byte
	}{
		{method: http.MethodPatch, path: "/v1/tasks/" + task.Task.ID, body: []byte(`{"name":"changed"}`)},
		{method: http.MethodDelete, path: "/v1/tasks/" + task.Task.ID},
	} {
		response := performMCPOperation(t, s, request.method, request.path, session.Credential, "", request.body)
		if response.Code != http.StatusForbidden {
			t.Fatalf("%s %s status=%d body=%s", request.method, request.path, response.Code, response.Body.String())
		}
	}
}

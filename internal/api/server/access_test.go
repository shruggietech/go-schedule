package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"runtime"
	"testing"

	"github.com/shruggietech/go-schedule/internal/authorization"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestOperationCatalogCoversEveryRegisteredManagementRoute(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate test source")
	}
	source, err := os.ReadFile(filename[:len(filename)-len("access_test.go")] + "server.go")
	if err != nil {
		t.Fatal(err)
	}
	matches := regexp.MustCompile(`HandleFunc\("([A-Z]+) ([^"]+)"`).FindAllSubmatch(source, -1)
	registered := map[string]bool{}
	for _, match := range matches {
		registered[string(match[1])+" "+string(match[2])] = true
	}
	cataloged := map[string]bool{}
	for _, operation := range authorization.Catalog() {
		cataloged[operation.Method+" "+operation.Pattern] = true
	}
	for route := range registered {
		if !cataloged[route] {
			t.Errorf("registered route absent from catalog: %s", route)
		}
	}
	for route := range cataloged {
		if !registered[route] {
			t.Errorf("catalog route is not registered: %s", route)
		}
	}
}

func TestAuthorizationReloadsRevokedActorAndAuditsDenial(t *testing.T) {
	s := newTestServer(t)
	actor, err := s.store.CreateActor(domain.ActorKindCLI, "Automation", domain.CapabilityObserve, nil)
	if err != nil {
		t.Fatal(err)
	}
	s.SetActorResolver(func(*http.Request) (string, error) { return actor.ID, nil })
	handler := s.Handler()
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
	if first.Code != http.StatusOK {
		t.Fatalf("first status=%d", first.Code)
	}
	if _, err := s.store.RevokeActor(actor.ID); err != nil {
		t.Fatal(err)
	}
	second := httptest.NewRecorder()
	handler.ServeHTTP(second, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
	if second.Code != http.StatusForbidden {
		t.Fatalf("second status=%d body=%s", second.Code, second.Body.String())
	}
	events, err := s.store.ListAudit(domain.AuditQuery{ActorID: actor.ID, Result: domain.AuditResultDenied, Limit: 10})
	if err != nil || len(events) != 1 || events[0].Operation != "health.read" {
		t.Fatalf("events=%+v err=%v", events, err)
	}
}

func TestActorAPIUsesIntentFirstAuditAndProtectsLocalActor(t *testing.T) {
	s := newTestServer(t)
	body := bytes.NewBufferString(`{"kind":"cli","display_name":"Release operator","capability":"manage"}`)
	created := httptest.NewRecorder()
	s.Handler().ServeHTTP(created, httptest.NewRequest(http.MethodPost, "/v1/access/actors", body))
	if created.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", created.Code, created.Body.String())
	}
	var actor domain.Actor
	if err := json.Unmarshal(created.Body.Bytes(), &actor); err != nil {
		t.Fatal(err)
	}
	events, err := s.store.ListAudit(domain.AuditQuery{Operation: "actors.create", Limit: 10})
	if err != nil || len(events) != 1 || events[0].Result != domain.AuditResultSucceeded {
		t.Fatalf("events=%+v err=%v", events, err)
	}
	local, _ := s.store.LocalActor()
	revoke := httptest.NewRecorder()
	s.Handler().ServeHTTP(revoke, httptest.NewRequest(http.MethodPost, "/v1/access/actors/"+local.ID+"/revoke", nil))
	if revoke.Code != http.StatusConflict {
		t.Fatalf("revoke status=%d body=%s", revoke.Code, revoke.Body.String())
	}
}

func TestCurrentActorReportsServerOwnedAuthorityToObserveClients(t *testing.T) {
	s := newTestServer(t)
	actor, err := s.store.CreateActor(domain.ActorKindDesktop, "Remote desktop", domain.CapabilityOperate, nil)
	if err != nil {
		t.Fatal(err)
	}
	s.SetActorResolver(func(*http.Request) (string, error) { return actor.ID, nil })
	response := httptest.NewRecorder()
	s.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/access/current", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var current domain.Actor
	if err := json.Unmarshal(response.Body.Bytes(), &current); err != nil || current.ID != actor.ID || current.Capability != domain.CapabilityOperate {
		t.Fatalf("actor=%+v err=%v", current, err)
	}
}

func TestAuditExportIsNDJSONAndDoesNotExposePayloads(t *testing.T) {
	s := newTestServer(t)
	response := httptest.NewRecorder()
	s.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/audit/export?limit=10", nil))
	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "application/x-ndjson" {
		t.Fatalf("status=%d content-type=%q", response.Code, response.Header().Get("Content-Type"))
	}
	for _, forbidden := range []string{"request_body", "response_body", "authorization", "raw_error", "command", "environment", "stdin", "path", "credential", "secret"} {
		if bytes.Contains(bytes.ToLower(response.Body.Bytes()), []byte(forbidden)) {
			t.Fatalf("export contains prohibited field %q: %s", forbidden, response.Body.String())
		}
	}
}

func TestUnavailableAuditStorageBlocksProtectedOperation(t *testing.T) {
	s := newTestServer(t)
	if err := s.store.Close(); err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	s.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/v1/access/actors", bytes.NewBufferString(`{"kind":"cli","display_name":"Blocked","capability":"observe"}`)))
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	if response.Body.String() == "" || !bytes.Contains(response.Body.Bytes(), []byte(CodeAuditUnavailable)) {
		t.Fatalf("body=%s", response.Body.String())
	}
}

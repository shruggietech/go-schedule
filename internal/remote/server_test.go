package remote

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	apis "github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
	"github.com/shruggietech/go-schedule/internal/enrollment"
	"github.com/shruggietech/go-schedule/internal/store"
)

func TestRemoteAllowlistAndBearerLifecycle(t *testing.T) {
	st, _ := store.Open(":memory:")
	defer st.Close()
	local, _ := st.LocalActor()
	api := apis.New(st, nil, nil, nil, "", nil)
	api.SetActorResolver(ActorID)
	service := enrollment.New(st)
	handler := NewHandler(api.Handler(), service, local.ID)
	public := httptest.NewRecorder()
	handler.ServeHTTP(public, httptest.NewRequest(http.MethodGet, "/api/v1/health", nil))
	if public.Code != http.StatusOK {
		t.Fatalf("health=%d %s", public.Code, public.Body.String())
	}
	excluded := httptest.NewRecorder()
	handler.ServeHTTP(excluded, httptest.NewRequest(http.MethodGet, "/api/v1/runtime-info", nil))
	if excluded.Code != http.StatusNotFound {
		t.Fatalf("excluded=%d", excluded.Code)
	}
	missing := httptest.NewRecorder()
	handler.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil))
	if missing.Code != http.StatusUnauthorized {
		t.Fatalf("missing=%d", missing.Code)
	}
	for name, header := range map[string]string{"wrong scheme": "Basic value", "empty bearer": "Bearer ", "malformed bearer": "Bearer one two"} {
		t.Run(name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
			request.Header.Set("Authorization", header)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusUnauthorized || response.Header().Get("WWW-Authenticate") != "Bearer" {
				t.Fatalf("status=%d header=%q", response.Code, response.Header().Get("WWW-Authenticate"))
			}
		})
	}
	pairing, _ := service.Create("remote JSON", domain.ActorKindJSON, domain.CapabilityObserve)
	body, _ := json.Marshal(map[string]any{"pairing_id": pairing.ID, "phrase": pairing.Phrase, "daemon_id": pairing.DaemonID, "display_name": pairing.DisplayName, "kind": pairing.Kind, "capability": pairing.Capability})
	enrollReq := httptest.NewRequest(http.MethodPost, "/api/v1/enroll", bytes.NewReader(body))
	enrollReq.Header.Set("Content-Type", "application/json")
	enrolled := httptest.NewRecorder()
	handler.ServeHTTP(enrolled, enrollReq)
	if enrolled.Code != http.StatusCreated {
		t.Fatalf("enroll=%d %s", enrolled.Code, enrolled.Body.String())
	}
	var issued domain.IssuedCredential
	json.NewDecoder(enrolled.Body).Decode(&issued)
	task := &domain.Task{Name: "visible name", Command: "secret-command", Env: map[string]string{"TOKEN": "secret-env"}, Stdin: "secret-stdin", WorkingDir: "secret-directory", RunAs: "secret-user", Timezone: "UTC", OverlapPolicy: domain.OverlapQueueOne, CatchupPolicy: domain.CatchupOne, MissingDatePolicy: domain.MissingDateSkip, TimeBasis: domain.TimeBasisWallClock, DSTGapPolicy: domain.DSTGapNextValid, DSTOverlapPolicy: domain.DSTOverlapFirst, State: domain.TaskActive}
	if err := st.CreateTask(task); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	request.Header.Set("Authorization", "Bearer "+issued.Token)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("tasks=%d %s", response.Code, response.Body.String())
	}
	for _, secret := range []string{"secret-command", "secret-env", "secret-stdin", "secret-directory", "secret-user"} {
		if bytes.Contains(response.Body.Bytes(), []byte(secret)) {
			t.Fatalf("task list leaked %q: %s", secret, response.Body.String())
		}
	}
	detailRequest := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+task.ID, nil)
	detailRequest.Header.Set("Authorization", "Bearer "+issued.Token)
	detail := httptest.NewRecorder()
	handler.ServeHTTP(detail, detailRequest)
	if detail.Code != http.StatusOK || !bytes.Contains(detail.Body.Bytes(), []byte("visible name")) {
		t.Fatalf("task detail=%d %s", detail.Code, detail.Body.String())
	}
	for _, secret := range []string{"secret-command", "secret-env", "secret-stdin", "secret-directory", "secret-user"} {
		if bytes.Contains(detail.Body.Bytes(), []byte(secret)) {
			t.Fatalf("task detail leaked %q: %s", secret, detail.Body.String())
		}
	}
	deniedRequest := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString(`{}`))
	deniedRequest.Header.Set("Content-Type", "application/json; charset=utf-8")
	deniedRequest.Header.Set("Authorization", "Bearer "+issued.Token)
	denied := httptest.NewRecorder()
	handler.ServeHTTP(denied, deniedRequest)
	if denied.Code != http.StatusForbidden {
		t.Fatalf("underprivileged mutation=%d %s", denied.Code, denied.Body.String())
	}
	origin := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	origin.Header.Set("Origin", "https://example.test")
	blocked := httptest.NewRecorder()
	handler.ServeHTTP(blocked, origin)
	if blocked.Code != http.StatusForbidden {
		t.Fatalf("origin=%d", blocked.Code)
	}
}

func TestOversizedBodyIsRejectedAtRemoteBoundary(t *testing.T) {
	handler := NewHandler(http.NotFoundHandler(), nil, "public")
	request := httptest.NewRequest(http.MethodPost, "/api/v1/enroll", bytes.NewReader(bytes.Repeat([]byte("x"), 4097)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusRequestEntityTooLarge || !bytes.Contains(response.Body.Bytes(), []byte(`"request_too_large"`)) {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestRemoteStreamUsesObservationProjectionOutsideRequestCapacity(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	local, _ := st.LocalActor()
	service := enrollment.New(st)
	pairing, err := service.Create("observer", domain.ActorKindJSON, domain.CapabilityObserve)
	if err != nil {
		t.Fatal(err)
	}
	issued, err := service.Exchange(pairing.ID, pairing.Phrase, pairing.DaemonID, pairing.DisplayName, pairing.Kind, pairing.Capability)
	if err != nil {
		t.Fatal(err)
	}
	projected := false
	handler := NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projected = r.URL.Query().Get("observation") == "true"
		w.WriteHeader(http.StatusNoContent)
	}), service, local.ID)
	for range cap(handler.concurrent) {
		handler.concurrent <- struct{}{}
	}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/events", nil)
	request.Header.Set("Authorization", "Bearer "+issued.Token)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent || !projected {
		t.Fatalf("status=%d projected=%v body=%s", response.Code, projected, response.Body.String())
	}
}

func TestEveryOperationHasUniqueDocumentedIdentity(t *testing.T) {
	contract, err := os.ReadFile("../../api/openapi/remote-v1.yaml")
	if err != nil {
		t.Fatal(err)
	}
	documented := map[string]bool{}
	for _, match := range regexp.MustCompile(`operationId: ([a-z0-9._]+)`).FindAllStringSubmatch(string(contract), -1) {
		documented[match[1]] = true
	}
	seen := map[string]bool{}
	for _, operation := range Operations() {
		if operation.ID == "" || operation.RemotePath == "" || seen[operation.Method+" "+operation.RemotePath] {
			t.Fatalf("invalid operation %+v", operation)
		}
		if operation.TargetKind == "" || operation.Audit == "" || len(operation.SecretExclusions) == 0 {
			t.Errorf("operation %s lacks security classification", operation.ID)
		}
		seen[operation.Method+" "+operation.RemotePath] = true
		if !documented[operation.ID] {
			t.Errorf("operation %s is missing from OpenAPI", operation.ID)
		}
	}
	if len(documented) != len(seen) {
		t.Fatalf("OpenAPI operations=%d runtime operations=%d", len(documented), len(seen))
	}
}

func TestTLSConfigurationRejectsLegacyProtocols(t *testing.T) {
	if got := TLSConfig(tls.Certificate{}).MinVersion; got != tls.VersionTLS13 {
		t.Fatalf("minimum TLS version=%x", got)
	}
}

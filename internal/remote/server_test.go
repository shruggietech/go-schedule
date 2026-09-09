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
	request := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	request.Header.Set("Authorization", "Bearer "+issued.Token)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("tasks=%d %s", response.Code, response.Body.String())
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

package remote

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apis "github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/authorization"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
	"github.com/shruggietech/go-schedule/v2/internal/enrollment"
	"github.com/shruggietech/go-schedule/v2/internal/store"
)

func TestS144RemoteAutomationOperations(t *testing.T) {
	for _, resource := range []string{"chains", "joins"} {
		for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete} {
			route := "/api/v1/" + resource
			if method == http.MethodPatch || method == http.MethodDelete {
				route += "/owned"
			}
			op, ok := lookup(method, route)
			if !ok {
				t.Fatalf("missing %s %s", method, route)
			}
			if method == http.MethodGet {
				if op.Retry != RetrySafe || op.Capability != domain.CapabilityObserve {
					t.Fatalf("read=%+v", op)
				}
			} else if op.Retry != RetryNever || op.Capability != domain.CapabilityManage || op.Audit != authorization.AuditMutation {
				t.Fatalf("mutation=%+v", op)
			}
		}
	}
	for _, route := range []string{"/api/v1/join-instances", "/api/v1/join-instances/exact"} {
		op, ok := lookup(http.MethodGet, route)
		if !ok || op.Capability != domain.CapabilityObserve || op.Retry != RetrySafe {
			t.Fatalf("instance route %s=%+v found=%v", route, op, ok)
		}
	}
}

func TestS144RemoteJoinCRUDHonorsAuthorityAndDaemon(t *testing.T) {
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	local, err := st.LocalActor()
	if err != nil {
		t.Fatal(err)
	}
	tasks := make([]domain.Task, 3)
	for i := range tasks {
		tasks[i] = domain.Task{Name: "task", Timezone: "UTC", State: domain.TaskActive}
		if err := st.CreateTask(&tasks[i]); err != nil {
			t.Fatal(err)
		}
	}
	api := apis.New(st, nil, nil, nil, "", nil)
	api.SetActorResolver(ActorID)
	service := enrollment.New(st)
	handler := NewHandler(api.Handler(), service, local.ID)
	issue := func(capability domain.Capability) domain.IssuedCredential {
		t.Helper()
		pair, err := service.Create("remote joins", domain.ActorKindJSON, capability)
		if err != nil {
			t.Fatal(err)
		}
		credential, err := service.Exchange(pair.ID, pair.Phrase, pair.DaemonID, pair.DisplayName, pair.Kind, pair.Capability)
		if err != nil {
			t.Fatal(err)
		}
		return credential
	}
	manager, observer := issue(domain.CapabilityManage), issue(domain.CapabilityObserve)
	request := func(method, route, token, daemon string, body any) *httptest.ResponseRecorder {
		t.Helper()
		var data []byte
		if body != nil {
			var err error
			data, err = json.Marshal(body)
			if err != nil {
				t.Fatal(err)
			}
		}
		req := httptest.NewRequest(method, route, bytes.NewReader(data))
		req.Header.Set("Authorization", "Bearer "+token)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if daemon != "" {
			req.Header.Set(apis.ExpectedDaemonHeader, daemon)
		}
		out := httptest.NewRecorder()
		handler.ServeHTTP(out, req)
		return out
	}
	body := apis.JoinCreateRequest{Name: strings.Repeat("N", 160000), EntryTaskID: tasks[0].ID, TargetTaskID: tasks[2].ID, Branches: []domain.ParallelJoinBranch{{ID: "a", Name: strings.Repeat("B", 140000), TaskID: tasks[1].ID}}}
	if denied := request(http.MethodPost, "/api/v1/joins", observer.Token, "", body); denied.Code != 403 {
		t.Fatalf("observe mutation=%d %s", denied.Code, denied.Body.String())
	}
	if wrong := request(http.MethodPost, "/api/v1/joins", manager.Token, "wrong", body); wrong.Code != 409 {
		t.Fatalf("wrong daemon=%d %s", wrong.Code, wrong.Body.String())
	}
	created := request(http.MethodPost, "/api/v1/joins", manager.Token, "", body)
	var join domain.ParallelJoin
	if created.Code != 201 || json.Unmarshal(created.Body.Bytes(), &join) != nil {
		t.Fatalf("create=%d %s", created.Code, created.Body.String())
	}
	for _, route := range []string{"/api/v1/joins", "/api/v1/joins/" + join.ID, "/api/v1/join-instances?join_id=" + join.ID, "/api/v1/chains"} {
		if read := request(http.MethodGet, route, observer.Token, "", nil); read.Code != 200 {
			t.Fatalf("read %s=%d %s", route, read.Code, read.Body.String())
		}
	}
	definitionPath := "/api/v1/joins/" + join.ID + "/definition"
	if denied := request(http.MethodGet, definitionPath, observer.Token, "", nil); denied.Code != http.StatusForbidden {
		t.Fatalf("observer authoring=%d", denied.Code)
	}
	definition := request(http.MethodGet, definitionPath, manager.Token, "", nil)
	var editable domain.ParallelJoin
	if definition.Code != http.StatusOK || json.Unmarshal(definition.Body.Bytes(), &editable) != nil || editable.Name != body.Name || len(editable.Branches) != 1 || editable.Branches[0].Name != body.Branches[0].Name {
		t.Fatalf("remote authoring text status=%d name length=%d", definition.Code, len(editable.Name))
	}
	if operation, ok := lookup(http.MethodGet, definitionPath); !ok || operation.Capability != domain.CapabilityManage || operation.Retry != RetrySafe || operation.Audit != authorization.AuditPrivilegedRead {
		t.Fatalf("authoring operation=%+v found=%v", operation, ok)
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	wrongMIME := httptest.NewRequest(http.MethodPost, "/api/v1/joins", bytes.NewReader(rawBody))
	wrongMIME.Header.Set("Authorization", "Bearer "+manager.Token)
	wrongMIME.Header.Set("Content-Type", "text/plain")
	mimeResponse := httptest.NewRecorder()
	handler.ServeHTTP(mimeResponse, wrongMIME)
	if mimeResponse.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("join MIME=%d", mimeResponse.Code)
	}
	oversized := body
	oversized.Name = strings.Repeat("X", 2<<20)
	if rejected := request(http.MethodPost, "/api/v1/joins", manager.Token, "", oversized); rejected.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("remote join oversize=%d", rejected.Code)
	}
	for _, raw := range []string{`{"name":"changed","unknown":true}`, `{"name":"changed"} {}`} {
		bad := httptest.NewRequest(http.MethodPatch, "/api/v1/joins/"+join.ID, strings.NewReader(raw))
		bad.Header.Set("Authorization", "Bearer "+manager.Token)
		bad.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, bad)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("remote invalid JSON=%d", response.Code)
		}
	}
	name := "changed"
	if changed := request(http.MethodPatch, "/api/v1/joins/"+join.ID, manager.Token, "", apis.JoinUpdateRequest{Name: &name, ExpectedUpdatedAt: &join.UpdatedAt}); changed.Code != 200 {
		t.Fatalf("patch=%d %s", changed.Code, changed.Body.String())
	}
	if removed := request(http.MethodDelete, "/api/v1/joins/"+join.ID, manager.Token, "", nil); removed.Code != 204 {
		t.Fatalf("delete=%d %s", removed.Code, removed.Body.String())
	}
}

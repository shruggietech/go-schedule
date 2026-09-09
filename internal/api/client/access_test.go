package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/domain"
)

func TestTypedActorAndAuditClientContracts(t *testing.T) {
	var requests []string
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request.Method+" "+request.URL.RequestURI())
		var body string
		switch request.URL.Path {
		case "/v1/access/actors":
			if request.Method == http.MethodGet {
				body = `{"actors":[]}`
			} else {
				body = `{"id":"actor-1"}`
			}
		case "/v1/audit":
			body = `{"events":[]}`
		case "/v1/audit/export":
			body = "{\"id\":\"event-1\"}\n"
		default:
			body = `{"id":"actor-1"}`
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	c := &Client{http: &http.Client{Transport: transport}}
	ctx := context.Background()
	if _, err := c.ListActors(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := c.CreateActor(ctx, server.ActorCreateRequest{Kind: domain.ActorKindCLI, DisplayName: "CLI", Capability: domain.CapabilityObserve}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.UpdateActor(ctx, "actor-1", server.ActorUpdateRequest{}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.RevokeActor(ctx, "actor-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.ListAudit(ctx, domain.AuditQuery{Operation: "tasks.create", Limit: 25}); err != nil {
		t.Fatal(err)
	}
	exported, err := c.ExportAudit(ctx, domain.AuditQuery{Result: domain.AuditResultDenied, Limit: 10})
	if err != nil || string(exported) != "{\"id\":\"event-1\"}\n" {
		t.Fatalf("export=%q err=%v", exported, err)
	}
	if len(requests) != 6 || !strings.Contains(requests[4], "operation=tasks.create") || !strings.Contains(requests[5], "result=denied") {
		t.Fatalf("requests=%v", requests)
	}
}

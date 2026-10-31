package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
)

func TestUpdateTaskTransportsExplicitEnvironmentIntent(t *testing.T) {
	for _, test := range []struct {
		name  string
		req   server.TaskUpdateRequest
		clear bool
		env   bool
	}{
		{name: "omit", req: server.TaskUpdateRequest{Name: "renamed"}},
		{name: "clear", req: server.TaskUpdateRequest{ClearEnv: true}, clear: true},
		{name: "replace", req: server.TaskUpdateRequest{Env: map[string]string{"EXAMPLE": "value"}}, env: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			var wire map[string]any
			c := &Client{http: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				if err := json.NewDecoder(r.Body).Decode(&wire); err != nil {
					t.Fatal(err)
				}
				return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"task":{"id":"fixture"}}`))}, nil
			})}}
			if _, err := c.UpdateTask(context.Background(), "fixture", test.req); err != nil {
				t.Fatal(err)
			}
			_, hasEnv := wire["env"]
			if (wire["clear_env"] == true) != test.clear || hasEnv != test.env {
				t.Fatalf("wire=%v, want clear=%t env=%t", wire, test.clear, test.env)
			}
		})
	}
}

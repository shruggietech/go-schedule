package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func TestUpdateGroupEscapesIdentityAndDecodesGroup(t *testing.T) {
	var method, path string
	client := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		method, path = request.Method, request.URL.EscapedPath()
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"id":"group one","name":"renamed","enabled":true}`))}, nil
	})}}
	group, err := client.UpdateGroup(context.Background(), "group one", server.GroupUpdateRequest{Name: "renamed"})
	if err != nil || method != http.MethodPatch || path != "/v1/groups/group%20one" || group.Name != "renamed" {
		t.Fatalf("method=%q path=%q group=%+v err=%v", method, path, group, err)
	}
}

package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func TestFilesystemWatcherClientUsesEscapedLifecyclePaths(t *testing.T) {
	var method, path string
	client := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		method, path = request.Method, request.URL.EscapedPath()
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"id":"watch one","name":"incoming","kind":"file","debounce":"250ms","stability":"500ms","health":{"state":"active"}}`))}, nil
	})}}
	item, err := client.UpdateFilesystemWatcher(context.Background(), "watch one", server.FilesystemWatcherUpdateRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if method != http.MethodPatch || path != "/v1/filesystem-watchers/watch%20one" || item.ID != "watch one" {
		t.Fatalf("method=%q path=%q item=%+v", method, path, item)
	}
}

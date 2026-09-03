package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRuntimeInfoUsesTypedEndpoint(t *testing.T) {
	client := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodGet || request.URL.Path != "/v1/runtime-info" {
			t.Fatalf("request = %s %s", request.Method, request.URL.Path)
		}
		body := `{"data_dir":"/custom/data","database_path":"/custom/data/goschedule.db","config_path":"/etc/custom.json","log_path":"/logs/events.log","lock_path":"/custom/data/goschedd.lock"}`
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}}

	got, err := client.RuntimeInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got.DataDir != "/custom/data" || got.ConfigPath != "/etc/custom.json" || got.LogPath != "/logs/events.log" {
		t.Fatalf("runtime info = %+v", got)
	}
}

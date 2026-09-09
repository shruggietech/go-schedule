package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

func TestTypedManifestClientLifecycle(t *testing.T) {
	var requests []string
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request.Method+" "+request.URL.Path)
		if request.Method != http.MethodGet {
			var body map[string]string
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil || len(body) != 1 {
				t.Fatalf("request body=%v err=%v", body, err)
			}
		}
		body := `{"installation_id":"daemon-1","display_name":"Workshop","product_version":"1.2.0","local_api_versions":["v1"],"remote_api_versions":[],"operating_mode":"local_only","capabilities":["tasks"],"platform":{"os":"windows","architecture":"amd64"}}`
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})
	c := &Client{http: &http.Client{Transport: transport}}
	ctx := context.Background()
	manifest, err := c.Manifest(ctx)
	if err != nil || manifest.InstallationID != "daemon-1" {
		t.Fatalf("Manifest=%+v err=%v", manifest, err)
	}
	if _, err := c.RenameDaemon(ctx, "Workshop"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.ResetDaemonIdentity(ctx, "daemon-1"); err != nil {
		t.Fatal(err)
	}
	want := "GET /v1/manifest|PATCH /v1/manifest|POST /v1/manifest/reset"
	if strings.Join(requests, "|") != want {
		t.Fatalf("requests=%v", requests)
	}
}

func TestTypedManifestClientPreservesConflictEnvelope(t *testing.T) {
	c := &Client{http: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		body := `{"error":{"code":"conflict","field":"confirm_installation_id","message":"identity changed"}}`
		return &http.Response{StatusCode: http.StatusConflict, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}}
	_, err := c.ResetDaemonIdentity(context.Background(), "old")
	status, ok := err.(*StatusError)
	if !ok || status.Code != server.CodeConflict || status.Field != "confirm_installation_id" {
		t.Fatalf("error=%T %+v", err, err)
	}
}

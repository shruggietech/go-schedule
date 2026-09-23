package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
	"github.com/shruggietech/go-schedule/internal/bundle"
)

func bundleTestResponse(t *testing.T, status int, value any) *http.Response {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(string(data)))}
}

func TestBundleClientValidationReportsMissingTargetRequirementsWithoutForwarding(t *testing.T) {
	requests := []string{}
	manifest := server.ManifestResponse{InstallationID: "daemon-a", Capabilities: []string{"tasks"}, Platform: server.ManifestPlatform{OS: "unknown"}}
	c := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request.Method+" "+request.URL.Path)
		if request.URL.Path != "/v1/manifest" {
			t.Fatalf("incompatible target received %s", request.URL.Path)
		}
		return bundleTestResponse(t, http.StatusOK, manifest), nil
	})}}
	doc := bundle.Document{Schema: bundle.SchemaV2, Tasks: []bundle.Task{{PortableID: "task", Name: "Task", Timezone: "UTC", OverlapPolicy: "queue_one", CatchupPolicy: "one", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"}}, Watchers: []bundle.Watcher{{PortableID: "watcher", Name: "Watcher", TargetTaskID: "task", Kind: "file", Debounce: "250ms", Stability: "500ms"}}}
	result, err := c.ValidateBundle(context.Background(), doc)
	if err != nil || result.Valid || len(result.Issues) != 4 {
		t.Fatalf("validation=%+v err=%v", result, err)
	}
	if _, err := c.PreviewBundle(context.Background(), doc); err == nil || !strings.Contains(err.Error(), "selected daemon") {
		t.Fatalf("preview error=%v", err)
	}
	if _, err := c.ApplyBundle(context.Background(), bundle.Plan{ID: "plan", TargetDaemonID: "daemon-a"}); err == nil || !strings.Contains(err.Error(), "selected daemon") {
		t.Fatalf("apply error=%v", err)
	}
	if !reflect.DeepEqual(requests, []string{"GET /v1/manifest", "GET /v1/manifest", "GET /v1/manifest"}) {
		t.Fatalf("requests=%v", requests)
	}
}

func TestBundleClientFreezesSelectedTargetAcrossDiscoveryAndPreview(t *testing.T) {
	requestsA := []string{}
	requestsB := []string{}
	var switchable *Client
	profile := server.ManifestResponse{InstallationID: "daemon-a", Capabilities: []string{"bundles", "groups"}, Platform: server.ManifestPlatform{OS: "windows"}}
	a := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requestsA = append(requestsA, request.Method+" "+request.URL.Path)
		if request.URL.Path == "/v1/manifest" {
			switchable.Use(&Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				requestsB = append(requestsB, request.Method+" "+request.URL.Path)
				return bundleTestResponse(t, http.StatusInternalServerError, nil), nil
			})}})
			return bundleTestResponse(t, http.StatusOK, profile), nil
		}
		if request.URL.Path != "/v1/bundles/preview" {
			t.Fatalf("unexpected request to A: %s", request.URL.Path)
		}
		return bundleTestResponse(t, http.StatusOK, bundle.Plan{ID: "plan-a", TargetDaemonID: "daemon-a"}), nil
	})}}
	switchable = NewSwitchable(a)
	plan, err := switchable.PreviewBundle(context.Background(), bundle.Document{Schema: bundle.SchemaV1, Groups: []bundle.Group{{PortableID: "group", Name: "Group"}}})
	if err != nil || plan.TargetDaemonID != "daemon-a" {
		t.Fatalf("plan=%+v err=%v", plan, err)
	}
	if !reflect.DeepEqual(requestsA, []string{"GET /v1/manifest", "POST /v1/bundles/preview"}) || len(requestsB) != 0 {
		t.Fatalf("A=%v B=%v", requestsA, requestsB)
	}
}

func TestBundleClientFailsClosedWhenManifestUnavailable(t *testing.T) {
	requests := []string{}
	c := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request.URL.Path)
		return bundleTestResponse(t, http.StatusInternalServerError, server.APIError{}), nil
	})}}
	if _, err := c.ExportBundle(context.Background()); err == nil {
		t.Fatal("export accepted missing manifest")
	}
	if _, err := c.CompareBundle(context.Background(), bundle.Document{Schema: bundle.SchemaV1}); err == nil {
		t.Fatal("compare accepted missing manifest")
	}
	if !reflect.DeepEqual(requests, []string{"/v1/manifest", "/v1/manifest"}) {
		t.Fatalf("requests=%v", requests)
	}
}

func TestBundleClientRejectsDaemonIdentityChangeDuringPreview(t *testing.T) {
	requests := []string{}
	c := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request.URL.Path)
		if request.URL.Path == "/v1/manifest" {
			return bundleTestResponse(t, http.StatusOK, server.ManifestResponse{InstallationID: "before", Capabilities: []string{"bundles"}}), nil
		}
		return bundleTestResponse(t, http.StatusOK, bundle.Plan{ID: "plan", TargetDaemonID: "after"}), nil
	})}}
	plan, err := c.PreviewBundle(context.Background(), bundle.Document{Schema: bundle.SchemaV1})
	if err == nil || plan.ID != "" || !strings.Contains(err.Error(), "changed during preview") {
		t.Fatalf("plan=%+v err=%v", plan, err)
	}
	if !reflect.DeepEqual(requests, []string{"/v1/manifest", "/v1/bundles/preview"}) {
		t.Fatalf("requests=%v", requests)
	}
}

func TestBundleClientRejectsApplyWhenReviewedFamilyBecomesUnavailable(t *testing.T) {
	requests := []string{}
	c := &Client{http: &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requests = append(requests, request.URL.Path)
		if request.URL.Path != "/v1/manifest" {
			t.Fatalf("apply reached incompatible daemon: %s", request.URL.Path)
		}
		return bundleTestResponse(t, http.StatusOK, server.ManifestResponse{InstallationID: "daemon-a", Capabilities: []string{"bundles"}, Platform: server.ManifestPlatform{OS: "linux"}}), nil
	})}}
	_, err := c.ApplyBundle(context.Background(), bundle.Plan{ID: "plan", TargetDaemonID: "daemon-a", Items: []bundle.Item{{Kind: "watcher", Action: bundle.ActionCreate}}})
	if err == nil || !strings.Contains(err.Error(), "watchers") || !reflect.DeepEqual(requests, []string{"/v1/manifest"}) {
		t.Fatalf("err=%v requests=%v", err, requests)
	}
}

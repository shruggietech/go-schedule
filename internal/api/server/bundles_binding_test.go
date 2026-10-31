package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/bundle"
)

func TestBundleUnknownWatcherBindingRejectedBeforeTargetAccess(t *testing.T) {
	for _, endpoint := range []string{"/v1/bundles/preview", "/v1/bundles/compare"} {
		for _, test := range []struct {
			name     string
			document bundle.Document
			paths    map[string]string
			wantID   string
		}{
			{"empty", bundle.Document{Schema: bundle.SchemaV2}, map[string]string{"missing": "/unused"}, "missing"},
			{"nonempty-watchers", watcherBindingTestDocument(), map[string]string{"missing": "/unused"}, "missing"},
			{"non-watcher", bundle.Document{Schema: bundle.SchemaV2, Groups: []bundle.Group{{PortableID: "group", Name: "Group"}}}, map[string]string{"group": "/unused"}, "group"},
			{"deterministic", bundle.Document{Schema: bundle.SchemaV2}, map[string]string{"z-last": "/unused", "a-first": "/unused"}, "a-first"},
		} {
			t.Run(endpoint+"/"+test.name, func(t *testing.T) {
				// A nil target store makes accidental reads during refusal observable.
				defer func() {
					if got := recover(); got != nil {
						t.Errorf("binding refusal accessed target state: %v", got)
					}
				}()
				data, err := json.Marshal(BundleRequest{Bundle: test.document, WatcherPaths: test.paths})
				if err != nil {
					t.Fatal(err)
				}
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(data)))
				(&Server{}).handleBundlePlan(rec, req, endpoint == "/v1/bundles/preview")
				requireWatcherBindingRefusal(t, rec, test.wantID, "not a watcher")
			})
		}
	}
}

func TestBundleDuplicateWatcherBindingsRejectedBeforeMapOverwrite(t *testing.T) {
	for _, endpoint := range []string{"/v1/bundles/validate", "/v1/bundles/preview", "/v1/bundles/compare"} {
		for _, values := range []string{`"watcher":"/first","watcher":"/second"`, `"watcher":"/same","watcher":"/same"`, `"watcher":"/first","\u0077atcher":"/second"`} {
			t.Run(endpoint+"/"+values, func(t *testing.T) {
				defer func() {
					if got := recover(); got != nil {
						t.Errorf("duplicate binding refusal accessed target state: %v", got)
					}
				}()
				body := `{"bundle":{"schema":"go-schedule.bundle/v2"},"watcher_paths":{` + values + `}}`
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(body))
				s := &Server{}
				if endpoint == "/v1/bundles/validate" {
					s.handleBundleValidate(rec, req)
				} else {
					s.handleBundlePlan(rec, req, endpoint == "/v1/bundles/preview")
				}
				requireWatcherBindingRefusal(t, rec, "watcher", "duplicate")
			})
		}
	}
}

func requireWatcherBindingRefusal(t *testing.T, rec *httptest.ResponseRecorder, id, reason string) {
	t.Helper()
	var response APIError
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("status=%d body=%s decode=%v", rec.Code, rec.Body.String(), err)
	}
	if rec.Code != http.StatusBadRequest || response.Error.Code != CodeValidation || response.Error.Field != "watcher_paths" || !strings.Contains(response.Error.Message, id) || !strings.Contains(response.Error.Message, reason) {
		t.Fatalf("status=%d error=%+v; want watcher_paths validation refusal for %q (%s)", rec.Code, response.Error, id, reason)
	}
}

func TestBundleBindingDecoderPreservesValidPaths(t *testing.T) {
	req, err := decodeBundleRequest(strings.NewReader(`{"bundle":{"schema":"go-schedule.bundle/v2"},"watcher_paths":{"watcher":"C:\\owned\\folder=part"}}`))
	if err != nil || len(req.WatcherPaths) != 1 || req.WatcherPaths["watcher"] != `C:\owned\folder=part` {
		t.Fatalf("request=%+v err=%v", req, err)
	}
	for _, body := range []string{`{"bundle":{"schema":"go-schedule.bundle/v2"}}`, `{"bundle":{"schema":"go-schedule.bundle/v2"},"watcher_paths":null}`, `{"bundle":{"schema":"go-schedule.bundle/v2"},"watcher_paths":{}}`} {
		if _, err := decodeBundleRequest(strings.NewReader(body)); err != nil {
			t.Fatalf("body=%s err=%v", body, err)
		}
	}
}

func watcherBindingTestDocument() bundle.Document {
	return bundle.Document{Schema: bundle.SchemaV2, Tasks: []bundle.Task{{PortableID: "task", Name: "Task", Timezone: "UTC", OverlapPolicy: "queue_one", CatchupPolicy: "one", MissingDatePolicy: "skip", TimeBasis: "wall_clock", DSTGapPolicy: "next_valid", DSTOverlapPolicy: "first"}}, Watchers: []bundle.Watcher{{PortableID: "watcher", Name: "Watcher", Kind: "directory", Pattern: "*.txt", Debounce: "0.25s", Stability: "0.5s", TargetTaskID: "task"}}}
}

func TestBundleRepeatedBindingObjectsRefused(t *testing.T) {
	for _, first := range []string{`{"watcher":"/first"}`, `null`, `{}`} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/bundles/preview", strings.NewReader(`{"bundle":{"schema":"go-schedule.bundle/v2"},"watcher_paths":`+first+`,"watcher_paths":{"watcher":"/second"}}`))
		s := newTestServer(t)
		s.handleBundlePreview(rec, req)
		requireWatcherBindingRefusal(t, rec, "watcher_paths", "duplicate")
		if len(s.bundlePlans) != 0 {
			t.Fatal("duplicate binding objects retained a plan")
		}
	}
}

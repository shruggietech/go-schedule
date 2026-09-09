package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

func TestManifestReportsBoundedDeterministicDaemonFacts(t *testing.T) {
	s := newTestServer(t)
	first := requestManifest(t, s, http.MethodGet, "/v1/manifest", nil)
	second := requestManifest(t, s, http.MethodGet, "/v1/manifest", nil)
	if first.InstallationID == "" || first.InstallationID != second.InstallationID || first.DisplayName != "go-schedule daemon" {
		t.Fatalf("manifest identity is not stable: first=%+v second=%+v", first, second)
	}
	if !reflect.DeepEqual(first.LocalAPIVersions, []string{"v1"}) || first.RemoteAPIVersions == nil || len(first.RemoteAPIVersions) != 0 || first.OperatingMode != "local_only" {
		t.Fatalf("protocol contract = %+v", first)
	}
	if !sort.StringsAreSorted(first.Capabilities) || len(first.Capabilities) != 11 || !contains(first.Capabilities, "actor-authorization") || !contains(first.Capabilities, "management-audit") || first.Platform.OS == "" || first.Platform.Architecture == "" {
		t.Fatalf("capability/platform contract = %+v", first)
	}
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/manifest", nil))
	var raw map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	wantKeys := []string{"capabilities", "display_name", "installation_id", "local_api_versions", "operating_mode", "platform", "product_version", "remote_api_versions"}
	gotKeys := make([]string, 0, len(raw))
	for key := range raw {
		gotKeys = append(gotKeys, key)
	}
	sort.Strings(gotKeys)
	if !reflect.DeepEqual(gotKeys, wantKeys) {
		t.Fatalf("manifest keys = %v, want %v", gotKeys, wantKeys)
	}
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func TestManifestRenameAndResetLifecycle(t *testing.T) {
	s := newTestServer(t)
	before := requestManifest(t, s, http.MethodGet, "/v1/manifest", nil)
	renamed := requestManifest(t, s, http.MethodPatch, "/v1/manifest", map[string]string{"display_name": "  Workshop scheduler  "})
	if renamed.DisplayName != "Workshop scheduler" || renamed.InstallationID != before.InstallationID {
		t.Fatalf("renamed manifest = %+v", renamed)
	}
	for _, body := range []any{map[string]string{"display_name": ""}, map[string]string{"unknown": "value"}} {
		rec := performJSON(s, http.MethodPatch, "/v1/manifest", body)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("invalid rename status = %d, body=%s", rec.Code, rec.Body)
		}
	}
	conflict := performJSON(s, http.MethodPost, "/v1/manifest/reset", map[string]string{"confirm_installation_id": "wrong"})
	if conflict.Code != http.StatusConflict {
		t.Fatalf("reset conflict status = %d, body=%s", conflict.Code, conflict.Body)
	}
	after := requestManifest(t, s, http.MethodPost, "/v1/manifest/reset", map[string]string{"confirm_installation_id": renamed.InstallationID})
	if after.InstallationID == renamed.InstallationID || after.DisplayName != renamed.DisplayName {
		t.Fatalf("reset manifest = %+v", after)
	}
	stale := performJSON(s, http.MethodPost, "/v1/manifest/reset", map[string]string{"confirm_installation_id": renamed.InstallationID})
	if stale.Code != http.StatusConflict {
		t.Fatalf("stale reset status = %d", stale.Code)
	}
}

func requestManifest(t *testing.T, s *Server, method, path string, body any) ManifestResponse {
	t.Helper()
	rec := performJSON(s, method, path, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("%s %s status = %d, body=%s", method, path, rec.Code, rec.Body)
	}
	var manifest ManifestResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &manifest); err != nil {
		t.Fatal(err)
	}
	return manifest
}

func performJSON(s *Server, method, path string, body any) *httptest.ResponseRecorder {
	var payload bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&payload).Encode(body)
	}
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, httptest.NewRequest(method, path, &payload))
	return rec
}

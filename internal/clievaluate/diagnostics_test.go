package clievaluate

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/client"
)

func TestSnapshotRetainsCollectionsBeforeTheFailedRequest(t *testing.T) {
	api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/credentials") {
			w.WriteHeader(503)
			_, _ = w.Write([]byte(`{"error":{"code":"internal","message":"owned unavailable"}}`))
			return
		}
		name := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		if name == "tasks" {
			_, _ = w.Write([]byte(`{"tasks":[{"id":"retained-task","name":"owned"}]}`))
			return
		}
		_, _ = w.Write([]byte(`{"` + name + `":[]}`))
	}))
	s := &session{api: api}
	partial, err := s.snapshot(context.Background())
	if err == nil {
		t.Fatal("incomplete snapshot accepted")
	}
	if _, ok := fieldAt(partial, "tasks.retained-task"); !ok {
		t.Fatalf("partial observation discarded: %v %v", partial, err)
	}
}

func TestCaseDiagnosticPhasesRetainRawCommandResultAndCapRequests(t *testing.T) {
	d := &executionDiagnostics{}
	ctx := d.phase(context.Background(), "after")
	api := client.NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`{"status":"ok"}`)) }))
	for i := 0; i < 300; i++ {
		_, _ = api.Health(ctx)
	}
	original := execution{Cleanup: errors.New("owned cleanup failure"), Code: 0, Out: []byte("command succeeded"), State: []byte(`{"after":{"tasks":{}},"observation_error":"timed out"}`)}
	result := d.finish(original)
	if result.Code != 0 || string(result.Out) != "command succeeded" {
		t.Fatal("inspection replaced command evidence")
	}
	var state struct {
		Diagnostics struct {
			Requests     []requestTiming `json:"requests"`
			Omitted      int             `json:"requests_omitted"`
			CleanupError string          `json:"process_cleanup_error"`
		}
		ObservationError string `json:"observation_error"`
	}
	if err := json.Unmarshal(result.State, &state); err != nil {
		t.Fatal(err)
	}
	if len(state.Diagnostics.Requests) != 256 || state.Diagnostics.Omitted != 44 || state.ObservationError != "timed out" || state.Diagnostics.CleanupError != "owned cleanup failure" {
		t.Fatalf("requests=%d omitted=%d observation_error=%q cleanup_error=%q", len(state.Diagnostics.Requests), state.Diagnostics.Omitted, state.ObservationError, state.Diagnostics.CleanupError)
	}
}

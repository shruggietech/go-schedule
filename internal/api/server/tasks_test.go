package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func doJSON(t *testing.T, s *Server, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatal(err)
		}
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, &buf)
	s.Handler().ServeHTTP(rec, req)
	return rec
}

func TestCreateTask_Recurring(t *testing.T) {
	s := newTestServer(t)
	rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "nightly", Command: "/bin/true", Schedule: "every day at 09:00", Timezone: "UTC",
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var resp TaskResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Task.ID == "" || resp.Schedule.HumanSummary != "Every day at 09:00" {
		t.Fatalf("unexpected task detail: %+v", resp)
	}
	if len(resp.NextRuns) == 0 {
		t.Fatal("expected computed next runs")
	}
}

func TestCreateTask_OneOffPastRejected(t *testing.T) {
	s := newTestServer(t)
	past := time.Now().UTC().Add(-time.Hour)
	rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "bday", Command: "/bin/true", At: &past,
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	var e APIError
	_ = json.Unmarshal(rec.Body.Bytes(), &e)
	if e.Error.Field != "at" {
		t.Fatalf("expected field 'at', got %+v", e)
	}
}

func TestCreateTask_OneOffFuture(t *testing.T) {
	s := newTestServer(t)
	future := time.Now().UTC().Add(24 * time.Hour)
	rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "bday", Command: "/bin/true", At: &future,
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestCreateTask_InvalidScheduleAndTimezone(t *testing.T) {
	s := newTestServer(t)
	if rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "x", Command: "/bin/true", Schedule: "every blorp",
	}); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad schedule: status %d", rec.Code)
	}
	if rec := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "x", Command: "/bin/true", Schedule: "every day", Timezone: "Mars/Phobos",
	}); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad timezone: status %d", rec.Code)
	}
}

func TestPreview(t *testing.T) {
	s := newTestServer(t)
	rec := doJSON(t, s, http.MethodPost, "/v1/schedules/preview", PreviewRequest{
		Schedule: "3rd wednesday monthly at 14:00", Timezone: "UTC",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	var resp PreviewResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.HumanSummary != "The 3rd Wednesday of every month at 14:00" || len(resp.NextRuns) == 0 {
		t.Fatalf("unexpected preview: %+v", resp)
	}
}

func TestRunNow_NotFound(t *testing.T) {
	s := newTestServer(t)
	rec := doJSON(t, s, http.MethodPost, "/v1/tasks/missing/run-now", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestListAndDeleteTask(t *testing.T) {
	s := newTestServer(t)
	create := doJSON(t, s, http.MethodPost, "/v1/tasks", TaskCreateRequest{
		Name: "x", Command: "/bin/true", Schedule: "every 15 minutes",
	})
	var resp TaskResponse
	_ = json.Unmarshal(create.Body.Bytes(), &resp)

	list := doJSON(t, s, http.MethodGet, "/v1/tasks", nil)
	if list.Code != http.StatusOK {
		t.Fatalf("list status %d", list.Code)
	}

	del := doJSON(t, s, http.MethodDelete, "/v1/tasks/"+resp.Task.ID, nil)
	if del.Code != http.StatusNoContent {
		t.Fatalf("delete status %d", del.Code)
	}
	get := doJSON(t, s, http.MethodGet, "/v1/tasks/"+resp.Task.ID, nil)
	if get.Code != http.StatusNotFound {
		t.Fatalf("get after delete status %d, want 404", get.Code)
	}
}

package server

import (
	"net/http"
	"testing"
)

func TestUpdateTaskEnvironmentClearIntent(t *testing.T) {
	s := newTestServer(t)
	created := newTaskFor(t, s, TaskCreateRequest{Name: "environment", Command: "echo", Env: map[string]string{"EXAMPLE": "old"}})
	path := "/v1/tasks/" + created.Task.ID
	for _, step := range []struct {
		name string
		body any
		want string
	}{
		{"omit", TaskUpdateRequest{Name: "renamed"}, "old"},
		{"replace", TaskUpdateRequest{Env: map[string]string{"EXAMPLE": "new"}}, "new"},
		{"clear", map[string]any{"clear_env": true}, ""},
	} {
		t.Run(step.name, func(t *testing.T) {
			rec := doJSON(t, s, http.MethodPatch, path, step.body)
			if rec.Code != http.StatusOK {
				t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
			}
			got := getTask(t, s, created.Task.ID).Task.Env
			if step.want == "" && len(got) != 0 || step.want != "" && got["EXAMPLE"] != step.want {
				t.Fatalf("environment=%v, want EXAMPLE=%q", got, step.want)
			}
		})
	}
}

func TestUpdateTaskEnvironmentClearRejectsReplacementBeforeMutation(t *testing.T) {
	s := newTestServer(t)
	created := newTaskFor(t, s, TaskCreateRequest{Name: "original", Command: "echo", Env: map[string]string{"EXAMPLE": "old"}})
	rec := doJSON(t, s, http.MethodPatch, "/v1/tasks/"+created.Task.ID, map[string]any{
		"clear_env": true, "env": map[string]string{"OTHER": "new"}, "name": "must not persist",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	got := getTask(t, s, created.Task.ID).Task
	if got.Name != "original" || len(got.Env) != 1 || got.Env["EXAMPLE"] != "old" {
		t.Fatalf("rejected request mutated task: %+v", got)
	}
}

package taskgroup

import (
	"encoding/json"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/commandline"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestRemovingFinalEnvironmentRowPreservesClearOnWire(t *testing.T) {
	current := server.TaskResponse{Task: domain.Task{Env: map[string]string{"EXAMPLE": "old"}}}
	request, err := updateRequest(TaskDraft{TaskDetail: TaskDetail{Mode: "manual"}}, current, commandline.Invocation{}, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]any
	if err := json.Unmarshal(encoded, &wire); err != nil {
		t.Fatal(err)
	}
	if wire["clear_env"] != true {
		t.Fatalf("clear intent lost during shared client encoding: %s", encoded)
	}
}

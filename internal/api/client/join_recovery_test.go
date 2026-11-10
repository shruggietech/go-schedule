package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestJoinRecoveryClientPreservesEmptyAndAbsentPackets(t *testing.T) {
	var packets []map[string]json.RawMessage
	c := NewInProcess(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var p map[string]json.RawMessage
		if err := json.Unmarshal(raw, &p); err != nil {
			t.Fatal(err)
		}
		packets = append(packets, p)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(domain.ParallelJoin{ID: "join", RecoveryTaskID: "repair", RecoveryOn: []domain.JoinRecoveryCondition{}, CleanupTaskID: "clean"})
	}))
	plain, err := c.CreateJoin(context.Background(), server.JoinCreateRequest{Name: "Base"})
	if err != nil || plain.RecoveryOn == nil {
		t.Fatalf("readback lost []: %+v %v", plain, err)
	}
	for _, k := range []string{"recovery_task_id", "recovery_on", "cleanup_task_id"} {
		if _, ok := packets[0][k]; ok {
			t.Errorf("legacy create gained %s", k)
		}
	}
	empty := []domain.JoinRecoveryCondition{}
	_, err = c.CreateJoin(context.Background(), server.JoinCreateRequest{RecoveryTaskID: "repair", RecoveryOn: &empty, CleanupTaskID: "clean"})
	if err != nil {
		t.Fatal(err)
	}
	if string(packets[1]["recovery_on"]) != "[]" {
		t.Fatalf("explicit empty=%s", packets[1]["recovery_on"])
	}
	clear := ""
	_, err = c.UpdateJoin(context.Background(), "join", server.JoinUpdateRequest{RecoveryTaskID: &clear, RecoveryOn: &empty, CleanupTaskID: &clear})
	if err != nil {
		t.Fatal(err)
	}
	if string(packets[2]["recovery_task_id"]) != `""` || string(packets[2]["cleanup_task_id"]) != `""` || string(packets[2]["recovery_on"]) != "[]" {
		t.Fatalf("clear=%v", packets[2])
	}
}

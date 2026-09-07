package taskgroup

import (
	"context"
	"testing"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

type recordingBackend struct {
	Backend
	detailed bool
}

func (b *recordingBackend) ListTaskDetails(context.Context, string, string) ([]server.TaskResponse, error) {
	b.detailed = true
	return []server.TaskResponse{{}}, nil
}

func TestLocalBackendUsesDetailedProtectedClientContract(t *testing.T) {
	daemon := &recordingBackend{}
	backend := NewLocalBackend(daemon)
	values, err := backend.ListTaskDetails(context.Background(), "", "")
	if err != nil || !daemon.detailed || len(values) != 1 {
		t.Fatalf("values=%+v called=%v err=%v", values, daemon.detailed, err)
	}
}

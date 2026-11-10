package mcpmanage

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func (f *fakeManageClient) CreateJoin(context.Context, server.JoinCreateRequest) (domain.ParallelJoin, error) {
	f.calls++
	return domain.ParallelJoin{ID: "join-created"}, f.err
}
func (f *fakeManageClient) UpdateJoin(context.Context, string, server.JoinUpdateRequest) (domain.ParallelJoin, error) {
	f.calls++
	return domain.ParallelJoin{ID: "join-updated"}, f.err
}
func (f *fakeManageClient) DeleteJoin(context.Context, string) error { f.calls++; return f.err }
func TestJoinManageUsesExactRetryAndAuthority(t *testing.T) {
	f := &fakeManageClient{}
	e := New(f, false)
	input := JoinInput{Envelope: Envelope{DaemonID: uuid.NewString(), RequestID: uuid.NewString()}, Action: "create", Create: &JoinCreate{Name: "Parallel", EntryTaskID: "entry", TargetTaskID: "target", Branches: []JoinBranch{{ID: "a", Name: "A", TaskID: "same"}, {ID: "b", Name: "B", TaskID: "same"}}}}
	_, first, err := e.joinHandler()(context.Background(), nil, input)
	if err != nil || first.Outcome != OutcomeAccepted {
		t.Fatalf("first=%+v err=%v", first, err)
	}
	_, again, _ := e.joinHandler()(context.Background(), nil, input)
	if again != first || f.calls != 1 {
		t.Fatalf("retry=%+v calls=%d", again, f.calls)
	}
}

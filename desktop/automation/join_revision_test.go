package automation

import (
	"context"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/api/client"
	"github.com/shruggietech/go-schedule/v2/internal/api/server"
	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

type revisionBackend struct {
	fakeBackend
	request server.JoinUpdateRequest
	err     error
}

func (b *revisionBackend) GetJoinDefinition(context.Context, string) (domain.ParallelJoin, error) {
	return domain.ParallelJoin{ID: "join", Name: "Previous", EntryTaskID: "previous-entry", TargetTaskID: "previous-target", UpdatedAt: time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC)}, nil
}

func (b *revisionBackend) UpdateJoin(_ context.Context, _ string, r server.JoinUpdateRequest) (domain.ParallelJoin, error) {
	b.request = r
	return domain.ParallelJoin{ID: "join"}, b.err
}
func TestJoinSaveRetainsExactBranchIDsAndAtomicRevision(t *testing.T) {
	b := &revisionBackend{}
	d := JoinDraft{ID: "join", Name: "Exports", EntryTaskID: "entry", TargetTaskID: "target", Branches: []domain.ParallelJoinBranch{{ID: "a", Name: "A", TaskID: "same"}, {ID: "b", Name: "B", TaskID: "same"}}, OriginalUpdatedAt: "2026-11-10T00:00:00Z"}
	r := NewService(b).SaveJoin(context.Background(), d)
	if r.Outcome != "accepted" || b.request.ExpectedUpdatedAt == nil || !b.request.ExpectedUpdatedAt.Equal(time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC)) || (*b.request.Branches)[0].ID == (*b.request.Branches)[1].ID {
		t.Fatalf("result=%+v request=%+v", r, b.request)
	}
	b.err = &client.StatusError{Code: server.CodeConflict}
	r = NewService(b).SaveJoin(context.Background(), d)
	if r.Outcome != "stale" {
		t.Fatalf("revision conflict=%+v", r)
	}
	d.OverwriteStale = true
	NewService(b).SaveJoin(context.Background(), d)
	if b.request.ExpectedUpdatedAt != nil {
		t.Fatal("explicit overwrite retained stale precondition")
	}
}

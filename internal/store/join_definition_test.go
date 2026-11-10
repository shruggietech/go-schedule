package store

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func TestS144AuthoringDefinitionPreservesAllAcceptedUnicodeNamesAcrossPages(t *testing.T) {
	st := openMem(t)
	entry := chainTask(t, st, "authoring entry")
	branch := chainTask(t, st, "authoring branch")
	target := chainTask(t, st, "authoring target")
	j := domain.ParallelJoin{Name: strings.Repeat("計画", 32769), EntryTaskID: entry.ID, TargetTaskID: target.ID}
	for i := 0; i < 101; i++ {
		j.Branches = append(j.Branches, domain.ParallelJoinBranch{Name: fmt.Sprintf("Branch %03d", i), TaskID: branch.ID})
	}
	for i, n := range []int{4097, 16385, 65537} {
		j.Branches[i].Name = strings.Repeat(string(rune('一'+i)), n)
	}
	if err := st.CreateParallelJoin(&j); err != nil {
		t.Fatal(err)
	}
	// Readonly task labels may be huge and are not part of editable definition intent.
	if _, err := st.db.Exec("UPDATE tasks SET name=? WHERE id IN (?,?,?)", strings.Repeat("display-only", 100000), entry.ID, branch.ID, target.ID); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	bounded, err := st.GetParallelJoinPageContext(ctx, j.ID, 0, 100, 4096)
	if err != nil || len(bounded.Name) > 4096 || len(bounded.TruncatedFields) == 0 || len(bounded.Branches[0].TruncatedFields) == 0 {
		t.Fatalf("observe projection lost bounds: length=%d fields=%v err=%v", len(bounded.Name), bounded.TruncatedFields, err)
	}
	first, err := st.GetParallelJoinDefinitionPageContext(ctx, j.ID, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	if first.Name != j.Name || first.ID != j.ID || first.EntryTaskID != j.EntryTaskID || first.TargetTaskID != j.TargetTaskID || !first.CreatedAt.Equal(j.CreatedAt) || !first.UpdatedAt.Equal(j.UpdatedAt) {
		t.Fatal("authoring read changed stored definition intent or revision")
	}
	if first.BranchCount != 101 || first.BranchOffset != 0 || len(first.Branches) != 100 || !first.HasMoreBranches || len(first.TruncatedFields) != 0 {
		t.Fatalf("authoring page metadata=%d/%d offset=%d more=%v fields=%v", len(first.Branches), first.BranchCount, first.BranchOffset, first.HasMoreBranches, first.TruncatedFields)
	}
	if first.EntryTaskName != "" || first.TargetTaskName != "" {
		t.Fatal("authoring projection includes readonly task labels")
	}
	for i, b := range first.Branches {
		if b.ID != j.Branches[i].ID || b.Name != j.Branches[i].Name || b.TaskID != j.Branches[i].TaskID || b.TaskName != "" || len(b.TruncatedFields) != 0 {
			t.Fatalf("authoring branch %d changed editable fields or included readonly labels", i)
		}
	}
	last, err := st.GetParallelJoinDefinitionPageContext(ctx, j.ID, 100, 100)
	if err != nil || last.BranchOffset != 100 || last.BranchCount != 101 || len(last.Branches) != 1 || last.HasMoreBranches {
		t.Fatalf("last authoring page offset=%d count=%d branches=%d more=%v err=%v", last.BranchOffset, last.BranchCount, len(last.Branches), last.HasMoreBranches, err)
	}
	first.Branches = append(first.Branches, last.Branches...)
	first.HasMoreBranches = false
	replacement := chainTask(t, st, "replacement target")
	first.TargetTaskID = replacement.ID
	if err := st.UpdateParallelJoin(&first); err != nil {
		t.Fatal(err)
	}
	after, err := st.GetParallelJoinDefinitionPageContext(ctx, j.ID, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if after.Name != j.Name || after.TargetTaskID != replacement.ID || !reflect.DeepEqual(after.Branches, j.Branches) {
		t.Fatal("target-only edit changed accepted definition or branch names")
	}
}

func TestS144AuthoringDefinitionCancellationCoversQueuedPoolAcquisition(t *testing.T) {
	st := openMem(t)
	j, _ := joinFixture(t, st)
	held, err := st.db.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer held.Close()
	baseline := st.db.Stats().WaitCount
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { _, err := st.GetParallelJoinDefinitionPageContext(ctx, j.ID, 0, 100); done <- err }()
	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()
	for st.db.Stats().WaitCount == baseline {
		select {
		case <-tick.C:
		case <-timeout.C:
			t.Fatal("authoring read did not enter SQL pool wait")
		}
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("authoring read lost cancellation: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("canceled authoring read waited for pool release")
	}
}

func TestS144AuthoringDefinitionMissingAndClosedStorage(t *testing.T) {
	st, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetParallelJoinDefinitionPageContext(context.Background(), "missing", 0, 100); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing definition=%v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetParallelJoinDefinitionPageContext(context.Background(), "missing", 0, 100); err == nil {
		t.Fatal("closed storage authoring read reported success")
	}
}

package mcpobserve

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shruggietech/go-schedule/v2/internal/domain"
)

func (f *fakeReader) ListJoinsPage(ctx context.Context, offset, limit, branchOffset, branchLimit, textLimit int) ([]domain.ParallelJoin, bool, error) {
	f.seenCtx = ctx
	return []domain.ParallelJoin{}, false, f.err
}
func (f *fakeReader) GetJoinPage(ctx context.Context, id string, offset, limit, textLimit int) (domain.ParallelJoin, error) {
	f.seenCtx = ctx
	return domain.ParallelJoin{ID: id}, f.err
}
func (f *fakeReader) ListJoinInstancesPage(ctx context.Context, filter domain.ParallelJoinInstanceFilter) ([]domain.ParallelJoinInstance, bool, error) {
	f.seenCtx = ctx
	return []domain.ParallelJoinInstance{}, false, f.err
}
func (f *fakeReader) GetJoinInstancePage(ctx context.Context, id string, filter domain.ParallelJoinInstanceFilter) (domain.ParallelJoinInstance, error) {
	f.seenCtx = ctx
	return domain.ParallelJoinInstance{ID: id}, f.err
}

type joinBoundReader struct {
	*fakeReader
	filter  domain.ParallelJoinInstanceFilter
	exactID string
}

func (f *joinBoundReader) ListJoinInstancesPage(ctx context.Context, filter domain.ParallelJoinInstanceFilter) ([]domain.ParallelJoinInstance, bool, error) {
	f.filter = filter
	return []domain.ParallelJoinInstance{{ID: "instance", EntryRunID: "entry-run", BranchCount: 101, HasMoreBranches: true, Branches: []domain.ParallelJoinInvocation{{ID: "invocation-a", Name: "Untrusted <instructions>", TaskID: "same", RunID: "exact-run-a"}}}}, true, nil
}
func (f *joinBoundReader) GetJoinInstancePage(ctx context.Context, id string, filter domain.ParallelJoinInstanceFilter) (domain.ParallelJoinInstance, error) {
	f.filter = filter
	f.exactID = id
	return domain.ParallelJoinInstance{ID: id, EntryRunID: "entry-run", BranchOffset: filter.BranchOffset, BranchCount: 101, Branches: []domain.ParallelJoinInvocation{{ID: "invocation-b", TaskID: "same", RunID: "exact-run-b"}}}, nil
}
func TestJoinResourcesPageExactInvocationsAtTheAPI(t *testing.T) {
	r := &joinBoundReader{fakeReader: &fakeReader{}}
	a := &adapter{client: r, now: time.Now, timeout: time.Second}
	e, err := a.envelope(context.Background(), "join_instances", joinInstancesURI, 0)
	if err != nil {
		t.Fatal(err)
	}
	if r.filter.Limit != PageLimit || r.filter.BranchLimit != PageLimit || r.filter.TextLimit != TextLimit {
		t.Fatalf("unbounded query=%+v", r.filter)
	}
	raw, _ := json.Marshal(e)
	if !strings.Contains(string(raw), "branches_next_uri") || !strings.Contains(string(raw), "entry-run") || !strings.Contains(string(raw), TrustNotice) {
		t.Fatalf("missing evidence or continuation: %s", raw)
	}
	uri := joinInstancesURI + "/instance/branches/" + encodeCursor(100)
	result, err := a.readJoinExact(joinInstancesURI)(context.Background(), &mcp.ReadResourceRequest{Params: &mcp.ReadResourceParams{URI: uri}})
	if err != nil || r.exactID != "instance" || r.filter.BranchOffset != 100 || !strings.Contains(result.Contents[0].Text, "exact-run-b") {
		t.Fatalf("exact=%+v filter=%+v err=%v", result, r.filter, err)
	}
	if strings.Contains(result.Contents[0].Text, "exact-run-a") {
		t.Fatal("mixed distinct branch pages")
	}
}

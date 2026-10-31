package cli

import (
	"context"
	"testing"

	"github.com/shruggietech/go-schedule/v2/internal/api/server"
)

func TestGroupEditRejectsCycleWithoutPartialRename(t *testing.T) {
	c := isolatedCommandBackend(t)
	ctx := context.Background()
	root, err := c.CreateGroup(ctx, server.GroupCreateRequest{Name: "original"})
	if err != nil {
		t.Fatal(err)
	}
	child, err := c.CreateGroup(ctx, server.GroupCreateRequest{Name: "child", ParentID: root.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runFixtureCommand(c, "", "group", "edit", root.ID, "--name", "must not persist", "--parent", child.ID); err == nil {
		t.Fatal("cycle accepted")
	}
	groups, err := c.ListGroups(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, got := range groups {
		if got.ID == root.ID && (got.Name != "original" || got.ParentID != "") {
			t.Fatalf("rejected edit partially persisted: %+v", got)
		}
	}
}

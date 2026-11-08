package maintainercheck

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func lifecycleFixture(t testing.TB, root, state, delivery, tasks, disposition string) {
	t.Helper()
	put(t, root, "specs/001-example/spec.md", "# Example\n\n**Status**: "+state+"\n\n**Delivery**: "+delivery+"\n"+disposition)
	put(t, root, "specs/001-example/tasks.md", tasks)
	put(t, root, "specs/README.md", "| 001-example | "+state+" | record |\n")
}

func TestLifecycleStates(t *testing.T) {
	for _, tc := range []struct{ state, delivery, tasks, disposition, want string }{
		{"Draft", "", "- [ ] T001 work\n", "", ""},
		{"Ready", "", "- [ ] T001 work\n", "", ""},
		{"In Progress", "", "- [x] T001 work\n- [ ] T002 work\n", "", ""},
		{"Implemented", "commit abc1234", "- [x] T001 work\n", "", ""},
		{"Implemented", "Pending implementation", "- [x] T001 work\n", "", "invalid delivery evidence"},
		{"Implemented", "commit abc1234", "- [ ] T001 work\n", "", "unresolved task"},
		{"In Progress", "", "- [x] T001 work\n", "", "completed and actionable"},
		{"Deferred", "", "", "**Disposition**: Owner deferred\n", ""},
		{"Superseded", "", "", "**Disposition**: New design\n", ""},
		{"Abandoned", "", "", "", "missing a disposition"},
		{"Finished", "", "", "", "invalid status"},
	} {
		t.Run(tc.state+tc.want, func(t *testing.T) {
			root := t.TempDir()
			lifecycleFixture(t, root, tc.state, tc.delivery, tc.tasks, tc.disposition)
			_, findings := Lifecycle(root)
			got := strings.Join(findings, "\n")
			if tc.want == "" && got != "" {
				t.Fatal(got)
			}
			if tc.want != "" && !strings.Contains(got, tc.want) {
				t.Fatalf("want %s got %s", tc.want, got)
			}
		})
	}
}

func TestLifecycleLiteralRoot(t *testing.T) {
	for _, name := range []string{"checkout[1]", "checkout[unclosed"} {
		t.Run(name, func(t *testing.T) {
			root := filepath.Join(t.TempDir(), name)
			lifecycleFixture(t, root, "Implemented", "commit abc1234", "- [x] T001 work\n", "")
			count, findings := Lifecycle(root)
			if count != 1 || len(findings) != 0 {
				t.Fatal(count, findings)
			}
		})
	}
}

func BenchmarkLifecycle133(b *testing.B) {
	root := b.TempDir()
	inventory := ""
	for i := 1; i <= 133; i++ {
		name := fmt.Sprintf("%03d-example", i)
		put(b, root, "specs/"+name+"/spec.md", "**Status**: Implemented\n**Delivery**: commit abc1234\n")
		put(b, root, "specs/"+name+"/tasks.md", "- [x] T001 work\n")
		inventory += "| " + name + " | Implemented | record |\n"
	}
	put(b, root, "specs/README.md", inventory)
	b.ResetTimer()
	for b.Loop() {
		count, findings := Lifecycle(root)
		if count != 133 || len(findings) != 0 {
			b.Fatal(count, findings)
		}
	}
}

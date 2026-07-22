package cli

import "testing"

// TestGroupIntent_ThreeWay pins the CLI half of FR-014/FR-015. Before this
// change `--group ""` was indistinguishable from omitting the flag, so no
// client could take a task back out of a group.
func TestGroupIntent_ThreeWay(t *testing.T) {
	tests := []struct {
		name string
		argv []string
		want *string // nil = leave membership unchanged
	}{
		{"omitted leaves membership unchanged", []string{"some-id"}, nil},
		{"explicit empty removes from group", []string{"some-id", "--group", ""}, strp("")},
		{"named group assigns", []string{"some-id", "--group", "g-123"}, strp("g-123")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := taskEdit()
			// Parse flags only; RunE would need a live daemon.
			if err := cmd.ParseFlags(tt.argv); err != nil {
				t.Fatalf("parse %v: %v", tt.argv, err)
			}
			group, err := cmd.Flags().GetString("group")
			if err != nil {
				t.Fatal(err)
			}

			got := groupIntent(cmd, group)
			switch {
			case tt.want == nil && got != nil:
				t.Errorf("got %q, want nil (membership must be left unchanged)", *got)
			case tt.want != nil && got == nil:
				t.Errorf("got nil, want %q", *tt.want)
			case tt.want != nil && got != nil && *got != *tt.want:
				t.Errorf("got %q, want %q", *got, *tt.want)
			}
		})
	}
}

func strp(s string) *string { return &s }

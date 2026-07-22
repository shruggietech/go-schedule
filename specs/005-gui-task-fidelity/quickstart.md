# Quickstart validation: GUI task fidelity

**Feature**: [spec.md](spec.md) · **Date**: 2026-07-22

Runnable checks that prove the feature end to end. Automated first (these are the
merge gate), then the manual walkthrough that covers what only a human at the GUI
can confirm.

## Prerequisites

- Go toolchain per `go.mod`, cgo available for the race build
- Windows for the GUI walkthrough (the reported defects are Windows/.msi)
- A built daemon and clients: `go build ./...`

## Automated — CI parity

Run in the foreground and watch to completion. Never background these: `go test`
buffers a package's output until the package finishes, so a backgrounded run
cannot be told apart from a dead one.

```bash
gofmt -l internal cmd test gui
```

```bash
go vet ./...
```

```bash
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run ./...
```

```bash
CGO_ENABLED=1 go test -race $(go list ./... | grep -vE '/cmd/gosched-gui|/gui$')
```

```bash
go test ./gui/...
```

Expected: `gofmt` prints nothing; everything else passes; core packages stay at
or above 80% coverage.

## Automated — feature-specific gates

| Check | Command | Proves |
|---|---|---|
| Phrase round-trip | `go test ./internal/schedule/ -run 'Render'` | FR-004: `Parse → Render → Parse` preserves the RRULE for every supported phrase |
| Migration safety | `go test ./internal/store/ -run 'Migrat'` | FR-002/FR-024: a v3 database upgrades with every schedule row intact, and re-opening is a no-op |
| Group tri-state | `go test ./internal/api/server/ -run 'Group'` | FR-014/FR-016: omitted / empty / named behave as specified, unknown IDs are 400 |
| Editor prefill | `go test ./gui/ -run 'Prefill\|Editor'` | FR-006/FR-007/FR-008/FR-011b: mode, phrase, one-off date+time, clean-dirty state, mode-switch gating |
| Group views | `go test ./gui/ -run 'Group'` | FR-018/FR-019/FR-020: membership tree, ungrouped node, move action |

## Manual — GUI walkthrough

Start the daemon, then `gosched-gui`.

### A. Pre-existing database *(the reported defect, SC-001)*

Use a database created before this change — either a real v0.3.0 profile or one
produced by checking out `main`, creating a few tasks, and returning to this
branch.

1. Open a recurring task for editing.
   **Expect**: Mode `Recurring`, Schedule populated with an equivalent phrase.
   **Fail if**: Schedule is blank.
2. Save without touching anything. Note the task's next run before and after.
   **Expect**: identical (SC-002, SC-007).

### B. Round-trip fidelity

3. Create a recurring task `every 15 minutes starting at 09:00`. Reopen it.
   **Expect**: Schedule `every 15 minutes`, Start at `09:00`, not doubled and not
   dropped.
4. Create a one-off for a future date. Reopen it.
   **Expect**: Mode `One-off`, Date and Time showing that instant **in the task's
   timezone** (set a non-local timezone to make this meaningful).
5. On a recurring task, switch Mode to `One-off` and leave date and time empty.
   **Expect**: Save disabled (FR-011b).
6. Open any task, change nothing, press Cancel.
   **Expect**: closes immediately, no discard prompt (FR-008).
7. On a recurring task, change only the timezone and save.
   **Expect**: next runs move to the new zone (FR-011).

### C. Group assignment

8. Create group `Parent`, then `Child` under it.
9. Create a task and assign it to `Child` from the editor.
   **Expect**: the choice list shows `Child` with its path under `Parent`
   (FR-013); the task list shows the group (FR-017).
10. In the Groups tab, expand `Child`.
    **Expect**: the task is listed beneath it, visually distinct from groups
    (FR-018).
11. Select the task there and move it to `Parent`.
    **Expect**: it relocates without a manual refresh (FR-020, FR-022).
12. Move it to `(none)`.
    **Expect**: it appears under the ungrouped area (FR-019, Story 3).
13. With no ungrouped tasks left, look at the hierarchy.
    **Expect**: the ungrouped area is still shown, empty (FR-019).
14. Select a task and press a group-only action (New Group / Enable-Disable /
    Delete).
    **Expect**: it does not act on the task (FR-021).
15. Assign the task back to `Parent`, disable `Parent`.
    **Expect**: the task is suppressed by the existing cascade (Story 2).

### D. Cross-client agreement

16. `gosched task list` — group membership matches the GUI.
17. `gosched task edit <id> --group ""` — the task ungroups, and the GUI reflects
    it live (FR-015, FR-022).
18. `gosched task edit <id>` with no `--group` — membership unchanged.

## Done when

Every automated command above passes in the foreground, and every "Expect" in
the walkthrough holds. Issues
[#3](https://github.com/shruggietech/go-schedule/issues/3) and
[#4](https://github.com/shruggietech/go-schedule/issues/4) are then no longer
reproducible (SC-008).

# Research: GUI Navigation and Information

## Decision: Keep the current Activity terminology

**Rationale**: Issue #29 predates PR #46, which renamed the combined log and alert
surface to Activity and added its bounded alert badge. Git history shows that the
tab reference is intentionally retained so its label can update without
rebuilding the collection. The new order is therefore Tasks, Groups, Schedule,
Activity, Info.

**Alternatives considered**:

- Restore Logs to match the old issue text. Rejected because it reverses a
  completed clarity improvement.
- Rename the surface again. Rejected because no new naming problem is in scope.

## Decision: Use one new package-local view file

**Rationale**: Existing GUI views are separated by concern (`tasks.go`,
`groups.go`, `schedule.go`, `logs.go`). A small `info.go` follows that convention
and can consume package-local `appIcon` plus `internal/buildinfo.Version` without
changing the Backend interface.

**Alternatives considered**:

- Inline all Info construction in `app.go`. Rejected because it makes navigation
  assembly own unrelated presentation detail.
- Add an Info viewmodel or daemon endpoint. Rejected because the data is immutable
  and already process-local.

## Decision: Use Fyne's standard image, scroll, and hyperlink primitives

**Rationale**: The pinned Fyne 2.7.4 source confirms that `canvas.Image` supports
aspect-preserving containment, `container.NewVScroll` bounds tall local content,
and `widget.Hyperlink` is focusable and delegates activation to `fyne.OpenURL`.
These existing primitives meet the accessibility and platform-integration needs
without custom browser logic.

**Alternatives considered**:

- Custom URL buttons. Rejected because they would recreate styling, focus, and
  activation behavior already maintained by Fyne.
- Rich-text markdown links. Rejected because independently testable hyperlink
  widgets provide clearer labels and URL contracts for this small fixed set.

## Decision: Construct canonical URLs as values, not parsed strings at runtime

**Rationale**: The three destinations are compile-time constants. Explicit
`url.URL` values avoid a fallible parse path and keep invalid hard-coded
destinations visible in review and tests.

**Alternatives considered**:

- Parse strings and ignore errors. Rejected by the code-quality constitution.
- Parse strings and panic on error. Rejected because the values can be expressed
  without either error handling or a panic.

## Decision: Test the public presentation contract headlessly

**Rationale**: The GUI package already shares one Fyne test application and
inspects its widget tree. Regression tests can assert the exact tab sequence,
Activity badge position, local mark resource, version label, link labels, and
URL values without a display server or external browser.

**Alternatives considered**:

- Screenshot-only testing. Rejected because it is platform-sensitive and weaker
  for exact text and destination contracts.
- Manual-only verification. Rejected because both issues describe stable,
  deterministic behavior suitable for automated regression coverage.

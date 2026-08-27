# Implementation Plan: GUI Activity Clarity

**Branch**: `codex/014-gui-activity-clarity` | **Date**: 2026-08-27 |
**Spec**: `specs/014-gui-activity-clarity/spec.md`

## Summary

Rename the GUI's mixed log-and-alert destination to Activity, centralize its
bounded badge label, and replace the destructive-looking "Dismiss All" control
with a clear-content action plus persistent explanatory copy. Preserve the
existing view cutoff and alert acknowledgement semantics. Close issues #26,
#28, and #30 without expanding into the log-path API work in #31.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Fyne v2.7.4; existing GUI view model and Backend
interface

**Storage**: No changes; existing daemon persistence remains authoritative

**Testing**: Go unit tests with Fyne's headless test driver, followed by
`sh scripts/verify.sh all`

**Target Platform**: Existing Linux, macOS, and Windows desktop GUI targets

**Project Type**: Go desktop application backed by a local daemon

**Performance Goals**: Constant-time tab-label formatting; no change to data
loading or rendering volume

**Constraints**: No API, daemon, persistence, dependency, workflow, or pinned
artifact changes; explanations must not depend on hover

**Scale/Scope**: One GUI navigation label, one badge formatter, one Activity
toolbar/action explanation, and focused regression tests

## Constitution Check

| Gate | Result | Evidence |
| --- | --- | --- |
| Code quality | PASS | A small pure label helper isolates the only boundary logic; low-value internal renames are excluded. |
| Testing | PASS | Tests are written first for labels, boundaries, and control presentation; the full aggregate remains mandatory. |
| UX consistency | PASS | Activity accurately names mixed content, the badge is bounded, and clear behavior is described without hover. |
| Performance | PASS | Formatting remains constant-time and the data path is unchanged. |
| Autopilot | PASS | The slice is traceable to #26, #28, and #30 and will halt once before publication. |
| Integration | PASS | Work stays on a review branch and will use `Closes` for all three completed issues. |

No constitution exception is required.

## Project Structure

### Documentation (this feature)

```text
specs/014-gui-activity-clarity/
├── checklists/ux.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
└── tasks.md
```

### Source Code

```text
gui/
├── app.go
├── app_test.go
├── logs.go
└── logs_test.go

CHANGELOG.md
CLAUDE.md
```

**Structure Decision**: Keep the existing GUI package and internal Logs names.
Only user-facing terminology changes. This avoids a broad mechanical rename
that would not improve behavior or comprehension.

## Design

1. Add a pure `activityTabLabel(int) string` helper in `gui/app.go`. Return
   `Activity` for counts at or below zero, the exact count for 1 through 99,
   and `Activity (99+)` from 100 upward.
2. Use the helper for both initial tab construction and badge refresh.
3. In `gui/logs.go`, rename the visible action to `Clear View`, replace the
   delete icon with Fyne's `ContentClearIcon`, and add subdued persistent copy:
   `Hides current activity and acknowledges visible alerts. Records are not deleted.`
4. Retain `clearedAt`, the current filtered-row alert ID collection, and the
   existing backend acknowledgement loop unchanged except for terminology in
   comments.
5. Add focused headless GUI tests before implementation, then run the full
   canonical aggregate.

## Post-Design Constitution Re-check

PASS. The design adds no architecture, dependencies, persistence behavior, or
governance machinery. Persistent help improves accessibility compared with the
issue's hover-tooltip suggestion, while the scope remains proportional.

## Explicit Deviation from Issue #28's Proposed Tooltip

Issue #28 proposed an on-hover tooltip. This plan uses always-visible explanatory
copy instead because the project has no tooltip pattern, Fyne's core button API
does not expose one, and hover-only disclosure would be unavailable on touch
devices. The result satisfies the issue's actual need with less machinery and
better discoverability.

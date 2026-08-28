# Implementation Plan: Activity Diagnostics Clarity

**Branch**: `codex/022-activity-diagnostics` | **Date**: 2026-08-28 | **Spec**: [spec.md](spec.md)

## Summary

Make Activity honest and useful during troubleshooting by carrying the daemon's
resolved log path with its bounded recent-record response, presenting that
metadata without normalizing or probing it, and replacing the ambiguous daemon
listening message with one testable startup-completion event. Preserve existing
CLI output by explicitly unwrapping the richer response at the CLI boundary.

## Technical Context

**Language/Version**: Go 1.25.0, Markdown
**Primary Dependencies**: Go standard library, existing Fyne v2 UI and slog/logbus stack
**Storage**: Existing rotating JSONL log only; no schema, retention, or persistence changes
**Testing**: Go unit and GUI tests plus all eight canonical verification gates
**Target Platform**: Windows and Linux local desktop/daemon environments
**Project Type**: Local daemon, IPC API, CLI, and desktop GUI
**Performance Goals**: No additional request, filesystem access, or unbounded state
**Constraints**: Exact path preservation, stable CLI output, UTF-8 without BOM, no file-opening behavior
**Scale/Scope**: One additive response field, one view-model field, one Activity label, and one startup record

## Constitution Check

| Principle | Gate | Design evidence |
| --- | --- | --- |
| I. Code Quality | PASS | The resolved path has one source of truth and passes unchanged through typed layers. |
| II. Testing Standards | PASS | Server, view-model, presentation, GUI integration, and startup behavior receive focused tests before implementation. |
| III. UX Consistency | PASS | Activity explains its bounded scope; startup and API fields use consistent structured names; CLI output remains stable. |
| IV. Performance | PASS | Metadata rides the existing logs request; no extra I/O, polling, or growing collection is introduced. |
| V. Autonomous Execution | PASS | S022 follows the full Spec-Kit sequence, review branch, canonical gates, local commit, and mandatory pre-push halt. |

### Post-design re-check

All gates remain satisfied. No dependency, migration, benchmark-sensitive path,
network surface, or governance mechanism is added.

## Architecture and Decision Log

### Extend the existing recent-activity response

`GET /v1/logs` gains an explicit `log_path` string alongside `logs`. The server
receives the already resolved `Config.LogPath()` value and returns it even when
the ring is nil or empty. An empty string is retained in the JSON shape and
means metadata is unavailable, which is also compatible with older daemons.

This avoids an extra health request and keeps record and path metadata from the
same daemon response. Changing the small internal server constructor directly
is preferable to introducing an options abstraction for six call sites.

### Preserve CLI compatibility at its boundary

The API client returns the complete typed logs response, consistent with other
client methods. The CLI deliberately formats and serializes only its `Logs`
member, preserving both the current table and bare-array JSON output.

### Keep Activity presentation semantic and passive

The view-model stores records and path atomically. Activity uses a refreshed,
word-wrapped diagnostics label that says the view is limited, points operators
to older daemon records, and renders either the exact path or an explicit
unavailable state. It does not duplicate the numeric ring limit or offer file
actions.

### Treat startup as a discrete event

A small helper emits exactly one `daemon startup complete` record at the current
pre-serve location with structured `endpoint`, `db`, and `log_path` fields. No
uptime or lifecycle state is introduced.

## Project Structure

```text
specs/022-activity-diagnostics/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── verification.md
├── contracts/activity-diagnostics.md
├── checklists/
└── tasks.md

cmd/goschedd/
gui/
├── logs.go
├── app_test.go
└── viewmodel/
internal/
├── api/client/
├── api/server/
└── cli/
docs/cli.md
specs/004-rebrand-gui-overhaul/contracts/api-logs.md
CHANGELOG.md
```

**Structure Decision**: Extend the existing local API, client, view-model, GUI,
and daemon seams. No new package or dependency is warranted.

## Implementation Phases

1. Add failing response and view-model propagation tests for exact log-path metadata.
2. Add failing presentation and GUI tests for bounded-view, exact-path, and unavailable wording.
3. Add a failing structured startup-record test.
4. Implement the narrow API, client, CLI, view-model, UI, and daemon changes.
5. Align the current API/user docs and changelog.
6. Analyze, run focused and canonical verification, audit encoding, commit locally, and halt before publication.

## Complexity Tracking

No constitutional violations or architecture exceptions are required.

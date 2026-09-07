# S062 Verification

**Date**: 2026-09-07

**Scope**: Complete production Wails Tasks and Groups workflow (#153), plus the shared safe platform-native first-task path (#192).

## Cross-artifact analysis

Spec Kit analysis found no critical, high, medium, ambiguity, duplication, or unmapped-requirement findings. All 30 functional requirements and 10 measurable outcomes map to 56 chronological tasks, and both requirements checklists are complete.

## Focused evidence

| Verification | Result |
| --- | --- |
| Focused server and client API tests | PASS, opt-in detailed reads, default compatibility, disabled draft groups, safe clear semantics, and escaped group updates |
| `go test ./internal/commandexample ./gui` | PASS, exact platform allowlist, bounded native execution, Fyne insertion, traversal, reinsertion, and edit isolation |
| `go test ./test/integration -run TestGuidedPlatformExampleCreatesRunsAndRecordsRecognizableActivity -count=1` | PASS on Windows, real `cmd.exe /d /c ver` execution persisted a successful recognizable manual Activity run without a visible child console |
| `go test -race -count=20 ./...` from `desktop/` | PASS, application, connection, task/group service, local adapter, stale-write, safe-error, hierarchy, differential-write, and 100-operation coverage |
| `npm test` from `desktop/frontend/` | PASS, 10 files and 21 tests |
| `npm run build` from `desktop/frontend/` | PASS, TypeScript and Vite production build |
| `npm run test:e2e` from `desktop/frontend/` | PASS, 8 Chromium contracts including 100 tasks, 20 group levels, WCAG checks, narrow width, 200 percent zoom, reduced motion, offline assets, and one-shot keyboard insertion |
| `go run github.com/wailsapp/wails/v2/cmd/wails@v2.14.0 build -clean` from `desktop/` | PASS, native Windows amd64 package built as `go-schedule.exe` |
| `go run ./scripts/github-format` | PASS, no em dashes or hard-wrapped Markdown prose |
| `sh scripts/spec-lifecycle-check.sh .` | PASS, S062 lifecycle and inventory are consistent |
| `git diff --check` | PASS |

## Acceptance audit

### #153

- The production Tasks route now provides search, effective-state filtering, stable selection, full-path groups, schedule and policy summaries, detail context, and honest empty, loading, and unavailable states.
- Task creation and editing expose name, group, direct command line, working directory, ordered environment input, standard input, run identity, enabled intent, timezone, recurring, one-off, and manual timing, plus every current scheduling policy.
- Writes are pessimistic and daemon-authoritative. The service parses direct command boundaries, validates environment rows, sends differential updates, preserves schedules on untouched saves, recompiles recurring schedules when timezone intent changes, rejects stale timestamps, and requires explicit overwrite.
- Groups support disabled draft creation, rename, reparent, root movement, declared state changes, cascade-aware counts, full duplicate-safe paths, descendant exclusion, task reassignment, and target-aware deletion.
- Run now, enable, disable, and delete controls name their task and This computer where confirmation is required, suppress duplicate activation, surface safe outcomes, preserve selection, and hand accepted Run now actions to the honest Activity route.
- Automated Chromium evidence exercises one hundred tasks and a twenty-level hierarchy with keyboard, WCAG, compact layout, zoom, reduced-motion, and local-only asset contracts.

### #192

- `internal/commandexample` is the one exact mapping for Windows `cmd.exe /d /c ver`, macOS `/usr/bin/sw_vers`, and Linux `uname -a`. Both supported desktop experiences consume it directly or are guarded by exact consistency tests.
- A fresh command field displays the execution-host suggestion, explains what it prints and the Run now to Activity path, inserts once on forward Tab while empty, keeps focus with the caret at the end, restores ordinary traversal after insertion, never overwrites content, and offers insertion again after clearing.
- Existing-task editors receive no suggestion, including blank legacy drafts. Shift+Tab remains ordinary reverse traversal because Fyne captures only the eligible unmodified empty-field Tab and React checks the modifier explicitly.
- Exact parsing, allowlisting, current-host execution under five seconds, recognizable captured output, documentation consistency, and retired Python listener and file-writing guidance absence are automated.
- The production desktop matrix executes the guided example on Windows, macOS, and Linux before building each native package. Hosted results are publication evidence recorded on the pull request.

## Deliberate implementation consolidation

The planned standalone `GroupEditor.tsx` was consolidated into `GroupsPanel.tsx` because the parent-choice exclusion set, hierarchy counts, selected identity, and cascade confirmation share one tightly coupled authoritative snapshot. This avoids a second hierarchy state owner without expanding scope. The existing shared table, field, notice, dialog, and state primitives were reused instead of adding synonymous abstractions to `components/index.tsx`.

## Canonical repository verification

The canonical Windows Bash compatibility command completed successfully:

```text
GO="/mnt/c/Program Files/Go/bin/go.exe" GOFMT="/mnt/c/Program Files/Go/bin/gofmt.exe" sh scripts/verify.sh all
```

All eight gates passed in order: format, vet, lint, race, GUI, coverage, documentation, and automation. Core coverage remained engine 81.9 percent, schedule 89.2 percent, timezone 91.3 percent, store 80.1 percent, catchup 88.9 percent, and logbus 91.1 percent.

## Hosted evidence boundary

Pull-request CI must pass the production desktop build and real guided example on Windows, macOS, and Linux, the Chromium contract, root race and coverage gates, CodeQL, and packaging checks. No hosted or third-party review result is claimed before it arrives on the pull request.

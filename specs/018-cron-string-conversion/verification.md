# Verification Evidence: Pure Schedule String Conversion

**Date**: 2026-08-28

**Branch**: `codex/018-cron-string-conversion`

## Baseline

- `go test ./internal/cron ./internal/cli -count=1` passed before test or
  production changes.
- `internal/cron.Explain` already performs local cron-to-human translation and
  named fidelity refusal.
- `internal/cron.Export` renders recurring schedules only through a task-aware
  wrapper that also rejects disabled tasks.
- `schedule.Parse` still rejects cron as task-authoring input, and existing
  import reaches task creation only through a human phrase.
- `cron explain` contacts the daemon only for upcoming-run preview; there is no
  pure symmetric one-string CLI operation.
- CLI results use stdout, diagnostics use stderr, and `errUsage` maps validation
  failures to exit code 2.

## Test-first evidence

### User Story 1 red phase

- Added automatic/forced cron classification, exact output, normalization,
  malformed input, unsupported input, no-fallback, invalid-destination, and
  daemon-independent command regressions.
- `go test ./internal/cron ./internal/cli -run
  'Test(Convert|CronConvert)' -count=1` failed as expected because the
  conversion types, service, and command did not yet exist.

### User Story 1 green phase

- Implemented stable syntax/result values, single-pass structural detection,
  local cron-to-human conversion, destination validation, and exact text CLI
  output.
- The same focused command passed with no daemon client in the conversion path.

### User Story 2 red phase

- Added canonical human output, forced direction, explicit-phase fidelity,
  carve-out, schedule-only export preservation, and calendar-parity regressions.
- `go test ./internal/cron -run
  'Test(ConvertHuman|ExportSchedule|Export_Expressible)' -count=1` failed as
  expected because the schedule-only renderer did not yet exist.

### User Story 2 green phase

- Extracted schedule-only rendering behind the unchanged task-state wrapper,
  canonicalized weekdays as `1-5`, rejected implicit/misaligned anchors, and
  implemented deterministic human-to-cron conversion.
- Corrected already-supported sub-daily cron phrases to retain their `00:00`
  phase rather than inheriting task-creation time.
- `go test ./internal/cron ./internal/cli -run
  'Test(Convert|CronConvert|Export|RoundTrip)' -count=1` passed.

### User Story 3 red phase

- Added stable five-field JSON success/refusal, exact stdout/stderr, exit-class,
  and duplicate-diagnostic regressions.
- `go test ./internal/cli -run
  'Test(CronConvert_JSON|HandleExecuteError)' -count=1` failed as expected
  because structured rendering and reported-error handling did not yet exist.

### User Story 3 green phase

- Implemented stable JSON through Cobra writers, structured refusals on stderr,
  usage-class exit 2, and a reported-error marker that suppresses duplicate
  root diagnostics while preserving every existing unreported error path.
- `go test ./internal/cli -run
  'Test(CronConvert|HandleExecuteError)' -count=1` and the combined focused
  cron/CLI suite passed.

### Classification refinement

- Added regressions for existing five-field human phrases after recognizing
  that field count alone conflicts with the supported human grammar.
- The focused conversion suite failed by classifying both phrases as cron.
- Narrowed automatic cron detection to `@` prefixes or five fields with a
  cron-shaped minute field. Forced destination behavior is unchanged.

## Focused verification

- `go test ./internal/cron ./internal/cli -count=1` passed after implementation.
- A freshly built `gosched` binary produced:
  - cron to human: exit 0, stdout `weekdays at 09:00\n`, empty stderr;
  - human to cron: exit 0, stdout `0 9 * * 1-5\n`, empty stderr;
  - structured success: exit 0, all five fields on stdout, empty stderr;
  - structured refusal: exit 2, empty stdout, all five fields and a named
    explicit-phase refusal on stderr.
- Both directions ran through the local command without daemon, IPC, API,
  network, storage, configuration, or task mutation.
- Requirements checklist: 16/16 complete.
- Conversion contract checklist: 24/24 complete.

## Repository verification

- `git diff --check` passed.
- Strict UTF-8 decoding, no-BOM, and mojibake audits passed across all 26
  changed and untracked files.
- `sh scripts/verify.sh all` passed in the foreground:
  1. format;
  2. vet;
  3. lint (0 issues);
  4. race;
  5. GUI;
  6. coverage (engine 88.1%, schedule 91.5%, timezone 88.9%, store 87.0%,
     catchup 87.5%, logbus 91.1%);
  7. docs (11 pages and all documentation contracts clean);
  8. automation (approved actions, CodeQL contract, and gate manifest clean).
- Final Spec-Kit analysis found 22/22 requirements and success outcomes covered,
  26/26 valid unique task IDs, and no ambiguity, duplication, or constitution
  conflict.

## Issue disposition

- #51: complete and eligible for `Closes #51` after merge.
- #50: this slice supplies a reusable pure conversion boundary but does not
  accept or retain cron in tasks; the pull request will use `Refs #50` and leave
  it open.

## Review follow-up

- AI review identified that the new hourly phase guard rejected a nonzero
  minute even though cron preserves it in the minute field.
- New export and conversion regressions reproduced the failure for `00:30` and
  `08:30` phases before the guard was corrected.
- Hourly conversion now checks only hour-step alignment after retaining the
  existing sub-minute precision check. A complementary `09:30` regression
  confirms that a genuinely misaligned two-hour phase is still refused.
- Focused cron/CLI tests and all eight canonical repository gates passed after
  the correction.

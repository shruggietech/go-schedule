# Verification: Natural Command Entry

## Baseline

- Branch `codex/046-natural-command-entry` started clean from `main` commit `96e061475052ce37b6a9827c5e583fd48db10bdc`.
- Git remote is `git@github.com:shruggietech/go-schedule.git`; issue [#110](https://github.com/shruggietech/go-schedule/issues/110) is open in milestone `v1.0.0 - Release readiness`.
- Pre-change `taskEditor` owns separate `command` and `args` widgets. Arguments are a multiline input with `one argument per line`; `splitArgs` trims lines and drops blanks, so empty, surrounding-space, and literal-newline arguments cannot round-trip through the GUI.
- The pre-change preview reconstructs one display string and cannot distinguish every exact argument boundary.
- Baseline focused command-preview tests passed: `go test ./gui -run 'TestEditor_CommandPreview|TestCommandLinePreview' -count=1`.
- `.gitignore` already covers Go binaries/tests, generated Windows resources, coverage, runtime data, vendor/workspace data, editor files, and `dist/`; no ignore edit is required.
- Spec Kit cross-artifact analysis covered 26 functional requirements, 8 success criteria, and 44 tasks with zero critical/high findings. Both requirements checklists pass (39/39).

## Test-First Evidence

- Foundational red run: `go test ./internal/commandline -count=1` failed to compile because `Invocation`, `Parse`, `Format`, and `QuoteDisplay` did not exist. This was observed after the grammar, error, canonical identity, exact display, and fuzz-seed tests were added and before `internal/commandline/commandline.go` was created.
- Foundational green run: `go test ./internal/commandline -count=1` passed after implementing the portable parser and formatter.
- User Story 1 red run: the focused GUI selection failed to compile because `taskEditor.commandLine` and `taskEditor.leftContent` did not exist. This was observed before replacing the old fields and layout.
- User Story 1 green run: `go test ./gui -run 'TestEditor_CommandField|TestEditor_CommandValidation|TestEditor_InvalidCommand|TestEditor_SaveGating|TestEditor_CommandPreview' -count=1` passed after the one-field editor, six-row minimum, stretch layout, validation, and exact create submission were implemented.
- User Story 2 red run: the focused GUI preview selection failed to compile because `commandPreviewText` did not exist. This was observed after exact-boundary, stale-preview, recovery, literal-punctuation, and help tests replaced the ambiguous reconstructed-preview expectations.
- User Story 2 green run: `go test ./gui -run 'TestCommandPreview|TestEditor_CommandPreview|TestEditor_InvalidCommandClears|TestEditor_CommandHelp' -count=1` passed after exact preview composition and portable direct-execution help were implemented.
- Compatibility testing exposed and corrected one formatter defect before delivery: a bare argument ending in backslash escaped the formatter's separating space and merged with the following argument. The canonical formatter now quotes trailing-backslash values, and a core regression keeps such a value before another argument.
- User Story 3 compatibility run initially failed on that trailing-backslash merge and on test-helper JSON decoding of the test runner's trailing `PASS` marker. After correcting the formatter and decoding one JSON value from helper output, `go test ./internal/commandline ./internal/api/server ./internal/store ./internal/executor -count=1` passed on native Windows, followed by the exact existing-task editor regression and the complete headless `gui` package (81.789s).

## Focused Verification

- Formatting plus focused packages passed: `gofmt` over every changed Go file, then `go test ./internal/commandline ./internal/api/server ./internal/store ./internal/executor -count=1`.
- The complete headless editor package passed: `go test ./gui -count=1` (77.795s on the final focused run before lifecycle-only edits).
- The seeded formatter/parser fuzz property passed for 3 seconds after 3,843,542 executions. Fuzzing first found an invalid-UTF-8 input that Go rune conversion would silently replace; `Parse` and `Format` now reject invalid UTF-8, and the generated corpus seed remains as a regression.
- Native Windows helper-process coverage passed with exact empty, spaced, Unicode, tab, newline, backslash, variable-marker, pipe, redirect, wildcard, separator, and ampersand arguments. `platform.HideConsole` remains applied to the Windows child process.
- Existing-task editor prefill/update, structured API, SQLite persistence, and executor helper-process tests all passed with exact unusual value and ordering comparisons.

## Full Verification

- `go test -race ./internal/commandline ./internal/api/server ./internal/store ./internal/executor -count=1` passed.
- `go test ./... -count=1` passed, including the complete GUI package (103.368s) and integration package (65.142s).
- The canonical `scripts/verify.sh all` run passed all eight gates in order: format, vet, lint (zero issues), race, GUI, coverage, docs, and automation. Core coverage remained above policy: engine 86.4%, schedule 89.2%, timezone 91.3%, store 84.1%, catchup 88.9%, and logbus 91.1%.
- The first canonical invocation revealed three lint findings in new code (one unchecked cache call and two `WriteString(fmt.Sprintf(...))` uses). They were corrected before the clean all-gates pass; no lint suppressions were added.
- The second and final external Codex review found two valid edge cases. The update contract now distinguishes omitted arguments from an explicit empty list so GUI edits can clear stale arguments, and exact preview rendering now escapes non-ASCII Unicode space separators. Focused regressions and the full canonical gates passed after both corrections.

## Requirements Audit

- **FR-001 through FR-003, SC-003**: `gui/editor.go` owns one multiline **Command line** entry, sets six visible rows, and uses a surplus-height layout. `TestEditor_CommandFieldIsOneRoomyGrowingMultilineInput` proves the minimum and growth behavior.
- **FR-004 through FR-010**: `internal/commandline` implements one portable lexical grammar, Unicode-whitespace splitting, quotation/adjacency, path-safe backslashes, empty values, strict UTF-8, NUL rejection, and position-aware quote errors. Table, identity, and fuzz tests cover each boundary.
- **FR-011 through FR-014, SC-004 and SC-005**: preview composition separately labels Program and numbered Arguments with exact escaped display; editor tests prove live valid, invalid, stale-clear, recovery, and Save-gating states. Inline Help contains Windows/POSIX, empty, repeated, Unicode, and literal-newline examples without `argv` vocabulary.
- **FR-015 through FR-020, SC-001, SC-002, SC-006, and SC-007**: create/update still submit the existing structured program plus ordered arguments. Existing-task, API, store, executor, and native-host tests prove exact values. No API, persistence, CLI, daemon, schema, or implicit-shell path changed. The existing CI race matrix executes host-specific native tests on Windows, macOS, and Linux.
- **FR-021 through FR-024**: inline Help documents both host families, quoting, spaces, empty/repeated/Unicode/multiline values, no expansion, and explicit-shell responsibility. The implementation retains Fyne's standard focusable multiline Entry, validation, preview, and Button controls and adds no pointer-only interaction.
- **FR-025, SC-008**: unit, fuzz, API, persistence, executor, full headless GUI, native Windows, race, integration, and all eight repository gates passed.
- **FR-026**: diff inspection confirms S046 contains no failure-message or output-capture policy change; issue #102 remains separate.
- All 26 functional requirements and all 8 measurable outcomes have implementation and verification evidence. No criterion remains partial or waived.

## Encoding and Diff Integrity

- Every changed and added file decoded as strict UTF-8 and none contained a UTF-8 BOM (30 files inspected).
- The mojibake signature scan found no suspect sequences.
- `git diff --check` passed with no whitespace errors.
- The final specification lifecycle and documentation link gates passed after S046 was marked Implemented and all 44 tasks were resolved.

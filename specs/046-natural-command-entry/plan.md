# Implementation Plan: Natural Command Entry

**Branch**: `codex/046-natural-command-entry` | **Date**: 2026-09-03 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/046-natural-command-entry/spec.md`

## Summary

Resolve [#110](https://github.com/shruggietech/go-schedule/issues/110) by replacing the task editor's separate executable and one-argument-per-line controls with one vertically roomy **Command line** editor. A small internal package will parse one documented, platform-independent direct-command grammar into the existing `Command` plus `Args` model and generate a canonical lossless editor representation for existing tasks. The GUI will show an exact separately labeled program/ordered-argument preview, actionable syntax errors, and clear no-shell guidance. API, persistence, CLI, executor, and issue #102 diagnostic behavior remain unchanged.

## Technical Context

**Language/Version**: Go 1.25.0

**Primary Dependencies**: Go standard library; Fyne 2.8.1 for the existing desktop editor

**Storage**: Existing SQLite task columns (`command` and JSON argument array); no migration

**Testing**: Go unit/integration tests, headless Fyne tests, platform-tagged Windows/POSIX tests, race detector, canonical `scripts/verify.sh all`

**Target Platform**: Windows, macOS, and Linux desktop clients; existing daemon execution on all supported hosts

**Project Type**: Cross-platform desktop application backed by a local daemon/API

**Performance Goals**: Parse, validate, format, and refresh realistic command lines within one synchronous edit cycle; linear time and memory in input size

**Constraints**: Direct execution only; identical grammar on every platform; exact round trip for every process-valid string; no new module dependency; no persistence or API shape change; at least six visible editor lines; no issue #102 diagnostic expansion

**Scale/Scope**: One editor surface, one parser/formatter package, existing task authoring/execution paths, representative command inputs up to the existing GUI/API practical limits

## Constitution Check

*GATE: Passed before research and again after design.*

| Principle / constraint | Design response | Gate |
| --- | --- | --- |
| I. Reliability First | The stored invocation remains the authority; formatter/parser identity tests prevent silent mutation | PASS |
| II. Determinism & Correctness | One portable grammar yields identical boundaries on every host; invalid quotation is rejected | PASS |
| III. Observability | The preview exposes the exact program and each numbered argument, including escaped invisible characters | PASS |
| IV. Performance | Parser and formatter are single-pass/bounded-linear operations; preview updates only the edited invocation | PASS |
| V. Autonomous Build-Phase Execution | S046 is issue-backed (#110), uses a review branch, runs all Spec Kit phases, and publishes only under explicit user authorization | PASS |
| Go/toolchain and dependency constraint | Standard library plus the pinned existing Fyne dependency is sufficient | PASS |
| Supported platforms | Portable parser tests run on the three-host CI matrix; platform-tagged tests prove native process boundaries | PASS |
| Backward compatibility | Domain, wire, SQLite, CLI, and executor shapes stay unchanged; existing tasks are formatted then reparsed losslessly | PASS |
| Security | No implicit shell, expansions, pipelines, redirects, or environment interpolation; shell use requires naming one explicitly | PASS |

### Deliberate Design Deviation

Issue #110 lists platform-aware parsing as one candidate. S046 deliberately uses one portable grammar instead because host-dependent parsing would allow the same persisted editor text to mean different argument boundaries on Windows and POSIX. The authoritative stored representation remains one program plus an ordered argument list, and the editor grammar is only a reversible authoring projection of that model.

### Workflow Tooling Deviation

The installed checklist prerequisite script requires `plan.md`, contradicting the mandated `specify → clarify → checklist → plan` order. S046 used the feature directory returned by the successful clarification prerequisite command to generate `checklists/command-entry.md` before plan setup. This preserves the governing sequence without modifying unrelated Spec Kit tooling.

## Project Structure

### Documentation (this feature)

```text
specs/046-natural-command-entry/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── portable-command-line.md
├── checklists/
│   ├── requirements.md
│   └── command-entry.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
internal/commandline/
├── commandline.go              # portable parser, canonical formatter, exact preview values
├── commandline_test.go         # grammar, errors, round trips, fuzz seeds
├── commandline_unix_test.go    # POSIX native process-boundary proof
└── commandline_windows_test.go # Windows native process-boundary proof

gui/
├── editor.go                   # one multiline command editor, validation, preview, submission
├── editor_data.go              # display-only invocation formatting
├── editor_test.go              # headless editor, sizing, validation, help, submission
└── editor_prefill_test.go      # existing-task lossless edit regressions

docs/
├── gui-fields.md               # portable syntax and exact preview contract
└── INSTALL-windows.md          # installed-service Windows authoring example
```

**Structure Decision**: Put the reusable, UI-independent grammar at `internal/commandline` because it needs focused cross-platform and native-process tests. Keep widget layout and plain-language preview composition in the existing `gui` package. Do not move or duplicate the authoritative task fields in domain, API, storage, CLI, or executor code.

## Phase 0: Research Decisions

Research findings and rejected alternatives are recorded in [research.md](research.md). The decisive choices are:

1. Use one portable direct-command grammar rather than Windows/POSIX-dependent parsing.
2. Preserve the existing structured program-plus-arguments storage and executor boundary.
3. Canonically format stored invocations using reversible quotation, without rewriting live user text.
4. Preview the program and numbered arguments independently with quoted escaped values.
5. Give the multiline editor a layout-enforced six-line minimum and vertical growth.
6. Prove behavior at parser, GUI, API/persistence, executor, POSIX, and native Windows boundaries.

## Phase 1: Design

### Grammar boundary

- `commandline.Parse` converts editor text into one non-empty program plus ordered arguments.
- Unicode whitespace separates values only outside quotes.
- Single and double quotes group values. Empty quotes create an empty value; adjacent segments compose one value.
- Backslash escapes whitespace or a quote outside quotes and a double quote inside double quotes. Other backslashes remain literal, preserving Windows paths and POSIX path text.
- Shell punctuation has no special handling.
- Invalid UTF-8 is rejected without replacement. An unmatched quote or NUL returns a position-aware error. A trailing ordinary backslash is literal.
- `commandline.Format` emits a canonical text form satisfying `Parse(Format(invocation)) == invocation` for every process-valid invocation.

### GUI boundary

- Replace `command` and `args` widgets with one `commandLine` multiline entry.
- Apply a minimum-height wrapper derived from six text rows while allowing the left-pane scroll/layout to allocate more height.
- Parse once per edit refresh and use that result for validation, preview, and submission.
- Invalid input clears the invocation preview and shows a plain-language error. Valid input shows `Program` and `Arguments in order`, with every value rendered through an escaped quoted display.
- Prefill existing tasks with `commandline.Format`; unchanged save reparses to identical structured values.

### Compatibility boundary

- `domain.Task`, API requests, SQLite schema, CLI flags, and `executor.Run` do not change.
- GUI create/update calls still send `Command` and `Args` separately.
- The executor still calls the platform's direct process API through Go's structured command facility and retains its no-visible-console rule.
- Shells remain runnable only when explicitly entered as the program. S046 adds no mode, inference, or acknowledgement state.

### Verification strategy

- Red-first parser tables and seeded fuzz tests cover grammar, invalid syntax, canonical identity, Unicode, spaced paths, repeated flags, empty values, quotes, backslashes, CR/LF/tab, and shell punctuation.
- Headless GUI tests cover one-field construction, six-line minimum, vertical growth contract, live preview, stale-preview clearing, keyboard-editable multiline configuration, Save gating, create/update submission, help, and existing task round trips.
- Existing API/store tests gain an exact unusual-argument round trip only if current coverage does not already pin it.
- Native tagged tests launch a helper process with representative values and compare the received argument vector on Windows and POSIX.
- Executor tests retain direct invocation and add shell-punctuation literals without entering issue #102 output presentation.

## Complexity Tracking

No constitution violation requires justification. One small internal package is the minimum coherent boundary because placing parsing in GUI callbacks would make platform tests and formatter identity hard to isolate, while changing the stored model would add unnecessary migration and compatibility risk.

# Feature Specification: Natural Command Entry

**Feature Branch**: `codex/046-natural-command-entry`

**Created**: 2026-09-03

**Status**: Implemented

<!-- Allowed states and transition evidence: specs/README.md -->

**Delivery**: Review branch `codex/046-natural-command-entry`; focused, race, full-suite, fuzz, native Windows, and canonical eight-gate verification passed 2026-09-03; pull-request review remains publication workflow.

**Input**: GitHub issue [#110](https://github.com/shruggietech/go-schedule/issues/110) and the approved S046 work-slice request

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Enter a Familiar Command Line (Priority: P1)

As a task author, I want one roomy command-line field where I can type a normal-looking program and its arguments, so I do not have to translate a familiar command into an unexplained one-argument-per-line format.

**Why this priority**: Command entry is the primary blocked workflow. If it remains confusing, users cannot reliably create tasks even when scheduling and execution are correct.

**Independent Test**: Enter representative Windows and POSIX commands containing flags, quoted paths, Unicode, empty arguments, repeated flags, and a quoted literal newline. The editor identifies one program and the exact ordered arguments, allows Save only for valid input, and submits those exact values.

**Acceptance Scenarios**:

1. **Given** a blank task editor, **When** the user enters `python -m http.server --bind "127.0.0.1"`, **Then** the preview identifies `python` as the program and the remaining four values as ordered arguments.
2. **Given** a Windows executable path containing spaces, **When** the user quotes the path and adds arguments on the same line, **Then** the preview preserves the path as one program value without invoking a command shell.
3. **Given** an empty quoted argument, repeated flags, Unicode text, or a quoted literal newline, **When** the command is parsed, **Then** each value and its position remain visibly distinguishable in the preview.
4. **Given** unmatched quotation or unsupported NUL input, **When** the user reviews the editor, **Then** a plain-language validation message identifies the problem and Save remains unavailable.

---

### User Story 2 - Understand Exactly What Will Run (Priority: P1)

As a cautious operator, I want the preview and help to distinguish the launched program from each argument in order, so I can verify the task without knowing process API terminology or guessing how a reconstructed command string will behave.

**Why this priority**: A more natural editor is only safe when its translation into the existing direct-execution model is exact, reviewable, and understandable.

**Independent Test**: Compare the editor preview with the submitted program and arguments for valid and invalid inputs. Every argument, including empty strings, embedded whitespace, quotes, backslashes, and newlines, is represented unambiguously and the help explains the syntax and non-shell behavior.

**Acceptance Scenarios**:

1. **Given** valid command input, **When** the preview refreshes, **Then** it separately labels the program and numbered arguments and uses an escaped display that makes boundaries and invisible characters unambiguous.
2. **Given** text containing `|`, `>`, `*`, `$`, `%`, or similar shell punctuation, **When** it is parsed, **Then** the characters remain literal argument content and no expansion, redirection, pipeline, or environment substitution occurs.
3. **Given** an operator who needs shell behavior, **When** they read the help, **Then** it explains that a shell must be named explicitly as the program and that the shell, not go-schedule, then owns interpretation and security.

---

### User Story 3 - Edit Existing Tasks Without Loss (Priority: P1)

As an existing user, I want every previously stored task to open in the new editor and save without changing its program or argument values, so the usability improvement does not require migration or alter what my tasks execute.

**Why this priority**: Existing tasks may contain values that the former line editor could not faithfully display or re-enter. Backward compatibility is a release gate, not a follow-up enhancement.

**Independent Test**: Open and save stored task fixtures containing spaces, quotes, empty strings, Unicode, repeated arguments, backslashes, tabs, carriage returns, and literal newlines. The resulting persisted and executed program-plus-argument vector is byte-for-byte equivalent in UTF-8 text values and positional order.

**Acceptance Scenarios**:

1. **Given** an existing task with an arbitrary valid program and ordered arguments, **When** it opens, **Then** the combined command line is a canonical lossless representation of those stored values.
2. **Given** an unchanged existing task, **When** it is saved, **Then** its program and arguments are identical to the pre-edit values and no persistence migration occurs.
3. **Given** the same canonical editor text on Windows, macOS, or Linux, **When** it is parsed, **Then** it produces the same program and ordered arguments on every platform.

### Edge Cases

- Leading, trailing, and repeated whitespace outside quoted content separates values but does not create accidental empty arguments.
- Empty quoted text creates an intentional empty program value or argument; an empty program is rejected while empty arguments are preserved.
- Single-quoted content preserves every enclosed character literally, including backslashes and line breaks.
- Double-quoted content can contain escaped quotation marks and backslashes; other backslashes, including Windows path separators, remain literal.
- Outside quotes, escaping whitespace, a quotation mark, or a backslash keeps that character in the current value; backslashes before ordinary characters remain literal.
- Adjacent quoted and unquoted segments form one value, allowing a quotation mark or whitespace-bearing segment inside a larger value.
- A literal newline outside quotation marks is a separator; a literal newline inside quotation marks is part of that value.
- Shell metacharacters, environment markers, glob characters, comments, and redirection characters have no special meaning.
- Extremely long values remain editable and inspectable without forcing a horizontal-scroll-only workflow.
- Existing values containing process-valid tabs, carriage returns, or line feeds remain losslessly representable and visibly escaped in the preview; invalid UTF-8 and NUL are rejected because the editor contract is UTF-8 text and supported process APIs cannot launch NUL.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The task editor MUST replace the separate Command and one-argument-per-line Arguments inputs with one primary field labeled **Command line**.
- **FR-002**: The Command line field MUST accept a program followed by zero or more arguments in a familiar whitespace-and-quotation form.
- **FR-003**: The field MUST show at least six text lines at the default task-dialog size and MUST grow vertically when additional dialog height is available.
- **FR-004**: The editor MUST parse command text with the same documented portable direct-command grammar on Windows, macOS, and Linux.
- **FR-005**: Outside quotation marks, Unicode whitespace MUST delimit values; repeated or surrounding delimiters MUST NOT create arguments.
- **FR-006**: Single and double quotation marks MUST group whitespace and literal newlines within one value, and the grouping quotation marks MUST NOT become part of the value.
- **FR-007**: An empty quoted segment MUST create an intentional empty value, and adjacent quoted or unquoted segments without a delimiter MUST compose one value.
- **FR-008**: Backslash treatment MUST preserve ordinary Windows and POSIX paths while supporting literal whitespace, quotation marks, and backslashes as defined in the user help and command-entry contract.
- **FR-009**: The parser MUST reject invalid UTF-8, unmatched quotation marks, and NUL content with plain-language errors that identify the problem; character-level syntax errors MUST locate the failing character or line and explain the expected correction.
- **FR-010**: The editor MUST reject an absent, empty, invalid-UTF-8, or NUL-containing program while preserving intentional empty arguments and rejecting invalid-UTF-8 or NUL-containing arguments.
- **FR-011**: The live preview MUST separately label the exact program and every argument in positional order without using `argv` as required user vocabulary.
- **FR-012**: Preview values MUST use an unambiguous escaped representation that distinguishes empty strings, spaces, tabs, carriage returns, line feeds, quotation marks, and backslashes.
- **FR-013**: Preview and validation MUST refresh as the user edits and MUST never display a valid invocation derived from stale text after the current text becomes invalid.
- **FR-014**: Save MUST remain unavailable until both the command line and all other currently relevant task fields are valid.
- **FR-015**: Creating or updating a task MUST submit the parsed program and ordered arguments through the existing structured task contract.
- **FR-016**: The daemon MUST continue launching tasks as one explicit program plus an ordered argument list without implicitly invoking `cmd.exe`, PowerShell, or a POSIX shell.
- **FR-017**: Characters commonly interpreted by shells, including pipes, redirects, wildcards, variable markers, separators, and comments, MUST remain literal unless the user explicitly names a shell program and passes those characters to it.
- **FR-018**: Existing stored tasks MUST open as a canonical command line that parses back to the identical program and ordered arguments, including empty and multiline arguments.
- **FR-019**: Opening and saving an unchanged existing task MUST NOT require or perform a persistence-schema migration.
- **FR-020**: The API, persistence, CLI, and executor contracts MUST remain compatible with the existing separate `command` and ordered `args` values.
- **FR-021**: Inline guidance MUST include Windows and POSIX examples for executable paths, spaces, quotation marks, empty arguments, repeated flags, Unicode, and literal newlines.
- **FR-022**: Help MUST state that the command-entry grammar is portable, performs no shell expansion, and has identical token boundaries on every supported platform.
- **FR-023**: Help MUST explain how to request shell behavior explicitly and warn that shell-specific quoting, expansion, and security then belong to the named shell.
- **FR-024**: Keyboard focus, selection, editing, validation, preview inspection, and Save gating MUST remain usable without a pointer.
- **FR-025**: Unit, API round-trip, persistence, executor, headless GUI, POSIX, and native Windows tests MUST cover the accepted command model and the issue #110 edge-case matrix.
- **FR-026**: S046 MUST NOT change task failure-message presentation or output-capture policy tracked separately by issue #102.

### Key Entities

- **Command Draft**: The user-authored portable command-line text, which may be valid, incomplete, or invalid while editing.
- **Direct Invocation**: One non-empty program and an ordered list of zero or more arguments. This remains the authoritative stored and executed form.
- **Canonical Command Line**: A lossless textual representation generated from a Direct Invocation and accepted by the same portable parser on every supported platform.
- **Launch Preview**: A user-readable, unambiguous view of the Direct Invocation with a separately labeled program and numbered argument values.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All issue #110 boundary cases (spaces, quotes, empty arguments, Unicode, spaced paths, repeated flags, and literal newlines) complete parse, canonical-format, reparse, create, edit, persistence, API, and execution round trips with zero value or ordering changes.
- **SC-002**: The same canonical text produces an identical program and ordered arguments in automated Windows, macOS, and Linux verification.
- **SC-003**: At the default task-dialog size, the command editor displays at least six text lines and gains additional usable height when the dialog height increases.
- **SC-004**: Every syntactically invalid command case disables Save and replaces any stale launch preview with actionable guidance within one edit-refresh cycle.
- **SC-005**: A user can enter each documented Windows and POSIX example in one field and verify the resulting program and argument order from the preview without consulting external documentation.
- **SC-006**: One hundred percent of pre-S046 task fixtures open and save with identical stored program and argument values and without a database migration.
- **SC-007**: Automated inspection finds zero implicit shell-launch paths introduced by S046, and shell punctuation remains literal in all direct-launch tests.
- **SC-008**: The complete repository verification suite passes with no regression in task creation, editing, scheduling, persistence, or execution.

## Assumptions

- The existing stored `command` plus ordered `args` representation is the correct execution boundary and remains authoritative.
- A single portable grammar is safer and more predictable than host-dependent parsing because scheduled tasks and configuration data can move between supported platforms.
- The editor may canonically rewrite quotation style when loading an existing task; only the represented program and arguments, not the cosmetic spelling, require round-trip identity.
- Directly naming `cmd`, PowerShell, `sh`, or another shell remains possible because shells are ordinary executables, but S046 adds no shell mode and performs no shell interpretation itself.
- Issue #102 remains the authority for richer failure messages and captured process output and is excluded from this slice.

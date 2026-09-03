# Research: Natural Command Entry

## Decision 1: Use one portable direct-command grammar

**Decision**: Parse the task editor with one documented lexical grammar that behaves identically on Windows, macOS, and Linux.

**Rationale**: Microsoft documents Windows command-line backslash/quote rules as a platform convention, while POSIX shells use different rules. A task can be persisted, backed up, or edited on a different host. Host-dependent parsing would therefore make identical editor text deceptively change argument boundaries. A portable grammar preserves the issue's stronger cross-platform requirement and still accepts familiar whitespace and quotation input.

**Alternatives considered**:

- Host-native Windows and POSIX parsing: familiar in isolation but unsafe for portable task data and difficult to preview as one stable contract.
- POSIX shell-word parsing everywhere: widely recognized but treats ordinary unquoted Windows backslashes as escapes and suggests shell semantics the executor does not provide.
- Keep one argument per line: lossless but explicitly rejected by the issue evidence as unnatural and too cramped.

**Primary references**:

- [Microsoft `CommandLineToArgvW` documentation](https://learn.microsoft.com/en-us/windows/win32/api/shellapi/nf-shellapi-commandlinetoargvw)
- [Go `os/exec` package documentation](https://pkg.go.dev/os/exec)

## Decision 2: Retain direct structured execution

**Decision**: Treat the parser as an authoring adapter only. Continue storing and launching one program plus an ordered argument list.

**Rationale**: Go's official `os/exec` documentation states that it intentionally does not invoke a system shell or perform globbing, pipelines, redirects, or expansions. The current executor already follows that model. Preserving it avoids injection surprises and maintains all API, CLI, storage, and execution compatibility.

**Alternatives considered**:

- Pass the entire field to a shell: superficially literal, but changes quoting, environment, redirection, security, and portability semantics.
- Store only the typed string: requires a migration and reparsing at run time, which could change task meaning after parser updates.
- Add a shell mode: not needed to resolve #110 because users can explicitly name a shell executable under the existing model.

## Decision 3: Make formatting a lossless projection

**Decision**: Generate a canonical editor string from stored program/arguments. Emit simple values bare, prefer double quotes for whitespace-bearing values without double quotes, prefer single quotes when needed, and compose quoted/escaped segments for values containing both quote styles.

**Rationale**: Existing tasks need a one-field representation that can always be reparsed without changing value boundaries. Canonical formatting avoids storing duplicate command text, permits stable tests, and preserves Windows path separators without noisy universal backslash doubling.

**Alternatives considered**:

- Preserve original spelling: unavailable for existing tasks and would require a redundant persisted field.
- Use a reconstructed display string as the authority: ambiguous for empty values and literal newlines.
- JSON array editing: lossless but less natural than the existing UI and issue request.

## Decision 4: Separate exact preview values

**Decision**: Present a separately labeled program and numbered arguments using a quoted escaped display representation.

**Rationale**: A combined string cannot visibly distinguish an empty argument, literal newline, or some quote/backslash boundaries. Separate values explain the launch without requiring `argv` vocabulary and make stale or invalid state obvious.

**Alternatives considered**:

- Continue the current reconstructed code block: readable for simple input but not exact for quotes, backslashes, empties, or newlines.
- Token chips: visually clear but expensive for multiline selection, keyboard editing, copying, and long arguments.
- Preview only the canonical command line: reversible for machines, less direct for users auditing boundaries.

## Decision 5: Enforce editor height through layout

**Decision**: Wrap the multiline entry in a tiny layout whose minimum height represents at least six text rows and whose assigned height expands with its container.

**Rationale**: Calling `Resize` on a widget managed by a parent layout is not a stable size contract. A layout-level minimum participates correctly in Fyne measurement and remains headless-testable.

**Alternatives considered**:

- Fixed pixel resize on the Entry: parent layouts can immediately overwrite it.
- Add a nested scroll pane: the Entry already owns multiline scrolling and a second scroll surface would degrade wheel and keyboard behavior.
- Enlarge the entire dialog only: does not guarantee the field receives the added space.

## Decision 6: Verify both representation and native launch boundaries

**Decision**: Combine pure parser/formatter tests, headless GUI tests, existing API/store coverage, executor tests, and build-tagged helper-process tests on Windows and POSIX.

**Rationale**: Pure parsing tests cannot prove the executor passes the resulting vector to a native process, while native tests alone do not cover invalid edit states or existing-task formatting. The CI race matrix already runs on Windows, macOS, and Ubuntu, so tagged tests provide real host evidence without a new workflow.

**Alternatives considered**:

- GUI tests only: miss persistence and native process boundaries.
- Windows attended testing only: cannot prove POSIX parity and is unnecessarily manual for argument vectors.
- New end-to-end workflow: existing CI already supplies the required operating systems.

## Decision 7: Preserve Spec Kit phase order around the checklist helper defect

**Decision**: Generate the requirements-quality checklist from the feature path returned by `check-prerequisites.ps1 -Json -PathsOnly`, then run plan setup.

**Rationale**: The installed checklist prerequisite requires a plan despite the command skill and autopilot protocol placing checklist before plan. The established workaround preserves the mandated order without unrelated workflow changes.

**Alternatives considered**:

- Plan before checklist: violates the governing sequence.
- Patch Spec Kit during S046: unrelated scope.
- Skip the checklist: violates the autopilot protocol.

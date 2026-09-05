# Verification: Persisted Adjustable Columns

## Baseline

- 2026-09-05: `git status --short --branch` showed clean synchronized `main`
  before branch creation.
- 2026-09-05: `go test ./gui` passed before implementation.
- Existing `.gitignore` covers Go/Fyne binaries, tests, coverage, local runtime
  data, editor metadata, and local build output. No ignore change was needed.

## Test-driven evidence

- Red: focused GUI compilation failed before implementation because
  `Preferred`, `newColumnProfile`, `newAdjustableStructuredList`, and
  `columnKeyboardStep` did not exist.
- Green: preference, allocation, pointer, keyboard, accessibility, alignment,
  view-isolation, reset, and default-allocation focused tests pass.

## Final verification

- `go test ./gui -count=1`: PASS (87.301 seconds).
- `go test -race ./gui -run "Test(ColumnProfile|ProfileWidths)" -count=1`:
  PASS (1.083 seconds).
- A broader `go test -race ./gui/...` attempt was stopped after it exceeded six
  minutes without output. The plan was corrected to the focused pure-state race
  surface because project policy explicitly excludes the Fyne widget package
  from the canonical race gate due to upstream font-cache behavior.
- `go test ./gui -run TestColumnBoundariesRenderInBothThemesAtSupportedWidths
  -count=1`: PASS (0.101 seconds), covering light/dark separator visibility and
  exact allocation at 560- and 1200-unit widths.

## Canonical iteration

- The first canonical run passed format, vet, lint, race, GUI, coverage, and
  documentation, then correctly rejected a lifecycle mismatch because the
  inventory said **In Progress** while the feature spec still said **Draft**.
  The spec status was corrected before rerunning the gate.
- Corrected `scripts/verify.sh automation`: PASS.
- Corrected `scripts/verify.sh all`: PASS. Format printed no files; vet and
  lint passed with zero issues; the project-defined race suite passed; GUI
  passed; core coverage was engine 86.5%, schedule 89.2%, timezone 91.3%, store
  84.3%, catchup 88.9%, and logbus 91.1%; docs and automation passed.
- Final implemented-lifecycle `scripts/verify.sh all`: PASS on the exact code
  candidate; the uncached GUI gate completed in 90.464 seconds with the same
  coverage results and green documentation and automation gates.

## Traceability audit

- FR-001 through FR-016: satisfied by the shared boundary widget, normalized
  profile, per-view constructors and reset actions, and focused regression tests.
- SC-001 through SC-007: satisfied by conservation, persistence/fallback,
  isolation/reset, practical default, appearance, full GUI, and canonical gate
  evidence above.
- No dependency, daemon schema, IPC contract, scheduler behavior, or pinned
  process artifact changed.

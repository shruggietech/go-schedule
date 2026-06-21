---
description: "Task list for GUI Editor Refinements"
---

# Tasks: GUI Editor Refinements

**Input**: Design documents from `specs/003-gui-editor-refinements/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md),
[data-model.md](data-model.md), [contracts/editor-ui.md](contracts/editor-ui.md),
[quickstart.md](quickstart.md)

**Tests**: INCLUDED — constitution Principle II (testing NON-NEGOTIABLE). Headless Fyne tests +
pure-function unit tests.

**Organization**: By user story. Most editor stories share `gui/editor.go` and are **sequential**;
the window helper (US5) and per-tab cursor swaps (US6) are independent files and can run in
parallel with the editor work.

## Path Conventions

Single Go module; all work under `gui/`.

---

## Phase 1: Setup (Shared widgets)

**Purpose**: Reusable primitives used by multiple stories.

- [ ] T001 [P] Add cursor-aware toolbar button constructors in `gui/widgets.go` (`newToolbarButton(label string, icon fyne.Resource, tapped func()) *cursorButton` and a no-icon variant), reusing the existing `cursorButton`; doc comments.
- [ ] T002 [P] Add a custom `collapsible` widget in `gui/widgets.go`: a header `cursorButton` with `theme.NavigateNextIcon()` (▶) when collapsed and `theme.MenuDropDownIcon()` (▼) when expanded, toggling a content container's visibility; starts collapsed; exposes `SetOpen`/`open` state (per [research.md](research.md) §5).
- [ ] T003 [P] Add tests in `gui/widgets_test.go`: `collapsible` starts collapsed showing the ▶ icon, toggles to ▼ + shows content, and the header reports `desktop.PointerCursor`; `newToolbarButton` reports the pointer cursor and fires its callback.

**Checkpoint**: Shared widgets compile and are tested.

---

## Phase 2: Foundational (Two-pane editor shell)

**Purpose**: Restructure `gui/editor.go` from a single-column form into a wider two-pane modal
(left form, right pane holder) without changing 002 behavior. Blocks US1–US4.

**⚠️ CRITICAL**: Blocks US1, US2, US3, US4.

- [ ] T004 Widen the editor dialog (~1180×720) and rebuild `build()` in `gui/editor.go` to a two-pane layout: left = the existing scrollable form (What to run / When / Advanced); right = a pane holder (`container.NewStack`) that will hold Preview or Help. Remove the "Preview" `FormItem` from the When form; move `schedPreview`+`cmdPreview` into the right pane (per [contracts/editor-ui.md](contracts/editor-ui.md)).
- [ ] T005 Update `rebuildWhen()` so it no longer appends the Preview row; verify mode visibility, Start-at, validation, and anchor still work. Update `gui/editor_test.go` helpers (e.g. `whenLabels`) for the new layout so existing 002 tests pass.

**Checkpoint**: Editor renders two panes; all existing 002 tests pass; create/edit still work.

---

## Phase 3: User Story 1 — Roomier two-pane editor (Priority: P1) 🎯 MVP

**Goal**: ~2× width, fields left, Preview right.

**Independent Test**: Open editor; confirm width and that schedule summary, next runs, and command
preview render in the right pane.

- [ ] T006 [US1] Confirm the right pane shows the Preview content by default and renders schedule summary + next runs + command preview together; adjust spacing/wrapping in `gui/editor.go`.
- [ ] T007 [US1] Add tests in `gui/editor_test.go`: the right pane contains the preview widgets (not the form), and the When form no longer has a "Preview" item.

**Checkpoint**: US1 functional and independently testable.

---

## Phase 4: User Story 2 — In-modal Help (Priority: P1)

**Goal**: Help toggle swaps the right pane to per-field guidance; Examples button removed.

**Independent Test**: Click Help → guidance with examples; toggle back → Preview intact; no Examples
button by Schedule.

- [ ] T008 [US2] Build a Help view in `gui/editor.go` (scrollable `widget.RichText`) with a section + example per field (Name, Command, Arguments, Timezone, Mode, Schedule incl. `starting at` anchor, One-off time, Overlap, Catch-up).
- [ ] T009 [US2] Add a Help toggle (cursor button) that swaps the right-pane holder between Preview and Help without rebuilding the form (preserves input + computed preview, FR-005).
- [ ] T010 [US2] Remove `examplesButton()` and its use beside the Schedule field in `gui/editor.go` (FR-006); fold its content into the Help Schedule section.
- [ ] T011 [US2] Add tests in `gui/editor_test.go`: toggling shows/hides Help; input persists across a Preview→Help→Preview toggle; no Examples button exists in the Schedule row.

**Checkpoint**: US2 functional.

---

## Phase 5: User Story 3 — Code-block command preview (Priority: P2)

**Goal**: Monospace code block, no "Will run:" prefix.

**Independent Test**: Command preview is monospace without the prefix.

- [ ] T012 [US3] Render the command preview via `widget.RichText` with a monospace/code-block segment in `gui/editor.go`; drop the "Will run: " prefix; keep empty-state guidance text. Reuse `commandLinePreview` for the text.
- [ ] T013 [US3] Update tests in `gui/editor_test.go` to assert the command preview text has no "Will run:" prefix and is rendered via the monospace/code-block widget.

**Checkpoint**: US3 functional.

---

## Phase 6: User Story 4 — Disclosure arrow, right footer, cancel-confirm (Priority: P2)

**Goal**: ▶/▼ Advanced Settings, right-aligned Save/Cancel, dirty-guarded Cancel.

**Independent Test**: Arrow direction toggles; buttons right-aligned; dirty Cancel prompts, clean
Cancel closes immediately.

- [ ] T014 [US4] Replace the `widget.Accordion` Advanced Settings with the custom `collapsible` (T002) holding overlap + catch-up, collapsed by default, in `gui/editor.go` (FR-009).
- [ ] T015 [US4] Right-align the footer: replace the non-expanding label spacer with `layout.NewSpacer()` (or a Border trailing slot) so Cancel/Save sit right (FR-010).
- [ ] T016 [US4] Capture a baseline at open and add `isDirty()` comparing current field values to baseline (empty for New, prefilled for Edit; unchanged defaults not dirty) in `gui/editor.go`.
- [ ] T017 [US4] Add `confirmCancel()`: if dirty, show `dialog.NewConfirm("Discard changes?", …)` and only close on confirm; else close immediately. Wire Cancel button (and dialog dismiss/Escape where feasible) through it (FR-011/FR-012).
- [ ] T018 [US4] Add tests in `gui/editor_test.go`: collapsible icon per state; `isDirty()` false on untouched New and on unchanged Edit baseline, true after editing a field; footer right-alignment (button ordering/spacer present).

**Checkpoint**: US4 functional.

---

## Phase 7: User Story 5 — Maximized window on launch (Priority: P3)

**Goal**: Main window opens at the screen work area.

**Independent Test**: Launch → window fills the work area, taskbar visible.

- [ ] T019 [P] [US5] Add `gui/screen_windows.go` (`//go:build windows`): query the work area via `user32.dll!SystemParametersInfoW(SPI_GETWORKAREA)` using stdlib `syscall`, returning width/height in pixels (0,0 on failure).
- [ ] T020 [P] [US5] Add `gui/screen_other.go` (`//go:build !windows`): return `(0,0)` (unknown).
- [ ] T021 [US5] Add a pure `windowSizeFor(workW, workH int, scale float32) fyne.Size` helper (pixels→units minus chrome margin; generous default when work-area is 0) and call it from `NewUI`/`Run` in `gui/app.go` to `Resize` + `CenterOnScreen` at startup (FR-001).
- [ ] T022 [P] [US5] Add `gui/screen_test.go`: unit-test `windowSizeFor` (scale conversion, margin, and fallback when work-area unknown).

**Checkpoint**: US5 functional (verify manually on Windows).

---

## Phase 8: User Story 6 — App-wide pointer cursor (Priority: P3)

**Goal**: Every clickable button shows the pointer cursor.

**Independent Test**: Hover toolbar buttons on each tab + dialog buttons → hand cursor.

- [ ] T023 [P] [US6] Swap `widget.NewButton*` → `newToolbarButton`/no-icon variant in `gui/tasks.go`.
- [ ] T024 [P] [US6] Swap in `gui/groups.go`.
- [ ] T025 [P] [US6] Swap in `gui/triggers.go`.
- [ ] T026 [P] [US6] Swap in `gui/schedule.go`.
- [ ] T027 [P] [US6] Swap in `gui/alerts.go`.
- [ ] T028 [US6] Add a test (e.g. `gui/widgets_test.go` or a tab test) asserting a representative toolbar button is a `cursorButton` / reports `desktop.PointerCursor`. Note the table/row out-of-scope boundary in a comment.

**Checkpoint**: All user stories independently functional.

---

## Phase 9: Polish & Cross-Cutting

- [ ] T029 [P] Update [docs/gui-fields.md](../../docs/gui-fields.md): note the two-pane editor, in-modal Help (replacing Examples), code-block command preview, ▶/▼ Advanced Settings, right-aligned/cancel-confirm footer, and maximized-on-launch.
- [ ] T030 Run `gofmt -l .`, `go vet ./...`, `go test ./gui/...`, the CI `-race` set, and `CGO_ENABLED=0 go build ./...`; fix findings; confirm no 002 regressions.
- [ ] T031 Execute [quickstart.md](quickstart.md) scenarios manually against the windowed GUI (incl. maximized launch and pointer cursors) and record results.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: no deps; T001–T003 parallel.
- **Foundational (Phase 2)**: depends on Setup; **blocks** US1–US4.
- **US1–US4**: share `gui/editor.go` → sequential, in order US1 → US2 → US3 → US4.
- **US5 (window)** and **US6 (cursor swaps)**: independent of `gui/editor.go`; can run in parallel
  with the editor stories (US6 depends on T001's constructors).
- **Polish (Phase 9)**: after all targeted stories.

### Within Each Story

- Tests written with implementation; commit after each task or logical group.
- US4's collapsible swap (T014) depends on T002; US6 swaps depend on T001.

### Parallel Opportunities

- Setup: T001, T002, T003.
- US5: T019, T020, T022 parallel; T021 after T019/T020.
- US6: T023–T027 parallel (different files).
- Polish docs T029 parallel.

---

## Parallel Example: User Story 6

```bash
Task: "Cursor buttons in gui/tasks.go"      # T023
Task: "Cursor buttons in gui/groups.go"     # T024
Task: "Cursor buttons in gui/triggers.go"   # T025
Task: "Cursor buttons in gui/schedule.go"   # T026
Task: "Cursor buttons in gui/alerts.go"     # T027
```

---

## Implementation Strategy

### MVP First (P1 stories)

1. Setup (widgets) → Foundational (two-pane shell).
2. US1 (two-pane preview) → US2 (Help) → validate quickstart 1–3.

### Incremental Delivery

3. US3 (code-block preview) → US4 (arrow/footer/cancel-confirm).
4. US5 (maximize) and US6 (cursor) in parallel.
5. Phase 9 polish, docs, full verification.

---

## Notes

- `[P]` = different files, no incomplete-task dependency. Editor-file tasks are NOT `[P]`.
- No daemon/CLI/API/store changes; keep all 002 behavior intact (FR-014).
- No new module dependency — the window helper uses stdlib `syscall`.
- Run `go test ./gui/...` (CI excludes `/gui` from `-race` due to Fyne font-cache races); the
  pure `windowSizeFor` test also runs under the `-race` set if placed in a non-GUI-driver path.

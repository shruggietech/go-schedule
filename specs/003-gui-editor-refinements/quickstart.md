# Quickstart: Validating GUI Editor Refinements

Run/validation guide. Contracts: [editor-ui.md](contracts/editor-ui.md).

## Prerequisites

- Go + cgo GUI toolchain (WinLibs MinGW GCC on Windows) per `CLAUDE.md`.
- Daemon (`goschedd`) running locally, or the in-process fake backend for headless tests.

## Automated checks

```bash
go test ./gui/...                 # headless: dirty-detect, pane toggle, collapsible, cursor, layout
go test -race ./gui/viewmodel/... # plus the CI -race set excludes /gui (Fyne font-cache races)
gofmt -l . && go vet ./...
go build ./...                    # CGO=0 build must still pass (gui stays cgo-free)
```

Expected: all green; no 002 regressions.

## Manual validation (windowed GUI)

1. **Maximized launch (US5)**: start the GUI → main window fills the work area, taskbar still
   visible, title bar + icon intact.
2. **Two-pane editor (US1)**: open New Task → modal ~2× wider; fields on the left, Preview on the
   right. Type a command + valid schedule → schedule summary, next runs, and the command code block
   all render on the right.
3. **Help (US2)**: click Help → right pane shows per-field guidance + examples; toggle back →
   Preview returns with inputs intact. Confirm there is no Examples button by Schedule.
4. **Command preview (US3)**: command preview is monospace, no "Will run:" prefix.
5. **Controls (US4)**:
   - Advanced Settings shows ▶ collapsed, ▼ expanded.
   - Save/Cancel are right-aligned.
   - Type into a field, click Cancel → confirmation prompt. On a fresh untouched form, Cancel closes
     immediately.
6. **Pointer cursor (US6)**: hover toolbar buttons on Tasks/Schedule/Groups/Triggers/Alerts and the
   dialog buttons → hand cursor everywhere.

## Regression (002)

- Mode toggle hides the irrelevant time field and preserves values; Save gating still blocks empty
  required fields; `every 15 minutes starting at 09:00` still aligns the preview; advanced labels
  still save the correct wire values.

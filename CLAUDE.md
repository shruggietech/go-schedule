<!-- SPECKIT START -->
# go-schedule — Active Plan

Cross-platform (Linux/macOS/Windows) task scheduler in **Go**. Architecture: a system-wide
**daemon** (`goschedd`) hosts the scheduling engine + SQLite store + executor; the **CLI**
(`gosched`) and **Go-native Fyne GUI** (`gosched-gui`) are thin clients over a local IPC API
(Unix socket / Windows named pipe). The GUI is built windowless (`-H windowsgui`) and tasks
spawn with no console window.

Governing documents:
- Constitution: `.specify/memory/constitution.md` (v1.0.0 — code quality, testing, UX, performance)
- Spec: `specs/001-task-scheduler/spec.md`
- Plan: `specs/001-task-scheduler/plan.md`
- Design: `specs/001-task-scheduler/research.md`, `data-model.md`, `contracts/`, `quickstart.md`

Active feature:
- Plan: `specs/004-rebrand-gui-overhaul/plan.md` (rename go-scheduler→go-schedule; Windows .msi install w/ auto-start service; Alerts→unified Logs view w/ filters + on-disk JSONL + detail; remove Triggers entirely (migration v3); real-time GUI via broker task/group/log events (drop manual Refresh); toggleable calendar view under Schedule)
- Prior: `specs/003-gui-editor-refinements/plan.md` (GUI editor refinements: maximized window, two-pane modal + Help, code-block preview, custom collapsible, cancel-confirm, app-wide pointer cursor)
- Prior: `specs/002-gui-task-editor-ux/plan.md` (GUI task-editor UX overhaul + interval anchor)

Key conventions: internal scheduling in UTC; per-task IANA timezone with DST (next-valid /
first-occurrence); recurrence via RFC 5545 RRULE (rrule-go) behind a human-readable layer;
injected `Clock` interface (no direct `time.Now()` in engine code); `log/slog` structured logs;
`go test -race`; dispatch latency p99 < 100ms.
<!-- SPECKIT END -->

# Verification: GUI Dual-Syntax Scheduling

## Chronological evidence

### 2026-08-28 - Baseline

- `go test ./gui ./internal/scheduleinput ./internal/api/server -count=1` passed before production changes (`gui` 5.828s, `scheduleinput` 0.033s, `api/server` 0.145s).
- The starting GUI used human-only `schedule.Parse` for validity and preview, hard-coded `schedule_syntax: human`, prefilled retained expressions, omitted recurring fields for one-offs, and preserved blank expressionless legacy timing on compatible edits.

### 2026-08-28 - Test-first implementation

- Added GUI tests for cron preview/create, exact cron prefill/update, both syntax switches, five-word human classification, invalid cron, and faithful refusal.
- The focused red run failed at the intended old boundaries: supported cron kept Save disabled, preview/update remained `human`, and a DOM/DOW fidelity refusal fell through to the human-parser message. The existing five-word human regression remained green.
- Replaced the two GUI human-only parse paths with one helper over `scheduleinput.Parse`. Preview, create, and update now carry its normalized expression and selected syntax; current field text overrides stale response metadata.
- `go test ./gui ./internal/scheduleinput ./internal/api/server -count=1` passed after implementation (`gui` 8.051s, `scheduleinput` 0.029s, `api/server` 0.138s).
- `scripts/docs-check.sh` passed all 11 pages, links, front matter, fences, and theme contracts after the GUI-local guidance update.
- Spec-Kit prerequisite discovery found all required design artifacts. Both checklists are complete, traceability covers all 16 functional requirements, and the post-design analysis found no critical inconsistency or constitution violation.

### 2026-08-28 - Final quality gates

- `git diff --check` passed.
- Strict UTF-8 decoding, no-BOM, and mojibake scans passed across all 18 changed and untracked paths.
- `scripts/verify.sh all` passed in the foreground with all eight gates: format, vet, lint (0 issues), race, GUI, coverage, docs, and automation.
- Core coverage passed: engine 88.1%, schedule 91.5%, timezone 88.9%, store 87.0%, catchup 87.5%, and logbus 91.1%.
- The docs gate passed 11 pages and the automation gate passed the approved action, CodeQL, and eight-gate manifest contracts.

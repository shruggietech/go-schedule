# Tasks: S103 native desktop notifications and release preparation

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [data-model.md](data-model.md), [contract](contracts/desktop-popups.md)

## Phase 1: Foundations

- [x] T001 Add versioned desktop-local popup preferences, validation, restore behavior, bridge methods, and persistence tests in `desktop/settings` (FR-001, FR-002).
- [x] T002 Add injectable Wails native notification adapter with DOM-ready initialization, contextual authorization, response callback, cleanup, and unsupported status (FR-005, FR-006).
- [x] T003 Add identity-pinned registered-daemon read-only observation with bounded lifecycle and secret-free exact reads (FR-003, FR-004).

## Phase 2: End-to-end behavior

- [x] T004 Implement candidate filtering, live-only subscription boundary, terminal-run and alert mapping, bounded dedupe, and no mutation, with tests for two daemons, reconnect, mute, and exclusions (FR-001 through FR-004).
- [x] T005 Wire activation intents through source-identity-checked Activity drilldown, with tests for switch, stale profile, and missing record (FR-005).
- [x] T006 Build Notifications-page controls and capability explanation for mute, conditions, severity, daemon scope, denial, and closed-app limit; add frontend interaction tests (FR-001, FR-002, FR-006, FR-007).
- [x] T007 Update Windows/macOS/Linux packaged-notification integration as needed, and add platform build/behavior checks (FR-005, FR-006). S103 PR and exact-main CI passed; no additional package change was needed.

## Phase 3: Release and review

- [x] T008 Update user documentation, changelog, v1.5.0 draft release notes, and GitHub planning dependencies without closing SMTP or parent issues (FR-007 through FR-009).
- [x] T009 Run spec-kit analyze and resolve gaps, canonical verification, formatting, targeted tests, and audit issue-level acceptance (all FRs).
- [x] T010 Commit, push, publish structured PR, resolve CI and up to two review rounds, and hand off for merge. PR #260 merged; no tag or public release occurred in S103 (FR-008, FR-009).

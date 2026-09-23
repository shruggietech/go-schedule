# Tasks: Windows local service tray and controls

## Phase 1: Setup

- [x] T001 Confirm issue #255 scope and Windows build/installer conventions in `specs/101-windows-tray-controls/plan.md`
- [x] T002 Record native API and elevation choices in `specs/101-windows-tray-controls/research.md`

## Phase 2: Foundation

- [x] T003 Add lossless read-only SCM state query and tests in `internal/service/status_windows.go` and `internal/service/status_windows_test.go`
- [x] T004 Add shared bounded local-health state projection and tests in `internal/desktopcontrol/`
- [x] T005 Add shared narrow service action orchestration, cancellation, and observed-result tests in `internal/desktopcontrol/`

## Phase 3: User Story 1, truthful tray presence

**Goal**: One icon per session reports real local service state independent of GUI.
**Independent test**: GUI close, service transitions, Explorer restart, duplicate launch, and logoff.

- [x] T006 [US1] Implement windowless per-session companion ownership and native notification-area event loop in `cmd/gosched-tray/`
- [x] T007 [US1] Add approved small-size icon variants and status tooltip/menu in `cmd/gosched-tray/` and `brand/platform/windows/`
- [x] T008 [US1] Add Explorer restart recovery, icon cleanup, and lifecycle tests in `cmd/gosched-tray/`

## Phase 4: User Story 2, service controls

**Goal**: Tray and GUI can perform observed Start, Stop, Restart with narrow UAC and no console flash.
**Independent test**: Admin and standard-user actions, cancel, error, timeout, remote-selected GUI.

- [x] T009 [US2] Add confirm/progress/action menu handling and narrow elevated helper mode in `cmd/gosched-tray/`
- [x] T010 [US2] Expose local service snapshot/actions through Wails bridge in `desktop/app.go`
- [x] T011 [US2] Add always-visible local service UI and interaction tests in `desktop/frontend/src/settings/ConnectionsPage.tsx`
- [x] T012 [US2] Preserve deliberate installed-service stops in `desktop/main.go` and add regression tests

## Phase 5: User Story 3, GUI activation

**Goal**: Open reveals/focuses one GUI and never duplicates it.
**Independent test**: Absent, visible, minimized, hidden GUI and repeated Open.

- [x] T013 [US3] Implement session-scoped GUI singleton and activation listener in `desktop/main.go` and `internal/desktopcontrol/`
- [x] T014 [US3] Wire tray Open to GUI activation or windowless GUI launch in `cmd/gosched-tray/`
- [x] T015 [US3] Add GUI activation regression tests in `desktop/` and `internal/desktopcontrol/`

## Phase 6: Packaging and polish

- [x] T016 Install, register, upgrade, and remove the companion in `build/windows/goschedule.wxs`
- [x] T017 Update MSI contract tests and build staging in `build/windows/verify_wxs.ps1`, `test/integration/windows_installer_contract_test.go`, `.github/workflows/ci.yml`, `.github/workflows/release.yml`, and `scripts/automation-check.sh`
- [x] T018 Document use/elevation and architecture decision in `docs/` and `CHANGELOG.md`
- [x] T019 Run `go run ./scripts/github-format`, `scripts/verify.sh all`, Windows MSI checks, and inspect built artifacts
- [x] T020 Complete independent spec-kit analysis, resolve blockers, and update `specs/101-windows-tray-controls/tasks.md`

## Dependencies and delivery

T003-T005 precede the tray and GUI state surfaces. T006-T008 make Story 1 viable. T009-T012 build on the shared state and satisfy Story 2. T013-T015 satisfy Story 3. T016-T019 close packaging and verification. Windows native code and frontend presentation can proceed independently after the shared contract stabilizes, but the final PR must deliver all three stories and installer lifecycle together.
